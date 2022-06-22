package server

import (
	"fmt"

	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/router"
)

func New() {
	config := config.Get()
	router := router.New()
	port := fmt.Sprintf("%v:%v", config.Api.Host, config.Api.Port)
	router.Run(port)
}
