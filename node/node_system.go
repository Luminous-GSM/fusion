package node

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/model/domain"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

type MyCustomClaims struct {
	Roles []string `json:"fusionRoles"`
	jwt.RegisteredClaims
}

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

func (node NodeService) GetSystemLoad() (*domain.SystemLoadModel, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		zap.S().Errorw("node: could not read virtual memory", "error", err)
		return nil, err
	}
	cpu, err := cpu.Percent(time.Millisecond, false)
	if err != nil {
		zap.S().Errorw("node: could not read cpu percentage", "error", err)
		return nil, err
	}
	usage, err := disk.Usage(config.Get().DataDirectory)
	if err != nil {
		zap.S().Errorw("node: could not read data directory size", "error", err)
		return nil, err
	}

	return &domain.SystemLoadModel{
		RamLoad:  fmt.Sprintf("%.2f", mem.UsedPercent),
		CpuLoad:  fmt.Sprintf("%.2f", cpu[0]),
		HddUsage: fmt.Sprintf("%.2f", usage.UsedPercent),
	}, nil
}

func (node NodeService) TemporaryAuthentication() (string, error) {
	// Create the claims
	claims := MyCustomClaims{
		[]string{"websocket"},
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fusion",
			Subject:   "websocket",
		},
	}

	securedBytes := []byte(config.Get().ApiSecurityToken)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(securedBytes)
	if err != nil {
		zap.S().Errorw("jwt-token: failed jwt signing", "error", err, "token", token)
		return "", err
	}

	// Create JWT
	return signedToken, nil
}
