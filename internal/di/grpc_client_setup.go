package di

import (
	"log"

	"google.golang.org/grpc"
)

func GRPCClientSetup(port string) *grpc.ClientConn {
	// setup grpc client
	doconn, err := grpc.NewClient("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to doctor service: %v", err)
	}

	return doconn
}
