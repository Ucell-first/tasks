package errors

import errs "errors"

var (
	// ErrRabbitMQConfigurationURIInvalid appears when trying to set main
	// configuration for RabbitMQ provider with URI being an empty string.
	ErrRabbitMQConfigurationURIInvalid = errs.New("invalid URI for RabbitMQ connection")

	// ErrRabbitMQConnectionEstablishingFailed appears when trying to
	// establish connection to RabbitMQ and failed to do so.
	ErrRabbitMQConnectionEstablishingFailed = errs.New("failed to establish connection to RabbitMQ server")

	// ErrRabbitMQNoURIDefined appears when trying to set configuration
	// for connection with URIs undefined.
	ErrRabbitMQNoURIDefined = errs.New("no RabbitMQ URIs defined")
)
