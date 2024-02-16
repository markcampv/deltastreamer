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
	Run:   monitorServicesRegistrations,
}

var (
	consulAddr   string
	pollInterval int
	startIndex   uint64
	serviceName  string
	mode         string
)

func init() {
	rootCmd.AddCommand(monitorCmd)

	//Define your flags and configuration settings.
	monitorCmd.Flags().StringVarP(&consulAddr, "consul-addr", "a", "http://localhost:8500", "Address of the Consul server")
	monitorCmd.Flags().IntVarP(&pollInterval, "poll-interval", "p", 10, "Polling interval in seconds")
	monitorCmd.Flags().Uint64Var(&startIndex, "start-index", 0, "Initial index to start watching for changes")
	monitorCmd.Flags().StringVar(&serviceName, "service-name", "", "The name of the service to monitor")
	monitorCmd.Flags().StringVar(&mode, "mode", "service", "Monitoring mode: 'service' for service registration/deregistration, 'instance' for service instances and health")
}

func monitorCommandHandler(cmd *cobra.Command, args []string) {
	client, err := api.NewClient((&api.Config{Address: consulAddr})
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	switch mode {
	case "service":
		monitorServicesRegistrations(client)
	case "instance":
		monitorServiceInstances(client)
	default:
		log.Fatalf("Invalid mode specified: %s. Valid modes are 'service' or 'instance'.")
	}
}

func monitorServicesRegistrations(client *api.Client) {
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

func monitorServiceInstances (client *api.Client) {
	client, err := api.NewClient((&api.Config{Address: consulAddr}))
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	
	lastIndex := startIndex // Use the startIndex specified by the flag
	var previousInstances []*api.ServiceEntry


	fmt.Println("Monitoring instances of service:", serviceName, "at", consulAddr)
	ticker :=  time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Fetching current state of service instances...")
		instances, newIndex, err := fetchServiceInstances(client, serviceName, lastIndex)
		if err != nil {
			log.Printf("Error fetching service instances: %v\n", err)
			continue
		}

		if newIndex != lastIndex {
			logInstanceDeltas(previousInstances, instances) // Use 'instances' directly as fetched
			previousInstances = make([]*api.ServiceEntry, len(instances)) //Reinitialize previousInstances to match instances length
			copy(previousInstances, instances) //Deep copy instances to previousInstances
		}

		// count health instances
		totalInstances := len(instances)
		healthyInstances := 0
		for _, instance := range instances {
			for _, check := range instance.Checks {
				if check.Status == "passing" {
					healthyInstances++
					break // This accounts for on health check per instance for now
				}
			}
		}

		fmt.Printf("Out of %d instances, %d are healthy.\n", totalInstances, healthyInstances)

		lastIndex = newIndex
		fmt.Printf("Last index update to: %d\n", lastIndex)
	}
}

func fetchServiceInstances(client *api.Client, serviceName string, lastIndex uint64) ([]*api.ServiceEntry, uint64, error) {
	queryOpts := &api.QueryOptions{
		WaitIndex: lastIndex,       // Use the last known index to wait for changes
		WaitTime:  5 * time.Minute, // Max wait time; may make a flag out of this
	}

	instances, meta, err := client.Health().Service(serviceName, "", true, queryOpts)
	if err != nil {
		return nil, 0, err
	}

	return instances, meta.LastIndex, nil
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

func logInstanceDeltas(prev, curr []*api.ServiceEntry) {
	prevMap := make(map[string]*api.ServiceEntry)
	currMap := make(map[string]*api.ServiceEntry)

	//Populate the previous instances map
	for _, instance := range prev {
		prevMap[instance.Service.ID] = instance
	}

	// Populate the current instances map and check for additions or health changes
	for _, instance := range curr {
	     currMap[instance.Service.ID] = instance
	     if _, exists := prevMap[instance.Service.ID]; !exists {
			 fmt.Printf("Added instance %s\n", instance.Service.ID)
		 } else {
			 //Check if health status has changed
			 if !healthStatusEquals(prevMap[instance.Service.ID], instance) {
				 fmt.Printf(("Health status changed for instance: %s\n", instance.Service.ID)
			 }
		 }
	}

	// Check for removed instances
	for _, instance := range prev {
		if _, exists := currMap[instance.Service.ID]; !exists {
			fmt.Printf("Removed instance: %s\n", instance.Service.ID)
		}
	}
}




