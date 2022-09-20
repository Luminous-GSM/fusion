package domain

type NodeDescriptionModel struct {
	Ip                 string              `json:"ip"`
	NodeUniqueId       string              `json:"nodeUniqueId"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	NodeStatus         NodeStatusType      `json:"nodeStatus"`
	NodeStatusExpected NodeStatusType      `json:"nodeStatusExpected"`
	Version            string              `json:"version"`
	HostingPlatform    HostingPlatformType `json:"hostingPlatform"`
	ActivePods         int                 `json:"activePods"`
	Token              string              `json:"token"`
	Warnings           []FusionWarning     `json:"warnings"`
}

type FusionContainerModel struct {
	Command string               `json:"command"`
	Created int                  `json:"created"`
	Id      string               `json:"id"`
	Image   string               `json:"image"`
	ImageId string               `json:"imageId"`
	Names   []string             `json:"names"`
	Status  string               `json:"status"`
	State   FusionContainerState `json:"state"`
	Ports   []ContainerPort      `json:"ports"`
}

type ContainerPort struct {
	Ip          string `json:"ip"`
	PrivatePort string `json:"privatePort"`
	PublicPort  string `json:"publicPort"`
	Type        string `json:"type"`
}

type FusionImageModel struct {
	Created    int    `json:"created"`
	Id         string `json:"id"`
	Size       int    `json:"size"`
	Containers int    `json:"containers"`
}

type HostingPlatformType string

const (
	SELF HostingPlatformType = "self"
	AWS  HostingPlatformType = "aws"
)

type FusionContainerState string

const (
	CONTAINER_CREATED    FusionContainerState = "created"
	CONTAINER_RESTARTING FusionContainerState = "restarting"
	CONTAINER_RUNNING    FusionContainerState = "running"
	CONTAINER_REMOVING   FusionContainerState = "removing"
	CONTAINER_PAUSED     FusionContainerState = "paused"
	CONTAINER_EXITED     FusionContainerState = "exited"
	CONTAINER_DEAD       FusionContainerState = "dead"
)

type NodeStatusType string

const (
	NODE_RUNNING    NodeStatusType = "running"
	NODE_PENDING    NodeStatusType = "pending"
	NODE_TERMINATED NodeStatusType = "terminated"
	NODE_INACTIVE   NodeStatusType = "inactive"
)

type SystemLoadModel struct {
	CpuLoad  string `json:"cpuLoad"`
	RamLoad  string `json:"ramLoad"`
	HddUsage string `json:"hddUsage"`
}

type FusionWarning struct {
	Severity FusionWarningSeverity `json:"severity"`
	Service  FusionWarningService  `json:"service"`
	Message  string                `json:"message"`
}

type FusionWarningSeverity string

const (
	HIGH   FusionWarningSeverity = "high"
	MEDIUM FusionWarningSeverity = "medium"
	LOW    FusionWarningSeverity = "low"
)

type FusionWarningService string

const (
	DOCKER FusionWarningService = "docker"
	EVENT  FusionWarningService = "event"
)
