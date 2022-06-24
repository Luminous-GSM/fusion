package cmd

import (
	"fmt"

	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/environment"
	"github.com/luminous-gsm/fusion/router"
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
	cfg := config.Get()
	env, err := environment.NewEnvironment()
	if err != nil {
		panic("EXPAND THIS")
	}
	mgr, err := server.NewManager(env)
	if err != nil {
		panic("EXPAND THIS")
	}
	router := router.NewRouter(mgr)

	port := fmt.Sprintf("%v:%v", cfg.Api.Host, cfg.Api.Port)
	log.Infof("cmd: started API server on %v", port)
	router.Run(port)
}

func initConfig() {
	config.Load(configPath)
}
