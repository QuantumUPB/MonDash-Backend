package middlewares

import (
	"bytes"
	"net/http"
	"strings"

	"mondash-backend/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

// LoggingMiddleware logs incoming requests and outgoing responses at debug level.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debugf("Request: %s %s", r.Method, r.URL.Path)
		lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lrw, r)
		logger.Log.Debugf("Response: %d %s", lrw.status, strings.TrimSpace(lrw.body.String()))
	})
}
