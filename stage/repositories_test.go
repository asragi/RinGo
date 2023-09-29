package stage

import (
	"time"

	"github.com/asragi/RinGo/core"
)

var (
	itemMasterRepo        = CreateMockItemMasterRepo()
	itemStorageRepo       = CreateMockItemStorageRepo()
	itemStorageUpdateRepo = createMockItemStorageUpdateRepo()
	userExploreRepo       = createMockUserExploreRepo()
	exploreMasterRepo     = createMockExploreMasterRepo()
	skillMasterRepo       = createMockSkillMasterRepo()
	userSkillRepo         = createMockUserSkillRepo()
	skillGrowthUpdateRepo = createMockSkillUpdateRepo()
	userStageRepo         = createMockUserStageRepo()
	stageMasterRepo       = createMockStageMasterRepo()
	skillGrowthDataRepo   = createMockSkillGrowthDataRepo()
	earningItemRepo       = createMockEarningItemRepo()
	consumingItemRepo     = createMockConsumingItemRepo()
	requiredSkillRepo     = createMockRequiredSkillRepo()
	reductionSkillRepo    = createMockReductionStaminaSkillRepo()
)

type MockItemMaster struct {
	ItemId      core.ItemId
	Price       core.Price
	DisplayName core.DisplayName
	Description core.Description
	MaxStock    core.MaxStock
	CreatedAt   core.CreatedAt
	UpdatedAt   core.UpdatedAt
	Explores    []ExploreId
}

type MockItemMasterRepo struct {
	Items map[core.ItemId]MockItemMaster
}

var MockItemIds []core.ItemId = []core.ItemId{
	"0000-ringo", "0001-burned", "0002-strick",
}

var t = time.Unix(1648771200, 0)

var MockItems []MockItemMaster = []MockItemMaster{
	{
		ItemId:      MockItemIds[0],
		Price:       200,
		DisplayName: "リンゴ",
		Description: "ごくふつうのリンゴ",
		MaxStock:    1000,
		CreatedAt:   core.CreatedAt(t),
		UpdatedAt:   core.UpdatedAt(t),
		Explores:    []ExploreId{mockExploreIds[0], mockExploreIds[1]},
	},
	{
		ItemId:      MockItemIds[1],
		Price:       500,
		DisplayName: "焼きリンゴ",
		Description: "リンゴを加熱したもの",
		MaxStock:    100,
		CreatedAt:   core.CreatedAt(t),
		UpdatedAt:   core.UpdatedAt(t),
		Explores:    []ExploreId{},
	},
	{
		ItemId:      MockItemIds[2],
		Price:       50,
		DisplayName: "木の枝",
		Description: "よく乾いた手頃なサイズの木の枝",
		MaxStock:    500,
		CreatedAt:   core.CreatedAt(t),
		UpdatedAt:   core.UpdatedAt(t),
		Explores:    []ExploreId{},
	},
}

func (m *MockItemMasterRepo) Get(itemId core.ItemId) (GetItemMasterRes, error) {
	item := m.Items[itemId]
	return GetItemMasterRes{
		ItemId:      itemId,
		Price:       item.Price,
		DisplayName: item.DisplayName,
		Description: item.Description,
		MaxStock:    item.MaxStock,
	}, nil
}

func (m *MockItemMasterRepo) BatchGet(ids []core.ItemId) ([]GetItemMasterRes, error) {
	result := make([]GetItemMasterRes, len(ids))
	for i, v := range ids {
		result[i], _ = m.Get(v)
	}
	return result, nil
}

func CreateMockItemMasterRepo() *MockItemMasterRepo {
	itemMasterRepo := MockItemMasterRepo{}
	items := make(map[core.ItemId]MockItemMaster)
	for _, v := range MockItems {
		items[v.ItemId] = v
	}
	itemMasterRepo.Items = items
	return &itemMasterRepo
}

type MockItemStorageMaster struct {
	UserId core.UserId
	ItemId core.ItemId
	Stock  core.Stock
}

type MockItemStorageRepo struct {
	Data map[core.UserId]map[core.ItemId]MockItemStorageMaster
}

