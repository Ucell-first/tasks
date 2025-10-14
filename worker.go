//tasks/worker.go

package tasks

import (
	"container/heap"
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"gitlab.local.iti.domain/mc2/golibs/tasks/logger"
	"gitlab.local.iti.domain/mc2/golibs/tasks/models"
)

const (
	defaultScheduledTaskDuration = 5 * time.Second
	minTimerDuration             = 100 * time.Millisecond
)

func (t *Tasks) startWorkers(ctx context.Context) {
	// Start regular task workers
	for i := 0; i < t.opts.numWorkers; i++ {
		t.wg.Add(1)
		go t.taskWorker(ctx, i+1, t.taskQueue)
	}

	// Start retry worker
	t.wgRetry.Add(1)
	go t.retryTaskWorker(ctx)

	// Start delayed task worker
	t.wgDelayed.Add(1)
	go t.delayedTaskWorker(ctx)

	// Start scheduled task worker
	go t.scheduledTaskWorker(ctx)
}

func (t *Tasks) taskWorker(ctx context.Context, workerID int, taskQueue <-chan models.Task) {
	defer t.wg.Done()

	for {
		select {
		case <-ctx.Done():
			t.opts.logger.Logf(logger.LogLevelInfo, "worker %d shutting down", nil, workerID)
			return
		case task, ok := <-taskQueue:
			if !ok {
				t.opts.logger.Logf(logger.LogLevelInfo, "worker %d: task queue closed", nil, workerID)
				return
			}

			if err := t.processTask(ctx, task); err != nil {
				t.opts.logger.Logf(logger.LogLevelError, "worker %d: processTask error: %s",
					map[string]interface{}{"task_name": task.Name}, workerID, err.Error())
			}
		}
	}
}

func (t *Tasks) processTask(ctx context.Context, task models.Task) error {
	t.tasksHandlersMutex.RLock()
	handler, ok := t.tasksHandlers[task.Name]
	t.tasksHandlersMutex.RUnlock()

	if !ok {
		return fmt.Errorf("%w: %w: task_name=%s", errProcessTask, ErrTaskNameNotRegistered, task.Name)
	}

	isScheduled := task.Params["scheduled"] == "true"
	isDelayed := task.Params["delayed"] == "true"

	t.opts.logger.Logf(logger.LogLevelInfo, "processing task: %s",
		map[string]interface{}{
			"task_name": task.Name,
			"scheduled": isScheduled,
			"delayed":   isDelayed,
		}, task.Name)

	if err := handler(task.Params); err != nil {
		// Don't retry scheduled tasks, they will run again on schedule
		if !isScheduled {
			t.addToRetryQueue(task)
		}

		return fmt.Errorf("%w: %w", errProcessTask, err)
	}

	return nil
}

