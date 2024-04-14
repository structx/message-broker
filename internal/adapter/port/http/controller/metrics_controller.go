package controller

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InvokeMetricsController start prometheus handler with default metrics
func InvokeMetricsController() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":2112", nil)
	}()
}
