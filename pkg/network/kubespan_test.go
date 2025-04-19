
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

func TestCloudClusterConnection(t *testing.T) {
	manager := NewKubeSpanManager("10.0.0.1:50000")
	
	err := manager.AddCloudCluster("cloud-1", "192.168.1.1:50000")
	if err != nil {
		t.Errorf("AddCloudCluster() error = %v", err)
	}

	if len(manager.CloudClusters) != 1 {
		t.Errorf("Expected 1 cloud cluster, got %d", len(manager.CloudClusters))
	}
}
