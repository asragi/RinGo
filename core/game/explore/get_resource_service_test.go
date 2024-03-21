package explore

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/test"
	"reflect"
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateGetUserResourceService(t *testing.T) {
	type testCase struct {
		mockResourceRes  *game.GetResourceRes
		getResourceError error
		userId           core.UserId
		expectedError    error
	}

	testCases := []testCase{
		{
			mockResourceRes: &game.GetResourceRes{
				UserId:             "userId",
				MaxStamina:         3000,
				StaminaRecoverTime: core.StaminaRecoverTime(test.MockTime()),
				Fund:               1000,
			},
			getResourceError: nil,
			userId:           "id",
			expectedError:    nil,
		},
	}

	for _, v := range testCases {
		var passedUserIdToResource core.UserId
		getResource := func(ctx context.Context, id core.UserId) (*game.GetResourceRes, error) {
			passedUserIdToResource = id
			return v.mockResourceRes, v.getResourceError
		}
		getFunc := CreateGetUserResourceService(getResource)
		ctx := test.MockCreateContext()
		res, err := getFunc(ctx, v.userId)
		if !errors.Is(err, v.expectedError) {
			t.Errorf("expected err: %s, got: %s", v.expectedError, err)
		}
		if !reflect.DeepEqual(v.mockResourceRes, res) {
			t.Errorf("expected: %+v, got:%+v", v.mockResourceRes, res)
		}
		if v.userId != passedUserIdToResource {
			t.Errorf("expected: %s, got: %s", v.userId, passedUserIdToResource)
		}
	}

}
