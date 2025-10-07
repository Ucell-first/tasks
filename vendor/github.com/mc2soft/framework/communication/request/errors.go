package request

import "errors"

var (
	// ErrCantUseBothDataAndFile indicates that we can't use both data from file and passed bytes.
	ErrCantUseBothDataAndFile = errors.New("can't use simultaneously data and file")
	// ErrFailedToMarshalRequestToJSON indicates that we are unable to marshal JSON from passed data.
	ErrFailedToMarshalRequestToJSON = errors.New("failed to marshal request to JSON")
	// ErrFailedToUnmarshalJSONtoRequest indicates that we are unable to parse JSON.
	ErrFailedToUnmarshalJSONtoRequest = errors.New("failed to unmarshal JSON to request")
)
