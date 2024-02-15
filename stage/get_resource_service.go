package stage

import (
	"fmt"
	"github.com/asragi/RinGo/core"
)

type CreateGetUserResourceServiceFunc func(
	resourceFunc GetResourceFunc,
) GetUserResourceServiceFunc
type GetUserResourceServiceFunc func(core.UserId) (GetResourceRes, error)

func CreateGetUserResourceService(
	getResource GetResourceFunc,
) GetUserResourceServiceFunc {
	get := func(
		userId core.UserId,
	) (GetResourceRes, error) {
		handleError := func(err error) (GetResourceRes, error) {
			return GetResourceRes{}, fmt.Errorf("error on get user resource: %w", err)
		}
		err := userId.IsValid()
		if err != nil {
			return handleError(err)
		}
		res, err := getResource(userId)
		if err != nil {
			return handleError(err)
		}
		return res, nil
	}

	return get
}
