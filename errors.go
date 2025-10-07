package tasks

import "errors"

var (
	ErrTasks = errors.New("tasks")
	// ErrRegisterHandler указывает на возникновение ошибки при регистрации обработчика задачи.
	ErrRegisterHandler = errors.New("RegisterHandler method")

	// ErrTaskNameAlreadyRegistered указывает на присутствие ранее зарегистрированного обработчика задачи.
	ErrTaskNameAlreadyRegistered = errors.New("task name already registered")

	// ErrCreate указывает на возникновение ошибки при попытке создать задачу на обработку.
	ErrCreate = errors.New("Create method")

	ErrUnknownProvider = errors.New("unknown provider")
	ErrUnknownContext  = errors.New("unknown context")

	ErrTaskNameNotRegistered = errors.New("task name not registered")
	ErrCreateScheduled       = errors.New("CreateScheduled method")

	ErrEmptyTopic = errors.New("empty topic")

	errProcessTask = errors.New("processTask method")

	errHandler = errors.New("handleTask method")
)
