
package network

import (
	"fmt"
	"net"
)

// KubeSpanManager handles mesh networking between clusters
type KubeSpanManager struct {
	HomeClusterEndpoint string
	CloudClusters       []string
}

func NewKubeSpanManager(homeEndpoint string) *KubeSpanManager {
	return &KubeSpanManager{
		HomeClusterEndpoint: homeEndpoint,
		CloudClusters:       make([]string, 0),
	}
}

func (k *KubeSpanManager) ValidateEndpoint() error {
	if k.HomeClusterEndpoint == "" {
		return fmt.Errorf("home cluster endpoint is required")
	}
	
	host, port, err := net.SplitHostPort(k.HomeClusterEndpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint format: %v", err)
	}

	if host == "" || port == "" {
		return fmt.Errorf("both host and port are required")
	}

	return nil
}

func (k *KubeSpanManager) AddCloudCluster(name, endpoint string) error {
	if name == "" || endpoint == "" {
		return fmt.Errorf("both name and endpoint are required")
	}

	_, _, err := net.SplitHostPort(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint format: %v", err)
	}

	k.CloudClusters = append(k.CloudClusters, endpoint)
	return nil
}
