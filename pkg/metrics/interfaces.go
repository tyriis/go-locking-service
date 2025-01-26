package metrics

// MetricsRecorder defines the interface for recording metrics
type MetricsRecorder interface {
	ObserveHTTPRequest(method, path string, statusCode int, duration float64)
	RecordUserAction(action string)
	IncrementErrorCount(errorType string)
}
