package explore

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core/game"
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
		mockExplore         []*game.UserExplore
		mockArgs            getItemDetailArgs
		mockCompensatedArgs *game.CompensatedMakeUserExploreArgs
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
				masterRes: &game.GetItemMasterRes{
					ItemId:      "",
					Price:       0,
					DisplayName: "",
					Description: "",
					MaxStock:    0,
				},
				storageRes: &game.StorageData{
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
			[]*game.ExploreStaminaPair,
			[]*game.GetExploreMasterRes,
			game.compensatedMakeUserExploreFunc,
		) []*game.UserExplore {
			return v.mockExplore
		}
		makeUserExplore := func(args *game.makeUserExploreArrayArgs) []*game.UserExplore {
			return v.mockExplore
		}
		var passedExploreIds []game.ExploreId
		fetchUserExploreArgs := func(
			ctx context.Context,
			id core.UserId,
			ids []game.ExploreId,
		) (*game.CompensatedMakeUserExploreArgs, error) {
			passedExploreIds = ids
			return v.mockCompensatedArgs, nil
		}
		compensatedMakeUserExplore := func(
			repoArgs *game.CompensatedMakeUserExploreArgs,
			currentTimer core.GetCurrentTimeFunc,
			execNum int,
			makeUserExplore game.OldMakeUserExploreFunc,
		) game.compensatedMakeUserExploreFunc {
			return func(*game.makeUserExploreArgs) []*game.UserExplore {
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
		expectedPassedExploreIds := game.UserExploreToIdArray(v.mockExplore)
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
		mockGetItemMasterRes   *game.GetItemMasterRes
		mockGetItemStorageRes  *game.StorageData
		mockGetExploreRes      []*game.GetExploreMasterRes
		mockItemExplore        []game.ExploreId
		mockExploreStaminaPair []*game.ExploreStaminaPair
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
			mockGetItemMasterRes: &game.GetItemMasterRes{
				ItemId:      itemId,
				Price:       300,
				DisplayName: "TestItem",
				Description: "TestDesc",
				MaxStock:    100,
			},
			mockGetItemStorageRes: &game.StorageData{
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
		mockGetItemMaster := func(_ context.Context, itemId []core.ItemId) ([]*game.GetItemMasterRes, error) {
			mockGetItemArgs = itemId[0]
			return []*game.GetItemMasterRes{v.mockGetItemMasterRes}, nil
		}

		var passedStorageArg core.ItemId
		mockGetItemStorage := func(_ context.Context, userId core.UserId, itemId []core.ItemId) (
			game.BatchGetStorageRes,
			error,
		) {
			passedStorageArg = itemId[0]
			return game.BatchGetStorageRes{
				UserId:   userId,
				ItemData: []*game.StorageData{v.mockGetItemStorageRes},
			}, nil
		}
		var passedExploreArgs []game.ExploreId
		mockExploreMaster := func(_ context.Context, exploreIds []game.ExploreId) ([]*game.GetExploreMasterRes, error) {
			passedExploreArgs = exploreIds
			return v.mockGetExploreRes, nil
		}
		var passedStaminaArgs []game.ExploreId
		consumingStamina := func(
			ctx context.Context,
			userId core.UserId,
			ids []game.ExploreId,
		) ([]*game.ExploreStaminaPair, error) {
			passedStaminaArgs = ids
			return v.mockExploreStaminaPair, nil
		}
		var passedItemRelationArg core.ItemId
		mockItemExplore := func(ctx context.Context, itemId core.ItemId) ([]game.ExploreId, error) {
			passedItemRelationArg = itemId
			return v.mockItemExplore, nil
		}
		req := v.request
		ctx := test.MockCreateContext()
		res, err := CreateGenerateGetItemDetailArgs(
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
