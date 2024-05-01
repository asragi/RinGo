package game

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCreateGetItemListService(t *testing.T) {
	type testCase struct {
		request        core.UserId
		mockStorage    []*StorageData
		mockItemMaster []*GetItemMasterRes
		expectedError  error
	}

	testCases := []testCase{
		{
			request:     "userId",
			mockStorage: []*StorageData{},
			mockItemMaster: []*GetItemMasterRes{
				{
					ItemId:      "A",
					Price:       100,
					DisplayName: "NameA",
					Description: "DescA",
					MaxStock:    200,
				},
			},
			expectedError: nil,
		},
		{
			request: "userId",
			mockStorage: []*StorageData{
				{
					UserId:  "userId",
					ItemId:  "A",
					Stock:   199,
					IsKnown: true,
				},
			},
			mockItemMaster: []*GetItemMasterRes{
				{
					ItemId:      "A",
					Price:       100,
					DisplayName: "NameA",
					Description: "DescA",
					MaxStock:    200,
				},
				{
					ItemId:      "B",
					Price:       1000,
					DisplayName: "NameB",
					Description: "DescB",
					MaxStock:    2000,
				},
			},
			expectedError: nil,
		},
	}

	for i, v := range testCases {
		req := v.request
		expect := func() []*ItemListRow {
			storageMap := func() map[core.ItemId]core.Stock {
				result := map[core.ItemId]core.Stock{}
				for _, v := range v.mockStorage {
					result[v.ItemId] = v.Stock
				}
				return result
			}()
			var res []*ItemListRow
			for _, item := range v.mockItemMaster {
				stock := func() core.Stock {
					if _, ok := storageMap[item.ItemId]; !ok {
						return 0
					}
					return storageMap[item.ItemId]
				}()
				res = append(
					res, &ItemListRow{
						ItemId:      item.ItemId,
						DisplayName: item.DisplayName,
						Stock:       stock,
						MaxStock:    item.MaxStock,
						Price:       item.Price,
					},
				)
			}
			return res
		}()
		mockGetAllStorage := func(ctx context.Context, id core.UserId) ([]*StorageData, error) {
			return v.mockStorage, nil
		}
		mockFetchItemMaster := func(ctx context.Context, ids []core.ItemId) ([]*GetItemMasterRes, error) {
			return v.mockItemMaster, nil
		}
		f := CreateGetItemListService(mockGetAllStorage, mockFetchItemMaster)
		ctx := test.MockCreateContext()
		res, err := f(ctx, req)
		if !errors.Is(err, v.expectedError) {
			t.Fatalf("test case: %d, expected: %v, got: %v", i, v.expectedError, err)
		}
		if test.DeepEqual(expect, res) {
			t.Errorf("test case: %d, expected: %v, got: %v", i, expect, res)
		}
	}
}
