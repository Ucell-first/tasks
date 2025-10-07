package request

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	comcontext "github.com/mc2soft/framework/communication/context"
	frmwrkreq "github.com/mc2soft/framework/communication/request"
)

// DefaultRequest represents generic request structure which is used for inter-service
// communication.
// Use New() and setters for filling Request struct
type DefaultRequest struct {
	ctx            context.Context
	headers        comcontext.Headers
	file           *frmwrkreq.File
	from           string
	to             string
	host           string
	method         string
	path           string
	metricPath     string
	data           json.RawMessage
	sendOnlyData   bool
	isAsynchronous bool
	isRaw          bool
	noMetrics      bool
}

// New creates new request struct.
func New(
	ctx context.Context,
	method, path string,
	headers comcontext.Headers,
	data json.RawMessage,
) *DefaultRequest {
	return &DefaultRequest{
		ctx:     ctx,
		method:  method,
		path:    path,
		data:    data,
		headers: headers,
	}
}

// GetContext returns request's context.
func (r *DefaultRequest) GetContext() context.Context {
	return r.ctx
}

// GetData returns request's data.
func (r *DefaultRequest) GetData() json.RawMessage {
	return r.data
}

// GetFile returns request's file.
func (r *DefaultRequest) GetFile() *frmwrkreq.File {
	return r.file
}

// GetFrom returns service name from which request is arrived.
func (r *DefaultRequest) GetFrom() string {
	return r.from
}

// GetHeaders returns request's headers.
func (r *DefaultRequest) GetHeaders() comcontext.Headers {
	return r.headers
}

// GetHost returns host from which request is arrived.
func (r *DefaultRequest) GetHost() string {
	return r.host
}

// GetIsAsynchronous returns asynchronous message state.
func (r *DefaultRequest) GetIsAsynchronous() bool {
	return r.isAsynchronous
}

// GetMethod returns HTTP verb for request.
func (r *DefaultRequest) GetMethod() string {
	return r.method
}

// GetMetricPath returns path for Prometheus metrics.
func (r *DefaultRequest) GetMetricPath() string {
	return r.metricPath
}

// GetNoMetrics returns "no metrics" flag for request.
func (r *DefaultRequest) GetNoMetrics() bool {
	return r.noMetrics
}

// GetPath returns request's path.
func (r *DefaultRequest) GetPath() string {
	return r.path
}

// GetSendOnlyData returns "send only data" flag for request.
func (r *DefaultRequest) GetSendOnlyData() bool {
	return r.sendOnlyData
}

// GetTo returns destination name for request.
func (r *DefaultRequest) GetTo() string {
	return r.to
}

// HasMetricPath returns true if request has path for Prometheus metrics defined.
func (r *DefaultRequest) HasMetricPath() bool {
	return len(r.metricPath) > 0
}

// IsRaw returns true if only data should be sent.
func (r *DefaultRequest) IsRaw() bool {
	return r.isRaw
}

// MarshalJSON marshals request into JSON.
func (r *DefaultRequest) MarshalJSON() ([]byte, error) {
	data := &requestStruct{
		From:           r.GetFrom(),
		To:             r.GetTo(),
		Method:         r.GetMethod(),
		Path:           r.GetPath(),
		Data:           r.GetData(),
		IsAsynchronous: r.GetIsAsynchronous(),
		Headers:        r.GetHeaders(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", frmwrkreq.ErrFailedToMarshalRequestToJSON, err)
	}

	return jsonData, nil
}

// SetContext sets context for request.
func (r *DefaultRequest) SetContext(ctx context.Context) {
	r.ctx = ctx
}

// SetData sets data for request.
func (r *DefaultRequest) SetData(jsonData json.RawMessage) error {
	if r.file != nil {
		return frmwrkreq.ErrCantUseBothDataAndFile
	}

	r.data = jsonData

	return nil
}

// SetFile sets file for request.
func (r *DefaultRequest) SetFile(file *frmwrkreq.File) error {
	if r.data != nil {
		return frmwrkreq.ErrCantUseBothDataAndFile
	}

	r.file = file

	return nil
}

// SetFrom sets originator name for request.
func (r *DefaultRequest) SetFrom(From string) {
	r.from = From
}

// SetHeaders sets headers for request.
func (r *DefaultRequest) SetHeaders(headers comcontext.Headers) {
	r.headers = headers
}

// SetHost sets requester's host.
func (r *DefaultRequest) SetHost(Host string) {
	r.host = Host
}

// SetIsAsynchronous sets asynchronous flag for request.
func (r *DefaultRequest) SetIsAsynchronous(IsAsynchronous bool) {
	r.isAsynchronous = IsAsynchronous
}

// SetIsRaw sets raw flag for request.
func (r *DefaultRequest) SetIsRaw(raw bool) {
	r.isRaw = raw
}

// SetMethod sets HTTP verb for request.
func (r *DefaultRequest) SetMethod(method string) {
	r.method = method
}

// SetMetricPath sets Prometheus metric path for request.
func (r *DefaultRequest) SetMetricPath(MetricPath string) {
	r.metricPath = MetricPath
}

// SetNoMetrics sets "no metrics" flag for request
func (r *DefaultRequest) SetNoMetrics(noMetrics bool) {
	r.noMetrics = noMetrics
}

// SetPath sets path for request.
func (r *DefaultRequest) SetPath(path string) {
	r.path = path
}

// SetSendOnlyData sets "send only data" flag for request.
func (r *DefaultRequest) SetSendOnlyData(SendOnlyData bool) {
	r.sendOnlyData = SendOnlyData
}

// SetTo sets destination.
func (r *DefaultRequest) SetTo(To string) {
	r.to = To
}

// SkipMetrics returns "no metrics" flag for request.
func (r *DefaultRequest) SkipMetrics() bool {
	return r.noMetrics
}

// ToHTTPHeader converts headers to HTTP headers.
func (r *DefaultRequest) ToHTTPHeader() http.Header {
	return r.GetHeaders().ToHTTPHeader()
}

// ToInterfaceMap converts headers to interface map.
func (r *DefaultRequest) ToInterfaceMap(delimiter string) map[string]interface{} {
	return r.headers.ToInterfaceMap(delimiter)
}

// ToSaramaHeaders converts headers into Sarama format.
func (r *DefaultRequest) ToSaramaHeaders(delimiterKey, delimiter string) []sarama.RecordHeader {
	return r.headers.ToSaramaHeaders(delimiterKey, delimiter)
}

// UnmarshalJSON parses passed bytes into request structure.
func (r *DefaultRequest) UnmarshalJSON(rawData []byte) error {
	req := &requestStruct{}

	err := json.Unmarshal(rawData, req)
	if err != nil {
		return frmwrkreq.ErrFailedToUnmarshalJSONtoRequest
	}

	r.SetFrom(req.From)
	r.SetTo(req.To)
	r.SetIsAsynchronous(req.IsAsynchronous)
	r.SetMethod(req.Method)
	r.SetPath(req.Path)
	r.SetHeaders(req.Headers)

	err = r.SetData(req.Data)
	if err != nil {
		return err
	}

	return nil
}
