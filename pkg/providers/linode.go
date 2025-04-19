package providers

type LinodeProvider struct {
	config Provider
}

func NewLinodeProvider(config Provider) *LinodeProvider {
	return &LinodeProvider{
		config: config,
	}
}

func (l *LinodeProvider) CreateCluster(spec ClusterSpec) error {
	// TODO: Implement actual Linode API calls
	return nil
}

func (l *LinodeProvider) DeleteCluster(name string) error {
	// TODO: Implement actual Linode API calls
	return nil
}

func (l *LinodeProvider) GetClusterStatus(name string) (ClusterStatus, error) {
	// TODO: Implement actual Linode API calls
	return ClusterStatus{}, nil
}