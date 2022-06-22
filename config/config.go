package config

import (
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

const DefaultLocation = "fusion.yaml"

var (
	_config *Configuration
)

type ApiConfiguration struct {
	Host string `default:"0.0.0.0" yaml:"host"`
	Port int    `default:"8899" yaml:"port"`
}

type Configuration struct {
	// The location from which this configuration instance was instantiated.
	path string

	// Determines if wings should be running in debug mode. This value is ignored
	// if the debug flag is passed through the command line arguments.
	Debug bool `default:"false" json:"debug" yaml:"debug"`

	AppName string `default:"Fusion" json:"app_name" yaml:"app_name"`

	// A unique identifier for this node in Luminous.
	Uuid string

	Api ApiConfiguration `json:"api" yaml:"api"`
}

// NewAtPath creates a new struct and set the path where it should be stored.
// This function does not modify the currently stored global configuration.
func SetDefaults(path string) (*Configuration, error) {
	var c Configuration
	// Configures the default values for many of the configuration options present
	// in the structs. Values set in the configuration file take priority over the
	// default values.
	if err := defaults.Set(&c); err != nil {
		return nil, err
	}
	// Track the location where we created this configuration.
	c.path = path
	return &c, nil
}

// Load reads the configuration from the provided file and stores it in the
// global singleton for this instance.
func Load() error {
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

	// Store this configuration in the global state.
	Set(c)
	return nil
}

// Set the global configuration instance. This is a blocking operation such that
// anything trying to set a different configuration value, or read the configuration
// will be paused until it is complete.
func Set(c *Configuration) {
	// mu.Lock()
	// if _config == nil || _config.AuthenticationToken != c.AuthenticationToken {
	// 	_jwtAlgo = jwt.NewHS256([]byte(c.AuthenticationToken))
	// }
	_config = c
	// mu.Unlock()
}

// Get returns the global configuration instance. This is a thread-safe operation
// that will block if the configuration is presently being modified.
//
// Be aware that you CANNOT make modifications to the currently stored configuration
// by modifying the struct returned by this function. The only way to make
// modifications is by using the Update() function and passing data through in
// the callback.
func Get() *Configuration {
	// mu.RLock()
	// Create a copy of the struct so that all modifications made beyond this
	// point are immutable.
	//goland:noinspection GoVetCopyLock
	c := *_config
	// mu.RUnlock()
	return &c
}
