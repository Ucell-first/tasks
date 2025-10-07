package errors

import errs "errors"

var (
	// ErrKafkaConfigurationURIInvalid appears when trying to set main
	// configuration for Kafka provider with URI being an empty string.
	ErrKafkaConfigurationURIInvalid = errs.New("invalid URI for Kafka connection")

	// ErrKafkaNoURIDefined appears when trying to set configuration
	// for connection with URIs undefined.
	ErrKafkaNoURIDefined = errs.New("no Kafka URIs defined")

	// ErrKafkaDoNotSkipMessage appears when you need to leave a message
	// in topic on handler's error.
	// Use this error when return error form handler.
	// Example:
	//	return fmt.Errorf(
	//		"myFunction: %w; %w",
	// 		err,
	// 		frameworkErr.ErrKafkaDoNotSkipMessage,
	//	)
	ErrKafkaDoNotSkipMessage = errs.New("do not skip this message")
)
