package cmd

import (
	"github.com/sampiiiii-dev/anvil_server/anvil/db"
	"github.com/spf13/cobra"
)

var initDbCmd = &cobra.Command{
	Use:   "initdb",
	Short: "Initialize the database",
	Run: func(cmd *cobra.Command, args []string) {
		initializeDatabase()
	},
}

func init() {
	rootCmd.AddCommand(initDbCmd)
}

func initializeDatabase() {
	db.GetDBInstance()
}
