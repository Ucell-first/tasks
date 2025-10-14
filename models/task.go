package models

import "time"

// Task это структура данных о задаче.
type Task struct {
	StartTime      time.Time         `json:"start_time"`
	TimeOfNextExec time.Time         `json:"-"`
	Params         map[string]string `json:"params"`
	Name           string            `json:"name"`
	Host           string            `json:"host,omitempty"`
	Period         time.Duration     `json:"-"`
}

type RetryPolicy struct {
	//  Amount of time that must elapse before the first retry occurs. default value 1 second.
	InitialInterval time.Duration
	// How much the retry interval increases. default value is 2.0
	BackoffCoefficient float64
	// Specifies the maximum interval between retries. default  300 × Initial Interval
	MaximumInterval time.Duration
	// Specifies the maximum number of execution attempts. 0  means unlimited.
	MaximumAttempts int
}
