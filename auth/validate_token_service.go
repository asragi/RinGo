package auth

import "github.com/asragi/RinGo/core"

type ValidateTokenFunc func(core.UserId, AccessToken) error

type ValidateTokenServiceFunc func(ValidateTokenRepoFunc) ValidateTokenFunc

func CreateValidateTokenService(
	validateRepoFunc ValidateTokenRepoFunc,
) ValidateTokenFunc {
	validate := func(
		userId core.UserId,
		token AccessToken,
	) error {
		return validateRepoFunc(userId, token)
	}

	return validate
}
