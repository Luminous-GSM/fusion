package server

import (
	"context"

	"github.com/luminous-gsm/fusion/environment"
)

type ServerManager struct {
	ctx         context.Context
	ctxCancel   *context.CancelFunc
	environment *environment.Environment
}

func NewManager(env *environment.Environment) (*ServerManager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &ServerManager{
		ctx:         ctx,
		ctxCancel:   &cancel,
		environment: env,
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

func (s *ServerManager) Environment() *environment.Environment {
	return s.environment
}
