package delivery

import (
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()
		next.ServeHTTP(w, r)
		// Log request details and duration
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log panic and return 500
			}
		}()
		next.ServeHTTP(w, r)
	})
}
