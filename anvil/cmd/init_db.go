package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Run database initialization",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Making database")

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
