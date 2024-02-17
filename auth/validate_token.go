package auth

import "fmt"

type ValidateTokenFunc func(*AccessToken) (*AccessTokenInformation, error)

func CreateValidateToken(compare CompareToken, informationFunc GetTokenInformationFunc) ValidateTokenFunc {
	return func(token *AccessToken) (*AccessTokenInformation, error) {
		handleError := func(err error) (*AccessTokenInformation, error) {
			return nil, fmt.Errorf("validate token: %w", err)
		}
		err := token.IsValid(compare)
		if err != nil {
			return handleError(err)
		}
		info, err := token.GetInformation(informationFunc)
		if err != nil {
			return handleError(err)
		}
		return info, nil
	}
}
