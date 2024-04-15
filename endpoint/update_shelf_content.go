package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type UpdateShelfSizeEndpointFunc func(
	ctx context.Context,
	req *gateway.UpdateShelfSizeRequest,
) (*gateway.UpdateShelfSizeResponse, error)

func CreateUpdateShelfSizeEndpointFunc(
	updateShelfSize shelf.UpdateShelfSizeFunc,
	validateToken auth.ValidateTokenFunc,
) UpdateShelfSizeEndpointFunc {
	return func(
		ctx context.Context,
		req *gateway.UpdateShelfSizeRequest,
	) (*gateway.UpdateShelfSizeResponse, error) {
		handleError := func(err error) (*gateway.UpdateShelfSizeResponse, error) {
			return nil, fmt.Errorf("on update shelf size endpoint: %w", err)
		}
		sizeInt := req.Size
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		size := shelf.Size(sizeInt)
		err = updateShelfSize(ctx, userId, size)
		if err != nil {
			return handleError(err)
		}
		return &gateway.UpdateShelfSizeResponse{
			Error: nil,
			Size:  sizeInt,
		}, nil
	}
}
