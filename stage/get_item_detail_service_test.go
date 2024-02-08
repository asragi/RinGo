package stage

import (
	"errors"
	"github.com/asragi/RinGo/test"
	"reflect"
	"testing"
	"time"

	"github.com/asragi/RinGo/core"
)

func TestCreateGetItemDetailService(t *testing.T) {
	type testCase struct {
		req                 GetUserItemDetailReq
		expectedErr         error
		mockTime            time.Time
		mockExplore         []UserExplore
		mockArgs            getItemDetailArgs
		mockCompensatedArgs CompensatedMakeUserExploreArgs
	}

	testCases := []testCase{
		{
			req:         GetUserItemDetailReq{},
			expectedErr: nil,
		},
	}

	for _, v := range testCases {
		timer := func() time.Time {
			return v.mockTime
		}
		createArgs := func(
			GetUserItemDetailReq,
		) (getItemDetailArgs, error) {
			return v.mockArgs, nil
		}
		getAllItem := func(
			[]ExploreStaminaPair,
			[]GetExploreMasterRes,
			compensatedMakeUserExploreFunc,
		) []UserExplore {
			return v.mockExplore
		}
		makeUserExplore := func(args makeUserExploreArrayArgs) []UserExplore {
			return v.mockExplore
		}
		var passedExploreIds []ExploreId
		fetchUserExploreArgs := func(
			id core.UserId,
			token core.AccessToken,
			ids []ExploreId,
		) (CompensatedMakeUserExploreArgs, error) {
			passedExploreIds = ids
			return v.mockCompensatedArgs, nil
		}
		compensatedMakeUserExplore := func(
			repoArgs CompensatedMakeUserExploreArgs,
			currentTimer core.GetCurrentTimeFunc,
			execNum int,
			makeUserExplore MakeUserExploreArrayFunc,
		) compensatedMakeUserExploreFunc {
			return func(args makeUserExploreArgs) []UserExplore {
				return v.mockExplore
			}
		}
		getItemDetail := CreateGetItemDetailService(
			timer,
			createArgs,
			getAllItem,
			makeUserExplore,
			fetchUserExploreArgs,
			compensatedMakeUserExplore,
		)
		expectedPassedExploreIds := UserExploreToIdArray(v.mockExplore)
		expectedRes := getUserItemDetailRes{
			UserId:       v.mockArgs.storageRes.UserId,
			ItemId:       v.mockArgs.masterRes.ItemId,
			Price:        v.mockArgs.masterRes.Price,
			DisplayName:  v.mockArgs.masterRes.DisplayName,
			Description:  v.mockArgs.masterRes.Description,
			MaxStock:     v.mockArgs.masterRes.MaxStock,
			Stock:        v.mockArgs.storageRes.Stock,
			UserExplores: v.mockExplore,
		}
		res, err := getItemDetail(v.req)
		if !errors.Is(err, v.expectedErr) {
			t.Errorf("expect: %s, got: %s", v.expectedErr.Error(), err.Error())
		}
		if !test.DeepEqual(expectedRes, res) {
			t.Errorf("expect: %+v, got: %+v", expectedRes, res)
		}
		if !test.DeepEqual(passedExploreIds, expectedPassedExploreIds) {
			t.Errorf("expect: %+v, got: %+v", expectedPassedExploreIds, passedExploreIds)
		}
	}
}

