package communication

import (
	"fmt"

	"github.com/mc2soft/framework/errors"
	"github.com/mc2soft/framework/internal/handler"
)

// Handler handles all communications channels application will use.
// Communication channels should be registered before application's
// business logic initialization in order to properly register
// handlers.
type Handler struct {
	handler.BaseHandler
}

// Initialize performs preliminary initialization actions like internal
// structures initialization.
func (h *Handler) Initialize() {}

// GetProvider returns provider by name.
func (h *Handler) GetProvider(providerName string) (Provider, error) {
	providerRaw, err := h.GetRawProvider(providerName)
	if err != nil {
		return nil, fmt.Errorf("communications handler: %w", err)
	}

	provider, ok := providerRaw.(Provider)
	if !ok {
		return nil, errors.ErrInvalidProvider
	}

	return provider, nil
}

// RegisterHandler registers handler for all communication channels.
func (h *Handler) RegisterHandler(providerName, method, path string, handler HandlerFunc) error {
	provider, err := h.GetProvider(providerName)
	if err != nil {
		return err
	}

	if err := provider.RegisterHandler(method, path, handler); err != nil {
		return fmt.Errorf("RegisterHandler: %w", err)
	}

	return nil
}

// RegisterMiddleware registers middleware for specified provider.
func (h *Handler) RegisterMiddleware(providerName string, handler MiddlewareFunc) error {
	provider, err := h.GetProvider(providerName)
	if err != nil {
		return err
	}

	provider.RegisterMiddleware(handler)

	return nil
}

// SetHeaderDelimiters sets same headerDelimiter for all communication providers.
func (h *Handler) SetHeaderDelimiters(delimiter string) error {
	for _, providerName := range h.GetProvidersNames() {
		provider, err := h.GetProvider(providerName)
		if err != nil {
			return err
		}

		provider.SetHeaderDelimiter(delimiter)
	}

	return nil
}
