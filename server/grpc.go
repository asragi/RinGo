package server

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type StopDBFunc func()
type RegisterServer func(grpc.ServiceRegistrar)

func NewRPCServer(port int, registerOption RegisterServer) (Serve, StopDBFunc, error) {
	var stopDB StopDBFunc
	handleError := func(err error) (Serve, StopDBFunc, error) {
		return nil, stopDB, fmt.Errorf("new rpc server: %w", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return handleError(err)
	}
	s := grpc.NewServer()
	stopDB = s.GracefulStop
	registerOption(s)
	// TODO: for debug
	reflection.Register(s)

	serve := func() error {
		return s.Serve(listener)
	}
	return serve, stopDB, nil
}
