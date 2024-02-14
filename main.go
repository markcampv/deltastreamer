package main

import (
	"github.com/markcampv/deltastreamer/cmd"
	"log"
)

func main() {
	// Execute the root command defined in the cmd package
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing DeltaStreamer: %v", err)
	}

}
