package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/awnzl/echo-server/internal/config"
	"github.com/awnzl/echo-server/internal/handlers"
	"github.com/awnzl/echo-server/internal/logger"
)

func main() {
	conf := config.Get()
	log  := logger.NewZap(conf.LogLevel)
	defer log.Sync()

	router       := mux.NewRouter()
	httpHandlers := handlers.New(log)
	httpHandlers.RegisterHandlers(router)

	s := &http.Server{
		Addr: fmt.Sprintf(":%v", conf.Port),
		Handler: router,
	}

	log.Info("start listening", zap.String("port", conf.Port))
	if err := s.ListenAndServe(); err != nil {
		log.Panic("server error", zap.Error(err))
	}
}
