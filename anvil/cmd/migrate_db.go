package cmd

import (
	"github.com/sampiiiii-dev/anvil_server/anvil/db"
	"github.com/spf13/cobra"
)

var migrateDbCmd = &cobra.Command{
	Use:   "migratedb",
	Short: "Migrate the database",
	Run: func(cmd *cobra.Command, args []string) {
		migrateDatabase()
	},
}

func init() {
	rootCmd.AddCommand(migrateDbCmd)
}

func migrateDatabase() {
	err := db.MigrateDB(db.GetDBInstance())
	if err != nil {
		return
	}
}
