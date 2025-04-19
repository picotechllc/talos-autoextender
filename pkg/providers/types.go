
package providers

// Provider represents a cloud provider (Linode, Hetzner etc)
type Provider struct {
	Name        string
	Region      string
	Credentials map[string]string
}

func (p *Provider) Validate() error {
	if p.Name == "" || p.Region == "" {
		return fmt.Errorf("provider name and region are required")
	}
	if len(p.Credentials) == 0 {
		return fmt.Errorf("provider credentials are required")
	}
	return nil
}

// ClusterSpec defines the desired cluster state
type ClusterSpec struct {
	NodeCount    int
	NodeSize     string
	TalosVersion string
}

func (s *ClusterSpec) Validate() error {
	if s.NodeCount < 1 {
		return fmt.Errorf("node count must be at least 1")
	}
	if s.NodeSize == "" {
		return fmt.Errorf("node size is required")
	}
	if s.TalosVersion == "" {
		return fmt.Errorf("talos version is required")
	}
	return nil
}
