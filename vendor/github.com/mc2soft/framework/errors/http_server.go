package errors

import errs "errors"

var (

	// ErrConfigurationIsNotForEcho appears when passed pointer to configuration
	// structure is not for Echo HTTP server provider.
	ErrConfigurationIsNotForEcho = errs.New("passed configuration is not for Echo HTTP server provider")

	// ErrInvalidHTTPServerListeningAddress appears when setting configuration for
	// HTTP server with invalid string as listening address (not in "host:port" form).
	ErrInvalidHTTPServerListeningAddress = errs.New("invalid HTTP server listening address")

	// ErrInvalidHTTPServerWaitForSeconds appears when setting configuration for
	// HTTP server with WaitForSeconds parameter < 1.
	ErrInvalidHTTPServerWaitForSeconds = errs.New("invalid HTTP server WaitForSeconds parameter value")

	// ErrNoHTTPServerListeningAddressDefined appears when setting configuration
	// for HTTP server with empty string as listening address.
	ErrNoHTTPServerListeningAddressDefined = errs.New("empty HTTP server listening address")

	// ErrUnableToStartHTTPServer appears when HTTP server cannot be started
	// or if post-start availability check failed.
	ErrUnableToStartHTTPServer = errs.New("unable to start HTTP server")

	// ErrUnknownHTTPMethod appears when trying to register handler for unknown
	// HTTP method (verb).
	ErrUnknownHTTPMethod = errs.New("unknown HTTP method")
)
