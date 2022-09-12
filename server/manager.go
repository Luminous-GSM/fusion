package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type ServerManager struct {
	ctx       context.Context
	ctxCancel *context.CancelFunc
	validator *validator.Validate
}

func NewManager(ctx context.Context, cnl context.CancelFunc) (*ServerManager, error) {
	zap.S().Info("creating new server manager")
	return &ServerManager{
		ctx:       ctx,
		ctxCancel: &cnl,
		validator: validator.New(),
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

func (s *ServerManager) BindAndValidate(c *gin.Context, obj any) error {
	if err := c.BindJSON(&obj); err != nil {
		return err
	}

	if err := s.ValidateData(obj); err != nil {
		return err
	}
	return nil
}

func (s *ServerManager) ValidateData(data interface{}) error {
	// Validate the configuration according to validation tags in the structs.
	if err := s.validator.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			zap.S().Errorw("request field error",
				"field", err.Field(),
				"value", err.Value(),
				"validation_type", err.Tag(),
				"field_type", err.Type(),
			)
		}
		zap.S().Debugw("request error", "error", err)
		return err
	}
	return nil
}
