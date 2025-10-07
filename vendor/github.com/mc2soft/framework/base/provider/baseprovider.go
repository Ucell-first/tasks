package provider

import (
	"context"
	"reflect"
	"sync"

	"github.com/mc2soft/framework/communication/request"
)

// Func is a function that executed on provider's start or on shutdown.
// Takes no parameters and might return error.
type Func func() error

// Base is an embeddable structure which describes what every provider will contain.
// This structure should be initialized by calling BaseProviderInitialize() function to
// properly populate it's internal state. This call should reside right after SetName()
// and SetLoggerParams() (if they're needed).
// Before calling Start() ensure that you've registered starting functions with RegisterStartFunc().
// Also before shutdown you should register a shutdown functions with RegisterShutdownFunc().
type Base struct {
	name string
	ctx  context.Context

	initializePostFuncs []Func
	initializePreFuncs  []Func
	startFuncs          []Func
	shutdownFuncs       []Func

	providerStarted      bool
	providerStartedMutex sync.RWMutex

	defaultRequest request.Request
}

// BaseProviderInitialize initializes internal state.
func (bp *Base) BaseProviderInitialize() {
	bp.initializePostFuncs = make([]Func, 0)
	bp.initializePreFuncs = make([]Func, 0)
	bp.startFuncs = make([]Func, 0)
	bp.shutdownFuncs = make([]Func, 0)
	bp.defaultRequest = &request.DefaultRequest{}
}

// GetContext returns provider's context. Should be used by providers itself (mostly).
func (bp *Base) GetContext() context.Context {
	return bp.ctx
}

// GetName returns name of provider that was previously set by SetName().
func (bp *Base) GetName() string {
	return bp.name
}

// GetNewDefaultRequestStruct get a new instance of request.Request
func (bp *Base) GetNewDefaultRequestStruct() request.Request {
	return reflect.New(reflect.ValueOf(bp.defaultRequest).Elem().Type()).Interface().(request.Request)
}

// IsStarted returns true if Start() was called.
func (bp *Base) IsStarted() bool {
	bp.providerStartedMutex.RLock()
	defer bp.providerStartedMutex.RUnlock()

	return bp.providerStarted
}

// RegisterDefaultRequestStruct registers custom defaultRequest, used by kafka, nats, rabbitmq during unmarshalling
func (bp *Base) RegisterDefaultRequestStruct(
	r request.Request,
) {
	bp.defaultRequest = r
}

// RegisterStartFunc registers starting function to execute.
func (bp *Base) RegisterStartFunc(f Func) {
	bp.startFuncs = append(bp.startFuncs, f)
}

// RegisterStopFunc registers a function to execute on application shutdown.
func (bp *Base) RegisterStopFunc(f Func) {
	bp.shutdownFuncs = append(bp.shutdownFuncs, f)
}

// SetContext sets provider's context.
func (bp *Base) SetContext(ctx context.Context) {
	bp.ctx = ctx
}

// SetName sets provider's name. Should be called ASAP.
func (bp *Base) SetName(name string) {
	bp.name = name
}

// Shutdown calls all registered shutdown functions and return all errors appeared.
func (bp *Base) Shutdown() []error {
	var errs []error

	for _, f := range bp.shutdownFuncs {
		err := f()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// Start calls all registered starting functions and return appeared error immediately.
func (bp *Base) Start() error {
	for _, f := range bp.startFuncs {
		err := f()
		if err != nil {
			return err
		}
	}

	bp.providerStartedMutex.Lock()
	bp.providerStarted = true
	bp.providerStartedMutex.Unlock()

	return nil
}
