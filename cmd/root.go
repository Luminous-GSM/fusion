package cmd

import (
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/server"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var (
	configPath = config.DefaultLocation
)

var rootCommand = &cobra.Command{
	Use:   "fusion",
	Short: "Runs the fusion API server, allowing controller nodes from the Luminous console",
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	Run: rootRun,
}

func init() {
	rootCommand.PersistentFlags().StringVar(&configPath, "config", config.DefaultLocation, "Set the location for the configuration file")
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatalf("failed to execute command: %s", err)
	}
}

func rootRun(cmd *cobra.Command, _ []string) {
	server.New()
}

func initConfig() {
	config.Load(configPath)
}
