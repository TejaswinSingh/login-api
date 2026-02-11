package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TejaswinSingh/login-api/internal/metrics"
)

func Metrics(next http.Handler, stdMetrics *metrics.StdMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		next.ServeHTTP(w, r)

		end := time.Since(start)

		rw, ok := w.(*RequestLoggerResponseWriter)
		if ok {
			status := strconv.Itoa(rw.status)

			stdMetrics.RequestCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
			stdMetrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(end.Seconds())
		}

	})
}
