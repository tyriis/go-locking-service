package infrastructure

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ActiveLocks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_locks_total",
		Help: "The total number of active locks",
	})
)

// Handler returns the metrics http handler for prometheus
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
