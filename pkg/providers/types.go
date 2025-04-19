
package providers

// Provider represents a cloud provider (Linode, Hetzner etc)
type Provider struct {
	Name       string
	Region     string
	Credentials map[string]string
}

// ClusterSpec defines the desired cluster state
type ClusterSpec struct {
	NodeCount  int
	NodeSize   string
	TalosVersion string
}
