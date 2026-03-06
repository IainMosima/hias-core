package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hias_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hias_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "hias_http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		httpActiveConnections.Inc()
		start := time.Now()

		ctx.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ctx.Writer.Status())
		path := ctx.FullPath()
		if path == "" {
			path = "unknown"
		}

		httpRequestsTotal.WithLabelValues(ctx.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(ctx.Request.Method, path).Observe(duration)
		httpActiveConnections.Dec()
	}
}
