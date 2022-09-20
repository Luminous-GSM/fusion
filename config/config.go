package config

const DefaultLocation = "fusion.yaml"

var (
	_config *Configuration
)

type PodConfiguration struct {
	TmpfsSize         string   `default:"100" yaml:"tmpfs_size" json:"tmpfsSize"`
	ContainerPidLimit int64    `default:"512" yaml:"container_pid_limit" json:"containerPidLimit" `
	Dns               []string `default:"[\"1.1.1.1\", \"1.0.0.1\"]"`
}

type Configuration struct {
	Path            string `default:"config.yml" env:"FUSION_CONFIG_PATH,required" json:"path"`
	Debug           bool   `default:"false" env:"FUSION_DEBUG" json:"debug" yaml:"debug"`
	Version         string `default:"0.1" json:"version" yaml:"version" env:"FUSION_VERSION"`
	HostingPlatform string `default:"self" json:"hostingPlatformType" yaml:"hostingPlatformType" env:"FUSION_HOSTING_PLATFORM_TYPE" validate:"oneof=self aws"`

	// |-----> System CONFIGURATION <-----|
	// The root directory where fusion data is stored.
	RootDirectory string `default:"/var/lib/fusion/" validate:"endswith=/" env:"FUSION_SYSTEM_ROOT_DIRECTORY" yaml:"system_root_directory" json:"system_root_directory"`
	// Directory where logs and events are logged.
	LogDirectory string `default:"/var/log/fusion/" validate:"endswith=/" env:"FUSION_SYSTEM_LOG_DIRECTORY" yaml:"system_log_directory" json:"system_log_directory"`
	// Directory where the server data is stored at.
	DataDirectory  string `default:"/var/lib/fusion/volumes/" validate:"endswith=/" env:"FUSION_SYSTEM_DATA_DIRECTORY" yaml:"system_data_directory" json:"system_data_directory"`
	CertsDirectory string `default:"/var/lib/fusion/certs/" validate:"endswith=/" env:"FUSION_SYSTEM_CERTS_DIRECTORY" yaml:"system_certs_directory" json:"system_certs_directory"`
	// User IDs
	SystemUserUid int `default:"1000" env:"FUSION_SYSTEM_UID" yaml:"system_uid" json:"system_uid"`
	SystemUserGid int `default:"1000" env:"FUSION_SYSTEM_GID" yaml:"system_gid" json:"system_gid"`
	// |-----> System CONFIGURATION <-----|

	// |-----> NODE CONFIGURATION <-----|
	NodeUniqueId    string `validate:"required" env:"FUSION_NODE_UNIQUE_ID,required" yaml:"node_unique_id" json:"node_unique_id"`
	NodeHostname    string `default:"localhost" validate:"alphanum" env:"FUSION_NODE_HOSTNAME" yaml:"node_hostname" json:"node_hostname"`
	NodeName        string `default:"Fusion" env:"FUSION_NODE_NAME" json:"node_name" yaml:"node_name"`
	NodeDescription string `default:"Node Control Plane" env:"FUSION_NODE_DESCRIPTION" json:"node_description" yaml:"node_description"`
	// |-----> NODE CONFIGURATION <-----|

	// |-----> API CONFIGURATION <-----|
	ApiHost          string `default:"0.0.0.0" validate:"url|ip" env:"FUSION_API_HOST" yaml:"api_host" json:"api_host"`
	ApiPort          int    `default:"8899" validate:"numeric" env:"FUSION_API_PORT" yaml:"api_port" json:"api_port"`
	ApiSecurityToken string `validate:"required" env:"FUSION_API_SECURITY_TOKEN,required" yaml:"api_token" json:"api_token"`
	// |-----> API CONFIGURATION <-----|

	ConsoleUrl    string `validate:"required,url|ip" env:"FUSION_CONSOLE_URL,required" json:"console_url" yaml:"console_url"`
	ManagementUrl string `validate:"required,url|ip" env:"FUSION_MANAGEMENT_URL,required" json:"management_url" yaml:"management_url"`

	AllowPrivateNetwork bool `default:"false" env:"FUSION_ALLOW_PRIVATE_NETWORK" json:"allow_private_network" yaml:"allow_private_network"`

	// |-----> Pod CONFIGURATION <-----|
	Pod PodConfiguration `yaml:"pod" json:"pod"`
	// |-----> Pod CONFIGURATION <-----|
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
