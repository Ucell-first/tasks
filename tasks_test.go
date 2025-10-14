package tasks

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"gitlab.local.iti.domain/mc2/golibs/tasks/logger"
	"gitlab.local.iti.domain/mc2/golibs/tasks/mocks"
	"gitlab.local.iti.domain/mc2/golibs/tasks/models"

	"github.com/stretchr/testify/suite"
)

type TasksSuite struct {
	suite.Suite
}

func TestTasksSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(TasksSuite))
}

func (ts *TasksSuite) TestTasks_Create() {
	mockProvider := mocks.New()

	tasker, err := New(WithContext(context.Background()), WithProvider(mockProvider, "test"),
		WithRetryPolicy(models.RetryPolicy{
			BackoffCoefficient: 1.5,
			MaximumAttempts:    5,
		}),
		WithNumWorkers(4),
		WithLogger(logger.DefaultLogger{}),
	)
	ts.Require().NoError(err)

	err = tasker.Start()
	ts.Require().NoError(err)

	err = tasker.RegisterHandler("test", testTask)
	ts.Require().NoError(err)

	err = tasker.RegisterHandler("test_error", testTaskWithError)
	ts.Require().NoError(err)

	err = tasker.Create(context.Background(), "test", map[string]string{
		"data": "dummy data",
	})

	ts.Require().NoError(err)

	ts.Run("Tasks with retry", func() {
		for i := 0; i < 100; i++ {
			err = tasker.Create(context.Background(), "test_error", map[string]string{
				"data": fmt.Sprintf("dummy task-%d", i+1),
			})

			ts.Require().NoError(err)
		}
	})

	time.Sleep(5 * time.Second)
	tasker.Stop()
}

func (ts *TasksSuite) TestTasks_CreateScheduled() {
	mockProvider := mocks.New()

	tasker, err := New(WithContext(context.Background()), WithProvider(mockProvider, "test"),
		WithRetryPolicy(models.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumAttempts:    5,
		}),
		WithNumWorkers(4),
		WithQueueSize(10),
		WithLogger(logger.DefaultLogger{}),
	)
	ts.Require().NoError(err)

	err = tasker.Start()
	ts.Require().NoError(err)

	err = tasker.RegisterHandler("test", testTask)
	ts.Require().NoError(err)

	now := time.Now().UTC()
	tz, _ := time.LoadLocation("Asia/Tashkent")
	timeTask := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, tz)

	ts.Run("OK", func() {
		err = tasker.CreateScheduled(context.Background(), "test", nil, time.Now(), time.Second)
		ts.Require().NoError(err)
	})

	ts.Run("Error", func() {
		err = tasker.CreateScheduled(context.Background(), "test_err", nil, timeTask, time.Hour)
		ts.Require().NoError(err)

		err = tasker.CreateScheduled(context.Background(), "test_err", nil, timeTask, time.Minute)
		ts.Require().Error(err)
	})

	time.Sleep(defaultScheduledTaskDuration)
}

