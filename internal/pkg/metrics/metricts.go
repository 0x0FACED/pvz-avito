package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ------tech------
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	HTTPResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.1, 0.5, 1, 2, 5},
	}, []string{"method", "path"})

	// ------business------
	PvzCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pvz_created_total",
		Help: "Total number of created PVZ",
	})

	ReceptionCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "reception_created_total",
		Help: "Total number of created receptions",
	})

	ProductsAddedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_added_total",
		Help: "Total number of addes products",
	})
)