// retryTaskWorker processes tasks that need to be retried with exponential backoff.
// It uses a priority queue (heap) to efficiently manage tasks by their retry time.
func (t *Tasks) retryTaskWorker(ctx context.Context) {
	defer t.wgRetry.Done()

	retryQueue := &RetryQueue{}
	heap.Init(retryQueue)

	timer := time.NewTimer(time.Hour) // Start with a long duration
	timer.Stop()                      // Stop immediately, will be reset when needed

	defer timer.Stop()

	t.opts.logger.Log(logger.LogLevelInfo, "retry worker started", nil)

	for {
		select {
		case <-ctx.Done():
			t.opts.logger.Logf(logger.LogLevelInfo, "retry worker shutting down, pending tasks: %d",
				nil, retryQueue.Len())
			return

		case <-timer.C:
			// Timer fired, process ready tasks
			now := time.Now().UTC()

			for retryQueue.Len() > 0 {
				nextTask := (*retryQueue)[0]

				if nextTask.StartTime.After(now) {
					// Next task is not ready yet
					duration := nextTask.StartTime.Sub(now)
					if duration < minTimerDuration {
						duration = minTimerDuration
					}
					timer.Reset(duration)
					break
				}

				// Pop and process the task
				task := heap.Pop(retryQueue).(*RetryTask)

				t.opts.logger.Logf(logger.LogLevelInfo, "retrying task: %s (attempt %s)",
					map[string]interface{}{
						"task_name": task.Task.Name,
						"attempts":  task.Task.Params["attempts"],
					},
					task.Task.Name, task.Task.Params["attempts"])

				err := t.Create(ctx, task.Task.Name, task.Task.Params)
				if err != nil {
					t.opts.logger.Logf(logger.LogLevelError, "retry task create error: %s",
						map[string]interface{}{"task_name": task.Task.Name}, err.Error())
				}
			}

			// If queue is empty, don't reset timer
			if retryQueue.Len() == 0 {
				t.opts.logger.Log(logger.LogLevelDebug, "retry queue is empty", nil)
			}

		case task, ok := <-t.retryQueue:
			if !ok {
				// Channel closed, process remaining tasks and exit
				t.opts.logger.Logf(logger.LogLevelInfo, "retry queue channel closed, processing remaining %d tasks",
					nil, retryQueue.Len())

				for retryQueue.Len() > 0 {
					retryTask := heap.Pop(retryQueue).(*RetryTask)
					err := t.Create(ctx, retryTask.Task.Name, retryTask.Task.Params)
					if err != nil {
						t.opts.logger.Logf(logger.LogLevelError, "final retry task create error: %s", nil, err.Error())
					}
				}
				return
			}

			// Add new task to retry queue
			retryTask := &RetryTask{
				Task:      task,
				StartTime: task.StartTime,
			}

			heap.Push(retryQueue, retryTask)

			t.opts.logger.Logf(logger.LogLevelDebug, "added task to retry queue: %s at %s",
				map[string]interface{}{
					"task_name": task.Name,
					"retry_at":  task.StartTime.Format(time.RFC3339),
				},
				task.Name, task.StartTime.Format(time.RFC3339))

			// If this is the next task to process, reset timer
			if retryQueue.Len() == 1 || (*retryQueue)[0].StartTime.Equal(task.StartTime) {
				now := time.Now().UTC()
				duration := task.StartTime.Sub(now)

				if duration < 0 {
					duration = 0
				} else if duration < minTimerDuration {
					duration = minTimerDuration
				}

				timer.Reset(duration)
			}
		}
	}
}

// delayedTaskWorker processes tasks that are scheduled to run at a specific future time.
// Similar to retryTaskWorker but for delayed tasks.
func (t *Tasks) delayedTaskWorker(ctx context.Context) {
	defer t.wgDelayed.Done()

	delayedQueue := &RetryQueue{}
	heap.Init(delayedQueue)

	timer := time.NewTimer(time.Hour)
	timer.Stop()

	defer timer.Stop()

	t.opts.logger.Log(logger.LogLevelInfo, "delayed task worker started", nil)

	for {
		select {
		case <-ctx.Done():
			t.opts.logger.Logf(logger.LogLevelInfo, "delayed task worker shutting down, pending tasks: %d",
				nil, delayedQueue.Len())
			return

		case <-timer.C:
			now := time.Now().UTC()

			for delayedQueue.Len() > 0 {
				nextTask := (*delayedQueue)[0]

				if nextTask.StartTime.After(now) {
					duration := nextTask.StartTime.Sub(now)
					if duration < minTimerDuration {
						duration = minTimerDuration
					}
					timer.Reset(duration)
					break
				}

				task := heap.Pop(delayedQueue).(*RetryTask)

				t.opts.logger.Logf(logger.LogLevelInfo, "executing delayed task: %s",
					map[string]interface{}{
						"task_name":     task.Task.Name,
						"scheduled_for": task.StartTime.Format(time.RFC3339),
					},
					task.Task.Name)

				// Remove delayed flag before processing
				delete(task.Task.Params, "delayed")

				err := t.Create(ctx, task.Task.Name, task.Task.Params)
				if err != nil {
					t.opts.logger.Logf(logger.LogLevelError, "delayed task create error: %s",
						map[string]interface{}{"task_name": task.Task.Name}, err.Error())
				}
			}

			if delayedQueue.Len() == 0 {
				t.opts.logger.Log(logger.LogLevelDebug, "delayed queue is empty", nil)
			}

		case task, ok := <-t.delayedQueue:
			if !ok {
				t.opts.logger.Logf(logger.LogLevelInfo, "delayed queue channel closed, processing remaining %d tasks",
					nil, delayedQueue.Len())

				for delayedQueue.Len() > 0 {
					delayedTask := heap.Pop(delayedQueue).(*RetryTask)
					delete(delayedTask.Task.Params, "delayed")
					err := t.Create(ctx, delayedTask.Task.Name, delayedTask.Task.Params)
					if err != nil {
						t.opts.logger.Logf(logger.LogLevelError, "final delayed task create error: %s", nil, err.Error())
					}
				}
				return
			}

			delayedTask := &RetryTask{
				Task:      task,
				StartTime: task.StartTime,
			}

			heap.Push(delayedQueue, delayedTask)

			t.opts.logger.Logf(logger.LogLevelDebug, "added task to delayed queue: %s at %s",
				map[string]interface{}{
					"task_name":  task.Name,
					"execute_at": task.StartTime.Format(time.RFC3339),
				},
				task.Name, task.StartTime.Format(time.RFC3339))

			if delayedQueue.Len() == 1 || (*delayedQueue)[0].StartTime.Equal(task.StartTime) {
				now := time.Now().UTC()
				duration := task.StartTime.Sub(now)

				if duration < 0 {
					duration = 0
				} else if duration < minTimerDuration {
					duration = minTimerDuration
				}

				timer.Reset(duration)
			}
		}
	}
}

