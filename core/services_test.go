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
		Explores:    []ExploreId{mockExploreIds[0], mockExploreIds[1]},
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
		ExploreId: mockExploreIds[0],
		IsKnown:   true,
	},
}

var MockConditions map[ExploreId][]Condition = map[ExploreId][]Condition{
	mockExploreIds[0]: {
		{
			ConditionId:          "enough-stick",
			ConditionType:        ConditionTypeItem,
			ConditionTargetId:    ConditionTargetId(MockItems[2].ItemId),
			ConditionTargetValue: ConditionTargetValue(10),
		},
	},
	mockExploreIds[1]: {
		{
			ConditionId:          "enough-apple",
			ConditionType:        ConditionTypeItem,
			ConditionTargetId:    ConditionTargetId(MockItems[1].ItemId),
			ConditionTargetValue: ConditionTargetValue(10),
		},
		{
			ConditionId:          "enough-apple-burned",
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
	Data map[UserId]map[ExploreId]ExploreUserData
}

func (m *MockUserExploreRepo) GetActions(userId UserId, exploreIds []ExploreId, token AccessToken) (GetActionsRes, error) {
	result := make([]ExploreUserData, len(exploreIds))
	for i, v := range exploreIds {
		d := m.Data[userId][v]
		result[i] = d
	}
	return GetActionsRes{Explores: result, UserId: userId}, nil
}

var mockUserExploreData = map[UserId]map[ExploreId]ExploreUserData{
	MockUserId: {
		MockItems[0].Explores[0]: ExploreUserData{
			ExploreId: MockItems[0].Explores[0],
			IsKnown:   true,
		},
		MockItems[0].Explores[1]: ExploreUserData{
			ExploreId: mockExploreIds[1],
			IsKnown:   false,
		},
	},
}

func createMockUserExploreRepo() *MockUserExploreRepo {
	repo := MockUserExploreRepo{}
	repo.Data = mockUserExploreData
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

var mockExploreIds = []ExploreId{
	ExploreId("burn-apple"),
	ExploreId("make-sword"),
}

var mockExploreMaster = map[ItemId][]GetAllExploreMasterRes{
	MockItems[0].ItemId: {
		{
			ExploreId:   mockExploreIds[0],
			DisplayName: "りんごを焼く",
			Description: "りんごを火にかけてみよう",
		},
		{
			ExploreId:   mockExploreIds[1],
			DisplayName: "りんごの家を作る",
			Description: "りんごを使って家を建てます",
		},
	},
}

type MockExploreMasterRepo struct {
	Data map[ItemId][]GetAllExploreMasterRes
}

func (m *MockExploreMasterRepo) GetAllExploreMaster(itemId ItemId) ([]GetAllExploreMasterRes, error) {
	return m.Data[itemId], nil
}

func createMockExploreMasterRepo() *MockExploreMasterRepo {
	repo := MockExploreMasterRepo{}
	repo.Data = mockExploreMaster
	return &repo
}

type testRequest struct {
	userId UserId
	itemId ItemId
}

type testExplore struct {
	exploreId  ExploreId
	name       DisplayName
	isKnown    IsKnown
	isPossible IsPossible
}

type testExpect struct {
	price    Price
	stock    Stock
	explores []testExplore
}

type testCase struct {
	request testRequest
	expect  testExpect
}

func TestCreateItemService(t *testing.T) {
	itemMasterRepo := CreateMockItemMasterRepo()
	itemStorageRepo := CreateMockItemStorageRepo()
	userExploreRepo := createMockUserExploreRepo()
	conditionRepo := createMockExploreConditionRepo()
	exploreMasterRepo := createMockExploreMasterRepo()
	itemService := CreateItemService(itemMasterRepo, itemStorageRepo, exploreMasterRepo, userExploreRepo, conditionRepo)
	getUserItemDetail := itemService.GetUserItemDetail

	testCases := []testCase{
		{
			request: testRequest{
				itemId: MockItems[0].ItemId,
				userId: MockUserId,
			},
			expect: testExpect{
				price: MockItems[0].Price,
				stock: 20,
				explores: []testExplore{
					{
						exploreId:  mockExploreIds[0],
						name:       mockExploreMaster[MockItems[0].ItemId][0].DisplayName,
						isKnown:    true,
						isPossible: true,
					},
					{
						exploreId:  mockExploreIds[1],
						name:       mockExploreMaster[MockItems[0].ItemId][1].DisplayName,
						isKnown:    false,
						isPossible: false,
					},
				},
			},
		},
	}
	// test
	check := func(expect string, actual string) {
		if expect != actual {
			t.Errorf("want %s, actual %s", expect, actual)
		}
	}
	checkBool := func(title string, expect bool, actual bool) {
		if expect != actual {
			t.Errorf("%s: want %t, actual %t", title, expect, actual)
		}
	}
	for _, v := range testCases {
		targetId := v.request.itemId
		req := GetUserItemDetailReq{
			UserId: v.request.userId,
			ItemId: targetId,
		}
		res := getUserItemDetail(req)
		// check proper id
		if res.ItemId != targetId {
			t.Errorf("want %s, actual %s", targetId, res.ItemId)
		}

		// check proper master data
		expect := v.expect
		if res.Price != expect.price {
			t.Errorf("want %d, actual %d", expect.price, res.Price)
		}

		// check proper user storage data
		targetStock := expect.stock
		if res.Stock != targetStock {
			t.Errorf("want %d, actual %d", targetStock, res.Stock)
		}

		// check explore
		if len(res.UserExplores) != len(expect.explores) {
			t.Errorf("want %d, actual %d", len(expect.explores), len(res.UserExplores))
		}
		for j, w := range expect.explores {
			actual := res.UserExplores[j]
			if w.exploreId != actual.ExploreId {
				t.Errorf("want %s, actual %s", w.exploreId, actual.ExploreId)
			}
			check(string(w.name), string(actual.DisplayName))
			checkBool("isKnown", bool(w.isKnown), bool(actual.IsKnown))
			checkBool("isPossible", bool(w.isPossible), bool(actual.IsPossible))
		}
	}
}
