package config

var (
	_config *Configuration
)

type PodConfiguration struct {
	TmpfsSize         string   `default:"100" yaml:"tmpfs_size" json:"tmpfsSize"`
	ContainerPidLimit int64    `default:"512" yaml:"container_pid_limit" json:"containerPidLimit" `
	Dns               []string `default:"[\"1.1.1.1\", \"1.0.0.1\"]" json:"dns"`
}

type SystemInformation struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
}

type SystemRole struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Configuration struct {
	Path            string `default:"config.yml" env:"FUSION_CONFIG_PATH,required" json:"path"`
	Debug           bool   `default:"false" env:"FUSION_DEBUG" json:"debug" yaml:"debug"`
	Version         string `default:"0.1" json:"version" yaml:"version" env:"FUSION_VERSION"`
	HostingPlatform string `default:"self" json:"hostingPlatformType" yaml:"hostingPlatformType" env:"FUSION_HOSTING_PLATFORM_TYPE" validate:"oneof=self aws"`

	// |-----> System CONFIGURATION <-----|
	// The root directory where fusion data is stored.
	RootDirectory string `default:"/var/lib/fusion/" validate:"" env:"FUSION_SYSTEM_ROOT_DIRECTORY" yaml:"system_root_directory" json:"systemRootDirectory"`
	// Directory where logs and events are logged.
	LogDirectory string `default:"/var/lib/fusion/logs" validate:"" env:"FUSION_SYSTEM_LOG_DIRECTORY" yaml:"system_log_directory" json:"systemLogDirectory"`
	// Directory where the server data is stored at.
	DataDirectory  string `default:"/var/lib/fusion/data" validate:"" env:"FUSION_SYSTEM_DATA_DIRECTORY" yaml:"system_data_directory" json:"systemDataDirectory"`
	CertsDirectory string `default:"/var/lib/fusion/certs" validate:"" env:"FUSION_SYSTEM_CERTS_DIRECTORY" yaml:"system_certs_directory" json:"systemCertsDirectory"`
	// User IDs
	SystemUserUid int `default:"1000" env:"FUSION_SYSTEM_UID" yaml:"system_uid" json:"systemUid"`
	SystemUserGid int `default:"1000" env:"FUSION_SYSTEM_GID" yaml:"system_gid" json:"systemGid"`
	// |-----> System CONFIGURATION <-----|

	// |-----> NODE CONFIGURATION <-----|
	NodeUniqueId    string `default:"{{fusion.generated.id}}" validate:"" env:"FUSION_NODE_UNIQUE_ID" yaml:"node_unique_id" json:"nodeUniqueId"`
	NodeHostname    string `default:"localhost" validate:"alphanum" env:"FUSION_NODE_HOSTNAME" yaml:"node_hostname" json:"nodeHostname"`
	NodeName        string `default:"{{fusion.generated.name}}" env:"FUSION_NODE_NAME" json:"nodeName" yaml:"node_name"`
	NodeDescription string `default:"Fusion Node Control Agent" env:"FUSION_NODE_DESCRIPTION" json:"nodeDescription" yaml:"node_description"`
	// |-----> NODE CONFIGURATION <-----|

	// |-----> API CONFIGURATION <-----|
	ApiHost          string `default:"0.0.0.0" validate:"url|ip" env:"FUSION_API_HOST" yaml:"api_host" json:"apiHost"`
	ApiPort          int    `default:"8899" validate:"numeric" env:"FUSION_API_PORT" yaml:"api_port" json:"apiPort"`
	ApiSecurityToken string `default:"{{fusion.generated.password}}" validate:"" env:"FUSION_API_SECURITY_TOKEN" yaml:"api_token" json:"apiToken"`
	// |-----> API CONFIGURATION <-----|

	ConsoleUrl    string `default:"https://localhost:3000" validate:"required,url|ip" env:"FUSION_CONSOLE_URL" json:"consoleUrl" yaml:"console_url"`
	ManagementUrl string `default:"https://localhost:8898" validate:"required,url|ip" env:"FUSION_MANAGEMENT_URL" json:"managementUrl" yaml:"management_url"`

	AllowPrivateNetwork bool `default:"false" env:"FUSION_ALLOW_PRIVATE_NETWORK" json:"allowPrivateNetwork" yaml:"allow_private_network"`

	// |-----> Pod CONFIGURATION <-----|
	Pod PodConfiguration `yaml:"pod" json:"pod"`
	// |-----> Pod CONFIGURATION <-----|

	// |-----> System Information <-----|
	SystemInformation SystemInformation `yaml:"system_information" json:"systemInformation"`
	// |-----> System Information <-----|

	SystemRoles []SystemRole `yaml:"system_roles" json:"systemRoles"`
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
