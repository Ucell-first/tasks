//tasks/tasks.go

package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mc2soft/framework/communication"
	defaultrequest "gitlab.local.iti.domain/mc2/golibs/legacy-framework-request"
	"gitlab.local.iti.domain/mc2/golibs/tasks/logger"
	"gitlab.local.iti.domain/mc2/golibs/tasks/models"
)

const (
	defaultQueueSize   = 100
	defaultMaxInterval = 300
)

// TaskHandler handleTask func.
type TaskHandler func(params map[string]string) error

// Tasker is an interface for tasks.
type Tasker interface {
	RegisterHandler(taskName string, handler TaskHandler) error
	Create(ctx context.Context, taskName string, params map[string]string) error
	CreateScheduled(ctx context.Context, taskName string, params map[string]string,
		startAt time.Time, period time.Duration) error
	Start() error
	Stop()
}

type Tasks struct {
	provider           communication.Provider
	tasksHandlers      map[string]TaskHandler
	scheduledTasks     map[string]models.Task
	taskQueue          chan models.Task
	retryQueue         chan models.Task
	opts               *options
	wg                 sync.WaitGroup
	wgRetry            sync.WaitGroup
	tasksHandlersMutex sync.RWMutex
	scheduledTaskMutex sync.RWMutex
	AreConsumersActive atomic.Bool
}

func New(opts ...Option) (Tasker, error) {
	tasks := &Tasks{opts: &options{}}

	for _, opt := range opts {
		opt.apply(tasks.opts)
	}

	if err := tasks.initialize(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrTasks, err)
	}

	return tasks, nil
}

//nolint:funcorder
func (t *Tasks) initialize() error {
	if t.opts.numWorkers == 0 {
		t.opts.numWorkers = runtime.NumCPU()
	}

	if t.opts.queueSize == 0 {
		t.opts.queueSize = defaultQueueSize
	}

	if t.opts.logger == nil {
		t.opts.logger = new(logger.DefaultLogger)
	}

	if t.opts.retryPolicy.InitialInterval == 0 {
		t.opts.retryPolicy.InitialInterval = time.Second
	}

	if t.opts.retryPolicy.BackoffCoefficient == 0 {
		t.opts.retryPolicy.BackoffCoefficient = 2.0
	}

	if t.opts.retryPolicy.MaximumInterval == 0 {
		t.opts.retryPolicy.MaximumInterval = defaultMaxInterval * t.opts.retryPolicy.InitialInterval
	}

	if t.opts.provider == nil {
		return fmt.Errorf("initialization: %w", ErrUnknownProvider)
	}

	if t.opts.topic == "" {
		return fmt.Errorf("initialization: %w", ErrEmptyTopic)
	}

	t.provider = t.opts.provider

	if t.opts.ctx == nil {
		return fmt.Errorf("initialization: %w", ErrUnknownContext)
	}

	t.provider.RegisterHandlerNamingFunc(func(_, _, path string) string {
		return path
	})

	t.provider.RegisterDefaultRequestStruct(&defaultrequest.DefaultRequest{})

	t.tasksHandlers = make(map[string]TaskHandler)
	t.scheduledTasks = make(map[string]models.Task)
	t.taskQueue = make(chan models.Task, t.opts.queueSize)
	t.retryQueue = make(chan models.Task, t.opts.queueSize)

	return nil
}

func (t *Tasks) Start() error {
	err := t.provider.RegisterHandler("", t.opts.topic, t.handleTask)
	if err != nil {
		return fmt.Errorf("initialization: %w", err)
	}

	t.startWorkers(t.opts.ctx)
	t.AreConsumersActive.Store(true)

	return nil
}

func (t *Tasks) RegisterHandler(taskName string, handler TaskHandler) error {
	t.tasksHandlersMutex.Lock()
	defer t.tasksHandlersMutex.Unlock()

	_, ok := t.tasksHandlers[taskName]
	if ok {
		return fmt.Errorf("%w: %w: %s", ErrRegisterHandler, ErrTaskNameAlreadyRegistered, taskName)
	}

	t.tasksHandlers[taskName] = handler

	return nil
}

func (t *Tasks) Create(ctx context.Context, taskName string, params map[string]string) error {
	task := models.Task{
		Name:      taskName,
		Params:    params,
		StartTime: time.Now().UTC(),
	}

	if task.Params == nil {
		task.Params = map[string]string{}
	}

	taskRaw, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("%w :%w", ErrCreate, err)
	}

	event := defaultrequest.New(
		ctx,
		"",
		t.opts.topic,
		nil,
		taskRaw,
	)

	err = t.provider.Send(event)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreate, err)
	}

	return nil
}

func (t *Tasks) CreateScheduled(
	_ context.Context,
	taskName string,
	params map[string]string,
	startAt time.Time,
	period time.Duration,
) error {
	task := models.Task{
		Name:           taskName,
		Params:         params,
		Period:         period,
		TimeOfNextExec: startAt.Add(period),
	}

	if task.Params == nil {
		task.Params = map[string]string{}
	}

	task.Params["scheduled"] = "true"

	t.scheduledTaskMutex.Lock()
	defer t.scheduledTaskMutex.Unlock()

	_, ok := t.scheduledTasks[taskName]
	if ok {
		return fmt.Errorf("%w: %w: %s", ErrCreateScheduled, ErrTaskNameAlreadyRegistered, taskName)
	}

	t.scheduledTasks[taskName] = task

	return nil
}

// Stop stop all tasks processing.
func (t *Tasks) Stop() {
	// @TODO в будущем если в интерфейс provider будет добавлятся
	// метод типа UnRegisterHandler()  stop subscriptions kafka, то не нужно AreConsumersActive
	t.AreConsumersActive.Store(false)

	t.waitForTaskQueueFree(t.opts.ctx)

	close(t.taskQueue)

	t.wg.Wait()

	close(t.retryQueue)

	t.wgRetry.Wait()
}

func (t *Tasks) waitForTaskQueueFree(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.opts.logger.Logf(logger.LogLevelInfo, "waiting for task queue to free: %d", nil,
				len(t.taskQueue))

			if len(t.taskQueue) == 0 {
				return
			}
		}
	}
}
