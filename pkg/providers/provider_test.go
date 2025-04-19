package providers

import (
	"testing"
)

func TestProviderValidation(t *testing.T) {
	tests := []struct {
		name        string
		provider    Provider
		shouldError bool
	}{
		{
			name: "valid linode provider",
			provider: Provider{
				Name:   "linode",
				Region: "us-east",
				Credentials: map[string]string{
					"api_token": "test-token",
				},
			},
			shouldError: false,
		},
		{
			name: "missing credentials",
			provider: Provider{
				Name:   "hetzner",
				Region: "eu-west",
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.provider.Validate()
			if (err != nil) != tt.shouldError {
				t.Errorf("Provider.Validate() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

func TestClusterSpecValidation(t *testing.T) {
	tests := []struct {
		name        string
		spec        ClusterSpec
		shouldError bool
	}{
		{
			name: "valid spec",
			spec: ClusterSpec{
				NodeCount:    3,
				NodeSize:    "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			shouldError: false,
		},
		{
			name: "invalid node count",
			spec: ClusterSpec{
				NodeCount:    0,
				NodeSize:    "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if (err != nil) != tt.shouldError {
				t.Errorf("ClusterSpec.Validate() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

func TestClusterCreation(t *testing.T) {
	tests := []struct {
		name        string
		provider    Provider
		spec        ClusterSpec
		shouldError bool
	}{
		{
			name: "valid linode cluster",
			provider: Provider{
				Name:   "linode",
				Region: "us-east",
				Credentials: map[string]string{
					"api_token": "test-token",
				},
			},
			spec: ClusterSpec{
				NodeCount:    3,
				NodeSize:    "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewProviderFactory()
			provider, err := factory.CreateProvider(tt.provider)
			if err != nil {
				t.Fatalf("Failed to create provider: %v", err)
			}

			err = provider.CreateCluster(tt.spec)
			if (err != nil) != tt.shouldError {
				t.Errorf("CreateCluster() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}