func (m *MockItemStorageRepo) Get(userId core.UserId, itemId core.ItemId, token core.AccessToken) (GetItemStorageRes, error) {
	return GetItemStorageRes{UserId: userId, Stock: m.GetStock(userId, itemId)}, nil
}

func (m *MockItemStorageRepo) BatchGet(userId core.UserId, itemId []core.ItemId, token core.AccessToken) (BatchGetStorageRes, error) {
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

func (m *MockItemStorageRepo) GetStock(userId core.UserId, itemId core.ItemId) core.Stock {
	return m.Data[userId][itemId].Stock
}

var MockUserId = core.UserId("User")

func CreateMockItemStorageRepo() *MockItemStorageRepo {
	itemStorageRepo := MockItemStorageRepo{}
	data := make(map[core.UserId]map[core.ItemId]MockItemStorageMaster)
	for i, v := range MockItems {
		if _, ok := data[MockUserId]; !ok {
			data[MockUserId] = make(map[core.ItemId]MockItemStorageMaster)
		}
		data[MockUserId][v.ItemId] = MockItemStorageMaster{
			UserId: MockUserId,
			ItemId: v.ItemId,
			Stock:  core.Stock((i + 1) * 20),
		}
	}
	itemStorageRepo.Data = data
	return &itemStorageRepo
}

var MockExplores map[core.UserId]ExploreUserData = map[core.UserId]ExploreUserData{
	MockUserId: {
		ExploreId: mockExploreIds[0],
		IsKnown:   true,
	},
}

type MockUserExploreRepo struct {
	Data map[core.UserId]map[ExploreId]ExploreUserData
}

func (m *MockUserExploreRepo) GetActions(userId core.UserId, exploreIds []ExploreId, token core.AccessToken) (GetActionsRes, error) {
	result := make([]ExploreUserData, len(exploreIds))
	for i, v := range exploreIds {
		d := m.Data[userId][v]
		result[i] = d
	}
	return GetActionsRes{Explores: result, UserId: userId}, nil
}

var mockUserExploreData = map[core.UserId]map[ExploreId]ExploreUserData{
	MockUserId: {
		MockItems[0].Explores[0]: ExploreUserData{
			ExploreId: MockItems[0].Explores[0],
			IsKnown:   true,
		},
		MockItems[0].Explores[1]: ExploreUserData{
			ExploreId: mockExploreIds[1],
			IsKnown:   false,
		},
		mockStageExploreIds[0]: {
			ExploreId: mockStageExploreIds[0],
			IsKnown:   true,
		},
		mockStageExploreIds[1]: {
			ExploreId: mockStageExploreIds[1],
			IsKnown:   true,
		},
	},
}

func createMockUserExploreRepo() *MockUserExploreRepo {
	repo := MockUserExploreRepo{}
	repo.Data = mockUserExploreData
	return &repo
}

var mockExploreIds = []ExploreId{
	ExploreId("burn-apple"),
	ExploreId("make-sword"),
}

var mockExploreMaster = map[core.ItemId][]GetExploreMasterRes{
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

var mockStageExploreIds = []ExploreId{
	ExploreId("pick-up-apple"),
	ExploreId("alchemize-apple"),
}

var mockStageIds = []StageId{
	StageId("forest"),
	StageId("volcano"),
}

var mockStageExploreMaster = map[StageId][]GetExploreMasterRes{
	mockStageIds[0]: {
		{
			ExploreId:            mockStageExploreIds[0],
			DisplayName:          "りんごを拾いに行く",
			Description:          "木くずや石も拾えるかも",
			ConsumingStamina:     120,
			RequiredPayment:      0,
			StaminaReducibleRate: 0.5,
		},
		{
			ExploreId:            mockStageExploreIds[1],
			DisplayName:          "錬金術でりんごを金に変える",
			Description:          "黄金の精神を持ってりんごを金に変えます",
			ConsumingStamina:     720,
			RequiredPayment:      1000000,
			StaminaReducibleRate: 0.5,
		},
	},
}

type MockExploreMasterRepo struct {
	Data       map[core.ItemId][]GetExploreMasterRes
	StageData  map[StageId][]GetExploreMasterRes
	ExploreMap map[ExploreId]GetExploreMasterRes
}

func (m *MockExploreMasterRepo) GetAllExploreMaster(itemId core.ItemId) ([]GetExploreMasterRes, error) {
	return m.Data[itemId], nil
}

func (m *MockExploreMasterRepo) GetStageAllExploreMaster(stageIdArr []StageId) (BatchGetStageExploreRes, error) {
	result := []StageExploreMasterRes{}
	for _, v := range stageIdArr {
		exploreMasters := m.StageData[v]
		info := StageExploreMasterRes{
			StageId:  v,
			Explores: exploreMasters,
		}
		result = append(result, info)
	}
	return BatchGetStageExploreRes{result}, nil
}

func (m *MockExploreMasterRepo) Get(e ExploreId) (GetExploreMasterRes, error) {
	return m.ExploreMap[e], nil
}

func createMockExploreMasterRepo() *MockExploreMasterRepo {
	repo := MockExploreMasterRepo{}
	repo.ExploreMap = make(map[ExploreId]GetExploreMasterRes)
	repo.Data = mockExploreMaster
	repo.StageData = mockStageExploreMaster
	for _, v := range repo.Data {
		for _, w := range v {
			repo.ExploreMap[w.ExploreId] = w
		}
	}
	for _, v := range repo.StageData {
		for _, w := range v {
			repo.ExploreMap[w.ExploreId] = w
		}
	}
	return &repo
}

var mockSkillIds = []core.SkillId{
	"apple", "fire",
}

var MockSkillMaster = []SkillMaster{
	{
		SkillId:     mockSkillIds[0],
		DisplayName: "りんご愛好家",
	},
	{
		SkillId:     mockSkillIds[1],
		DisplayName: "火の祝福",
	},
}

type MockSkillMasterRepo struct {
	Skills map[core.SkillId]SkillMaster
}

func (m *MockSkillMasterRepo) BatchGet(skills []core.SkillId) (BatchGetSkillMasterRes, error) {
	result := make([]SkillMaster, len(skills))
	for i, id := range skills {
		result[i] = m.Skills[id]
	}
	return BatchGetSkillMasterRes{
		Skills: result,
	}, nil
}

func createMockSkillMasterRepo() *MockSkillMasterRepo {
	skills := map[core.SkillId]SkillMaster{}
	for _, v := range MockSkillMaster {
		skills[v.SkillId] = v
	}
	repo := MockSkillMasterRepo{Skills: skills}

	return &repo
}

var MockUserSkill = func() map[core.UserId]map[core.SkillId]UserSkillRes {
	result := make(map[core.UserId]map[core.SkillId]UserSkillRes)
	result[MockUserId] = map[core.SkillId]UserSkillRes{
		MockSkillMaster[0].SkillId: {
			UserId:   MockUserId,
			SkillId:  MockSkillMaster[0].SkillId,
			SkillExp: 35,
		},
		MockSkillMaster[1].SkillId: {
			UserId:   MockUserId,
			SkillId:  MockSkillMaster[0].SkillId,
			SkillExp: 0,
		},
	}
	return result
}()

type MockUserSkillRepo struct {
	Data map[core.UserId]map[core.SkillId]UserSkillRes
}

func (m *MockUserSkillRepo) BatchGet(userId core.UserId, skillIds []core.SkillId, token core.AccessToken) (BatchGetUserSkillRes, error) {
	list := m.Data[userId]
	result := make([]UserSkillRes, len(skillIds))
	for i, v := range skillIds {
		result[i] = UserSkillRes{
			UserId:   userId,
			SkillId:  v,
			SkillExp: list[v].SkillExp,
		}
	}
	return BatchGetUserSkillRes{
		UserId: userId,
		Skills: result,
	}, nil
}

func createMockUserSkillRepo() *MockUserSkillRepo {
	repo := MockUserSkillRepo{Data: MockUserSkill}
	return &repo
}

type mockUserStageRepo struct {
	Data map[core.UserId]map[StageId]UserStage
}

func (m *mockUserStageRepo) GetAllUserStages(userId core.UserId, ids []StageId) (GetAllUserStagesRes, error) {
	result := []UserStage{}
	for _, v := range ids {
		result = append(result, m.Data[userId][v])
	}
	return GetAllUserStagesRes{result}, nil
}

var mockUserStageData = map[core.UserId]map[StageId]UserStage{
	MockUserId: {
		mockStageIds[0]: {
			StageId: mockStageIds[0],
			IsKnown: true,
		},
		mockStageIds[1]: {
			StageId: mockStageIds[1],
			IsKnown: false,
		},
	},
}

func createMockUserStageRepo() *mockUserStageRepo {
	repo := mockUserStageRepo{}
	repo.Data = mockUserStageData
	return &repo
}

type mockStageMasterRepo struct {
	Data map[StageId]StageMaster
}

func (m *mockStageMasterRepo) Get(stageId StageId) (StageMaster, error) {
	return m.Data[stageId], nil
}

func (m *mockStageMasterRepo) GetAllStages() (GetAllStagesRes, error) {
	result := []StageMaster{}
	for _, v := range m.Data {
		result = append(result, v)
	}
	return GetAllStagesRes{Stages: result}, nil
}

var mockStageMasterData = map[StageId]StageMaster{
	mockStageIds[0]: {
		StageId:     mockStageIds[0],
		DisplayName: "ポムポムのもり",
		Description: "りんごの木がたくさんある森\nいつでもたくさんのりんごが採れる",
	},
	mockStageIds[1]: {
		StageId:     mockStageIds[1],
		DisplayName: "リゴーかざん",
		Description: "鉱石が採れるかも",
	},
}

func createMockStageMasterRepo() *mockStageMasterRepo {
	repo := mockStageMasterRepo{}
	repo.Data = mockStageMasterData
	return &repo
}

type MockSkillGrowthDataRepo struct {
	Data map[ExploreId][]SkillGrowthData
}

func (m *MockSkillGrowthDataRepo) BatchGet(exploreId ExploreId) []SkillGrowthData {
	return m.Data[exploreId]
}

var MockSkillGrowthData = map[ExploreId][]SkillGrowthData{
	mockExploreIds[0]: {
		{
			ExploreId:    mockExploreIds[0],
			SkillId:      mockSkillIds[0],
			GainingPoint: 10,
		},
		{
			ExploreId:    mockExploreIds[0],
			SkillId:      mockSkillIds[1],
			GainingPoint: 10,
		},
	},
	mockStageExploreIds[0]: {
		{
			ExploreId:    mockExploreIds[0],
			SkillId:      mockSkillIds[0],
			GainingPoint: 10,
		},
		{
			ExploreId:    mockExploreIds[0],
			SkillId:      mockSkillIds[1],
			GainingPoint: 10,
		},
	},
}

func createMockSkillGrowthDataRepo() *MockSkillGrowthDataRepo {
	repo := MockSkillGrowthDataRepo{Data: MockSkillGrowthData}
	return &repo
}

type mockEarningItemRepo struct {
	Data map[ExploreId][]EarningItem
}

func (m *mockEarningItemRepo) BatchGet(exploreId ExploreId) []EarningItem {
	return m.Data[exploreId]
}

var mockEarningItemData = map[ExploreId][]EarningItem{
	mockExploreIds[0]: {
		{
			ItemId:   MockItemIds[0],
			MinCount: 1,
			MaxCount: 100,
		},
		{
			ItemId:   MockItemIds[1],
			MinCount: 0,
			MaxCount: 1000,
		},
	},
	mockExploreIds[1]: {
		{
			ItemId:   MockItemIds[0],
			MinCount: 1,
			MaxCount: 10,
		},
		{
			ItemId:   MockItemIds[2],
			MinCount: 10,
			MaxCount: 100,
		},
	},
	mockStageExploreIds[0]: {
		{
			ItemId:   MockItemIds[0],
			MinCount: 10,
			MaxCount: 60,
		},
	},
}

func createMockEarningItemRepo() *mockEarningItemRepo {
	return &mockEarningItemRepo{Data: mockEarningItemData}
}

type mockConsumingItemRepo struct {
	Data map[ExploreId][]ConsumingItem
}

func (m *mockConsumingItemRepo) BatchGet(exploreId ExploreId) ([]ConsumingItem, error) {
	return m.Data[exploreId], nil
}

func (m *mockConsumingItemRepo) AllGet(exploreId []ExploreId) ([]BatchGetConsumingItemRes, error) {
	result := make([]BatchGetConsumingItemRes, len(exploreId))
	for i, v := range exploreId {
		result[i] = BatchGetConsumingItemRes{
			ExploreId:      v,
			ConsumingItems: m.Data[v],
		}
	}
	return result, nil
}

var mockConsumingItemData = map[ExploreId][]ConsumingItem{
	mockExploreIds[0]: {
		{
			ItemId:          MockItemIds[0],
			ConsumptionProb: 1,
			MaxCount:        10,
		},
		{
			ItemId:          MockItemIds[1],
			MaxCount:        15,
			ConsumptionProb: 0.5,
		},
	},
	mockExploreIds[1]: {
		{
			ItemId:          MockItemIds[0],
			ConsumptionProb: 1,
			MaxCount:        10,
		},
		{
			ItemId:          MockItemIds[2],
			ConsumptionProb: 0,
			MaxCount:        100,
		},
	},
	mockStageExploreIds[0]: {
		{
			ItemId:          MockItemIds[2],
			ConsumptionProb: 1,
			MaxCount:        1,
		},
	},
	mockStageExploreIds[1]: {
		{
			ItemId:          MockItemIds[0],
			ConsumptionProb: 1,
			MaxCount:        1000,
		},
	},
}

func createMockConsumingItemRepo() *mockConsumingItemRepo {
	return &mockConsumingItemRepo{Data: mockConsumingItemData}
}

type mockRequiredSkillRepo struct {
	Data map[ExploreId][]RequiredSkill
}

func (m *mockRequiredSkillRepo) Get(exploreId ExploreId) ([]RequiredSkill, error) {
	if _, ok := m.Data[exploreId]; !ok {
		return []RequiredSkill{}, nil
	}
	return m.Data[exploreId], nil
}

func (m *mockRequiredSkillRepo) BatchGet(ids []ExploreId) ([]RequiredSkillRow, error) {
	result := make([]RequiredSkillRow, len(ids))
	for i, v := range ids {
		row := RequiredSkillRow{
			ExploreId:      v,
			RequiredSkills: m.Data[v],
		}
		result[i] = row
	}
	return result, nil
}

func createMockRequiredSkillRepo() *mockRequiredSkillRepo {
	return &mockRequiredSkillRepo{}
}

type mockItemStorageUpdateRepo struct {
	Data map[core.UserId][]ItemStock
}

func (m *mockItemStorageUpdateRepo) Update(userId core.UserId, items []ItemStock, _ core.AccessToken) error {
	m.Data = make(map[core.UserId][]ItemStock)
	m.Data[userId] = items
	return nil
}

func (m *mockItemStorageUpdateRepo) Get(userId core.UserId) []ItemStock {
	return m.Data[userId]
}

func createMockItemStorageUpdateRepo() *mockItemStorageUpdateRepo {
	return &mockItemStorageUpdateRepo{}
}

type mockSkillUpdateRepo struct {
	Data map[core.UserId][]SkillGrowthPostRow
}

func (m *mockSkillUpdateRepo) Update(req SkillGrowthPost) error {
	m.Data = make(map[core.UserId][]SkillGrowthPostRow)
	m.Data[req.UserId] = req.SkillGrowth
	return nil
}

func (m *mockSkillUpdateRepo) Get(userId core.UserId) []SkillGrowthPostRow {
	return m.Data[userId]
}

func createMockSkillUpdateRepo() *mockSkillUpdateRepo {
	return &mockSkillUpdateRepo{}
}

type mockReductionStaminaSkillRepo struct {
	Data map[ExploreId][]core.SkillId
}

func (m *mockReductionStaminaSkillRepo) Get(exploreId ExploreId) ([]core.SkillId, error) {
	return m.Data[exploreId], nil
}

var mockStaminaReductionSkill = map[ExploreId][]core.SkillId{
	mockStageExploreIds[1]: {mockSkillIds[0]},
}

func createMockReductionStaminaSkillRepo() *mockReductionStaminaSkillRepo {
	return &mockReductionStaminaSkillRepo{Data: mockStaminaReductionSkill}
}
