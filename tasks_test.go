package tasks

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func testTask(params map[string]string) error {
	log.Println("task params", params)

	return nil
}

func testTaskWithError(params map[string]string) error {
	log.Println("task params", params)

	//nolint:err113
	return errors.New("task error")
}
