package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

type requestData struct {
	Word string `json:"word"`
}

type Server struct {
	router *mux.Router
	logger *zap.Logger
}

func New(logger *zap.Logger) *Server {
	s := &Server{
		router: mux.NewRouter(),
		logger: logger,
	}
	s.init()

	return s
}

func (s *Server) init() {
	s.router.HandleFunc("/", s.name)
	s.router.HandleFunc("/echo", s.echo)
	s.router.Use(s.loggingMiddleware)
	http.Handle("/", s.router)
	http.Handle("/echo", s.router)
}

func (s *Server) Run(port string) error {
	if s.router == nil || s.logger == nil {
		return fmt.Errorf("server isn't initialized")
	}

	return http.ListenAndServe(fmt.Sprintf(":%v", port), s.router)
}

func (s *Server) name(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if _, err := w.Write([]byte("{\n\t\"service_name\": \"echo-server\"\n}")); err != nil {
		s.logger.Error(fmt.Sprintf("error: %v", err))
	}
}

func (s *Server) echo(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	var requestData requestData
	if err := json.Unmarshal(data, &requestData); err != nil {
		s.logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	msg := fmt.Sprintf("{\n\t\"echo\": \"%v\"\n}", requestData.Word)
	if _, err := w.Write([]byte(msg)); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Server) loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(fmt.Sprintf("Requested URI: %v, from Remote address: %v", r.RequestURI, r.RemoteAddr))
	})
}
