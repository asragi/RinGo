package endpoint

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type GetResourceFunc func(request *gateway.GetResourceRequest) (*gateway.GetResourceResponse, error)

func CreateGetResourceEndpoint(
	serviceFunc stage.GetUserResourceServiceFunc,
	validateToken auth.ValidateTokenFunc,
) GetResourceFunc {
	get := func(req *gateway.GetResourceRequest) (*gateway.GetResourceResponse, error) {
		handleError := func(err error) (*gateway.GetResourceResponse, error) {
			return &gateway.GetResourceResponse{}, fmt.Errorf("error on get resource: %w", err)
		}
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		res, err := serviceFunc(userId)
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetResourceResponse{
			UserId:      string(res.UserId),
			MaxStamina:  int32(res.MaxStamina),
			RecoverTime: timestamppb.New(time.Time(res.StaminaRecoverTime)),
			Fund:        int32(res.Fund),
		}, nil
	}
	return get
}
