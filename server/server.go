package server

import (
	"fmt"

	"github.com/luminous-gsm/fusion/config"
)

func New() {
	config := config.Get()
	router := NewRouter()
	port := fmt.Sprintf("%v:%v", config.Api.Host, config.Api.Port)
	router.Run(port)
}
