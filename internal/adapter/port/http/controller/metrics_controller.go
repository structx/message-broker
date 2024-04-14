package controller

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InvokeMetricsController start prometheus handler with default metrics
func InvokeMetricsController() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			log.Fatalf("failed to start metrics server %v", err)
		}
	}()
}
