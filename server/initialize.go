package server

import (
	"context"
	"github.com/asragi/RinGo/debug"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/grpc"
)

func parseArgs() *debug.RunMode {
	return debug.NewRunMode()
}

type gRPCServer struct {
	gateway.UnimplementedRingoServer
	endpoints *endpoint.Endpoints
}

func newGrpcServer(endpoints *endpoint.Endpoints) *gRPCServer {
	return &gRPCServer{endpoints: endpoints}
}

func (s *gRPCServer) SignUp(
	ctx context.Context,
	req *gateway.SignUpRequest,
) (*gateway.SignUpResponse, error) {
	return s.endpoints.SignUp(ctx, req)
}

func (s *gRPCServer) Login(
	ctx context.Context,
	req *gateway.LoginRequest,
) (*gateway.LoginResponse, error) {
	return s.endpoints.Login(ctx, req)
}

func (s *gRPCServer) GetResource(
	ctx context.Context,
	req *gateway.GetResourceRequest,
) (*gateway.GetResourceResponse, error) {
	return s.endpoints.GetResource(ctx, req)
}

func (s *gRPCServer) UpdateUserName(
	ctx context.Context,
	req *gateway.UpdateUserNameRequest,
) (*gateway.UpdateUserNameResponse, error) {
	return s.endpoints.UpdateUserName(ctx, req)
}

func (s *gRPCServer) UpdateShopName(
	ctx context.Context,
	req *gateway.UpdateShopNameRequest,
) (*gateway.UpdateShopNameResponse, error) {
	return s.endpoints.UpdateShopName(ctx, req)
}

func (s *gRPCServer) GetMyShelf(
	ctx context.Context,
	req *gateway.GetMyShelfRequest,
) (*gateway.GetMyShelfResponse, error) {
	return s.endpoints.GetMyShelves(ctx, req)
}

func (s *gRPCServer) GetStageList(
	ctx context.Context,
	req *gateway.GetStageListRequest,
) (*gateway.GetStageListResponse, error) {
	return s.endpoints.GetStageList(ctx, req)
}

func (s *gRPCServer) GetStageActionDetail(
	ctx context.Context,
	req *gateway.GetStageActionDetailRequest,
) (*gateway.GetStageActionDetailResponse, error) {
	return s.endpoints.GetStageActionDetail(ctx, req)
}

func (s *gRPCServer) PostAction(
	ctx context.Context,
	req *gateway.PostActionRequest,
) (*gateway.PostActionResponse, error) {
	return s.endpoints.PostAction(ctx, req)
}

func (s *gRPCServer) GetItemList(
	ctx context.Context,
	req *gateway.GetItemListRequest,
) (*gateway.GetItemListResponse, error) {
	return s.endpoints.GetItemList(ctx, req)
}

func (s *gRPCServer) GetItemDetail(
	ctx context.Context,
	req *gateway.GetItemDetailRequest,
) (*gateway.GetItemDetailResponse, error) {
	return s.endpoints.GetItemDetail(ctx, req)
}

func (s *gRPCServer) GetItemActionDetail(
	ctx context.Context,
	req *gateway.GetItemActionDetailRequest,
) (*gateway.GetItemActionDetailResponse, error) {
	return s.endpoints.GetItemActionDetail(ctx, req)
}

func (s *gRPCServer) UpdateShelfContent(
	ctx context.Context,
	req *gateway.UpdateShelfContentRequest,
) (*gateway.UpdateShelfContentResponse, error) {
	return s.endpoints.UpdateShelfContent(ctx, req)
}

func (s *gRPCServer) UpdateShelfSize(
	ctx context.Context,
	req *gateway.UpdateShelfSizeRequest,
) (*gateway.UpdateShelfSizeResponse, error) {
	return s.endpoints.UpdateShelfSize(ctx, req)
}

func (s *gRPCServer) GetDailyRanking(
	ctx context.Context,
	req *gateway.GetDailyRankingRequest,
) (*gateway.GetDailyRankingResponse, error) {
	return s.endpoints.GetRankingUserList(ctx, req)
}

func SetUpServer(port int, endpoints *endpoint.Endpoints) (Serve, StopDBFunc, error) {
	grpcServer := newGrpcServer(endpoints)
	registerServer := func(s grpc.ServiceRegistrar) {
		gateway.RegisterRingoServer(s, grpcServer)
	}
	return NewRPCServer(port, registerServer)
}
