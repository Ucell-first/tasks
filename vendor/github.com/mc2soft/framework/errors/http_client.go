package errors

import errs "errors"

// ErrHTTPClientInvalidTimeout appears when invalid Timeout value in configuration
// struct was specified (<1).
var ErrHTTPClientInvalidTimeout = errs.New("invalid request timeout for HTTP client")
