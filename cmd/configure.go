package cmd

import (
	"errors"
	"os"
	"strconv"

	"github.com/luminous-gsm/fusion/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	forceConfigOverride bool
)

var configureCommand = &cobra.Command{
	Use:   "configure",
	Short: "Configure fusion based on environment variables",
	Run:   configureRun,
}

func init() {
	configureCommand.PersistentFlags().BoolVar(&forceConfigOverride, "force", false, "Should override configuration if file exist")
}

func configureRun(cmd *cobra.Command, _ []string) {
	Configure()
}

func Configure() {
	// Writes config if files does not exist, or it should be overriden.
	if !FileExist(configPath) || forceConfigOverride {
		zap.S().Infow("configuration file not found",
			"configPath", configPath,
			"forceConfigOverride", forceConfigOverride,
		)

		// Required environment variables
		apiSecretToken := getEnvironment("FUSION_API_SECRET_TOKEN", true)
		uniqueId := getEnvironment("FUSION_UNIQUE_ID", true)
		consoleLocation := getEnvironment("FUSION_CONSOLE_LOCATION", true)

		// Optional environment variables
		hostname := getEnvironment("FUSION_HOSTNAME", false)
		name := getEnvironment("FUSION_NAME", false)
		description := getEnvironment("FUSION_DESCRIPTION", false)
		logDirectory := getEnvironment("FUSION_LOG_DIRECTORY", false)
		dataDirectory := getEnvironment("FUSION_DATA_DIRECTORY", false)
		rootDirectory := getEnvironment("FUSION_ROOT_DIRECTORY", false)

		fusionIp := getEnvironment("FUSION_API_HOST", false)
		fusionPortStr := getEnvironment("FUSION_API_PORT", false)
		var fusionPort int
		if fusionPortStr != "" {
			if fusionPort, err := strconv.Atoi(fusionPortStr); err != nil {
				zap.S().Panicw("FUSION_API_PORT is not a number", "FUSION_API_PORT", fusionPort)
			}
		} else {
			fusionPort = 0
		}

		conf, err := config.SetDefaults(configPath)
		if err != nil {
			zap.S().Panicw("could not set configuration defaults", "configPath", configPath)
		}

		conf.ConsoleLocation = consoleLocation
		conf.Node.UniqueId = uniqueId
		conf.Api.Security.Token = apiSecretToken

		if hostname != "" {
			conf.Node.Hostname = hostname
		}
		if name != "" {
			conf.Node.Name = name
		}
		if description != "" {
			conf.Node.Description = description
		}
		if logDirectory != "" {
			conf.System.LogDirectory = logDirectory
		}
		if dataDirectory != "" {
			conf.System.DataDirectory = dataDirectory
		}
		if rootDirectory != "" {
			conf.System.RootDirectory = rootDirectory
		}
		if fusionIp != "" {
			conf.Api.Host = fusionIp
		}
		if fusionPort != 0 {
			conf.Api.Port = fusionPort
		}

		if err = config.ValidateConfig(conf); err != nil {
			zap.S().Panicw("could not set configuration defaults",
				"error", err,
				"configPath", configPath,
				"config", conf,
			)
		}

		WriteToDisk(conf, configPath)
	} else {
		zap.S().Infow("configuration file found", "configPath", configPath)
	}

}

func getEnvironment(envVariable string, fatal bool) string {
	value, ok := os.LookupEnv(envVariable)
	if !ok {
		if fatal {
			zap.S().Fatalf("%v environment variable not found, exiting now", envVariable)
		} else {
			zap.S().Infof("%v environment variable not found, using default value", envVariable)
		}

	}
	return value
}

func WriteToDisk(c *config.Configuration, configPath string) error {

	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configPath, b, 0o600); err != nil {
		return err
	}

	zap.S().Infow("created configuration file", "configPath", configPath)
	return nil
}

func FileExist(configPath string) bool {
	if _, err := os.Stat(configPath); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		return false

	} else {
		return false
	}
}
