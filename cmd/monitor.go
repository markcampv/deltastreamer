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
	startIndex   uint64
)

func init() {
	rootCmd.AddCommand(monitorCmd)

	//Define your flags and configuration settings.
	monitorCmd.Flags().StringVarP(&consulAddr, "consul-addr", "a", "http://localhost:8500", "Address of the Consul server")
	monitorCmd.Flags().IntVarP(&pollInterval, "poll-interval", "p", 10, "Polling interval in seconds")
	monitorCmd.Flags().Uint64Var(&startIndex, "start-index", 0, "Initial index to start watching for changes")
}

func monitorServices(cmd *cobra.Command, args []string) {
	client, err := api.NewClient((&api.Config{Address: consulAddr}))
	if err != nil {
		log.Fatalf("Failed to create Consul client: &v", err)
	}

	previousServices := make(map[string]struct{}) // track known services
	lastIndex := startIndex                       //Use the startIndex specified by the flag

	fmt.Println("Monitoring services at:", consulAddr)
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Fetching current state of services..")
		currentServices, newIndex, err := fetchServices(client, lastIndex)
		if err != nil {
			log.Printf("error fetching services: %v\n", err)
			continue
		}
		if newIndex != lastIndex {
			logDeltas(previousServices, currentServices)
			previousServices = currentServices // Update previous state for next comparison
		}
		lastIndex = newIndex
		fmt.Printf("Last index update to: %d\n", lastIndex)
	}
}

func fetchServices(client *api.Client, lastIndex uint64) (map[string]struct{}, uint64, error) {
	queryOpts := &api.QueryOptions{
		WaitIndex: lastIndex,       // Use the last known index to wait for changes
		WaitTime:  5 * time.Minute, // Max wait time; may make a flag out of this
	}

	services, meta, err := client.Catalog().Services(queryOpts)
	if err != nil {
		return nil, 0, err
	}

	serviceMap := make(map[string]struct{})
	for serviceName := range services {
		serviceMap[serviceName] = struct{}{}
	}

	return serviceMap, meta.LastIndex, nil
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
