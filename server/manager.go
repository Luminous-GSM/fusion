package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/luminous-gsm/fusion/service"
	"go.uber.org/zap"
)

type ServerManager struct {
	ctx            context.Context
	ctxCancel      *context.CancelFunc
	serviceManager *service.ServiceManager
}

func NewManager(ctx context.Context, cnl context.CancelFunc, srvMgr *service.ServiceManager) (*ServerManager, error) {
	zap.S().Info("creating new server manager")
	return &ServerManager{
		ctx:            ctx,
		ctxCancel:      &cnl,
		serviceManager: srvMgr,
	}, nil
}

// Cancels the context
func (s *ServerManager) CtxCancel() {
	if s.ctxCancel != nil {
		(*s.ctxCancel)()
	}
}

// Returns a context instance for the server.
func (s *ServerManager) Context() context.Context {
	return s.ctx
}

func (s *ServerManager) ServiceManager() *service.ServiceManager {
	return s.serviceManager
}

func (s *ServerManager) BindAndValidate(c *gin.Context, obj any) error {
	if err := c.BindJSON(&obj); err != nil {
		return err
	}

	if err := s.ServiceManager().ValidateData(obj); err != nil {
		return err
	}
	return nil
}
