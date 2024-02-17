package auth

import (
	"testing"
)

func TestStringToBase64(t *testing.T) {
	type testCase struct {
		text       string
		base64Text string
	}

	testCases := []*testCase{
		{
			text:       "text",
			base64Text: "dGV4dA==",
		},
	}

	for _, v := range testCases {
		base64String := stringToBase64(v.text)
		if base64String != v.base64Text {
			t.Errorf("expected: %s, got %s", v.base64Text, base64String)
		}
	}

	for _, v := range testCases {
		decodedText, err := Base64ToString(v.base64Text)
		if err != nil {
			t.Fatalf("error: %s", err.Error())
		}
		if decodedText != v.text {
			t.Errorf("expected: %s, got %s", v.text, decodedText)
		}
	}
}
