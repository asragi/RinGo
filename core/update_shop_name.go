package core

import "context"

type UpdateShopNameServiceFunc func(context.Context, UserId, Name) error

func CreateUpdateShopNameServiceFunc(updateShopName UpdateShopNameFunc) UpdateShopNameServiceFunc {
	return func(ctx context.Context, userId UserId, userName Name) error {
		return updateShopName(ctx, userId, userName)
	}
}
