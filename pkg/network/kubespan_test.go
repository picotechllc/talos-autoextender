
package network

import (
	"testing"
)

func TestKubeSpanManager(t *testing.T) {
	tests := []struct {
		name          string
		homeEndpoint  string
		cloudClusters []string
		wantErr       bool
	}{
		{
			name:         "valid home endpoint",
			homeEndpoint: "10.0.0.1:50000",
			wantErr:     false,
		},
		{
			name:         "invalid home endpoint",
			homeEndpoint: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeSpanManager(tt.homeEndpoint)
			err := manager.ValidateEndpoint()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKubeSpanNetworkTransitions(t *testing.T) {
	tests := []struct {
		name           string
		homeEndpoint   string
		initialClusters []struct {
			name     string
			endpoint string
			state    string
		}
		operations []struct {
			name    string
			cluster string
			action  string
		}
		expectedState map[string]string
		shouldError   bool
	}{
		{
			name:         "handle node failures gracefully",
			homeEndpoint: "10.0.0.1:50000",
			initialClusters: []struct {
				name     string
				endpoint string
				state    string
			}{
				{"cloud-1", "192.168.1.1:50000", "connected"},
				{"cloud-2", "192.168.1.2:50000", "connected"},
			},
			operations: []struct {
				name    string
				cluster string
				action  string
			}{
				{"simulate_failure", "cloud-1", "disconnect"},
				{"verify_failover", "cloud-2", "verify"},
				{"restore_node", "cloud-1", "connect"},
			},
			expectedState: map[string]string{
				"cloud-1": "connected",
				"cloud-2": "connected",
			},
			shouldError: false,
		},
		{
			name:         "network partition recovery",
			homeEndpoint: "10.0.0.1:50000",
			initialClusters: []struct {
				name     string
				endpoint string
				state    string
			}{
				{"cloud-1", "192.168.1.1:50000", "connected"},
				{"cloud-2", "192.168.1.2:50000", "connected"},
			},
			operations: []struct {
				name    string
				cluster string
				action  string
			}{
				{"partition_network", "cloud-1", "disconnect"},
				{"partition_network", "cloud-2", "disconnect"},
				{"restore_network", "cloud-1", "connect"},
				{"restore_network", "cloud-2", "connect"},
			},
			expectedState: map[string]string{
				"cloud-1": "connected",
				"cloud-2": "connected",
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeSpanManager(tt.homeEndpoint)

			// Setup initial clusters
			for _, cluster := range tt.initialClusters {
				err := manager.AddCloudCluster(cluster.name, cluster.endpoint)
				if err != nil {
					t.Fatalf("Failed to add initial cluster %s: %v", cluster.name, err)
				}
			}

			// Execute operations
			for _, op := range tt.operations {
				var err error
				switch op.action {
				case "disconnect":
					err = manager.DisconnectCluster(op.cluster)
				case "connect":
					cluster := findCluster(tt.initialClusters, op.cluster)
					err = manager.AddCloudCluster(cluster.name, cluster.endpoint)
				case "verify":
					err = manager.VerifyClusterConnectivity(op.cluster)
				}

				if (err != nil) != tt.shouldError {
					t.Errorf("Operation %s error = %v, shouldError %v", op.name, err, tt.shouldError)
				}
			}

			// Verify final state
			for clusterName, expectedState := range tt.expectedState {
				state, err := manager.GetClusterState(clusterName)
				if err != nil {
					t.Errorf("Failed to get cluster %s state: %v", clusterName, err)
				}
				if state != expectedState {
					t.Errorf("Cluster %s expected state %s, got %s", clusterName, expectedState, state)
				}
			}
		})
	}
}

func TestKubeSpanMeshConnectivity(t *testing.T) {
	tests := []struct {
		name          string
		homeEndpoint  string
		cloudClusters []struct {
			name     string
			endpoint string
		}
		operations  []string
		shouldError bool
	}{
		{
			name:         "full mesh connectivity",
			homeEndpoint: "10.0.0.1:50000",
			cloudClusters: []struct {
				name     string
				endpoint string
			}{
				{"cloud-1", "192.168.1.1:50000"},
				{"cloud-2", "192.168.1.2:50000"},
			},
			operations:  []string{"connect", "verify", "disconnect"},
			shouldError: false,
		},
		{
			name:         "handle connection failures",
			homeEndpoint: "10.0.0.1:50000",
			cloudClusters: []struct {
				name     string
				endpoint string
			}{
				{"cloud-1", "invalid:50000"},
			},
			operations:  []string{"connect"},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeSpanManager(tt.homeEndpoint)

			for _, op := range tt.operations {
				var err error
				switch op {
				case "connect":
					for _, cluster := range tt.cloudClusters {
						err = manager.AddCloudCluster(cluster.name, cluster.endpoint)
					}
				case "verify":
					err = manager.VerifyMeshConnectivity()
				case "disconnect":
					err = manager.DisconnectCluster("cloud-1")
				}

				if (err != nil) != tt.shouldError {
					t.Errorf("%s operation error = %v, shouldError %v", op, err, tt.shouldError)
				}
			}
		})
	}
}
