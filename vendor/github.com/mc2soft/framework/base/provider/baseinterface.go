package provider

import (
	"context"

	"github.com/mc2soft/framework/communication/request"
)

// Interface describes an interface every provider should conform. Any custom
// provider interface should embed this one.
type Interface interface {
	// BaseProviderInitialize initializes internal state for base provider. It should be
	// called by provider itself and not by external user.
	BaseProviderInitialize()
	// GetContext returns provider's context. Should be used by providers itself (mostly).
	GetContext() context.Context
	// GetName returns name of provider that was previously set by SetName().
	GetName() string
	// Initialize initializes provider's internal state.
	Initialize() error
	// RegisterStartFunc registers starting function to execute.
	RegisterStartFunc(f Func)
	// RegisterStopFunc registers a function to execute on application shutdown.
	RegisterStopFunc(f Func)
	// SetConfig sets provider's configuration.
	SetConfig(cfg interface{}) error
	// SetContext sets provider's context.
	SetContext(ctx context.Context)
	// SetName sets provider's name. Should be called ASAP.
	SetName(name string)
	// Shutdown calls all registered shutdown functions and return all errors appeared.
	Shutdown() []error
	// Start calls all registered starting functions and return appeared error immediately.
	Start() error
	// RegisterDefaultRequestStruct registers custom defaultRequest, used by kafka, nats, rabbitmq during unmarshalling
	RegisterDefaultRequestStruct(request.Request)
	// GetNewDefaultRequestStruct get a new instance of request.Request
	GetNewDefaultRequestStruct() request.Request
}
