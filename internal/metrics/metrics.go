package metrics

import "github.com/prometheus/client_golang/prometheus"

const namespace = "myapp"

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_requests_total",
		Help:      "Number of requests.",
	},
	[]string{"path", "method", "code"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_response_status",
		Help:      "Status of HTTP responses.",
	},
	[]string{"path", "method", "code"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "http_response_time_seconds",
		Help:      "Duration of HTTP requests.",
	},
	[]string{"path", "method", "code"},
)
