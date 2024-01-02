package core

type ValidateTokenFunc func(UserId, AccessToken) error

func CreateValidateTokenService(
	validateRepoFunc ValidateTokenRepoFunc,
) ValidateTokenFunc {
	validate := func(
		userId UserId,
		token AccessToken,
	) error {
		return validateRepoFunc(userId, token)
	}

	return validate
}
