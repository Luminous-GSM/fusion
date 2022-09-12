package node

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

var (
	_conce      sync.Once
	nodeService *NodeService
)

type NodeService struct {
	ctx context.Context
}

// Returns a docker client.
func InitNodeService(context context.Context) *NodeService {
	_conce.Do(func() {
		nodeService = &NodeService{
			ctx: context,
		}
		zap.S().Info("node: configured node client")
	})
	return nodeService
}

func (ns NodeService) InitEventListeners() {

}

func Instance() *NodeService {
	return nodeService
}
