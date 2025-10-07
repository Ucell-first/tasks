package context

import (
	"context"
	"io"
	"sync"
)

// OutFunc is a function signature for data output.
type OutFunc func(code int, data interface{}, additionals ...interface{}) error

// DefaultContext is an incoming request context. If you're implementing your
// own request context structure - do not forget to embed it into your structure!
type DefaultContext struct {
	providerCtx     interface{}
	ctx             context.Context
	body            io.ReadCloser
	store           map[string]interface{}
	requestHeaders  Headers
	responseHeaders Headers
	genericOutFunc  OutFunc
	jsonOutFunc     OutFunc
	stringOutFunc   OutFunc
	htmlOutFunc     OutFunc
	protocol        string
	path            string
	from            string
	lock            sync.RWMutex
}

// NewDefaultContext creates new default context structure.
func NewDefaultContext() *DefaultContext {
	dc := DefaultContext{
		requestHeaders:  make(Headers),
		responseHeaders: make(Headers),
	}

	return &dc
}

// AddResponseHeader appends new value for response header.
func (dc *DefaultContext) AddResponseHeader(key, value string) { dc.responseHeaders.Add(key, value) }

// Body is a body of request.
func (dc *DefaultContext) Body() io.ReadCloser { return dc.body }

// Context returns context.Context instance for later use with other
// parts of application (like database queries).
func (dc *DefaultContext) Context() context.Context {
	return dc.ctx
}

// From returns origin identificator. This could be service name
// or complete PROTO://ADDR line, depending on communication provider
// in use.
func (dc *DefaultContext) From() string { return dc.from }

// Get gets value from a store by key
func (dc *DefaultContext) Get(key string) interface{} {
	dc.lock.RLock()
	defer dc.lock.RUnlock()

	return dc.store[key]
}

// GetRequestHeaderValue returns value from single header for request.
func (dc *DefaultContext) GetRequestHeaderValue(key string) string { return dc.requestHeaders.Get(key) }

// GetRequestHeaderValues returns all values from single header.
func (dc *DefaultContext) GetRequestHeaderValues(key string) []string {
	return dc.requestHeaders.Values(key)
}

// Path returns complete path on which request was made. For example
// it can be /api/v1/YOUNAMEIT.
func (dc *DefaultContext) Path() string { return dc.path }

// Protocol returns protocol name for request. For example it can be set
// to "http" or "rabbitmq".
func (dc *DefaultContext) Protocol() string { return dc.protocol }

// ProviderContext returns interface to provider Context object.
func (dc *DefaultContext) ProviderContext() interface{} { return dc.providerCtx }

// PutRequestHeader sets new values to header overwriting all other values.
func (dc *DefaultContext) PutRequestHeader(key string, value []string) {
	dc.requestHeaders.Put(key, value)
}

// Set sets key and its value in a store
func (dc *DefaultContext) Set(key string, val interface{}) {
	dc.lock.Lock()
	defer dc.lock.Unlock()

	if dc.store == nil {
		dc.store = make(map[string]interface{})
	}

	dc.store[key] = val
}

// SetBody sets body field.
func (dc *DefaultContext) SetBody(body io.ReadCloser) {
	dc.body = body
}

// SetContext sets request's context.
func (dc *DefaultContext) SetContext(ctx context.Context) { dc.ctx = ctx }

// SetFrom sets "From" field.
func (dc *DefaultContext) SetFrom(from string) { dc.from = from }

// SetRequestHeader sets single new value to header overwriting all other values.
func (dc *DefaultContext) SetRequestHeader(key, value string) { dc.requestHeaders.Set(key, value) }

// SetRequestHeaders sets "Headers" field.
func (dc *DefaultContext) SetRequestHeaders(headers Headers) { dc.requestHeaders = headers }

// SetResponseHeaders sets "Headers" field.
func (dc *DefaultContext) SetResponseHeaders(headers Headers) { dc.responseHeaders = headers }

// SetPath sets "Path" field.
func (dc *DefaultContext) SetPath(path string) { dc.path = path }

// SetProtocol sets "Protocol" field.
func (dc *DefaultContext) SetProtocol(protocol string) { dc.protocol = protocol }

// SetProviderContext sets provider-specific context for later usage.
func (dc *DefaultContext) SetProviderContext(ctx interface{}) {
	dc.providerCtx = ctx
}

// RequestHeaders returns all headers for request.
func (dc *DefaultContext) RequestHeaders() Headers { return dc.requestHeaders }

// ResponseHeaders returns all headers for request.
func (dc *DefaultContext) ResponseHeaders() Headers { return dc.responseHeaders }

// SetWriteFunc sets generic data writing function from provider.
func (dc *DefaultContext) SetWriteFunc(outputFunc OutFunc) { dc.genericOutFunc = outputFunc }

// SetWriteHTMLFunc sets HTML writing function from provider.
func (dc *DefaultContext) SetWriteHTMLFunc(outputFunc OutFunc) { dc.htmlOutFunc = outputFunc }

// SetWriteJSONFunc sets JSON writing function from provider.
func (dc *DefaultContext) SetWriteJSONFunc(outputFunc OutFunc) { dc.jsonOutFunc = outputFunc }

// SetWriteStringFunc sets string writing function from provider.
func (dc *DefaultContext) SetWriteStringFunc(outputFunc OutFunc) { dc.stringOutFunc = outputFunc }

// Write writes generic data to client, if applicable by provider. Writer should specify valid
// content type in headers ("content-type") otherwise error will be returned.
func (dc *DefaultContext) Write(code int, data []byte, additionals ...interface{}) error {
	return dc.genericOutFunc(code, data, additionals...)
}

// WriteHTML writes HTML back to client, if applicable by provider.
func (dc *DefaultContext) WriteHTML(code int, html string, additionals ...interface{}) error {
	return dc.htmlOutFunc(code, html, additionals...)
}

// WriteJSON writes JSON back, if applicable by provider.
func (dc *DefaultContext) WriteJSON(code int, data interface{}, additionals ...interface{}) error {
	return dc.jsonOutFunc(code, data, additionals...)
}

// WriteString writes string back to client, if applicable by provider.
func (dc *DefaultContext) WriteString(code int, data string, additionals ...interface{}) error {
	return dc.stringOutFunc(code, data, additionals...)
}
