package metrics

import (
	"net/http"
	"time"

	"github.com/tyriis/go-locking-service/pkg/metrics"
)

type MetricsMiddleware struct {
	recorder metrics.MetricsRecorder
}

func NewMetricsMiddleware(recorder metrics.MetricsRecorder) *MetricsMiddleware {
	return &MetricsMiddleware{recorder: recorder}
}

func (m *MetricsMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := &customResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(crw, r)

		duration := time.Since(start).Seconds()
		m.recorder.ObserveHTTPRequest(r.Method, r.URL.Path, crw.statusCode, duration)
		if crw.statusCode >= 500 {
			m.recorder.IncrementErrorCount("server_error")
		}
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}
