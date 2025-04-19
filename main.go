package main

import (
	"fmt"
	"os"

	"talos-autoextender/pkg/dns"
	"talos-autoextender/pkg/network"
	"talos-autoextender/pkg/providers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "talos-autoextender",
	Short: "A tool to extend Talos clusters into the cloud",
	Long: `talos-autoextender is a tool that helps bridge a home-based Talos cluster
to one or more cloud-based Talos clusters.

It enables secure, seamless exposure of home-hosted services to the
internet via cloud ingress, with dynamic DNS, automated provisioning,
and minimal manual intervention.`,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a cloud cluster extension",
	Long:  `Create a new Talos cluster in the cloud to extend your home cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		nodeCount, _ := cmd.Flags().GetInt("nodes")
		nodeSize, _ := cmd.Flags().GetString("size")
		talosVersion, _ := cmd.Flags().GetString("talos-version")
		apiKey, _ := cmd.Flags().GetString("api-key")

		fmt.Printf("Creating cluster with provider %s in region %s with %d nodes of size %s\n",
			provider, region, nodeCount, nodeSize)

		// Initialize provider configuration
		providerConfig := providers.Provider{
			Name:   provider,
			Region: region,
			Credentials: map[string]string{
				"api_key": apiKey,
			},
		}

		// Create cluster specification
		spec := providers.ClusterSpec{
			NodeCount:    nodeCount,
			NodeSize:     nodeSize,
			TalosVersion: talosVersion,
		}

		// Create provider factory
		factory := providers.NewProviderFactory()

		// Create cloud provider
		cloudProvider, err := factory.CreateProvider(providerConfig)
		if err != nil {
			fmt.Printf("Error creating provider: %v\n", err)
			return
		}

		// Create the cluster
		err = cloudProvider.CreateCluster(spec)
		if err != nil {
			fmt.Printf("Error creating cluster: %v\n", err)
			return
		}

		fmt.Println("Cluster created successfully")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cloud cluster extension",
	Long:  `Delete an existing Talos cluster in the cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		clusterName, _ := cmd.Flags().GetString("name")
		apiKey, _ := cmd.Flags().GetString("api-key")

		fmt.Printf("Deleting cluster %s with provider %s in region %s\n",
			clusterName, provider, region)

		// Initialize provider configuration
		providerConfig := providers.Provider{
			Name:   provider,
			Region: region,
			Credentials: map[string]string{
				"api_key": apiKey,
			},
		}

		// Create provider factory
		factory := providers.NewProviderFactory()

		// Create cloud provider
		cloudProvider, err := factory.CreateProvider(providerConfig)
		if err != nil {
			fmt.Printf("Error creating provider: %v\n", err)
			return
		}

		// Delete the cluster
		err = cloudProvider.DeleteCluster(clusterName)
		if err != nil {
			fmt.Printf("Error deleting cluster: %v\n", err)
			return
		}

		fmt.Println("Cluster deleted successfully")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a cloud cluster extension",
	Long:  `Get the status of an existing Talos cluster in the cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		clusterName, _ := cmd.Flags().GetString("name")
		apiKey, _ := cmd.Flags().GetString("api-key")

		fmt.Printf("Getting status for cluster %s with provider %s in region %s\n",
			clusterName, provider, region)

		// Initialize provider configuration
		providerConfig := providers.Provider{
			Name:   provider,
			Region: region,
			Credentials: map[string]string{
				"api_key": apiKey,
			},
		}

		// Create provider factory
		factory := providers.NewProviderFactory()

		// Create cloud provider
		cloudProvider, err := factory.CreateProvider(providerConfig)
		if err != nil {
			fmt.Printf("Error creating provider: %v\n", err)
			return
		}

		// Get cluster status
		status, err := cloudProvider.GetClusterStatus(clusterName)
		if err != nil {
			fmt.Printf("Error getting cluster status: %v\n", err)
			return
		}

		fmt.Printf("Cluster %s status:\n", clusterName)
		fmt.Printf("  State: %s\n", status.State)
		fmt.Printf("  Nodes: %d (%d ready)\n", status.NodeCount, status.ReadyNodeCount)
		if status.Error != "" {
			fmt.Printf("  Error: %s\n", status.Error)
		}
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect a cloud cluster to home cluster",
	Long:  `Connect a cloud cluster to your home Talos cluster using KubeSpan.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeEndpoint, _ := cmd.Flags().GetString("home-endpoint")
		cloudEndpoint, _ := cmd.Flags().GetString("cloud-endpoint")
		clusterName, _ := cmd.Flags().GetString("name")

		fmt.Printf("Connecting cluster %s at %s to home cluster at %s\n",
			clusterName, cloudEndpoint, homeEndpoint)

		// Create KubeSpan manager
		manager := network.NewKubeSpanManager(homeEndpoint)

		// Validate the home endpoint
		if err := manager.ValidateEndpoint(); err != nil {
			fmt.Printf("Invalid home endpoint: %v\n", err)
			return
		}

		// Add the cloud cluster
		if err := manager.AddCloudCluster(clusterName, cloudEndpoint); err != nil {
			fmt.Printf("Error adding cloud cluster: %v\n", err)
			return
		}

		fmt.Println("Cloud cluster connected successfully")
	},
}

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS records for the cluster",
	Long:  `Manage DNS records for exposing services from your home cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		domain, _ := cmd.Flags().GetString("domain")
		recordName, _ := cmd.Flags().GetString("record")
		recordType, _ := cmd.Flags().GetString("type")
		content, _ := cmd.Flags().GetString("content")

		fmt.Printf("Managing DNS record %s.%s of type %s with provider %s\n",
			recordName, domain, recordType, provider)

		// Create DNS manager
		manager := dns.NewDNSManager(provider, domain)

		// Validate configuration
		if err := manager.ValidateConfig(); err != nil {
			fmt.Printf("Invalid DNS configuration: %v\n", err)
			return
		}

		// Upsert DNS record
		if err := manager.UpsertRecord(recordName, recordType, content); err != nil {
			fmt.Printf("Error upserting DNS record: %v\n", err)
			return
		}

		fmt.Println("DNS record updated successfully")
	},
}

func init() {
	// Create command flags
	createCmd.Flags().String("provider", "linode", "Cloud provider to use (linode, hetzner)")
	createCmd.Flags().String("region", "us-east", "Region to deploy the cluster")
	createCmd.Flags().Int("nodes", 3, "Number of nodes in the cluster")
	createCmd.Flags().String("size", "g6-standard-2", "Size/type of the nodes")
	createCmd.Flags().String("talos-version", "v1.6.0", "Talos version to use")
	createCmd.Flags().String("api-key", "", "API key for the cloud provider")

	// Delete command flags
	deleteCmd.Flags().String("provider", "linode", "Cloud provider to use (linode, hetzner)")
	deleteCmd.Flags().String("region", "us-east", "Region of the cluster")
	deleteCmd.Flags().String("name", "", "Name of the cluster to delete")
	deleteCmd.Flags().String("api-key", "", "API key for the cloud provider")

	// Status command flags
	statusCmd.Flags().String("provider", "linode", "Cloud provider to use (linode, hetzner)")
	statusCmd.Flags().String("region", "us-east", "Region of the cluster")
	statusCmd.Flags().String("name", "", "Name of the cluster to get status for")
	statusCmd.Flags().String("api-key", "", "API key for the cloud provider")

	// Connect command flags
	connectCmd.Flags().String("home-endpoint", "", "Endpoint of the home cluster (IP:PORT)")
	connectCmd.Flags().String("cloud-endpoint", "", "Endpoint of the cloud cluster (IP:PORT)")
	connectCmd.Flags().String("name", "", "Name of the cloud cluster")

	// DNS command flags
	dnsCmd.Flags().String("provider", "cloudflare", "DNS provider to use")
	dnsCmd.Flags().String("domain", "", "Domain to manage records for")
	dnsCmd.Flags().String("record", "", "Record name (e.g., www)")
	dnsCmd.Flags().String("type", "A", "Record type (A, CNAME, etc.)")
	dnsCmd.Flags().String("content", "", "Record content (e.g., IP address)")

	// Add commands to root command
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(dnsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
