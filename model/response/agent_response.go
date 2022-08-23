package response

type NodeDescriptionResponse struct {
	Ip              string              `json:"ip"`
	NodeUniqueId    string              `json:"nodeUniqueId"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	NodeStatus      NodeStatusType      `json:"nodeStatus"`
	Version         string              `json:"version"`
	HostingPlatform HostingPlatformType `json:"hostingPlatform"`
	ActivePods      int                 `json:"activePods"`
	Token           string              `json:"token"`
}

type NodeStatusType string

const (
	RUNNING    NodeStatusType = "running"
	PENDING    NodeStatusType = "pending"
	TERMINATED NodeStatusType = "terminated"
	INACTIVE   NodeStatusType = "inactive"
)

type HostingPlatformType string

const (
	DAEMON HostingPlatformType = "daemon"
	LOCAL  HostingPlatformType = "local"
	AWS    HostingPlatformType = "aws"
)
