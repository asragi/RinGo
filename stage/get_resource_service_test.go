package stage

import (
	"reflect"
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateGetUserResourceService(t *testing.T) {
	type testCase struct {
		res              GetResourceRes
		validateToken    error
		getResourceError error
		userId           core.UserId
		token            core.AccessToken
		expectedError    error
	}

	testCases := []testCase{
		{
			res:              GetResourceRes{},
			validateToken:    nil,
			getResourceError: nil,
			userId:           "id",
			token:            "token",
			expectedError:    nil,
		},
	}

	for _, v := range testCases {
		var passedUserId core.UserId
		var passedToken core.AccessToken
		validateToken := func(id core.UserId, token core.AccessToken) error {
			passedUserId = id
			passedToken = token
			return v.validateToken
		}

		var passedUserIdToResource core.UserId
		var passedTokenToResource core.AccessToken
		getResource := func(id core.UserId, token core.AccessToken) (GetResourceRes, error) {
			passedUserIdToResource = id
			passedTokenToResource = token
			return v.res, v.getResourceError
		}

		getFunc := CreateGetUserResourceService(validateToken, getResource)
		res, err := getFunc(v.userId, v.token)
		if err != v.expectedError {
			t.Errorf("expected err: %s, got: %s", v.expectedError, err)
		}
		if !reflect.DeepEqual(v.res, res) {
			t.Errorf("expected: %+v, got:%+v", v.res, res)
		}
		if v.userId != passedUserId {
			t.Errorf("expected: %s, got: %s", v.userId, passedUserId)
		}
		if v.token != passedToken {
			t.Errorf("expected: %s, got: %s", v.token, passedToken)
		}
		if v.userId != passedUserIdToResource {
			t.Errorf("expected: %s, got: %s", v.userId, passedUserIdToResource)
		}
		if v.token != passedTokenToResource {
			t.Errorf("expected: %s, got: %s", v.token, passedTokenToResource)
		}
	}

}
