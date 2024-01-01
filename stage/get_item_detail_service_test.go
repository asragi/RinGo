package stage

import (
	"reflect"
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateGetItemDetailService(t *testing.T) {
	type testCase struct {
		req         GetUserItemDetailReq
		res         getUserItemDetailRes
		expectedErr error
	}

	testCases := []testCase{
		{
			req:         GetUserItemDetailReq{},
			res:         getUserItemDetailRes{},
			expectedErr: nil,
		},
	}

	for _, v := range testCases {
		createArgs := func(
			GetUserItemDetailReq,
		) (getItemDetailArgs, error) {
			return getItemDetailArgs{}, nil
		}
		getAllItem := func(
			[]ExploreStaminaPair,
			[]GetExploreMasterRes,
			compensatedMakeUserExploreFunc,
		) []UserExplore {
			return nil
		}
		calcBatchConsumingStaminaFunc := func(
			makeUserExploreArgs,
		) []UserExplore {
			return nil
		}
		getItemDetail := CreateGetItemDetailService(
			createArgs,
			getAllItem,
			calcBatchConsumingStaminaFunc,
		)

		res, err := getItemDetail(v.req)
		if v.expectedErr != err {
			t.Errorf("expect: %s, got: %s", v.expectedErr.Error(), err.Error())
		}
		if !reflect.DeepEqual(v.res, res) {
			t.Errorf("expect: %+v, got: %+v", v.res, res)
		}
	}
}

func TestCreateGetItemDetailArgs(t *testing.T) {
	type testCase struct {
		request       GetUserItemDetailReq
		expect        getItemDetailArgs
		expectedError error
	}

	var mockGetItemArgs core.ItemId
	getItemRes := GetItemMasterRes{}
	mockGetItemMaster := func(itemId core.ItemId) (GetItemMasterRes, error) {
		mockGetItemArgs = itemId
		return getItemRes, nil
	}
	getStorageRes := GetItemStorageRes{}
	var mockStorageArg core.ItemId
	mockGetItemStorage := func(userId core.UserId, itemId core.ItemId, token core.AccessToken) (GetItemStorageRes, error) {
		mockStorageArg = itemId
		return getStorageRes, nil
	}
	var mockExploreArgs []ExploreId
	getExploreMasterRes := []GetExploreMasterRes{}
	mockExploreMaster := func(exploreIds []ExploreId) ([]GetExploreMasterRes, error) {
		mockExploreArgs = exploreIds
		return getExploreMasterRes, nil
	}
	var mockItemRelationArg core.ItemId
	itemExploreRelation := []ExploreId{}
	mockItemExplore := func(itemId core.ItemId) ([]ExploreId, error) {
		mockItemRelationArg = itemId
		return itemExploreRelation, nil
	}
	var mockStaminaArgs []GetExploreMasterRes
	exploreStamina := []ExploreStaminaPair{}
	consumingStamina := func(userId core.UserId, token core.AccessToken, masters []GetExploreMasterRes) ([]ExploreStaminaPair, error) {
		mockStaminaArgs = masters
		return exploreStamina, nil
	}
	testCases := []testCase{}

	for i, v := range testCases {
		req := v.request
		res, err := createGetItemDetailArgs(
			v.request,
			mockGetItemMaster,
			mockGetItemStorage,
			mockExploreMaster,
			mockItemExplore,
			consumingStamina,
		)
		if err != v.expectedError {
			t.Fatalf("case: %d, expect error is: %s, got: %s", i, v.expectedError.Error(), err.Error())
		}

		if mockGetItemArgs != req.ItemId {
			t.Errorf("expect: %s, got: %s", req.ItemId, mockGetItemArgs)
		}
		if mockStorageArg != req.ItemId {
			t.Errorf("expect: %s, got: %s", req.ItemId, mockStorageArg)
		}
		if !reflect.DeepEqual(mockExploreArgs, itemExploreRelation) {
			t.Errorf("expect: %s, got: %s", itemExploreRelation, mockExploreArgs)
		}
		if mockItemRelationArg != req.ItemId {
			t.Errorf("expect: %s, got: %s", req.ItemId, mockItemRelationArg)
		}
		if !reflect.DeepEqual(mockStaminaArgs, getExploreMasterRes) {
			t.Errorf("stamina args and explore res not matched: mock args: %+v, res: %+v", mockStaminaArgs, getExploreMasterRes)
		}
		if !reflect.DeepEqual(v.expect, res) {
			t.Errorf("expect and actual are not matched:%+v, %+v", res, v.expect)
		}
	}
}
