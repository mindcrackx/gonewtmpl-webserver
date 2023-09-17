package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w)
		path := r.URL.Path
		method := r.Method
		start := time.Now()

		next.ServeHTTP(rw, r)

		statusCode := strconv.Itoa(rw.statusCode)
		totalRequests.WithLabelValues(path, method, statusCode).Inc()
		responseStatus.WithLabelValues(path, method, statusCode).Inc()
		httpDuration.WithLabelValues(path, method, statusCode).Observe(time.Since(start).Seconds())
	})
}
