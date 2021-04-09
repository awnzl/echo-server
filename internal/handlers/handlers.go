package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handlers struct {
	logger *zap.Logger
}

func New(log *zap.Logger) *Handlers {
	return &Handlers{
		logger: log,
	}
}

type echoRequest struct {
	Word string `json:"word"`
}

type echoResponse struct {
	EchoWord string `json:"echo"`
}

type errorResponse struct {
	Level string `json:"level"`
	Error string `json:"error"`
}

type nameResponse struct {
	ServiceName string `json:"service_name"`
}

func (h *Handlers) RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/", h.nameHandler)
	router.HandleFunc("/echo", h.echoHandler)
}

func (h *Handlers) nameHandler(w http.ResponseWriter, r *http.Request) {
	resp := nameResponse{
		ServiceName: "echo-server",
	}

	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeError("system", "internal server error", http.StatusInternalServerError, w)
		return
	}

	if err := h.writeResponse(b, w); err != nil {
		h.logger.Error("response writing error", zap.Error(err))
	}
}

func (h *Handlers) echoHandler(w http.ResponseWriter, r *http.Request) {
	var requestData echoRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeError("user","failed to read request body", http.StatusBadRequest, w)
		return
	}

	resp := echoResponse{
		EchoWord: requestData.Word,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeError("system", "internal server error", http.StatusInternalServerError, w)
		return
	}

	if err := h.writeResponse(b, w); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *Handlers) setRespHeaderJSONMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("Request", zap.String("URI", r.RequestURI), zap.String("Addr", r.RemoteAddr))
		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) writeResponse(b []byte, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (h *Handlers) writeError(lvl, msg string, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	resp := errorResponse{
		Level: lvl,
		Error: msg,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("failed to marshal", zap.Error(err))
		return
	}

	if _, err := w.Write(b); err != nil {
		h.logger.Error("failed to write response", zap.Error(err))
	}
}
