package context

import (
	"context"
	"io"
)

// Context is a generic request context that will contain everything
// every request might need.
// Every communication provider SHOULD have own functions for filling
// this context. Every handler SHOULD use this context.
// This context is for servers only, client communication providers
// should not implement it.
type Context interface {
	// AddResponseHeader appends new value for response header.
	AddResponseHeader(key, value string)
	// Body is a body of request.
	Body() io.ReadCloser
	// Context returns context.Context instance for later use with other
	// parts of application (like database queries).
	Context() context.Context
	// From returns origin identificator. This could be service name
	// or complete PROTO://ADDR line, depending on communication provider
	// in use.
	From() string
	// Gets value from a store by key
	Get(key string) interface{}
	// GetRequestHeaderValue returns value from single header for request.
	GetRequestHeaderValue(key string) string
	// GetRequestHeaderValues returns all values from single header.
	GetRequestHeaderValues(key string) []string
	// Path returns complete path on which request was made. For example
	// it can be /api/v1/YOUNAMEIT.
	Path() string
	// Protocol returns protocol name for request. For example it can be set
	// to "http" or "rabbitmq".
	Protocol() string
	// ProviderContext returns interface to provider Context object.
	ProviderContext() interface{}
	// RequestHeaders returns Headers from request structure which contains all known headers for
	// this request. It's API pretty similar to http.Header.
	RequestHeaders() Headers
	// ResponseHeaders returns Headers for response structure which contains all known headers for
	// this request. It's API pretty similar to http.Header.
	ResponseHeaders() Headers
	// Sets key and its value in a store
	Set(key string, val interface{})
	// SetContext sets request's context.
	SetContext(ctx context.Context)
	// Write writes generic data to client, if applicable by provider. Writer should specify valid
	// content type in headers ("content-type") otherwise error will be returned.
	Write(code int, data []byte, additionals ...interface{}) error
	// WriteHTML writes HTML back to client, if applicable by provider.
	WriteHTML(code int, html string, additionals ...interface{}) error
	// WriteJSON writes JSON back, if applicable by provider.
	WriteJSON(code int, data interface{}, additionals ...interface{}) error
	// WriteString writes string back to client, if applicable by provider.
	WriteString(code int, data string, additionals ...interface{}) error
}
