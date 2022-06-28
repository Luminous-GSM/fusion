package config

import (
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const DefaultLocation = "fusion.yaml"

var (
	_config *Configuration
)

// TODO Add validate to normal fields as well, such as only string, or integer
type NodeInformation struct {
	UniqueId    string `validate:"required" yaml:"unique_id" json:"unique_id"`
	Hostname    string `default:"localhost" validate:"alphanum" yaml:"hostname" json:"hostname"`
	Name        string `default:"Fusion" json:"name" yaml:"name"`
	Description string `default:"Node Control Plane" json:"description" yaml:"description"`
}

type ApiSecurity struct {
	Token string `validate:"required" yaml:"token" json:"token"`
}

type ApiConfiguration struct {
	Host     string      `default:"0.0.0.0" yaml:"host"`
	Port     int         `default:"8899" yaml:"port"`
	Security ApiSecurity `yaml:"security" json:"security"`
}

type PodConfiguration struct {
	TmpfsSize         string   `default:"100" yaml:"tmpfs_size" json:"tmpfsSize"`
	ContainerPidLimit int64    `default:"512" yaml:"container_pid_limit" json:"containerPidLimit" `
	Dns               []string `default:"[\"1.1.1.1\", \"1.0.0.1\"]"`
}

type SystemConfiguration struct {
	// The root directory where fusion data is stored.
	RootDirectory string `default:"/var/lib/fusion/" validate:"endswith=/" yaml:"root_directory" json:"root_directory"`

	// Directory where logs and events are logged.
	LogDirectory string `default:"/var/log/fusion/" validate:"endswith=/" yaml:"log_directory" json:"log_directory"`

	// Directory where the server data is stored at.
	DataDirectory string `default:"/var/lib/fusion/volumes/" validate:"endswith=/" yaml:"data_directory" json:"data_directory"`

	User struct {
		Uid int `default:"1000" yaml:"uid" json:"uid"`
		Gid int `default:"1000" yaml:"gid" json:"gid"`
	} `yaml:"user" json:"user"`
}

type Configuration struct {
	path  string
	Debug bool `default:"false" json:"debug" yaml:"debug"`

	System SystemConfiguration `yaml:"system" json:"system"`
	Node   NodeInformation     `yaml:"node" json:"node"`
	Api    ApiConfiguration    `yaml:"api" json:"api"`

	ConsoleLocation string `validate:"required,url|ip" json:"console_location" yaml:"console_location"`

	AllowPrivateNetwork bool `default:"false" json:"allow_private_network" yaml:"allow_private_network"`

	Pod PodConfiguration `yaml:"pod" json:"pod"`
}

func SetDefaults(path string) (*Configuration, error) {
	var c Configuration
	// Configures the default values for many of the configuration options present
	// in the structs. Values set in the configuration file take priority over the
	// default values.
	if err := defaults.Set(&c); err != nil {
		zap.S().Errorw("rrror setting default values", "error", err, "configuration", c)
		return nil, err
	}

	// Leave this false, if it's true,
	// the server will auto turn on debug mode on configuration generation.
	c.Debug = false
	c.path = path
	return &c, nil
}

// Load reads the configuration from the provided file and stores it in the
// global singleton for this node.
func Load(configLocation string) error {
	zap.S().Infow("loading configuration from file", "configfile", configLocation)

	b, err := os.ReadFile(configLocation)
	if err != nil {
		return err
	}
	c, err := SetDefaults(configLocation)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		return err
	}

	if err = ValidateConfig(c); err != nil {
		return err
	}

	zap.S().Debug("running in debug mode")

	// Store this configuration in the global state.
	Set(c)

	zap.S().Info("configuration completed")

	return nil
}

func ValidateConfig(c *Configuration) error {
	validate := validator.New()
	// Validate the configuration according to validation tags in the structs.
	if err := validate.Struct(c); err != nil {
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
