package auth

import "testing"

func TestCreateHashedPassword(t *testing.T) {
	type testCase struct {
		text     string
		expected string
	}

	testCases := []*testCase{
		{
			text:     "text",
			expected: "expected",
		},
	}

	for _, v := range testCases {
		var passedText string
		mockEncrypt := func(s string) (string, error) {
			passedText = s
			return v.expected, nil
		}

		f := createHashedPassword(mockEncrypt)
		res, _ := f(rowPassword(v.text))
		expectedPass := hashedPassword(v.expected)
		if res != expectedPass {
			t.Errorf("expected: %s, got %s", v.expected, string(res))
		}
		if passedText != v.text {
			t.Errorf("expected: %s, got: %s", v.text, passedText)
		}

	}
}
