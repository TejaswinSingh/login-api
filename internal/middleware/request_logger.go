package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/TejaswinSingh/login-api/internal/constants"
	"github.com/TejaswinSingh/login-api/internal/logging"
	"github.com/google/uuid"
)

func RequestLogger(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		r = addRequestIdToContext(r)
		rw := &RequestLoggerResponseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		end := time.Since(start)

		logger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.status),
			slog.Float64("duration", end.Seconds()),
		)

	})
}

type RequestLoggerResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *RequestLoggerResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Unwrap supports http.ResponseController.
func (rw *RequestLoggerResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func addRequestIdToContext(r *http.Request) *http.Request {
	var id uuid.UUID

	id, err := uuid.Parse(r.Header.Get("Request-Id"))
	if err != nil {
		id = uuid.New()
		r.Header.Set("Request-Id", id.String())
	}

	// ctx := context.WithValue(r.Context(), constants.RequestId_CtxKey, id.String())
	ctx := logging.AppendAttrToCtx(r.Context(), slog.String(string(constants.RequestIdCtxKey), id.String()))
	r = r.WithContext(ctx)
	return r
}
