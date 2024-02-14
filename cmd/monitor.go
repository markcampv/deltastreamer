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
		logDetlas(previousservices, currentServices)
		previousServices = currentServices // Update previous state for next comparison
	}
}
