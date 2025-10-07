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
)

func (t *Tasks) startWorkers(ctx context.Context) {
	for i := 0; i < t.opts.numWorkers; i++ {
		t.wg.Add(1)

		go t.taskWorker(ctx, i+1, t.taskQueue)
	}

	t.wgRetry.Add(1)

	go t.retryTaskWorker(ctx, t.retryQueue)

	go t.scheduledTaskWorker(ctx)
}

func (t *Tasks) taskWorker(ctx context.Context, _ int, taskQueue <-chan models.Task) {
	defer t.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-taskQueue:
			if !ok {
				return
			}

			if err := t.processTask(ctx, task); err != nil {
				t.opts.logger.Logf(logger.LogLevelError, "processTask error:%s", nil, err.Error())
			}
		}
	}
}

func (t *Tasks) processTask(_ context.Context, task models.Task) error {
	t.tasksHandlersMutex.RLock()
	handler, ok := t.tasksHandlers[task.Name]
	t.tasksHandlersMutex.RUnlock()

	if !ok {
		return fmt.Errorf("%w: %w", errProcessTask, ErrTaskNameNotRegistered)
	}

	isScheduled := task.Params["scheduled"] == "true"

	if err := handler(task.Params); err != nil {
		if !isScheduled {
			t.addToRetryQueue(task)
		}

		return fmt.Errorf("%w: %w", errProcessTask, err)
	}

	return nil
}

func (t *Tasks) retryTaskWorker(ctx context.Context, queueTasks <-chan models.Task) {
	defer t.wgRetry.Done()

	retryQueue := &RetryQueue{}
	heap.Init(retryQueue)

	timer := time.NewTimer(0) // Initialize with 0 to trigger immediately
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if retryQueue.Len() > 0 {
				//nolint:forcetypeassert
				nextTask := heap.Pop(retryQueue).(*RetryTask)

				err := t.Create(ctx, nextTask.Task.Name, nextTask.Task.Params)
				if err != nil {
					t.opts.logger.Logf(logger.LogLevelError, "retry task create :%s",
						nil, err.Error())
				}

				// Schedule the next task
				if retryQueue.Len() > 0 {
					nextStartTime := (*retryQueue)[0].StartTime
					timer.Reset(nextStartTime.Sub(time.Now().UTC()))
				}
			}

		case task, ok := <-queueTasks:
			if !ok {
				// Channel is closed, exit if the queue is empty
				if retryQueue.Len() == 0 {
					return
				} else {
					// Otherwise, process remaining tasks
					timer.Reset(0)

					continue
				}
			}

			// Add the new task to the retry queue
			heap.Push(retryQueue, &RetryTask{
				Task:      task,
				StartTime: task.StartTime,
			})

			// Reset the timer if this task is the next one to process
			if (*retryQueue)[0].StartTime.Equal(task.StartTime) {
				timer.Reset(task.StartTime.Sub(time.Now().UTC()))
			}
		}
	}
}

func (t *Tasks) scheduledTaskWorker(ctx context.Context) {
	ticker := time.NewTicker(defaultScheduledTaskDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.processScheduledTasks(ctx)
		}
	}
}

func (t *Tasks) processScheduledTasks(ctx context.Context) {
	t.scheduledTaskMutex.Lock()
	defer t.scheduledTaskMutex.Unlock()

	for name, task := range t.scheduledTasks {
		if task.TimeOfNextExec.Before(time.Now().UTC()) {
			task.TimeOfNextExec = task.TimeOfNextExec.Add(task.Period)
			t.scheduledTasks[name] = task

			err := t.Create(ctx, task.Name, task.Params)
			if err != nil {
				t.opts.logger.Logf(logger.LogLevelError, "scheduled task create :%s", nil, err.Error())
			}
		}
	}
}

func (t *Tasks) addToRetryQueue(task models.Task) {
	attempts, _ := strconv.Atoi(task.Params["attempts"])
	maxAttempts := t.opts.retryPolicy.MaximumAttempts

	if maxAttempts > 0 && attempts > maxAttempts {
		t.opts.logger.Logf(logger.LogLevelInfo, "max attempts exceeded, attempts:%d", nil, attempts)

		return
	}

	attempts++
	backoff := t.calculateBackoff(attempts)

	task.Params["attempts"] = strconv.Itoa(attempts)
	task.StartTime = time.Now().UTC().Add(backoff)

	t.retryQueue <- task
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
