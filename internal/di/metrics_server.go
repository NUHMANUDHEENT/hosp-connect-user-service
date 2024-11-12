package di

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Starts an HTTP server to serve metrics at /metrics
func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Prometheus metrics server starting on :2112")
	if err := http.ListenAndServe(":2112", nil); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}
