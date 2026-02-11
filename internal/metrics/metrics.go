package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type StdMetrics struct {
	RequestCount    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func NewStdMetrics(reg prometheus.Registerer) *StdMetrics {
	m := &StdMetrics{
		RequestCount: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of requests received",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: []float64{0.02, 0.05, 0.1, 0.2, 0.5, 1, 2},
			},
			[]string{"method", "path", "status"},
		),
	}
	return m
}
