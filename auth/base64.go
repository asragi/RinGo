package auth

import (
	"encoding/base64"
	"fmt"
)

type base64EncodeFunc func(string) string

func StringToBase64(text string) string {
	src := []byte(text)
	return base64.StdEncoding.EncodeToString(src)
}

type base64DecodeFunc func(string) (string, error)

func Base64ToString(base64Text string) (string, error) {
	dec, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}
	return string(dec), nil
}
