package main

import (
	"flag"
	"fmt"
	"os"

	"fusion/config"
	"fusion/server"
)

func main() {
	environment := flag.String("e", "fusion", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.New(*environment)
	server.New()
}
