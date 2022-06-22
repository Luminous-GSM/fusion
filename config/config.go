package config

import (
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"
)

const DefaultLocation = "fusion.yaml"

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
// global singleton for this instance.
func Load() error {
	validate = validator.New()

	b, err := os.ReadFile(DefaultLocation)
	if err != nil {
		return err
	}
	c, err := SetDefaults(DefaultLocation)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		return err
	}

	if err := validate.Struct(c); err != nil {
		return err
	}

	// Store this configuration in the global state.
	Set(c)
	return nil
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
