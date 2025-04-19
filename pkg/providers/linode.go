package providers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

// LinodeProvider implements the CloudProvider interface for Linode
type LinodeProvider struct {
	config  Provider
	client  linodego.Client
	context context.Context
}

// NewLinodeProvider creates a new Linode provider
func NewLinodeProvider(config Provider) (*LinodeProvider, error) {
	apiKey, ok := config.Credentials["api_key"]
	if !ok {
		return nil, fmt.Errorf("Linode API key not found in credentials")
	}

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})
	oauth2Client := oauth2.NewClient(ctx, tokenSource)

	client := linodego.NewClient(oauth2Client)
	client.SetDebug(false)

	return &LinodeProvider{
		config:  config,
		client:  client,
		context: ctx,
	}, nil
}

// CreateCluster creates a Talos cluster on Linode
func (l *LinodeProvider) CreateCluster(spec ClusterSpec) error {
	// Validate region
	regions, err := l.client.ListRegions(l.context, nil)
	if err != nil {
		return fmt.Errorf("failed to list Linode regions: %v", err)
	}

	validRegion := false
	for _, region := range regions {
		if region.ID == l.config.Region {
			validRegion = true
			break
		}
	}

	if !validRegion {
		return fmt.Errorf("invalid Linode region: %s", l.config.Region)
	}

	// Get Talos image ID (in a real implementation, we would find or upload the proper image)
	// For now, we'll use a Debian image as a placeholder
	images, err := l.client.ListImages(l.context, nil)
	if err != nil {
		return fmt.Errorf("failed to list Linode images: %v", err)
	}

	var imageID string
	for _, image := range images {
		if image.ID == "linode/debian11" {
			imageID = image.ID
			break
		}
	}

	if imageID == "" {
		return fmt.Errorf("Debian 11 image not found")
	}

	// Create nodes
	for i := 0; i < spec.NodeCount; i++ {
		nodeName := fmt.Sprintf("talos-node-%d", i)

		createOpts := linodego.InstanceCreateOptions{
			Region:   l.config.Region,
			Type:     spec.NodeSize,
			Label:    nodeName,
			Image:    imageID,
			RootPass: generateRandomPassword(), // In production, use proper secret management
			Tags:     []string{"talos-autoextender", "talos-node"},
		}

		instance, err := l.client.CreateInstance(l.context, createOpts)
		if err != nil {
			return fmt.Errorf("failed to create Linode instance: %v", err)
		}

		// Wait for instance to boot
		err = waitForInstanceStatus(l.context, &l.client, instance.ID, linodego.InstanceRunning, 300)
		if err != nil {
			return fmt.Errorf("instance failed to start: %v", err)
		}
	}

	return nil
}

// DeleteCluster deletes a Talos cluster from Linode
func (l *LinodeProvider) DeleteCluster(name string) error {
	// Use tags to find the instances rather than trying to use a map as filter
	options := linodego.NewListOptions(0, "")
	instances, err := l.client.ListInstances(l.context, options)
	if err != nil {
		return fmt.Errorf("failed to list instances: %v", err)
	}

	// Filter instances with the talos-autoextender tag
	for _, instance := range instances {
		hasTalosTag := false
		for _, tag := range instance.Tags {
			if tag == "talos-autoextender" {
				hasTalosTag = true
				break
			}
		}

		if hasTalosTag {
			err := l.client.DeleteInstance(l.context, instance.ID)
			if err != nil {
				return fmt.Errorf("failed to delete instance %d: %v", instance.ID, err)
			}
		}
	}

	return nil
}

// GetClusterStatus returns the status of a Talos cluster on Linode
func (l *LinodeProvider) GetClusterStatus(name string) (ClusterStatus, error) {
	// Use tags to find the instances rather than trying to use a map as filter
	options := linodego.NewListOptions(0, "")
	instances, err := l.client.ListInstances(l.context, options)
	if err != nil {
		return ClusterStatus{}, fmt.Errorf("failed to list instances: %v", err)
	}

	// Count instances with the talos-autoextender tag
	var filteredInstances []linodego.Instance
	for _, instance := range instances {
		for _, tag := range instance.Tags {
			if tag == "talos-autoextender" {
				filteredInstances = append(filteredInstances, instance)
				break
			}
		}
	}

	status := ClusterStatus{
		Name:      name,
		NodeCount: len(filteredInstances),
	}

	readyCount := 0
	for _, instance := range filteredInstances {
		if instance.Status == linodego.InstanceRunning {
			readyCount++
		}
	}

	status.ReadyNodeCount = readyCount

	if readyCount == len(filteredInstances) && readyCount > 0 {
		status.State = "ready"
	} else if readyCount > 0 {
		status.State = "partially_ready"
	} else if len(filteredInstances) > 0 {
		status.State = "provisioning"
	} else {
		status.State = "not_found"
	}

	return status, nil
}

// UpdateCluster updates an existing Talos cluster on Linode
func (l *LinodeProvider) UpdateCluster(spec ClusterSpec) error {
	// For now, just stub this out for the tests to compile
	return nil
}

// Helper functions

func waitForInstanceStatus(ctx context.Context, client *linodego.Client, id int, status linodego.InstanceStatus, timeoutSeconds int) error {
	start := time.Now()
	for {
		instance, err := client.GetInstance(ctx, id)
		if err != nil {
			return err
		}

		if instance.Status == status {
			return nil
		}

		if time.Since(start) >= time.Duration(timeoutSeconds)*time.Second {
			return fmt.Errorf("timed out waiting for instance %d to reach status %s", id, status)
		}

		time.Sleep(5 * time.Second)
	}
}

func generateRandomPassword() string {
	// In a real implementation, use a proper random password generator
	// This is just a placeholder
	return "Talos@" + strconv.FormatInt(time.Now().Unix(), 10)
}
