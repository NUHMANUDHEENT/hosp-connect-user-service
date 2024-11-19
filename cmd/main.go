package main

import (
	"log"
	"os"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/di"
)

func main() {
	di.LoadEnv()
	port := os.Getenv("USER_PORT")
	if port == "" {
		log.Fatal("USER_PORT environment variable not set")
	}

	listener, server := di.GRPCSetup(port)

	log.Printf("Starting gRPC server on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
