package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Query execution metrics - tracks end-to-end query execution time
	queryDuration *prometheus.HistogramVec

	// Query total counter - counts queries by type and status
	queryTotal *prometheus.CounterVec
)

// Init initializes and registers all Prometheus metrics with the default registry.
// This should be called once during application startup, before the metrics server starts.
func Init() {
	queryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "rsprandom",
			Name:      "query_duration_seconds",
			Help:      "Duration of request time for a random sign",
			Buckets:   prometheus.DefBuckets, // [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
		},
		[]string{"query_type"},
	)

	queryTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rsprandom",
			Name:      "query_total",
			Help:      "Total number of requests for a random sign by type and status",
		},
		[]string{"query_type", "status"}, // status: success, error
	)

}

// RecordQueryDuration records the duration of a query execution
func RecordQueryDuration(queryType string, duration float64) {
	queryDuration.WithLabelValues(queryType).Observe(duration)
}

// RecordQuerySuccess increments the success counter for a query type
func RecordQuerySuccess(queryType string) {
	queryTotal.WithLabelValues(queryType, "success").Inc()
}

// RecordQueryError increments the error counter for a query type
func RecordQueryError(queryType string) {
	queryTotal.WithLabelValues(queryType, "error").Inc()
}
