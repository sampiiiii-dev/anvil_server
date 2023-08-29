package cmd

import (
	"github.com/sampiiiii-dev/anvil_server/anvil"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Anvil server",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func startServer() {
	// Initialize and start the Anvil server
	a := anvil.Forge()
	a.Run()
}
