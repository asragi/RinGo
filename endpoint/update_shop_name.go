package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type UpdateShopNameEndpoint func(
	context.Context,
	*gateway.UpdateShopNameRequest,
) (*gateway.UpdateShopNameResponse, error)

func CreateUpdateShopNameEndpoint(
	updateShopName core.UpdateShopNameServiceFunc,
	validateToken auth.ValidateTokenFunc,
) UpdateShopNameEndpoint {
	return func(ctx context.Context, req *gateway.UpdateShopNameRequest) (*gateway.UpdateShopNameResponse, error) {
		handleError := func(err error) (*gateway.UpdateShopNameResponse, error) {
			return nil, fmt.Errorf("on update user name endpoint: %w", err)
		}
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		shopName, err := core.NewName(req.GetShopName())
		if err != nil {
			return handleError(err)
		}
		err = updateShopName(ctx, userId, shopName)
		if err != nil {
			return handleError(err)
		}
		return &gateway.UpdateShopNameResponse{
			Error:    nil,
			ShopName: shopName.String(),
		}, nil
	}
}
