package core

import (
	"testing"
	"time"
)

type MockItemMaster struct {
	ItemId      ItemId
	Price       Price
	DisplayName DisplayName
	Description Description
	MaxStock    MaxStock
	CreatedAt   CreatedAt
	UpdatedAt   UpdatedAt
}

type MockItemMasterRepo struct {
	Items map[ItemId]MockItemMaster
}

var t = time.Unix(1648771200, 0)

var MockItems [2]MockItemMaster = [2]MockItemMaster{
	{
		ItemId:      "0001-ringo",
		Price:       200,
		DisplayName: "リンゴ",
		Description: "ごくふつうのリンゴ",
		MaxStock:    1000,
		CreatedAt:   CreatedAt(t),
		UpdatedAt:   UpdatedAt(t),
	},
	{
		ItemId:      "0002-burned",
		Price:       500,
		DisplayName: "焼きリンゴ",
		Description: "リンゴを加熱したもの",
		MaxStock:    100,
		CreatedAt:   CreatedAt(t),
		UpdatedAt:   UpdatedAt(t),
	},
}

func (m *MockItemMasterRepo) Get(itemId ItemId) (GetItemMasterRes, error) {
	item := m.Items[itemId]
	return GetItemMasterRes{
		ItemId:      itemId,
		Price:       item.Price,
		DisplayName: item.DisplayName,
		Description: item.Description,
		MaxStock:    item.MaxStock,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func CreateMockItemMasterRepo() *MockItemMasterRepo {
	itemMasterRepo := MockItemMasterRepo{}
	items := make(map[ItemId]MockItemMaster)
	for _, v := range MockItems {
		items[v.ItemId] = v
	}
	itemMasterRepo.Items = items
	return &itemMasterRepo
}

type MockItemStorageMaster struct {
	UserId UserId
	ItemId ItemId
	Stock  Stock
}

type MockItemStorageRepo struct {
	Data map[UserId]map[ItemId]MockItemStorageMaster
}

func (m *MockItemStorageRepo) Get(userId UserId, itemId ItemId, token AccessToken) (GetItemStorageRes, error) {
	return GetItemStorageRes{}, nil
}

func (m *MockItemStorageRepo) GetStock(userId UserId, itemId ItemId) Stock {
	return m.Data[userId][itemId].Stock
}

var MockUserId = UserId("User")

func CreateMockItemStorageRepo() *MockItemStorageRepo {
	itemStorageRepo := MockItemStorageRepo{}
	data := make(map[UserId]map[ItemId]MockItemStorageMaster)
	for i, v := range MockItems {
		if _, ok := data[MockUserId]; !ok {
			data[MockUserId] = make(map[ItemId]MockItemStorageMaster)
		}
		data[MockUserId][v.ItemId] = MockItemStorageMaster{
			UserId: MockUserId,
			ItemId: v.ItemId,
			Stock:  Stock((i + 1) * 20),
		}
	}
	itemStorageRepo.Data = data
	return &itemStorageRepo
}

func TestCreateItemService(t *testing.T) {
	itemMasterRepo := CreateMockItemMasterRepo()
	itemStorageRepo := CreateMockItemStorageRepo()
	itemService := CreateItemService(itemMasterRepo, itemStorageRepo)
	getUserItemDetail := itemService.GetUserItemDetail

	// test
	for _, v := range MockItems {
		userId := UserId("hh")
		targetItem := v
		targetId := targetItem.ItemId
		req := GetUserItemDetailReq{
			UserId: userId,
			ItemId: targetId,
		}
		res := getUserItemDetail(req)

		// check proper id
		if res.ItemId != targetId {
			t.Errorf("want %s, actual %s", targetId, res.ItemId)
		}

		// check proper master data
		targetPrice := targetItem.Price
		if res.Price != targetPrice {
			t.Errorf("want %d, actual %d", targetPrice, res.Price)
		}

		// check proper user storage data
		targetStock := itemStorageRepo.GetStock(userId, targetId)
		if res.Stock != targetStock {
			t.Errorf("want %d, actual %d", targetPrice, res.Price)
		}
	}
}
