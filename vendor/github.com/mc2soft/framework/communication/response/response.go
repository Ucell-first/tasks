package response

import (
	"io"

	comcontext "github.com/mc2soft/framework/communication/context"
)

// Response represents generic service response structure for inter-service
// communication.
type Response struct {
	// Body is a response body.
	Body io.ReadCloser
	// Headers is a map of headers (like in HTTP).
	Headers comcontext.Headers
	// Code is a response code (like in HTTP).
	Code int
}
