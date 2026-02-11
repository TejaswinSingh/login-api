package middleware

import (
	"log/slog"
	"net/http"
)

func SampleMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("sample middleware (before next)")
		next.ServeHTTP(w, r)
		logger.Info("sample middleware (after next)")
	})
}
