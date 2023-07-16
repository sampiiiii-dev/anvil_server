package cmd

import (
	"github.com/sampiiiii-dev/anvil_server/anvil"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "anvil",
	Short: "Start the anvil server",
	Run:   startAnvil,
}

func startAnvil(cmd *cobra.Command, args []string) {
	a := anvil.Forge()
	a.Run("127.0.0.1:1323")
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
