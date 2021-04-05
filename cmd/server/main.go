package main

import (
	"fmt"

	"github.com/awnzl/echo-server/internal/config"
	"github.com/awnzl/echo-server/internal/logger"
	"github.com/awnzl/echo-server/internal/server"
)

func main() {
	conf, err := config.Get()
	if err != nil {
		conf = config.Config{
			Port:     "8080",
			LogLevel: "info",
		}
	}

	log := logger.NewZap(conf.LogLevel)
	defer log.Sync()

	s := server.New(log)
	if err := s.Run(conf.Port); err != nil {
		fmt.Println("Server run failure", err)
	}
}
