package auth

import (
	"fmt"
	"github.com/asragi/RinGo/core"
)

type createTokenFunc func(core.UserId) (AccessToken, error)
type sha256Func func(secretHashKey, string) string

func createToken(
	base64Encode base64EncodeFunc,
	getTime core.GetCurrentTimeFunc,
	secret secretHashKey,
	sha256 sha256Func,
) createTokenFunc {
	return func(userId core.UserId) (AccessToken, error) {
		nowTime := getTime()
		header := `{ alg: 'HS256', typ: 'JWT' }`
		payload := fmt.Sprintf(`{ sub: '%s', iat: %d}`, userId, nowTime.Unix())
		unsignedToken := fmt.Sprintf("%s.%s", base64Encode(header), base64Encode(payload))
		signature := sha256(secret, unsignedToken)
		jwt := fmt.Sprintf("%s.%s", unsignedToken, signature)
		return AccessToken(jwt), nil
	}
}

type AccessToken string

func (token *AccessToken) IsValid() error {
	if len(*token) <= 0 {
		return TokenIsInvalidError{token: *token}
	}
	return nil
}

type TokenIsInvalidError struct {
	token AccessToken
}

func (e TokenIsInvalidError) Error() string {
	return fmt.Sprintf("token is invalid: %s", e.token)
}
