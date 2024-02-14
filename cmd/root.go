package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deltastreamer",
	Short: "A tool to monitor and stream deltas of service changes in Consul",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Define any global flags here
}
