package explore

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type (
	CreateGetUserResourceServiceFunc func(
		resourceFunc game.GetResourceFunc,
	) GetUserResourceServiceFunc

	GetUserResourceServiceFunc func(
		context.Context,
		core.UserId,
	) (*game.GetResourceRes, error)
)

func CreateGetUserResourceService(
	getResource game.GetResourceFunc,
) GetUserResourceServiceFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
	) (*game.GetResourceRes, error) {
		handleError := func(err error) (*game.GetResourceRes, error) {
			return nil, fmt.Errorf("error on get user resource: %w", err)
		}
		err := userId.IsValid()
		if err != nil {
			return handleError(err)
		}
		res, err := getResource(ctx, userId)
		if err != nil {
			return handleError(err)
		}
		return res, nil
	}
}
