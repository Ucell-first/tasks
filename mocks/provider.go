package mocks

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/mc2soft/framework/base/provider"
	"github.com/mc2soft/framework/communication"
	comcontext "github.com/mc2soft/framework/communication/context"
	"github.com/mc2soft/framework/communication/request"
	"github.com/mc2soft/framework/communication/response"
)

var errNotFound = errors.New("handler not found")

type MockProvider struct {
	handlers map[string]communication.HandlerFunc
}

func New() MockProvider {
	return MockProvider{handlers: make(map[string]communication.HandlerFunc)}
}

func (m MockProvider) BaseProviderInitialize() {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) GetContext() context.Context {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) GetName() string {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) Initialize() error {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) RegisterStartFunc(f provider.Func) {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) RegisterStopFunc(f provider.Func) {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) SetConfig(i interface{}) error {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) SetContext(ctx context.Context) {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) SetName(s string) {
	// TODO implement me
	panic("implement me")
}

func (m MockProvider) Shutdown() []error {
	return nil
}

func (m MockProvider) Start() error {
	return nil
}

func (m MockProvider) IsClient() bool {
	return true
}

func (m MockProvider) IsServer() bool {
	return false
}

func (m MockProvider) SendAsync(request request.Request) error {
	return nil
}

func (m MockProvider) SendAsyncRaw(request request.Request) error {
	return nil
}

func (m MockProvider) SetHeaderDelimiter(delimiter string) {
}

func (m MockProvider) RegisterHandlerNamingFunc(f communication.HandlerNamingFunc) {
}

func (m MockProvider) RegisterMiddleware(middlewareFunc communication.MiddlewareFunc) {
}

func (m MockProvider) RegisterHandler(method, path string, handler communication.HandlerFunc) error {
	m.handlers[path] = handler

	return nil
}

func (m MockProvider) Request(request request.Request) (*response.Response, error) {
	return &response.Response{}, nil
}

//nolint:wrapcheck
func (m MockProvider) Send(request request.Request) error {
	cctx := comcontext.NewDefaultContext()

	handler, ok := m.handlers[request.GetPath()]
	if !ok {
		return errNotFound
	}

	data, err := request.GetData().MarshalJSON()
	if err != nil {
		return err
	}

	readCloser := io.NopCloser(bytes.NewReader(data))

	cctx.SetBody(readCloser)

	return handler(cctx)
}

//nolint:wrapcheck
func (m MockProvider) SendRaw(request request.Request) error {
	cctx := comcontext.NewDefaultContext()

	handler, ok := m.handlers[request.GetPath()]
	if !ok {
		return errNotFound
	}

	data, err := request.GetData().MarshalJSON()
	if err != nil {
		return err
	}

	readCloser := io.NopCloser(bytes.NewReader(data))

	cctx.SetBody(readCloser)

	return handler(cctx)
}

func (m MockProvider) GetNewDefaultRequestStruct() request.Request {
	return &request.DefaultRequest{}
}

func (m MockProvider) RegisterDefaultRequestStruct(
	r request.Request,
) {
}
