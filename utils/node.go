package utils

import (
	"encoding/base64"
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
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

func GenerateSecureRandomString(length int, allowSymbols bool) (string, error) {
	numberOfIntegers := int(math.Round(float64(length) / 3))
	numberOfSymbols := 0
	if allowSymbols {
		numberOfSymbols = int(math.Round(float64(length) / 4))
	}
	secureRandomString, err := password.Generate(length, numberOfIntegers, numberOfSymbols, false, true)
	if err != nil {
		zap.S().Errorw("utils: Could not generate secure random string")
		return "", err
	}
	return secureRandomString, nil
}

func HashPasswordBasedOnArgon2(password string, faultTolerant bool) (string, error) {

	var iterations uint32 = 3
	var memory uint32 = 32 * 1024
	var parallelism uint8 = 1
	var keyLength uint32 = 32

	salt, err := GenerateSecureRandomString(12, true)
	if err != nil {
		if !faultTolerant {
			return "", err
		}
		salt = "2*v0aFNv(0v&s1Me"
	}
	saltBytes := []byte(salt)
	passwordBytes := []byte(password)
	hashedPasswordBytes := argon2.Key(passwordBytes, saltBytes, iterations, memory, parallelism, keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(saltBytes)
	b64HashPassword := base64.RawStdEncoding.EncodeToString(hashedPasswordBytes)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, iterations, parallelism, b64Salt, b64HashPassword)

	return encodedHash, nil
}
