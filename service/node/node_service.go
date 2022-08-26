package node

import (
	"context"

	"go.uber.org/zap"
)

type NodeService struct {
	ctx context.Context
}

// Returns a docker client.
func NewNodeService(context context.Context) (*NodeService, error) {
	zap.S().Info("node: configured node client")
	return &NodeService{
		ctx: context,
	}, nil
}
