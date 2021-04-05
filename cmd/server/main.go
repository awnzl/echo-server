package main

import (
	"fmt"

	"github.com/awnzl/echo-server/internal/logger"
	"github.com/awnzl/echo-server/internal/server"
)

func main() {
	log := logger.NewZap("info")
	defer log.Sync()

	s := server.New(log)
	if err := s.Run("8080"); err != nil {
		fmt.Println("Server run failure", err)
	}
}