func (t *Tasks) scheduledTaskWorker(ctx context.Context) {
	ticker := time.NewTicker(defaultScheduledTaskDuration)
	defer ticker.Stop()

	t.opts.logger.Log(logger.LogLevelInfo, "scheduled task worker started", nil)

	for {
		select {
		case <-ctx.Done():
			t.opts.logger.Log(logger.LogLevelInfo, "scheduled task worker shutting down", nil)
			return
		case <-ticker.C:
			t.processScheduledTasks(ctx)
		}
	}
}

func (t *Tasks) processScheduledTasks(ctx context.Context) {
	t.scheduledTaskMutex.Lock()
	defer t.scheduledTaskMutex.Unlock()

	now := time.Now().UTC()

	for name, task := range t.scheduledTasks {
		if task.TimeOfNextExec.Before(now) || task.TimeOfNextExec.Equal(now) {
			t.opts.logger.Logf(logger.LogLevelDebug, "executing scheduled task: %s",
				map[string]interface{}{
					"task_name": task.Name,
					"next_exec": task.TimeOfNextExec.Format(time.RFC3339),
				},
				task.Name)

			task.TimeOfNextExec = task.TimeOfNextExec.Add(task.Period)
			t.scheduledTasks[name] = task

			err := t.Create(ctx, task.Name, task.Params)
			if err != nil {
				t.opts.logger.Logf(logger.LogLevelError, "scheduled task create error: %s",
					map[string]interface{}{"task_name": task.Name}, err.Error())
			}
		}
	}
}

func (t *Tasks) addToRetryQueue(task models.Task) {
	attemptsStr := task.Params["attempts"]
	attempts, _ := strconv.Atoi(attemptsStr)
	maxAttempts := t.opts.retryPolicy.MaximumAttempts

	if maxAttempts > 0 && attempts >= maxAttempts {
		t.opts.logger.Logf(logger.LogLevelInfo, "max retry attempts exceeded for task: %s (attempts: %d)",
			map[string]interface{}{
				"task_name":    task.Name,
				"attempts":     attempts,
				"max_attempts": maxAttempts,
			},
			task.Name, attempts)
		return
	}

	attempts++
	backoff := t.calculateBackoff(attempts)

	task.Params["attempts"] = strconv.Itoa(attempts)
	task.StartTime = time.Now().UTC().Add(backoff)

	t.opts.logger.Logf(logger.LogLevelInfo, "scheduling retry for task: %s (attempt %d, delay: %s)",
		map[string]interface{}{
			"task_name": task.Name,
			"attempts":  attempts,
			"backoff":   backoff.String(),
		},
		task.Name, attempts, backoff.String())

	select {
	case t.retryQueue <- task:
		// Successfully added to retry queue
	default:
		t.opts.logger.Logf(logger.LogLevelError, "retry queue is full, dropping task: %s",
			map[string]interface{}{"task_name": task.Name}, task.Name)
	}
}

// calculateBackoff computes the retry delay for exponential backoff.
func (t *Tasks) calculateBackoff(attempts int) time.Duration {
	retryPolicy := t.opts.retryPolicy

	backoff := float64(retryPolicy.InitialInterval) * math.Pow(retryPolicy.BackoffCoefficient, float64(attempts-1))

	if backoff > float64(retryPolicy.MaximumInterval) {
		backoff = float64(retryPolicy.MaximumInterval)
	}

	return time.Duration(backoff)
}
