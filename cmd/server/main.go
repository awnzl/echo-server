package main

import (
	"fmt"
	"github.com/awnzl/echo-server/internal/config"
	"github.com/awnzl/echo-server/internal/logger"
	"github.com/awnzl/echo-server/internal/server"
	"log"
)

const envFilepath = "../.env"

func main() {
	conf, err := config.Get(envFilepath)
	if err != nil {
		log.Fatal(err)
	}

	log := logger.NewZap(conf.LogLevel)
	defer log.Sync()

	s := server.New(log)
	if err := s.Run(conf.Port); err != nil {
		log.Fatal(fmt.Sprintf("Server run failure: %v", err))
	}
}
