package scenario

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectFunc func() (*grpc.ClientConn, error)

func Connect(address string) ConnectFunc {
	return func() (*grpc.ClientConn, error) {
		return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}

func closeConnection(conn *grpc.ClientConn) {
	err := conn.Close()
	if err == nil {
		return
	}
	fmt.Printf("error on close connection: %v\n", err)
}
