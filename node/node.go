package node

import (
	"context"
	"sync"

	"github.com/luminous-gsm/fusion/event"
	"github.com/luminous-gsm/fusion/model/domain"
	"github.com/vmware/transport-go/model"
	"go.uber.org/zap"
)

var (
	_conce      sync.Once
	nodeService *NodeService
)

type NodeService struct {
	ctx      context.Context
	warnings []domain.FusionWarning
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
	event.Instance().ListenWithPanic(event.EVENT_NODE_WARNING).Handle(
		func(m *model.Message) {
			event := m.Payload.(event.FusionEvent[domain.FusionWarning])
			ns.warnings = append(ns.warnings, event.Data)
		},
		event.Instance().DefaultErrorHandler,
	)
}

func Instance() *NodeService {
	return nodeService
}
