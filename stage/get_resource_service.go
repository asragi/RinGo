package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type GetUserResourceServiceFunc func(core.UserId, core.AccessToken) (GetResourceRes, error)

func CreateGetUserResourceService(
	validateToken core.ValidateTokenFunc,
	getResource GetResourceFunc,
) GetUserResourceServiceFunc {
	get := func(
		userId core.UserId,
		token core.AccessToken,
	) (GetResourceRes, error) {
		handleError := func(err error) (GetResourceRes, error) {
			return GetResourceRes{}, fmt.Errorf("error on get user resource: %w", err)
		}
		err := userId.IsValid()
		if err != nil {
			return handleError(err)
		}
		err = validateToken(userId, token)
		if err != nil {
			return handleError(err)
		}
		res, err := getResource(userId, token)
		if err != nil {
			return handleError(err)
		}
		return res, nil
	}

	return get
}
