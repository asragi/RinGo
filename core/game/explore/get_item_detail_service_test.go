package explore

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCreateGetItemDetailService(t *testing.T) {
	type testCase struct {
		req         GetUserItemDetailReq
		expectedErr error
		mockArgs    getItemDetailArgs
	}

	userId := core.UserId("user")
	itemId := core.ItemId("item")
	exploreId := game.ExploreId("explore")
	testCases := []testCase{
		{
			req: GetUserItemDetailReq{
				UserId: userId,
				ItemId: itemId,
			},
			expectedErr: nil,
			mockArgs: getItemDetailArgs{
				masterRes: &game.GetItemMasterRes{
					ItemId:      itemId,
					Price:       111,
					DisplayName: "display_name",
					Description: "desc",
					MaxStock:    222,
				},
				storageRes: &game.StorageData{
					UserId: userId,
					Stock:  120,
				},
				exploreStaminaPair: []*game.ExploreStaminaPair{
					{
						ExploreId:      exploreId,
						ReducedStamina: 10,
					},
				},
				explores: []*game.GetExploreMasterRes{
					{
						ExploreId:            exploreId,
						DisplayName:          "explore_name",
						Description:          "explore_desc",
						ConsumingStamina:     20,
						RequiredPayment:      90,
						StaminaReducibleRate: 0.4,
					},
				},
			},
		},
	}

	for _, v := range testCases {
		mockGenerateArgs := func(ctx context.Context, id core.UserId, itemId core.ItemId) (*getItemDetailArgs, error) {
			return &v.mockArgs, nil
		}
		getItemDetail := CreateGetItemDetailService(
			mockGenerateArgs,
		)
		expectedRes := &getUserItemDetailRes{
			UserId:       v.mockArgs.storageRes.UserId,
			ItemId:       v.mockArgs.masterRes.ItemId,
			Price:        v.mockArgs.masterRes.Price,
			DisplayName:  v.mockArgs.masterRes.DisplayName,
			Description:  v.mockArgs.masterRes.Description,
			MaxStock:     v.mockArgs.masterRes.MaxStock,
			Stock:        v.mockArgs.storageRes.Stock,
			UserExplores: nil,
		}
		ctx := test.MockCreateContext()
		res, err := getItemDetail(ctx, v.req.UserId, v.req.ItemId)
		if !errors.Is(err, v.expectedErr) {
			t.Errorf("expect: %s, got: %s", v.expectedErr.Error(), err.Error())
		}
		if !test.DeepEqual(expectedRes, res) {
			t.Errorf("expect: %+v, got: %+v", expectedRes, res)
			if expectedRes.UserId != res.UserId {
				t.Errorf("expect: %s, got: %s", expectedRes.UserId, res.UserId)
			}
			if expectedRes.ItemId != res.ItemId {
				t.Errorf("expect: %s, got: %s", expectedRes.ItemId, res.ItemId)
			}
			if expectedRes.Price != res.Price {
				t.Errorf("expect: %d, got: %d", expectedRes.Price, res.Price)
			}
			if expectedRes.DisplayName != res.DisplayName {
				t.Errorf("expect: %s, got: %s", expectedRes.DisplayName, res.DisplayName)
			}
			if expectedRes.Description != res.Description {
				t.Errorf("expect: %s, got: %s", expectedRes.Description, res.Description)
			}
			if expectedRes.MaxStock != res.MaxStock {
				t.Errorf("expect: %d, got: %d", expectedRes.MaxStock, res.MaxStock)
			}
			if expectedRes.Stock != res.Stock {
				t.Errorf("expect: %d, got: %d", expectedRes.Stock, res.Stock)
			}
			if !test.DeepEqual(expectedRes.UserExplores, res.UserExplores) {
				t.Errorf("expect: %+v, got: %+v", expectedRes.UserExplores, res.UserExplores)
			}
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
		mockUserExplore        []*game.UserExplore
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
		expectedRes := &getItemDetailArgs{
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
		mockMakeUserExplore := func(
			ctx context.Context,
			id core.UserId,
			exploreIds []game.ExploreId,
			execNum int,
		) ([]*game.UserExplore, error) {
			return v.mockUserExplore, nil
		}
		req := v.request
		ctx := test.MockCreateContext()
		generateGetItemDetailArgs := CreateGenerateGetItemDetailArgs(
			mockGetItemMaster,
			mockGetItemStorage,
			mockExploreMaster,
			mockItemExplore,
			consumingStamina,
			mockMakeUserExplore,
		)
		resArgs, err := generateGetItemDetailArgs(ctx, req.UserId, req.ItemId)
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
		if !test.DeepEqual(passedExploreArgs, v.mockItemExplore) {
			t.Errorf("expect: %s, got: %s", v.mockItemExplore, passedExploreArgs)
		}
		if passedItemRelationArg != req.ItemId {
			t.Errorf("expect: %s, got: %s", req.ItemId, passedItemRelationArg)
		}
		if !test.DeepEqual(passedStaminaArgs, v.mockItemExplore) {
			t.Errorf(
				"mockReducedStamina args and explore res not matched: mock args: %+v, res: %+v",
				v.mockItemExplore,
				passedStaminaArgs,
			)
		}
		if !test.DeepEqual(expectedRes, resArgs) {
			t.Errorf("expect:%+v, got:%+v", expectedRes, resArgs)
			if !test.DeepEqual(expectedRes.masterRes, resArgs.masterRes) {
				t.Errorf("masterRes expect: %+v, got: %+v", expectedRes.masterRes, resArgs.masterRes)
			}
			if !test.DeepEqual(expectedRes.storageRes, resArgs.storageRes) {
				t.Errorf("storageRes expect: %+v, got: %+v", expectedRes.storageRes, resArgs.storageRes)
			}
			if !test.DeepEqual(expectedRes.exploreStaminaPair, resArgs.exploreStaminaPair) {
				t.Errorf(
					"exploreStaminaPair expect: %+v, got: %+v",
					expectedRes.exploreStaminaPair,
					resArgs.exploreStaminaPair,
				)
			}
			if !test.DeepEqual(expectedRes.explores, resArgs.explores) {
				t.Errorf("explores expect: %+v, got: %+v", expectedRes.explores, resArgs.explores)
			}
		}
	}
}
