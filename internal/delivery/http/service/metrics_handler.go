package service

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns the metrics http handler for prometheus
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
