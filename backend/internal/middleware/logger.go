package middleware
package middleware

import (
	"net/http"
	"time"

	"github.com/apariciocch/psicosstcloud/internal/pkg/logger"
	"go.uber.org/zap"
)

// RequestLogger middleware para loguear requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Envolver response writer para capturar status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Ejecutar handler
		next.ServeHTTP(wrapped, r)

		// Log del request
		duration := time.Since(start).Milliseconds()
		userID := ""
		if uid, ok := r.Context().Value(UserIDKey).(string); ok {
			userID = uid
		}

		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.Int("status", wrapped.statusCode),
			zap.Int64("duration_ms", duration),
			zap.String("user_id", userID),
			zap.String("ip", getClientIP(r)),
		}

		logger.Info("HTTP Request", fields...)
	})
}

// responseWriter wrapper para capturar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}
