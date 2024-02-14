package main

import (
	"fmt"
	"log"
	// Assume "consul/api" is the Consul API client library
	"github.com/hashicorp/consul/api"
)

func main() {
	// Initialize Consul client
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	//  print services
	services, _, err := consulClient.Catalog().Services(nil)
	if err != nil {
		log.Fatalf("error querying Consul services: %v", err)
	}

	fmt.Println("Services", services)
}
