package server

import (
	"context"

	"go.uber.org/zap"
)

type ServerManager struct {
	ctx       context.Context
	ctxCancel *context.CancelFunc
}

func NewManager(ctx context.Context, cnl context.CancelFunc) (*ServerManager, error) {
	zap.S().Info("creating new server manager")
	return &ServerManager{
		ctx:       ctx,
		ctxCancel: &cnl,
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