func (ts *TasksSuite) TestTasks_CreateDelayed() {
	mockProvider := mocks.New()

	var executedTasks []string
	var mu sync.Mutex

	tasker, err := New(
		WithContext(context.Background()),
		WithProvider(mockProvider, "test"),
		WithRetryPolicy(models.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumAttempts:    3,
		}),
		WithNumWorkers(2),
		WithQueueSize(10),
		WithLogger(logger.DefaultLogger{}),
	)
	ts.Require().NoError(err)

	// Register handler that tracks execution
	err = tasker.RegisterHandler("delayed_test", func(params map[string]string) error {
		mu.Lock()
		defer mu.Unlock()
		executedTasks = append(executedTasks, params["task_id"])
		log.Printf("Delayed task executed: %s at %s", params["task_id"], time.Now().Format(time.RFC3339))
		return nil
	})
	ts.Require().NoError(err)

	err = tasker.Start()
	ts.Require().NoError(err)

	ts.Run("Success - Single delayed task", func() {
		executedTasks = []string{}
		startTime := time.Now().UTC().Add(2 * time.Second)

		err = tasker.CreateDelayed(
			context.Background(),
			"localhost",
			"delayed_test",
			map[string]string{"task_id": "task_1"},
			startTime,
		)
		ts.Require().NoError(err)

		// Wait for task to execute
		time.Sleep(3 * time.Second)

		mu.Lock()
		ts.Require().Len(executedTasks, 1)
		ts.Require().Equal("task_1", executedTasks[0])
		mu.Unlock()
	})

	ts.Run("Success - Multiple delayed tasks", func() {
		executedTasks = []string{}
		now := time.Now().UTC()

		// Create 3 tasks with different delays
		err = tasker.CreateDelayed(
			context.Background(),
			"localhost",
			"delayed_test",
			map[string]string{"task_id": "task_A"},
			now.Add(3*time.Second),
		)
		ts.Require().NoError(err)

		err = tasker.CreateDelayed(
			context.Background(),
			"localhost",
			"delayed_test",
			map[string]string{"task_id": "task_B"},
			now.Add(1*time.Second),
		)
		ts.Require().NoError(err)

		err = tasker.CreateDelayed(
			context.Background(),
			"localhost",
			"delayed_test",
			map[string]string{"task_id": "task_C"},
			now.Add(2*time.Second),
		)
		ts.Require().NoError(err)

		// Wait for all tasks to execute
		time.Sleep(4 * time.Second)

		mu.Lock()
		ts.Require().Len(executedTasks, 3)
		// Verify execution order (should be B, C, A)
		ts.Require().Equal("task_B", executedTasks[0])
		ts.Require().Equal("task_C", executedTasks[1])
		ts.Require().Equal("task_A", executedTasks[2])
		mu.Unlock()
	})

	ts.Run("Error - Past time", func() {
		pastTime := time.Now().UTC().Add(-1 * time.Hour)

		err = tasker.CreateDelayed(
			context.Background(),
			"localhost",
			"delayed_test",
			map[string]string{"task_id": "task_past"},
			pastTime,
		)
		ts.Require().Error(err)
		ts.Require().ErrorIs(err, ErrCreateDelayed)
	})

	ts.Run("Success - With host parameter", func() {
		executedTasks = []string{}
		startTime := time.Now().UTC().Add(1 * time.Second)

		err = tasker.CreateDelayed(
			context.Background(),
			"host-server-01",
			"delayed_test",
			map[string]string{"task_id": "task_with_host"},
			startTime,
		)
		ts.Require().NoError(err)

		time.Sleep(2 * time.Second)

		mu.Lock()
		ts.Require().Len(executedTasks, 1)
		ts.Require().Equal("task_with_host", executedTasks[0])
		mu.Unlock()
	})

	ts.Run("Success - Delayed task with retry on failure", func() {
		executedTasks = []string{}
		attemptCount := 0

		// Register handler that fails first 2 times
		err = tasker.RegisterHandler("delayed_test_retry", func(params map[string]string) error {
			mu.Lock()
			attemptCount++
			currentAttempt := attemptCount
			mu.Unlock()

			log.Printf("Attempt %d for task: %s", currentAttempt, params["task_id"])

			if currentAttempt < 3 {
				return errors.New("simulated failure")
			}

			mu.Lock()
			executedTasks = append(executedTasks, params["task_id"])
			mu.Unlock()

			return nil
		})
		ts.Require().NoError(err)

		startTime := time.Now().UTC().Add(1 * time.Second)

		err = tasker.CreateDelayed(
			context.Background(),
			"localhost",
			"delayed_test_retry",
			map[string]string{"task_id": "retry_task"},
			startTime,
		)
		ts.Require().NoError(err)

		// Wait for task execution and retries
		time.Sleep(10 * time.Second)

		mu.Lock()
		ts.Require().GreaterOrEqual(attemptCount, 3)
		ts.Require().Len(executedTasks, 1)
		mu.Unlock()
	})

	time.Sleep(1 * time.Second)
	tasker.Stop()
}

func testTask(params map[string]string) error {
	log.Println("task params", params)

	return nil
}

func testTaskWithError(params map[string]string) error {
	log.Println("task params", params)

	//nolint:err113
	return errors.New("task error")
}
