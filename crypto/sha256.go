package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type SHA256Func func(*string, *string) (*string, error)

func SHA256WithKey(key, msg *string) (*string, error) {
	mac := hmac.New(sha256.New, []byte(*key))
	_, err := mac.Write([]byte(*msg))
	if err != nil {
		return nil, fmt.Errorf("sha256 with key: %w", err)
	}
	str := hex.EncodeToString(mac.Sum(nil))
	return &str, nil
}
