package providers

import (
	"fmt"
)

// DefaultProviderFactory implements the ProviderFactory interface
type DefaultProviderFactory struct {
	registeredProviders map[string]func(Provider) (CloudProvider, error)
}

// NewProviderFactory creates a new provider factory with registered providers
func NewProviderFactory() ProviderFactory {
	factory := &DefaultProviderFactory{
		registeredProviders: make(map[string]func(Provider) (CloudProvider, error)),
	}

	// Register built-in providers
	factory.RegisterProvider("linode", func(config Provider) (CloudProvider, error) {
		return NewLinodeProvider(config)
	})

	return factory
}

// RegisterProvider adds a new provider to the factory
func (f *DefaultProviderFactory) RegisterProvider(name string, creator func(Provider) (CloudProvider, error)) {
	f.registeredProviders[name] = creator
}

// CreateProvider instantiates a cloud provider based on configuration
func (f *DefaultProviderFactory) CreateProvider(config Provider) (CloudProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid provider configuration: %v", err)
	}

	creator, ok := f.registeredProviders[config.Name]
	if !ok {
		return nil, fmt.Errorf("unknown provider type: %s", config.Name)
	}

	return creator(config)
}
