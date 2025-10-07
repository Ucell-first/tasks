package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/mc2soft/framework/base/provider"
	"github.com/mc2soft/framework/errors"
)

// BaseHandler is an embeddable structure for all subsystems handlers. It contains everything
// they have in common.
type BaseHandler struct {
	providers      map[string]provider.Interface
	globalCtx      context.Context
	providersMutex sync.RWMutex
}

// GetRawProvider returns requested provider by name. Returns raw provider which is
// provider.Interface and should be transformed into needed interface by actual handler.
// Should not be called by hands in application's code unless you're know what you're doing!
func (bh *BaseHandler) GetRawProvider(providerName string) (provider.Interface, error) {
	bh.providersMutex.RLock()
	defer bh.providersMutex.RUnlock()

	if len(bh.providers) == 0 {
		return nil, errors.ErrNoProvidersRegistered
	}

	provider, found := bh.providers[providerName]
	if !found {
		return nil, errors.ErrProviderNotFound
	}

	return provider, nil
}

// GetProvidersNames returns a slice of strings with names of registered providers.
func (bh *BaseHandler) GetProvidersNames() []string {
	bh.providersMutex.RLock()
	defer bh.providersMutex.RUnlock()

	names := make([]string, 0, len(bh.providers))

	for name := range bh.providers {
		names = append(names, name)
	}

	return names
}

// RegisterProvider registers new communication provider.
func (bh *BaseHandler) RegisterProvider(name string, provIface provider.Interface) error {
	bh.providersMutex.Lock()
	defer bh.providersMutex.Unlock()

	if bh.providers == nil {
		bh.providers = make(map[string]provider.Interface)
	}

	if _, found := bh.providers[name]; found {
		return errors.ErrProviderAlreadyRegistered
	}

	err := provIface.Initialize()
	if err != nil {
		return fmt.Errorf("basehandler: register provider: initialize provider: %w", err)
	}

	bh.providers[name] = provIface

	return nil
}

// SetConfig sets provider-specific configuration.
func (bh *BaseHandler) SetConfig(providerName string, config interface{}) error {
	bh.providersMutex.RLock()
	defer bh.providersMutex.RUnlock()

	if len(bh.providers) == 0 {
		return errors.ErrNoProvidersRegistered
	}

	provider, found := bh.providers[providerName]

	if !found {
		return errors.ErrProviderNotFound
	}

	if err := provider.SetConfig(config); err != nil {
		return fmt.Errorf("set config: %w", err)
	}

	return nil
}

// SetContext sets context that will be passed to every provider whenever possible.
func (bh *BaseHandler) SetContext(ctx context.Context) {
	bh.globalCtx = ctx
}

// Start passes global context to provider and tries to start it.
func (bh *BaseHandler) Start(providerName string) error {
	provider, err := bh.GetRawProvider(providerName)
	if err != nil {
		return err
	}

	provider.SetContext(bh.globalCtx)

	if err := provider.Start(); err != nil {
		return fmt.Errorf("Start: %w", err)
	}

	return nil
}

// StartAll starts all registered providers.
func (bh *BaseHandler) StartAll() error {
	providers := make([]string, 0)

	bh.providersMutex.RLock()
	for name := range bh.providers {
		providers = append(providers, name)
	}
	bh.providersMutex.RUnlock()

	for _, name := range providers {
		if err := bh.Start(name); err != nil {
			return fmt.Errorf("provider %s: %w", name, err)
		}
	}

	return nil
}

// Stop stops specific provider.
func (bh *BaseHandler) Stop(providerName string) []error {
	provider, err := bh.GetRawProvider(providerName)
	if err != nil {
		return []error{err}
	}

	errs := provider.Shutdown()
	if len(errs) > 0 {
		return errs
	}

	return nil
}

// StopAll stops all registered providers.
func (bh *BaseHandler) StopAll() []error {
	bh.providersMutex.RLock()
	defer bh.providersMutex.RUnlock()

	for _, provider := range bh.providers {
		errs := provider.Shutdown()
		if len(errs) > 0 {
			return errs
		}
	}

	return nil
}
