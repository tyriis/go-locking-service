package metrics

import (
	"fmt"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Add this function at the top level
func normalizePath(path string) string {
	// Replace lock key pattern with a placeholder
	re := regexp.MustCompile(`/api/v\d+/locks/[^/]+`)
	if re.MatchString(path) {
		return "/api/v1/locks/:key"
	}
	return path
}

type PrometheusMetricsService struct {
	httpRequestDuration *prometheus.HistogramVec
	httpRequestCounter  *prometheus.CounterVec
	userActionCounter   *prometheus.CounterVec
	errorCounter        *prometheus.CounterVec
	locksCounter        prometheus.Gauge
}

func NewPrometheusMetricsService() *PrometheusMetricsService {
	return &PrometheusMetricsService{
		httpRequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),

		httpRequestCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		}, []string{"method", "path", "status"}),

		locksCounter: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "locks_total",
			Help: "The total number of active locks",
		}),
	}
}

func (m *PrometheusMetricsService) ObserveHTTPRequest(method, path string, statusCode int, duration float64) {
	status := fmt.Sprintf("%d", statusCode)
	normalizedPath := normalizePath(path)
	m.httpRequestDuration.WithLabelValues(method, normalizedPath, status).Observe(duration)
	m.httpRequestCounter.WithLabelValues(method, normalizedPath, status).Inc()
}

func (m *PrometheusMetricsService) RecordUserAction(action string) {
	m.userActionCounter.WithLabelValues(action).Inc()
}

func (m *PrometheusMetricsService) IncrementErrorCount(errorType string) {
	m.errorCounter.WithLabelValues(errorType).Inc()
}

func (m *PrometheusMetricsService) SetLockCount(value float64) {
	m.locksCounter.Set(value)
}
