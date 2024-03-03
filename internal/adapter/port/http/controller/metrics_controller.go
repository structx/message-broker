package controller

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InvokeMetricsHandler start prometheus metrics handler on port 2112
func InvokeMetricsHandler() error {

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":2112", nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start metrics server %v", err)
		}
	}()

	return nil
}
