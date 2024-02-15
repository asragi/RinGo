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
		expectedError    error
	}

	testCases := []testCase{
		{
			res:              GetResourceRes{},
			validateToken:    nil,
			getResourceError: nil,
			userId:           "id",
			expectedError:    nil,
		},
	}

	for _, v := range testCases {
		var passedUserId core.UserId

		var passedUserIdToResource core.UserId
		getResource := func(id core.UserId) (GetResourceRes, error) {
			passedUserIdToResource = id
			return v.res, v.getResourceError
		}

		getFunc := CreateGetUserResourceService(getResource)
		res, err := getFunc(v.userId)
		if err != v.expectedError {
			t.Errorf("expected err: %s, got: %s", v.expectedError, err)
		}
		if !reflect.DeepEqual(v.res, res) {
			t.Errorf("expected: %+v, got:%+v", v.res, res)
		}
		if v.userId != passedUserId {
			t.Errorf("expected: %s, got: %s", v.userId, passedUserId)
		}
		if v.userId != passedUserIdToResource {
			t.Errorf("expected: %s, got: %s", v.userId, passedUserIdToResource)
		}
	}

}
