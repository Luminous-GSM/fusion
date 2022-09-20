package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

func BindAndValidate(c *gin.Context, obj any) error {
	if err := c.BindJSON(&obj); err != nil {
		return err
	}

	if err := ValidateData(obj); err != nil {
		return err
	}
	return nil
}

func ValidateData(data interface{}) error {
	// Validate the configuration according to validation tags in the structs.
	if err := validator.New().Struct(data); err != nil {
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
