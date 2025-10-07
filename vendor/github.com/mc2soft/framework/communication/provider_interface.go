package communication

import (
	"github.com/mc2soft/framework/base/provider"
	comcontext "github.com/mc2soft/framework/communication/context"
	"github.com/mc2soft/framework/communication/request"
	"github.com/mc2soft/framework/communication/response"
)

// Provider is an interface for communication channel provider.
type Provider interface {
	provider.Interface

	// IsClient returns true if communication provider is able to do
	// client-server communication.
	IsClient() bool
	// IsServer returns true if communication provider is able to serving
	// things.
	IsServer() bool
	// RegisterHandler registers handler for communication channel provider.
	// Some
	// providers might ignore "method" parameter (e.g. message queues).
	RegisterHandler(method, path string, handler HandlerFunc) error
	// Request sends request to remote thing in synchronous manner and returns
	// either a response or error.
	Request(request request.Request) (*response.Response, error)
	// Send sends request to remote thing in synchronous manner and does not
	// return anything but protocol error.
	Send(request request.Request) error
	// SendRaw sends request to remote thing in synchronous manner.
	// That method sends request's Data field as message and puts Headers into
	// provider's headers,
	// if applicable, all other fields are ignored.
	SendRaw(request request.Request) error
	// SendAsync sends request to remote thing in asynchronous manner. It
	// won't return any response or protocol errors, only internal for
	// framework.
	SendAsync(request request.Request) error
	// SendAsyncRaw sends request to remote thing in asynchronous manner. It
	// won't return any response or protocol errors, only internal for
	// framework. That method sends request's Data field as message and puts
	// Headers into provider's headers,
	// if applicable, all other fields are ignored.
	SendAsyncRaw(request request.Request) error
	// SetHeaderDelimiter sets delimiter var for separating same key headers if
	// array sending is impossible.
	SetHeaderDelimiter(delimiter string)
	// RegisterHandlerNamingFunc registers function for custom naming.
	RegisterHandlerNamingFunc(f HandlerNamingFunc)
	// RegisterMiddleware registers middleware which runs before actual request
	// processing with handler. Middlewares currently supported only for
	// providers which supports to serve things (a.k.a. server providers).
	RegisterMiddleware(mldwr MiddlewareFunc)
}

// HandlerFunc is a function signature for communications handler.
type HandlerFunc func(comcontext.Context) error

// HandlerNamingFunc is a function signature for generating queue name.
type HandlerNamingFunc func(appName, method, path string) string

// MiddlewareFunc is a function which provides middleware functionality and
// registered in providers. First parameter is a communication context, second
// is a next middleware to run.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc
