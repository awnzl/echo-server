package handlers

import (
	"encoding/json"
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
	router.Use(h.loggingMiddleware, h.JSONResponseMiddleware)
}

func (h *Handler) nameHandler(w http.ResponseWriter, r *http.Request) {
	resp := nameResponse{
		ServiceName: "echo-server",
	}

	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeError(
			[]byte(`{"level": "system", "error": "internal server error"}`),
			http.StatusInternalServerError,
			w,
		)
		return
	}

	if err := h.writeResponse(b, w); err != nil {
		h.logger.Error("response writing error", zap.Error(err))
	}
}

func (h *Handler) echoHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeError(
			[]byte(`{"level": "user", "error": "bad request. failed to read request body"}`),
			http.StatusBadRequest,
			w,
		)
		return
	}

	var requestData echoRequest
	if err := json.Unmarshal(data, &requestData); err != nil {
		h.logger.Error(err.Error())
		h.writeError(
			[]byte(`{"level": "user", "error": "bad request. failed to unmarshal request json data"}`),
			http.StatusBadRequest,
			w,
		)
		return
	}

	resp := echoResponse{
		EchoWord: requestData.Word,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeError(
			[]byte(`{"level": "system", "error": "internal server error"}`),
			http.StatusInternalServerError,
			w,
		)
		return
	}

	if err := h.writeResponse(b, w); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *Handler) JSONResponseMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

func (h *Handler) loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("Request", zap.String("URI", r.RequestURI), zap.String("Addr", r.RemoteAddr))
		handler.ServeHTTP(w, r)
	})
}

func (h *Handler) writeResponse(b []byte, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (h *Handler) writeError(b []byte, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		h.logger.Error(err.Error())
	}
}
