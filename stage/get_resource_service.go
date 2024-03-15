package stage

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type (
	CreateGetUserResourceServiceFunc func(
		resourceFunc GetResourceFunc,
	) GetUserResourceServiceFunc

	GetUserResourceServiceFunc func(
		context.Context,
		core.UserId,
	) (*GetResourceRes, error)
)

func CreateGetUserResourceService(
	getResource GetResourceFunc,
) GetUserResourceServiceFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
	) (*GetResourceRes, error) {
		handleError := func(err error) (*GetResourceRes, error) {
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
