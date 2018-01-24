package api

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Number of request counts.",
		},
		[]string{"appid"},
	)

	ErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_error_total",
			Help: "Number of request error counts.",
		},
		[]string{"error_code"},
	)
)

func RegistForwardMetrics() {
	prometheus.MustRegister(RequestTotal)
	prometheus.MustRegister(ErrorTotal)
}
