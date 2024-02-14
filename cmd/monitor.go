package cmd

import (
	"dev/nomad/api"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// monitorCmd represents the monitor command

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor changes in service states within Consul",
	Long:  "Monitors and streams deltas of service changes in the Consul cluster, focusing on delivering only the changes rather than the entire payload.",
	Run: func(cmd *cobra.Command, args []string) {
		consulAddr, _ := cmd.Flags().GetString("consul-addr")
		pollInterval, _ := cmd.Flags().GetInt("poll-interval")
		monitorServices(consulAddr, pollInterval)
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	//Define your flags and configuration settings.
	monitorCmd.Flags().String("consul-addr", "http://localhost:8500", "Consul address")
	monitorCmd.Flags().Int("poll-interval", 10, "Polling interval in seconds")
}

func monitorServices(consulAddr string, interval int) {
	// Initialize Consul client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddr
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal("error creating Consul client: %v", err)
	}

	fmt.Println("Monitoring services at:", consulAddr)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Fetching current state of services..")
		// Placeholder for actual monitoring logic
	}
}
