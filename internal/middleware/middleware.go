package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type Logger struct {
	handler http.Handler
	logger  *zap.Logger
}

func NewLogging(log *zap.Logger, handler http.Handler) *Logger {
	return &Logger{
		logger:  log,
		handler: handler,
	}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.logger.Info("Request", zap.String("URI", r.RequestURI), zap.String("Addr", r.RemoteAddr))
	l.handler.ServeHTTP(w, r)
}

type ContentTypeJSON struct {
	handler http.Handler
}

func NewContentTypeJSON(handler http.Handler) *ContentTypeJSON {
	return &ContentTypeJSON{
		handler: handler,
	}
}

func (h *ContentTypeJSON) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.handler.ServeHTTP(w, r)
}
