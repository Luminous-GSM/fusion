package server

import (
	"fmt"
	"fusion/config"
)

func New() {
	config := config.GetConfig()
	r := NewRouter()
	port := fmt.Sprintf(":%v", config.GetString("server.port"))
	r.Run(port)
}
