package cmd

import (
	"errors"
	"os"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator"
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/utils"
	"github.com/luminous-gsm/fusion/variables"
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

	// TODO : Get version from environment variable.

	configPath, found := os.LookupEnv("FUSION_CONFIG_PATH")
	if !found {
		configPath = "fusion.yml"
	}

	var generateNewConfig = false
	// Writes config if files does not exist, or it should be overriden.
	if !fileExist(configPath) || forceConfigOverride {
		zap.S().Infow("configuration file not found,creating one",
			"configPath", configPath,
			"forceConfigOverride", forceConfigOverride,
		)
		generateNewConfig = true
	} else {
		zap.S().Infow("configuration file found", "configPath", configPath)
		generateNewConfig = false
	}

	setupConfigurationOperations(configPath, generateNewConfig)

	zap.S().Info("configuration completed")

}

// Do basic configuration operations
func setupConfigurationOperations(configPath string, generateNewConfig bool) {
	// Step 1 : Create configuration with defaults
	conf, err := getDefaultedConfiguration(configPath)
	if err != nil {
		zap.S().Fatal("default configuration creation error")
	}

	// Set the default config as config to be able to use it in futher setup steps.
	// Mostly for replacing template tokens with values,
	// and then overriding them if they are found in ENV variables.
	config.Set(conf)

	if err = loadCalculatedDefaults(conf); err != nil {
		zap.S().Fatalw("calculated variables loading error", "errors", err)
	}

	if !generateNewConfig {
		// Step 2 : Load configuration values from file
		if err := loadFromFile(conf, configPath); err != nil {
			zap.S().Fatal("configuration file loading error")
		}
	}

	// Step 3 : Load configuration from environment variables
	if err = loadSystemEnvironments(conf); err != nil {
		zap.S().Fatalw("evironment variables loading error", "errors", err)
	}

	// Step 4 : Validate the config for validation errors
	if err = validateConfig(conf); err != nil {
		zap.S().Fatal("configuration valitation error")
	}

	// Step 5 : Write to disk
	if err := writeToDisk(conf, configPath); err != nil {
		zap.S().Fatal("configuration writing error")
	}

	// Store this configuration in the global state.
	config.Set(conf)

	// Refresh all values in variables with the updated, applied and correct values
	variables.Instance().RefreshAllVariables()
}

// Write the configuration file to disk
func writeToDisk(c *config.Configuration, configPath string) error {

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

// Check if the file exist
func fileExist(configPath string) bool {
	if _, err := os.Stat(configPath); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		return false

	} else {
		return false
	}
}

// Get the configuration struct with the default values set
func getDefaultedConfiguration(path string) (*config.Configuration, error) {
	var conf config.Configuration
	// Configures the default values for many of the configuration options present
	// in the structs. Values set in the configuration file take priority over the
	// default values.
	if err := defaults.Set(&conf); err != nil {
		zap.S().Errorw("rrror setting default values", "error", err, "configuration", conf)
		return nil, err
	}

	// Leave this false,
	// if it's true, the server will auto turn on debug mode on configuration generation.
	conf.Debug = false
	conf.Path = path
	return &conf, nil
}

// Load reads the configuration from the provided file
func loadFromFile(conf *config.Configuration, configPath string) error {
	zap.S().Infow("loading configuration from file", "configfile", configPath)

	fileByteArray, err := os.ReadFile(configPath)
	if err != nil {
		zap.S().Errorw("configuration read error",
			"error", err,
			"configPath", configPath,
		)
		return err
	}

	if err := yaml.Unmarshal(fileByteArray, conf); err != nil {
		zap.S().Errorw("yaml parsing error",
			"error", err,
			"configPath", configPath,
		)
		return err
	}

	return nil
}

func loadCalculatedDefaults(conf *config.Configuration) error {
	variables.Instance().RefreshAllVariables()

	// conf.SystemInformation.Os = runtime.GOOS
	conf.SystemInformation.Os = "linux"
	conf.SystemInformation.Arch = runtime.GOARCH

	// Create the default role username and password for the file browser extension.
	systemRoles := make([]config.SystemRole, 0)
	fileBrowserPassword, err := utils.GenerateSecureRandomString(24, true)
	if err != nil {
		return err
	}
	fileBrowserPasswordHashed, err := utils.HashPasswordBasedOnArgon2(fileBrowserPassword, false)
	if err != nil {
		return err
	}
	systemRoles = append(systemRoles, config.SystemRole{
		Username: "filebrowser",
		Password: fileBrowserPasswordHashed,
	})
	conf.SystemRoles = systemRoles

	// Apply all the config values that supports variables.
	conf.NodeUniqueId = variables.Instance().ReplaceGlobalVariablesInString(conf.NodeUniqueId)
	conf.NodeName = variables.Instance().ReplaceGlobalVariablesInString(conf.NodeName)
	conf.ApiSecurityToken = variables.Instance().ReplaceGlobalVariablesInString(conf.ApiSecurityToken)
	conf.NodeHostname = variables.Instance().ReplaceGlobalVariablesInString(conf.NodeHostname)

	return nil
}

// Validate the configuration struct based on specified validations
func validateConfig(conf *config.Configuration) error {
	validate := validator.New()
	// Validate the configuration according to validation tags in the structs.
	if err := validate.Struct(conf); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			zap.S().Errorw("configuration error: please ensure the following field is correct",
				"field", err.Field(),
				"value", err.Value(),
				"validation_type", err.Tag(),
				"field_type", err.Type(),
			)
		}
		return err
	}
	return nil
}

func loadSystemEnvironments(conf *config.Configuration) error {
	if err := env.Parse(conf); err != nil {
		return err
	}
	return nil
}
