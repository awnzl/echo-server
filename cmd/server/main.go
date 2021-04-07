package main

import (
	"fmt"
	"net/http"

	"github.com/awnzl/echo-server/internal/config"
	"github.com/awnzl/echo-server/internal/handlers"
	"github.com/awnzl/echo-server/internal/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	conf := config.Get()
	log  := logger.NewZap(conf.LogLevel)
	defer log.Sync()

	router   := mux.NewRouter()
	handlers := handlers.New(log)
	handlers.RegisterHandlers(router)

	s := &http.Server{
		Addr: fmt.Sprintf(":%v", conf.Port),
		Handler: router,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Panic("server error", zap.Error(err))
	}
}
