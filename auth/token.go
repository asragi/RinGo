package auth

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/utils"
	"strings"
)

type createTokenFunc func(*core.UserId) (*AccessToken, error)
type sha256Func func(*SecretHashKey, *string) (*string, error)

type AccessToken string
type ExpirationTime int
type AccessTokenInformation struct {
	UserId         core.UserId    `json:"user_id"`
	ExpirationTime ExpirationTime `json:"exp"`
}

func CreateTokenFuncEmitter(
	base64Encode base64EncodeFunc,
	getTime core.GetCurrentTimeFunc,
	jsonFunc utils.StructToJsonFunc[AccessTokenInformation],
	secret SecretHashKey,
	sha256 sha256Func,
) createTokenFunc {
	return func(userId *core.UserId) (*AccessToken, error) {
		handleError := func(err error) (*AccessToken, error) {
			return nil, fmt.Errorf("create token: %w", err)
		}
		nowTime := getTime()
		header := `{ alg: 'HS256', typ: 'JWT' }`
		info := &AccessTokenInformation{
			UserId:         *userId,
			ExpirationTime: ExpirationTime(nowTime.Unix()),
		}
		payload, err := jsonFunc(info)
		if err != nil {
			return handleError(err)
		}
		unsignedToken := fmt.Sprintf("%s.%s", base64Encode(header), base64Encode(*payload))
		signature, err := sha256(&secret, &unsignedToken)
		if err != nil {
			return handleError(err)
		}
		jwt := fmt.Sprintf("%s.%s", unsignedToken, *signature)
		token := AccessToken(jwt)
		return &token, nil
	}
}

type GetTokenInformationFunc func(token *AccessToken) (*AccessTokenInformation, error)

func CreateGetTokenInformation(
	decodeBase64 base64DecodeFunc,
	unmarshalJson utils.JsonToStructFunc[AccessTokenInformation],
) GetTokenInformationFunc {
	return func(token *AccessToken) (*AccessTokenInformation, error) {
		handleError := func(err error) (*AccessTokenInformation, error) {
			return nil, fmt.Errorf("get token info: %w", err)
		}
		tokenString := string(*token)
		splitToken := strings.Split(tokenString, ".")
		if len(splitToken) != 3 {
			return nil, TokenIsInvalidError{token: *token}
		}
		payloadString := splitToken[1]
		payloadJsonString, err := decodeBase64(payloadString)
		if err != nil {
			return handleError(err)
		}

		tokenInfo, err := unmarshalJson(payloadJsonString)
		if err != nil {
			return handleError(err)
		}
		return tokenInfo, nil
	}
}

type CompareToken func(token *AccessToken) error

func CreateCompareToken(key *SecretHashKey, sha256 sha256Func) CompareToken {
	return func(token *AccessToken) error {
		if len(*token) <= 0 {
			return TokenIsInvalidError{token: *token}
		}
		tokenString := string(*token)
		splitToken := strings.Split(tokenString, ".")
		if len(splitToken) != 3 {
			return TokenIsInvalidError{token: *token}
		}
		unsignedSignature := fmt.Sprintf("%s.%s", splitToken[0], splitToken[1])
		signature := splitToken[2]
		hashedUnsignedToken, err := sha256(key, &unsignedSignature)
		if err != nil {
			return fmt.Errorf("compare token: %w", err)
		}
		if *hashedUnsignedToken != signature {
			return TokenIsInvalidError{token: *token}
		}
		return nil
	}
}

func (token *AccessToken) GetInformation(getInfo GetTokenInformationFunc) (*AccessTokenInformation, error) {
	return getInfo(token)
}

func (token *AccessToken) IsValid(compare CompareToken) error {
	return compare(token)
}

type TokenIsInvalidError struct {
	token AccessToken
}

func (e TokenIsInvalidError) Error() string {
	return fmt.Sprintf("token is invalid: %s", e.token)
}
