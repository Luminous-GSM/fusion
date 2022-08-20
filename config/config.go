package config

const DefaultLocation = "fusion.yaml"

var (
	_config *Configuration
)

// TODO Add validate to normal fields as well, such as only string, or integer
type NodeInformation struct {
	UniqueId    string `validate:"required" env:"FUSION_NODE_UNIQUE_ID,required" yaml:"unique_id" json:"unique_id"`
	Hostname    string `default:"localhost" validate:"alphanum" env:"FUSION_NODE_HOSTNAME" yaml:"hostname" json:"hostname"`
	Name        string `default:"Fusion" env:"FUSION_NODE_NAME" json:"name" yaml:"name"`
	Description string `default:"Node Control Plane" env:"FUSION_NODE_DESCRIPTION" json:"description" yaml:"description"`
}

type ApiSecurity struct {
	Token string `validate:"required" env:"FUSION_API_TOKEN,required" yaml:"token" json:"token"`
}

type ApiConfiguration struct {
	Host     string      `default:"0.0.0.0" validate:"url|ip" env:"FUSION_API_HOST" yaml:"host" json:"host"`
	Port     int         `default:"8899" validate:"numeric" env:"FUSION_API_PORT" yaml:"port" json:"port"`
	Security ApiSecurity `yaml:"security" json:"security"`
}

type PodConfiguration struct {
	TmpfsSize         string   `default:"100" yaml:"tmpfs_size" json:"tmpfsSize"`
	ContainerPidLimit int64    `default:"512" yaml:"container_pid_limit" json:"containerPidLimit" `
	Dns               []string `default:"[\"1.1.1.1\", \"1.0.0.1\"]"`
}

type SystemConfiguration struct {
	// The root directory where fusion data is stored.
	RootDirectory string `default:"/var/lib/fusion/" validate:"endswith=/" env:"FUSION_ROOT_DIRECTORY" yaml:"root_directory" json:"root_directory"`

	// Directory where logs and events are logged.
	LogDirectory string `default:"/var/log/fusion/" validate:"endswith=/" env:"FUSION_LOG_DIRECTORY" yaml:"log_directory" json:"log_directory"`

	// Directory where the server data is stored at.
	DataDirectory string `default:"/var/lib/fusion/volumes/" validate:"endswith=/" env:"FUSION_DATA_DIRECTORY" yaml:"data_directory" json:"data_directory"`

	User struct {
		Uid int `default:"1000" env:"FUSION_UID" yaml:"uid" json:"uid"`
		Gid int `default:"1000" env:"FUSION_GID" yaml:"gid" json:"gid"`
	} `yaml:"user" json:"user"`
}

type Configuration struct {
	Path  string `default:"config.yml" env:"FUSION_CONFIG_PATH,required"`
	Debug bool   `default:"false" json:"debug" yaml:"debug"`

	System SystemConfiguration `yaml:"system" json:"system"`
	Node   NodeInformation     `yaml:"node" json:"node"`
	Api    ApiConfiguration    `yaml:"api" json:"api"`

	ConsoleLocation string `validate:"required,url|ip" env:"FUSION_CONSOLE_LOCATION,required" json:"console_location" yaml:"console_location"`

	AllowPrivateNetwork bool `default:"false" env:"FUSION_ALLOW_PRIVATE_NETWORK" json:"allow_private_network" yaml:"allow_private_network"`

	Pod PodConfiguration `yaml:"pod" json:"pod"`
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
