package main

import (
	"context"
	"github.com/asragi/RinGo/debug"
	"github.com/asragi/RinGo/initialize"
	"github.com/asragi/RinGo/server"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/grpc"
)

func parseArgs() *debug.RunMode {
	return debug.NewRunMode()
}

type gRPCServer struct {
	gateway.UnimplementedRegisterServer
	endpoints *initialize.Endpoints
}

func newGrpcServer(endpoints *initialize.Endpoints) *gRPCServer {
	return &gRPCServer{endpoints: endpoints}
}

func (s *gRPCServer) RegisterUser(
	ctx context.Context,
	req *gateway.RegisterUserRequest,
) (*gateway.RegisterUserResponse, error) {
	return s.endpoints.SignUp(ctx, req)
}

func setUpServer(port int, endpoints *initialize.Endpoints) (server.Serve, server.StopDBFunc, error) {
	grpcServer := newGrpcServer(endpoints)
	registerServer := func(s grpc.ServiceRegistrar) {
		gateway.RegisterRegisterServer(s, grpcServer)
	}
	return server.NewRPCServer(port, registerServer)
}
