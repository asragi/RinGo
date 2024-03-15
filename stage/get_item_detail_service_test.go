package stage

import (
	"context"
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
		mockExplore         []*UserExplore
		mockArgs            getItemDetailArgs
		mockCompensatedArgs *CompensatedMakeUserExploreArgs
	}

	testCases := []testCase{
		{
			req: GetUserItemDetailReq{
				UserId: "",
				ItemId: "",
			},
			expectedErr: nil,
			mockTime:    test.MockTime(),
			mockExplore: nil,
			mockArgs: getItemDetailArgs{
				masterRes: &GetItemMasterRes{
					ItemId:      "",
					Price:       0,
					DisplayName: "",
					Description: "",
					MaxStock:    0,
				},
				storageRes: &StorageData{
					UserId: "",
					Stock:  0,
				},
				exploreStaminaPair: nil,
				explores:           nil,
			},
			mockCompensatedArgs: nil,
		},
	}

	for _, v := range testCases {
		timer := func() time.Time {
			return v.mockTime
		}
		createArgs := func(
			context.Context,
			GetUserItemDetailReq,
		) (getItemDetailArgs, error) {
			return v.mockArgs, nil
		}
		getAllItem := func(
			[]*ExploreStaminaPair,
			[]*GetExploreMasterRes,
			compensatedMakeUserExploreFunc,
		) []*UserExplore {
			return v.mockExplore
		}
		makeUserExplore := func(args *makeUserExploreArrayArgs) []*UserExplore {
			return v.mockExplore
		}
		var passedExploreIds []ExploreId
		fetchUserExploreArgs := func(
			ctx context.Context,
			id core.UserId,
			ids []ExploreId,
		) (*CompensatedMakeUserExploreArgs, error) {
			passedExploreIds = ids
			return v.mockCompensatedArgs, nil
		}
		compensatedMakeUserExplore := func(
			repoArgs *CompensatedMakeUserExploreArgs,
			currentTimer core.GetCurrentTimeFunc,
			execNum int,
			makeUserExplore MakeUserExploreArrayFunc,
		) compensatedMakeUserExploreFunc {
			return func(*makeUserExploreArgs) []*UserExplore {
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
		ctx := test.MockCreateContext()
		res, err := getItemDetail(ctx, v.req)
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
		mockGetItemMasterRes   *GetItemMasterRes
		mockGetItemStorageRes  *StorageData
		mockGetExploreRes      []*GetExploreMasterRes
		mockItemExplore        []ExploreId
		mockExploreStaminaPair []*ExploreStaminaPair
	}

	userId := core.UserId("user")
	itemId := core.ItemId("item")

	testCases := []testCase{
		{
			request: GetUserItemDetailReq{
				UserId: userId,
				ItemId: itemId,
			},
			expectedError: nil,
			mockGetItemMasterRes: &GetItemMasterRes{
				ItemId:      itemId,
				Price:       300,
				DisplayName: "TestItem",
				Description: "TestDesc",
				MaxStock:    100,
			},
			mockGetItemStorageRes: &StorageData{
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
			masterRes:          v.mockGetItemMasterRes,
			storageRes:         v.mockGetItemStorageRes,
			exploreStaminaPair: v.mockExploreStaminaPair,
			explores:           v.mockGetExploreRes,
		}
		var mockGetItemArgs core.ItemId
		mockGetItemMaster := func(_ context.Context, itemId []core.ItemId) ([]*GetItemMasterRes, error) {
			mockGetItemArgs = itemId[0]
			return []*GetItemMasterRes{v.mockGetItemMasterRes}, nil
		}

		var passedStorageArg core.ItemId
		mockGetItemStorage := func(_ context.Context, userId core.UserId, itemId []core.ItemId) (
			BatchGetStorageRes,
			error,
		) {
			passedStorageArg = itemId[0]
			return BatchGetStorageRes{
				UserId:   userId,
				ItemData: []*StorageData{v.mockGetItemStorageRes},
			}, nil
		}
		var passedExploreArgs []ExploreId
		mockExploreMaster := func(_ context.Context, exploreIds []ExploreId) ([]*GetExploreMasterRes, error) {
			passedExploreArgs = exploreIds
			return v.mockGetExploreRes, nil
		}
		var passedStaminaArgs []ExploreId
		consumingStamina := func(
			ctx context.Context,
			userId core.UserId,
			ids []ExploreId,
		) ([]*ExploreStaminaPair, error) {
			passedStaminaArgs = ids
			return v.mockExploreStaminaPair, nil
		}
		var passedItemRelationArg core.ItemId
		mockItemExplore := func(ctx context.Context, itemId core.ItemId) ([]ExploreId, error) {
			passedItemRelationArg = itemId
			return v.mockItemExplore, nil
		}
		req := v.request
		ctx := test.MockCreateContext()
		res, err := FetchGetItemDetailArgs(
			ctx,
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
