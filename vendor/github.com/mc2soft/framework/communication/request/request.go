package request

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IBM/sarama"
	comcontext "github.com/mc2soft/framework/communication/context"
)

// Request is an interface for request data implementations.
type Request interface {
	GetContext() context.Context
	GetData() json.RawMessage
	GetFile() *File
	GetFrom() string
	GetHeaders() comcontext.Headers
	GetHost() string
	GetIsAsynchronous() bool
	GetMethod() string
	GetMetricPath() string
	GetNoMetrics() bool
	GetPath() string
	GetSendOnlyData() bool
	GetTo() string
	HasMetricPath() bool
	IsRaw() bool
	MarshalJSON() ([]byte, error)
	SetData(data json.RawMessage) error
	SetFile(file *File) error
	SetFrom(from string)
	SetHeaders(headers comcontext.Headers)
	SetHost(host string)
	SetIsAsynchronous(value bool)
	SetIsRaw(raw bool)
	SetMethod(method string)
	SetMetricPath(metricPath string)
	SetNoMetrics(value bool)
	SetPath(path string)
	SetSendOnlyData(value bool)
	SetTo(to string)
	SkipMetrics() bool
	ToHTTPHeader() http.Header
	ToInterfaceMap(delimiter string) map[string]interface{}
	ToSaramaHeaders(delimiterKey, delimiter string) []sarama.RecordHeader
	UnmarshalJSON(rawData []byte) error
}
