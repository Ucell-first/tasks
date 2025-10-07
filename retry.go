package tasks

import (
	"time"

	"gitlab.local.iti.domain/mc2/golibs/tasks/models"
)

// RetryTask represents a task in the retry queue.
type RetryTask struct {
	StartTime time.Time
	Task      models.Task
}

// RetryQueue implements heap.Interface for managing retry tasks.
//
//nolint:recvcheck
type RetryQueue []*RetryTask

func (rq RetryQueue) Len() int           { return len(rq) }
func (rq RetryQueue) Less(i, j int) bool { return rq[i].StartTime.Before(rq[j].StartTime) }
func (rq RetryQueue) Swap(i, j int)      { rq[i], rq[j] = rq[j], rq[i] }

//nolint:forcetypeassert
func (rq *RetryQueue) Push(x any) {
	*rq = append(*rq, x.(*RetryTask))
}

func (rq *RetryQueue) Pop() any {
	old := *rq
	n := len(old)
	task := old[n-1]
	*rq = old[:n-1]

	return task
}