func TestFetchGetItemDetailArgs(t *testing.T) {
	type testCase struct {
		request                GetUserItemDetailReq
		expectedError          error
		mockGetItemMasterRes   GetItemMasterRes
		mockGetItemStorageRes  ItemData
		mockGetExploreRes      []GetExploreMasterRes
		mockItemExplore        []ExploreId
		mockExploreStaminaPair []ExploreStaminaPair
	}

	userId := core.UserId("user")
	itemId := core.ItemId("item")

	testCases := []testCase{
		{
			request: GetUserItemDetailReq{
				UserId:      userId,
				ItemId:      itemId,
				AccessToken: "token",
			},
			expectedError: nil,
			mockGetItemMasterRes: GetItemMasterRes{
				ItemId:      itemId,
				Price:       300,
				DisplayName: "TestItem",
				Description: "TestDesc",
				MaxStock:    100,
			},
			mockGetItemStorageRes: ItemData{
				UserId:  userId,
				ItemId:  itemId,
				Stock:   50,
				IsKnown: true,
			},
			mockGetExploreRes:      nil,
			mockItemExplore:        nil,
			mockExploreStaminaPair: nil,
		},
	}

	for i, v := range testCases {
		expectedRes := getItemDetailArgs{
			masterRes: v.mockGetItemMasterRes,
			storageRes: GetItemStorageRes{
				UserId: userId,
				Stock:  v.mockGetItemStorageRes.Stock,
			},
			exploreStaminaPair: v.mockExploreStaminaPair,
			explores:           v.mockGetExploreRes,
		}
		var mockGetItemArgs core.ItemId
		mockGetItemMaster := func(itemId []core.ItemId) ([]GetItemMasterRes, error) {
			mockGetItemArgs = itemId[0]
			return []GetItemMasterRes{v.mockGetItemMasterRes}, nil
		}

		var passedStorageArg core.ItemId
		mockGetItemStorage := func(userId core.UserId, itemId []core.ItemId, token core.AccessToken) (
			BatchGetStorageRes,
			error,
		) {
			passedStorageArg = itemId[0]
			return BatchGetStorageRes{
				UserId:   userId,
				ItemData: []ItemData{v.mockGetItemStorageRes},
			}, nil
		}
		var passedExploreArgs []ExploreId
		mockExploreMaster := func(exploreIds []ExploreId) ([]GetExploreMasterRes, error) {
			passedExploreArgs = exploreIds
			return v.mockGetExploreRes, nil
		}
		var passedStaminaArgs []ExploreId
		consumingStamina := func(
			userId core.UserId,
			token core.AccessToken,
			ids []ExploreId,
		) ([]ExploreStaminaPair, error) {
			passedStaminaArgs = ids
			return v.mockExploreStaminaPair, nil
		}
		var passedItemRelationArg core.ItemId
		mockItemExplore := func(itemId core.ItemId) ([]ExploreId, error) {
			passedItemRelationArg = itemId
			return v.mockItemExplore, nil
		}
		req := v.request
		res, err := FetchGetItemDetailArgs(
			v.request,
			mockGetItemMaster,
			mockGetItemStorage,
			mockExploreMaster,
			mockItemExplore,
			consumingStamina,
		)
		if !errors.Is(err, v.expectedError) {
			t.Fatalf(
				"case: %d, expect error is: %s, got: %s",
				i,
				test.ErrorToString(v.expectedError),
				test.ErrorToString(err),
			)
		}
		if mockGetItemArgs != req.ItemId {
			t.Errorf("expect: %s, got: %s", req.ItemId, mockGetItemArgs)
		}
		if passedStorageArg != req.ItemId {
			t.Errorf("expect: %s, got: %s", req.ItemId, passedStorageArg)
		}
		if !reflect.DeepEqual(passedExploreArgs, v.mockItemExplore) {
			t.Errorf("expect: %s, got: %s", v.mockItemExplore, passedExploreArgs)
		}
		if passedItemRelationArg != req.ItemId {
			t.Errorf("expect: %s, got: %s", req.ItemId, passedItemRelationArg)
		}
		if !reflect.DeepEqual(passedStaminaArgs, v.mockItemExplore) {
			t.Errorf(
				"mockReducedStamina args and explore res not matched: mock args: %+v, res: %+v",
				v.mockItemExplore,
				passedStaminaArgs,
			)
		}
		if !reflect.DeepEqual(expectedRes, res) {
			t.Errorf("expect:%+v, got:%+v", expectedRes, res)
		}
	}
}
