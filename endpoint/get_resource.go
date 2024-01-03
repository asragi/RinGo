package endpoint

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type GetResourceFunc func(request *gateway.GetResourceRequest) (*gateway.GetResourceResponse, error)

func CreateGetResourceEndpoint(
	serviceFunc stage.GetUserResourceServiceFunc,
) GetResourceFunc {
	get := func(req *gateway.GetResourceRequest) (*gateway.GetResourceResponse, error) {
		handleError := func(err error) (*gateway.GetResourceResponse, error) {
			return &gateway.GetResourceResponse{}, fmt.Errorf("error on get resource: %w", err)
		}
		userId := core.UserId(req.UserId)
		token := core.AccessToken(req.Token)
		err := userId.IsValid()
		if err != nil {
			return handleError(err)
		}
		err = token.IsValid()
		if err != nil {
			return handleError(err)
		}
		res, err := serviceFunc(userId, token)
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
