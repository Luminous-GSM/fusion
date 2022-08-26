package node

import (
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/model/domain"
)

func (node NodeService) GetNodeDescription() domain.NodeDescriptionModel {
	config := config.Get()

	return domain.NodeDescriptionModel{
		Ip:              "0.0.0.0",
		NodeUniqueId:    config.NodeUniqueId,
		Name:            config.NodeName,
		Description:     config.NodeDescription,
		NodeStatus:      "running",
		Version:         config.Version,
		HostingPlatform: domain.HostingPlatformType(config.HostingPlatform),
		ActivePods:      0,
		Token:           config.ApiSecurityToken,
	}
}
