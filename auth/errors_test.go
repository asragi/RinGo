package auth

import (
	"errors"
	"fmt"
	"github.com/asragi/RinGo/core"
	"testing"
)

func TestCreateUserAlreadyExistsError(t *testing.T) {
	type testCase struct {
		userId      core.UserId
		expectedErr error
	}
	userId := core.UserId("user")
	testCases := []testCase{
		{userId: userId, expectedErr: UserAlreadyExistsError},
	}
	for _, v := range testCases {
		testErr := UserAlreadyExistsError
		testErrorWrapped := fmt.Errorf("wrapped: %w", testErr)
		if !errors.Is(testErr, v.expectedErr) {
			t.Errorf("Error didn't match: %s, %s", testErr, v.expectedErr)
		}
		if !errors.Is(testErrorWrapped, v.expectedErr) {
			t.Errorf("Error didn't match: %s, %s", testErrorWrapped, v.expectedErr)
		}
	}
}
