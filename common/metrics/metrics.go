package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	ErrorsTotal     *prometheus.CounterVec
	CircuitState    *prometheus.GaugeVec
}

// NewMetrics creates and registers Prometheus metrics
func NewMetrics(serviceName string) *Metrics {
	metrics := &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: serviceName + "_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    serviceName + "_request_duration_seconds",
				Help:    "HTTP request latencies in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: serviceName + "_errors_total",
				Help: "Total number of errors",
			},
			[]string{"type"},
		),
		CircuitState: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: serviceName + "_circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
			},
			[]string{"service"},
		),
	}

	return metrics
}

// Middleware returns a Gin middleware for collecting metrics
func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		m.RequestsTotal.WithLabelValues(method, path, status).Inc()
		m.RequestDuration.WithLabelValues(method, path).Observe(duration)

		// Record errors for 5xx status codes
		if c.Writer.Status() >= 500 {
			m.ErrorsTotal.WithLabelValues("server_error").Inc()
		} else if c.Writer.Status() >= 400 {
			m.ErrorsTotal.WithLabelValues("client_error").Inc()
		}
	}
}

// RecordError records an error metric
func (m *Metrics) RecordError(errorType string) {
	m.ErrorsTotal.WithLabelValues(errorType).Inc()
}

// SetCircuitState sets the circuit breaker state
// state: 0=closed, 1=half-open, 2=open
func (m *Metrics) SetCircuitState(service string, state float64) {
	m.CircuitState.WithLabelValues(service).Set(state)
}
