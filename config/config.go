package config

import (
	"encoding/json"
	"os"

	"github.com/apex/log"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator"

	// "github.com/goccy/go-json"
	"gopkg.in/yaml.v2"
)

const defaultLocation = "fusion.yaml"

var (
	_config *Configuration
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type NodeInformation struct {
	UniqueId    string `validate:"required" yaml:"unique_id" json:"unique_id"`
	Hostname    string `validate:"required" yaml:"hostname" json:"hostname"`
	Name        string `default:"Fusion" json:"name" yaml:"name"`
	Description string `default:"Node Control Plane" json:"description" yaml:"description"`
}

type ApiSecurity struct {
	Token string `yaml:"token" json:"token" validate:"required"`
}

type ApiConfiguration struct {
	Host     string      `default:"0.0.0.0" yaml:"host"`
	Port     int         `default:"8899" yaml:"port"`
	Security ApiSecurity `yaml:"security" json:"security"`
}

type Configuration struct {
	path  string
	Debug bool `default:"false" json:"debug" yaml:"debug"`

	Node NodeInformation  `yaml:"node" json:"node"`
	Api  ApiConfiguration `yaml:"api" json:"api"`

	ConsoleLocation string `validate:"required,url|ip" json:"console_location" yaml:"console_location"`

	AllowPrivateNetwork bool `default:"false" json:"allow_private_network" yaml:"allow_private_network"`
}

func SetDefaults(path string) (*Configuration, error) {
	var c Configuration
	// Configures the default values for many of the configuration options present
	// in the structs. Values set in the configuration file take priority over the
	// default values.
	if err := defaults.Set(&c); err != nil {
		return nil, err
	}
	c.path = path
	return &c, nil
}

// Load reads the configuration from the provided file and stores it in the
// global singleton for this node.
func Load() error {
	log.WithField("config_file", defaultLocation).Info("loading configuration from file")

	validate = validator.New()

	b, err := os.ReadFile(defaultLocation)
	if err != nil {
		return err
	}
	c, err := SetDefaults(defaultLocation)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		return err
	}

	if c.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("running in debug mode")

	// Validate the configuration according to validation tags in the structs.
	if err := validate.Struct(c); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			log.WithFields(
				log.Fields{
					"Field":          err.Field(),
					"Value":          err.Value(),
					"ValidationType": err.Tag(),
					"FieldType":      err.Type(),
				}).Error("Configuration Error: Please ensure the following field is correct.")

		}
		return err
	}

	// Store this configuration in the global state.
	Set(c)

	// Print the current configuration
	printConfig()

	return nil
}

func printConfig() {
	config_marshalled, _ := json.MarshalIndent(_config, "", "	")

	log.Debug(string(config_marshalled))
}

// Set the global configuration instance.
func Set(c *Configuration) {
	_config = c
}

// Get returns the global configuration instance.
// Be aware that you CANNOT make modifications to the currently stored configuration
// by modifying the struct returned by this function.
func Get() *Configuration {
	c := *_config
	return &c
}
