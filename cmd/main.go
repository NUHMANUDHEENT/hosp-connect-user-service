package main

import (
	"log"
	"os"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/di"
)

func main() {
	// Load environment variables
	di.LoadEnv()
	port := os.Getenv("USER_PORT")
	if port == "" {
		log.Fatal("USER_PORT environment variable not set")
	}

	// Start the gRPC service
	listener, server := di.GRPCSetup(port)

	// Start Prometheus metrics server in a separate goroutine
	go di.StartMetricsServer()

	// Serve the gRPC server
	log.Printf("Starting gRPC server on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
