package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func New(log *zap.Logger) *Handler {
	return &Handler{
		logger: log,
	}
}

type echoRequest struct {
	Word string `json:"word"`
}

type echoResponse struct{
	EchoWord string `json:"echo"`
}

type nameResponse struct{
	ServiceName string `json:"service_name"`
}

func (h *Handler) RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/", h.nameHandler)
	router.HandleFunc("/echo", h.echoHandler)
	router.Use(h.loggingMiddleware)
}

func (h *Handler) nameHandler(w http.ResponseWriter, r *http.Request) {
	resp := nameResponse{
		ServiceName: "echo-server",
	}
	b := h.marshal(resp)

	if err := h.writeResponse(b, w); err != nil {
		h.logger.Error("response writing error", zap.Error(err))
	}
}

func (h *Handler) echoHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}

	var requestData echoRequest
	if err := json.Unmarshal(data, &requestData); err != nil {
		h.logger.Error(err.Error())
		return
	}

	resp := echoResponse{
		EchoWord: requestData.Word,
	}
	b := h.marshal(resp)

	if err := h.writeResponse(b, w); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *Handler) loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info(fmt.Sprintf("Requested URI: %v, from Remote address: %v", r.RequestURI, r.RemoteAddr))
		handler.ServeHTTP(w, r)
	})
}

func (h *Handler) marshal(resp interface{}) []byte {
	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("marshalling error", zap.Error(err))
	}

	return b
}

func (h *Handler) writeResponse(b []byte, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}
