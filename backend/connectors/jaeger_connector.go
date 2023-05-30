package connectors

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func JaegerConnect(addr string) grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("jaeger connection failed: %v", err)
	}

	return *conn
}
