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
				NodeSize:     "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			shouldError: false,
		},
		{
			name: "invalid node count",
			spec: ClusterSpec{
				NodeCount:    0,
				NodeSize:     "g6-standard-2",
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

func TestClusterScaling(t *testing.T) {
	tests := []struct {
		name        string
		provider    Provider
		initialSpec ClusterSpec
		targetSpec  ClusterSpec
		shouldError bool
	}{
		{
			name: "scale up cluster nodes",
			provider: Provider{
				Name:   "linode",
				Region: "us-east",
				Credentials: map[string]string{
					"api_token": "test-token",
				},
			},
			initialSpec: ClusterSpec{
				NodeCount:    3,
				NodeSize:     "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			targetSpec: ClusterSpec{
				NodeCount:    5,
				NodeSize:     "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			shouldError: false,
		},
		{
			name: "scale down cluster nodes",
			provider: Provider{
				Name:   "linode",
				Region: "us-east",
				Credentials: map[string]string{
					"api_token": "test-token",
				},
			},
			initialSpec: ClusterSpec{
				NodeCount:    5,
				NodeSize:     "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			targetSpec: ClusterSpec{
				NodeCount:    3,
				NodeSize:     "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			shouldError: false,
		},
		{
			name: "upgrade node size",
			provider: Provider{
				Name:   "linode",
				Region: "us-east",
				Credentials: map[string]string{
					"api_token": "test-token",
				},
			},
			initialSpec: ClusterSpec{
				NodeCount:    3,
				NodeSize:     "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			targetSpec: ClusterSpec{
				NodeCount:    3,
				NodeSize:     "g6-standard-4",
				TalosVersion: "v1.6.0",
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Skip until provider is fully implemented")
			factory := NewProviderFactory()
			provider, err := factory.CreateProvider(tt.provider)
			if err != nil {
				t.Fatalf("Failed to create provider: %v", err)
			}

			err = provider.CreateCluster(tt.initialSpec)
			if err != nil {
				t.Fatalf("Failed to create initial cluster: %v", err)
			}

			err = provider.UpdateCluster(tt.targetSpec)
			if (err != nil) != tt.shouldError {
				t.Errorf("UpdateCluster() error = %v, shouldError %v", err, tt.shouldError)
			}

			status, err := provider.GetClusterStatus("test-cluster")
			if err != nil {
				t.Fatalf("Failed to get cluster status: %v", err)
			}

			if status.NodeCount != tt.targetSpec.NodeCount {
				t.Errorf("Expected node count %d, got %d", tt.targetSpec.NodeCount, status.NodeCount)
			}
		})
	}
}

func TestClusterLifecycle(t *testing.T) {
	tests := []struct {
		name        string
		provider    Provider
		spec        ClusterSpec
		operations  []string
		shouldError bool
	}{
		{
			name: "full cluster lifecycle",
			provider: Provider{
				Name:   "linode",
				Region: "us-east",
				Credentials: map[string]string{
					"api_token": "test-token",
				},
			},
			spec: ClusterSpec{
				NodeCount:    3,
				NodeSize:     "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			operations:  []string{"create", "scale", "delete"},
			shouldError: false,
		},
		{
			name: "cluster creation with invalid credentials",
			provider: Provider{
				Name:   "linode",
				Region: "us-east",
				Credentials: map[string]string{
					"api_token": "",
				},
			},
			spec: ClusterSpec{
				NodeCount:    3,
				NodeSize:     "g6-standard-2",
				TalosVersion: "v1.6.0",
			},
			operations:  []string{"create"},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Skip until provider is fully implemented")
			factory := NewProviderFactory()
			provider, err := factory.CreateProvider(tt.provider)
			if err != nil {
				t.Fatalf("Failed to create provider: %v", err)
			}

			for _, op := range tt.operations {
				var err error
				switch op {
				case "create":
					err = provider.CreateCluster(tt.spec)
				case "scale":
					newSpec := tt.spec
					newSpec.NodeCount = 5
					err = provider.UpdateCluster(newSpec)
				case "delete":
					err = provider.DeleteCluster("test-cluster")
				}

				if (err != nil) != tt.shouldError {
					t.Errorf("Operation %s error = %v, shouldError %v", op, err, tt.shouldError)
				}
			}
		})
	}
}
