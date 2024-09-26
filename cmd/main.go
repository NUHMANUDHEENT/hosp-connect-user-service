package main

import (
	"log"
	"os"

	"github.com/nuhmanudheent/hosp-connect-user-service/internal/di"
)

func main() {
	di.LoadEnv()
	port := os.Getenv("USERPORT")
	listener, server := di.GRPCSetup(port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
																																																																																																																																																																																																																																																																																																																																																															