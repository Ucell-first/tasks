package errors

import errs "errors"

var (
	// ErrTasksQueueChannelIsNil appears when worker requests channel that is
	// nil'd. This can happen if code requests specific worker's channel without
	// starting task itself.
	ErrTasksQueueChannelIsNil = errs.New("task worker incoming channel is nil")

	// ErrTasksQueueHandlerAlreadyRegistered appears when trying to register queue
	// handler with name that is already used.
	ErrTasksQueueHandlerAlreadyRegistered = errs.New("queue handler for passed queue name already registered")

	// ErrTasksQueueNotFound appears when passed queue name wasn't previously registered.
	ErrTasksQueueNotFound = errs.New("queue not found")

	// ErrTaskQueuesNotSet appears when trying to start task for execution but no queues
	// was registered
	ErrTaskQueuesNotSet = errs.New("task queues not registered")

	// ErrTasksWorkerNotSpecified appears when requesting worker-specific things like channel
	// for worker #3, but we're get out of bounds error.
	ErrTasksWorkerNotSpecified = errs.New("worker not specified")
)
