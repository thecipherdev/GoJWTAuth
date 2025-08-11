package main

import (
	"fmt"
	"github.com/thecipherdev/goauth/cmd/api"
	"github.com/thecipherdev/goauth/config"
	"log"
)

func main() {
	cfg := config.Get()
	PORT := fmt.Sprintf(":%v", cfg.Port)
	server := api.NewAPIServer(PORT)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
