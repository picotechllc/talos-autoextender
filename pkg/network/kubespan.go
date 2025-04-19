
package network

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
