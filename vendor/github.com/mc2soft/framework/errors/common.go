package errors

import errs "errors"

var (
	// ErrInvalidProvider returns only when subsystem handler failed to "cast"
	// registered provider into needed interface to perform subsystem-specific actions.
	ErrInvalidProvider = errs.New("invalid provider")

	// ErrNoProvidersRegistered appears when trying to do some modification
	// to provider (like registering handlers or setting configuration) but
	// no providers was registered.
	ErrNoProvidersRegistered = errs.New("no providers was registered")

	// ErrProviderAlreadyRegistered appears when trying to register a
	// provider with already taken name.
	ErrProviderAlreadyRegistered = errs.New("provider with such name already registered")

	// ErrProviderNotFound appears when trying to do something with a
	// provider which wasn't yet registered.
	ErrProviderNotFound = errs.New("provider wasn't found")

	// ErrShutdownInProgress appears when trying to do something but shutdown
	// procedure was initiated.
	ErrShutdownInProgress = errs.New("shutdown in progress")
)
