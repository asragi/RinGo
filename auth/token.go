package auth

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/crypto"
	"github.com/asragi/RinGo/utils"
	"strings"
)

type CreateTokenFunc func(core.UserId) (AccessToken, error)

// TODO: remove "*"
type Sha256Func func(SecretHashKey, string) (string, error)

func CryptWithSha256(key SecretHashKey, text string) (string, error) {
	keyString := string(key)
	return crypto.SHA256WithKey(keyString, text)
}

type AccessToken string

func NewAccessToken(token string) (AccessToken, error) {
	if len(token) <= 0 {
		return "", TokenIsInvalidError{token: AccessToken(token)}
	}
	return AccessToken(token), nil
}

func (token *AccessToken) String() string {
	return string(*token)
}

type ExpirationTime int
type AccessTokenInformation struct {
	UserId         core.UserId
	ExpirationTime ExpirationTime
}

type AccessTokenInformationFromJson struct {
	UserId         string `json:"UserId"`
	ExpirationTime int    `json:"ExpirationTime"`
}

func (info *AccessTokenInformationFromJson) ToInformation() (*AccessTokenInformation, error) {
	userId, err := core.NewUserId(info.UserId)
	if err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}
	return &AccessTokenInformation{
		UserId:         userId,
		ExpirationTime: ExpirationTime(info.ExpirationTime),
	}, nil
}

func CreateTokenFuncEmitter(
	base64Encode Base64EncodeFunc,
	getTime core.GetCurrentTimeFunc,
	jsonFunc utils.StructToJsonFunc[AccessTokenInformation],
	secret SecretHashKey,
	sha256 Sha256Func,
) CreateTokenFunc {
	return func(userId core.UserId) (AccessToken, error) {
		handleError := func(err error) (AccessToken, error) {
			return "", fmt.Errorf("create token: %w", err)
		}
		nowTime := getTime()
		header := `{ alg: 'HS256', typ: 'JWT' }`
		info := &AccessTokenInformation{
			UserId:         userId,
			ExpirationTime: ExpirationTime(nowTime.Unix()),
		}
		payload, err := jsonFunc(info)
		if err != nil {
			return handleError(err)
		}
		unsignedToken := fmt.Sprintf("%s.%s", base64Encode(header), base64Encode(*payload))
		signature, err := sha256(secret, unsignedToken)
		if err != nil {
			return handleError(err)
		}
		jwt := fmt.Sprintf("%s.%s", unsignedToken, signature)
		token := AccessToken(jwt)
		return token, nil
	}
}

type GetTokenInformationFunc func(token *AccessToken) (*AccessTokenInformation, error)

func CreateGetTokenInformation(
	decodeBase64 Base64DecodeFunc,
	unmarshalJson utils.JsonToStructFunc[AccessTokenInformationFromJson],
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
		info, err := tokenInfo.ToInformation()
		if err != nil {
			return handleError(err)
		}
		return info, nil
	}
}

type CompareToken func(token *AccessToken) error

func CreateCompareToken(key SecretHashKey, sha256 Sha256Func) CompareToken {
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
		hashedUnsignedToken, err := sha256(key, unsignedSignature)
		if err != nil {
			return fmt.Errorf("compare token: %w", err)
		}
		if hashedUnsignedToken != signature {
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
