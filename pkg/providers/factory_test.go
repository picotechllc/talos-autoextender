package providers

import (
	"testing"
)

func TestNewProviderFactory(t *testing.T) {
	factory := NewProviderFactory()
	if factory == nil {
		t.Fatal("Expected non-nil factory")
	}
}

func TestProviderFactoryCreateProvider(t *testing.T) {
	factory := NewProviderFactory()

	tests := []struct {
		name        string
		config      Provider
		shouldError bool
	}{
		{
			name: "valid linode provider",
			config: Provider{
				Name:   "linode",
				Region: "us-east",
				Credentials: map[string]string{
					"api_key": "test-api-key",
				},
			},
			shouldError: false,
		},
		{
			name: "unknown provider type",
			config: Provider{
				Name:   "unknown",
				Region: "us-east",
				Credentials: map[string]string{
					"api_key": "test-api-key",
				},
			},
			shouldError: true,
		},
		{
			name: "invalid provider config",
			config: Provider{
				Name:   "linode",
				Region: "",
				Credentials: map[string]string{
					"api_key": "test-api-key",
				},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := factory.CreateProvider(tt.config)
			if (err != nil) != tt.shouldError {
				t.Errorf("CreateProvider() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// Mock implementation for testing provider registration
type MockProvider struct{}

func (m *MockProvider) CreateCluster(spec ClusterSpec) error {
	return nil
}

func (m *MockProvider) DeleteCluster(name string) error {
	return nil
}

func (m *MockProvider) GetClusterStatus(name string) (ClusterStatus, error) {
	return ClusterStatus{}, nil
}

func (m *MockProvider) UpdateCluster(spec ClusterSpec) error {
	return nil
}

func TestRegisterProvider(t *testing.T) {
	factory := NewProviderFactory().(*DefaultProviderFactory)

	// Register a mock provider
	factory.RegisterProvider("mock", func(config Provider) (CloudProvider, error) {
		return &MockProvider{}, nil
	})

	// Try to create the mock provider
	provider, err := factory.CreateProvider(Provider{
		Name:   "mock",
		Region: "test-region",
		Credentials: map[string]string{
			"test": "value",
		},
	})

	if err != nil {
		t.Errorf("Failed to create mock provider: %v", err)
	}

	if provider == nil {
		t.Error("Expected non-nil provider")
	}

	// Verify it's the right type
	_, ok := provider.(*MockProvider)
	if !ok {
		t.Error("Expected *MockProvider type")
	}
}
