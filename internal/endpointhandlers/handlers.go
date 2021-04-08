package endpointhandlers

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
	router.Use(h.loggingMiddleware, h.setHeaderMiddleware)
}

func (h *Handler) nameHandler(w http.ResponseWriter, r *http.Request) {
	resp := nameResponse{
		ServiceName: "echo-server",
	}

	b, err := h.marshal(resp)
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
	b, err := h.marshal(resp)
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

func (h *Handler) setHeaderMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

func (h *Handler) loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info(fmt.Sprintf("Requested URI: %v, from Remote address: %v", r.RequestURI, r.RemoteAddr))
		handler.ServeHTTP(w, r)
	})
}

func (h *Handler) marshal(resp interface{}) ([]byte, error) {
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (h *Handler) writeResponse(b []byte, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (h *Handler) writeError(b []byte, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		h.logger.Error(err.Error())
	}
}
