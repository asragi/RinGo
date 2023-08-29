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
	Explores    []ExploreId
}

type MockItemMasterRepo struct {
	Items map[ItemId]MockItemMaster
}

var t = time.Unix(1648771200, 0)

var MockItems [3]MockItemMaster = [3]MockItemMaster{
	{
		ItemId:      "0001-ringo",
		Price:       200,
		DisplayName: "リンゴ",
		Description: "ごくふつうのリンゴ",
		MaxStock:    1000,
		CreatedAt:   CreatedAt(t),
		UpdatedAt:   UpdatedAt(t),
		Explores:    []ExploreId{"burn-apple"},
	},
	{
		ItemId:      "0002-burned",
		Price:       500,
		DisplayName: "焼きリンゴ",
		Description: "リンゴを加熱したもの",
		MaxStock:    100,
		CreatedAt:   CreatedAt(t),
		UpdatedAt:   UpdatedAt(t),
		Explores:    []ExploreId{},
	},
	{
		ItemId:      "0003-stick",
		Price:       50,
		DisplayName: "木の枝",
		Description: "よく乾いた手頃なサイズの木の枝",
		MaxStock:    500,
		CreatedAt:   CreatedAt(t),
		UpdatedAt:   UpdatedAt(t),
		Explores:    []ExploreId{},
	},
}

var MockExplores map[UserId]ExploreUserData = map[UserId]ExploreUserData{
	MockUserId: {
		ExploreId: MockItems[0].Explores[0],
		IsKnown:   true,
	},
}

var MockConditions map[ExploreId][]Condition = map[ExploreId][]Condition{
	MockItems[0].Explores[0]: {
		{
			ConditionId:          "enough-stick",
			ConditionType:        ConditionTypeItem,
			ConditionTargetId:    ConditionTargetId(MockItems[2].ItemId),
			ConditionTargetValue: ConditionTargetValue(100),
		},
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
	return GetItemStorageRes{UserId: userId, Stock: m.GetStock(userId, itemId)}, nil
}

func (m *MockItemStorageRepo) BatchGet(userId UserId, itemId []ItemId, token AccessToken) (BatchGetStorageRes, error) {
	result := make([]ItemData, len(itemId))
	for i, v := range itemId {
		itemData := ItemData{
			UserId: userId,
			ItemId: v,
			Stock:  m.Data[userId][v].Stock,
		}
		result[i] = itemData
	}
	res := BatchGetStorageRes{
		UserId:   userId,
		ItemData: result,
	}
	return res, nil
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

type MockUserExploreRepo struct {
	Data map[UserId]map[ItemId][]ExploreUserData
}

func (m *MockUserExploreRepo) GetActions(userId UserId, itemId ItemId, token AccessToken) (GetActionsRes, error) {
	return GetActionsRes{Explores: m.Data[userId][itemId], ItemId: itemId}, nil
}

func createMockUserExploreRepo() *MockUserExploreRepo {
	repo := MockUserExploreRepo{}
	data := make(map[UserId]map[ItemId][]ExploreUserData)
	data[MockUserId] = make(map[ItemId][]ExploreUserData)

	repo.Data = data
	return &repo
}

type MockExploreConditionRepo struct {
	Data map[ExploreId][]Condition
}

func (m *MockExploreConditionRepo) GetAllConditions(id []ExploreId) (GetAllConditionsRes, error) {
	result := make([]ExploreConditions, len(id))
	for i, v := range id {
		s := ExploreConditions{
			ExploreId:  v,
			Conditions: m.Data[v],
		}
		result[i] = s
	}
	return GetAllConditionsRes{Explores: result}, nil
}

func createMockExploreConditionRepo() *MockExploreConditionRepo {
	repo := MockExploreConditionRepo{}
	repo.Data = MockConditions
	return &repo
}

func TestCreateItemService(t *testing.T) {
	itemMasterRepo := CreateMockItemMasterRepo()
	itemStorageRepo := CreateMockItemStorageRepo()
	userExploreRepo := createMockUserExploreRepo()
	conditionRepo := createMockExploreConditionRepo()
	itemService := CreateItemService(itemMasterRepo, itemStorageRepo, userExploreRepo, conditionRepo)
	getUserItemDetail := itemService.GetUserItemDetail

	// test
	for _, v := range MockItems {
		targetItem := v
		targetId := targetItem.ItemId
		req := GetUserItemDetailReq{
			UserId: MockUserId,
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
		targetStock := itemStorageRepo.GetStock(MockUserId, targetId)
		if res.Stock != targetStock {
			t.Errorf("want %d, actual %d", targetStock, res.Stock)
		}

		// check improper user storage data
		req = GetUserItemDetailReq{
			UserId: UserId("ImproperUserNameTest"),
			ItemId: targetId,
		}
		res = getUserItemDetail(req)
		if res.Stock == targetStock {
			t.Errorf("don't want to be %d, actual %d", targetPrice, res.Price)
		}
	}
}
