package cmd

import (
	"github.com/sampiiiii-dev/anvil_server/anvil/config"
	"github.com/sampiiiii-dev/anvil_server/anvil/db"
	"github.com/sampiiiii-dev/anvil_server/anvil/logs"
	"github.com/sampiiiii-dev/anvil_server/anvil/workers"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the background workers",
	Run: func(cmd *cobra.Command, args []string) {
		startWorkers()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}

func startWorkers() {
	s := logs.HireScribe()
	s.Info("Starting workers")
	c := config.GetConfigInstance(s)
	rc := db.InitializeRedisClient(c)
	jq := workers.NewRedisJobQueue(rc, rc.Context())
	wp := workers.NewWorkerPool(5, 100, jq, s)
	wp.Start()
}
