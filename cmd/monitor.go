package cmd

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// monitorCmd represents the monitor command

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor changes in service states within Consul",
	Long:  "Monitors and streams deltas of service changes in the Consul cluster, focusing on delivering only the changes rather than the entire payload.",
	Run:   monitorServices,
}

var (
	consulAddr   string
	pollInterval int
)

func init() {
	rootCmd.AddCommand(monitorCmd)

	//Define your flags and configuration settings.
	monitorCmd.Flags().StringVarP(&consulAddr, "consul-addr", "addr", "http://localhost:8500", "Address of the Consul server")
	monitorCmd.Flags().IntVarP(&pollInterval, "poll-interval", "p", 10, "Polling interval in seconds")
}

func monitorServices(cmd *cobra.Command, args []string) {
	client, err := api.NewClient((&api.Config{Address: consulAddr}))
	if err != nil {
		log.Fatalf("Failed to create Consul client: &v", err)
	}

	previousServices := make(map[string]struct{}) // track known services

	fmt.Println("Monitoring services at:", consulAddr)
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Fetching current state of services..")
		currentServices, err := fetchServices(client)
		if err != nil {
			log.Printf("error fetching services: %v\n", err)
			continue
		}
		logDeltas(previousServices, currentServices)
		previousServices = currentServices // Update previous state for next comparison
	}
}

func fetchServices(client *api.Client) (map[string]struct{}, error) {
	catalog := client.Catalog()
	services, _, err := catalog.Services(nil)
	if err != nil {
		return nil, err
	}

	serviceMap := make(map[string]struct{})
	for serviceName := range services {
		serviceMap[serviceName] = struct{}{}
	}

	return serviceMap, nil
}

func logDeltas(prev, current map[string]struct{}) {
	// Log added services
	for service := range current {
		if _, exists := prev[service]; !exists {
			fmt.Printf("Service added: %s\n", service)
		}
	}

	// Log removed services
	for service := range prev {
		if _, exists := current[service]; !exists {
			fmt.Printf("Service removed: %s\n", service)
		}
	}
}
