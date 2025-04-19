package providers

// ProviderFactory creates cloud provider implementations
type ProviderFactory interface {
	CreateProvider(config Provider) (CloudProvider, error)
}

// CloudProvider defines common operations for cloud providers
type CloudProvider interface {
	CreateCluster(spec ClusterSpec) error
	DeleteCluster(name string) error
	GetClusterStatus(name string) (ClusterStatus, error)
	UpdateCluster(spec ClusterSpec) error
}

// ClusterStatus represents the current state of a cluster
type ClusterStatus struct {
	Name           string
	State          string // provisioning, ready, error, etc.
	NodeCount      int
	ReadyNodeCount int
	Error          string
}
