package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type UpdateShelfContentEndpointFunc func(
	ctx context.Context,
	req *gateway.UpdateShelfContentRequest,
) (*gateway.UpdateShelfContentResponse, error)

func CreateUpdateShelfSizeEndpointFunc(
	updateShelfContent shelf.UpdateShelfContentFunc,
	validateToken auth.ValidateTokenFunc,
) UpdateShelfContentEndpointFunc {
	return func(
		ctx context.Context,
		req *gateway.UpdateShelfContentRequest,
	) (*gateway.UpdateShelfContentResponse, error) {
		handleError := func(err error) (*gateway.UpdateShelfContentResponse, error) {
			return nil, fmt.Errorf("on update shelf content endpoint: %w", err)
		}
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		itemId := core.ItemId(req.ItemId)
		index := shelf.Index(req.Index)
		setPrice := shelf.SetPrice(req.SetPrice)
		err = updateShelfContent(ctx, userId, itemId, setPrice, index)
		if err != nil {
			return handleError(err)
		}
		return &gateway.UpdateShelfContentResponse{
			Error: nil,
		}, nil

	}
}
