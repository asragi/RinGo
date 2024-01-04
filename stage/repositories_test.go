package stage

import (
	"github.com/asragi/RinGo/core"
)

var (
	itemStorageRepo         = CreateMockItemStorageRepo()
	exploreMasterRepo       = createMockExploreMasterRepo()
	itemExploreRelationRepo = createMockItemExploreRelationRepo()
	skillMasterRepo         = createMockSkillMasterRepo()
	userSkillRepo           = createMockUserSkillRepo()
	earningItemRepo         = createMockEarningItemRepo()
	consumingItemRepo       = createMockConsumingItemRepo()
	requiredSkillRepo       = createMockRequiredSkillRepo()
	reductionSkillRepo      = createMockReductionStaminaSkillRepo()
)

type MockItemStorageMaster struct {
	UserId  core.UserId
	ItemId  core.ItemId
	Stock   core.Stock
	IsKnown core.IsKnown
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
			UserId:  userId,
			ItemId:  v,
			Stock:   m.Data[userId][v].Stock,
			IsKnown: m.Data[userId][v].IsKnown,
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

func (m *MockItemStorageRepo) Add(userId core.UserId, items []MockItemStorageMaster) {
	itemMap := func() map[core.ItemId]MockItemStorageMaster {
		result := make(map[core.ItemId]MockItemStorageMaster)
		for _, v := range items {
			result[v.ItemId] = v
		}
		return result
	}()
	m.Data[userId] = itemMap
}

// Deprecated: ID should be prepared for each test.
var MockUserId = core.UserId("User")

func CreateMockItemStorageRepo() *MockItemStorageRepo {
	itemStorageRepo := MockItemStorageRepo{}
	data := make(map[core.UserId]map[core.ItemId]MockItemStorageMaster)
	itemStorageRepo.Data = data
	return &itemStorageRepo
}

type MockExploreMasterRepo struct {
	Data map[ExploreId]GetExploreMasterRes
}

func (m *MockExploreMasterRepo) BatchGet(e []ExploreId) ([]GetExploreMasterRes, error) {
	result := make([]GetExploreMasterRes, len(e))
	for i, v := range e {
		result[i] = m.Data[v]
	}
	return result, nil
}

func (m *MockExploreMasterRepo) Get(e ExploreId) (GetExploreMasterRes, error) {
	return m.Data[e], nil
}

func (m *MockExploreMasterRepo) Add(e ExploreId, master GetExploreMasterRes) {
	m.Data[e] = master
}

type MockItemExploreRelationRepo struct {
	Data map[core.ItemId][]ExploreId
}

func (m *MockItemExploreRelationRepo) Get(id core.ItemId) ([]ExploreId, error) {
	return m.Data[id], nil
}

func (m *MockItemExploreRelationRepo) AddItem(itemId core.ItemId, e []ExploreId) {
	m.Data[itemId] = e
}

func createMockItemExploreRelationRepo() *MockItemExploreRelationRepo {
	return &MockItemExploreRelationRepo{Data: map[core.ItemId][]ExploreId{}}
}

func createMockExploreMasterRepo() *MockExploreMasterRepo {
	return &MockExploreMasterRepo{
		Data: map[ExploreId]GetExploreMasterRes{},
	}
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

func (m *MockSkillMasterRepo) Add(id core.SkillId, master SkillMaster) {
	m.Skills[id] = master
}

func createMockSkillMasterRepo() *MockSkillMasterRepo {
	return &MockSkillMasterRepo{Skills: map[core.SkillId]SkillMaster{}}
}

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

func (m *MockUserSkillRepo) Add(userId core.UserId, skills []UserSkillRes) {
	if _, ok := m.Data[userId]; !ok {
		m.Data[userId] = map[core.SkillId]UserSkillRes{}
	}
	for _, v := range skills {
		m.Data[userId][v.SkillId] = v
	}
}

func createMockUserSkillRepo() *MockUserSkillRepo {
	return &MockUserSkillRepo{Data: map[core.UserId]map[core.SkillId]UserSkillRes{}}
}

type mockEarningItemRepo struct {
	Data map[ExploreId][]EarningItem
}

func (m *mockEarningItemRepo) BatchGet(exploreId ExploreId) []EarningItem {
	return m.Data[exploreId]
}

func (m *mockEarningItemRepo) Add(e ExploreId, items []EarningItem) {
	m.Data[e] = items
}

func createMockEarningItemRepo() *mockEarningItemRepo {
	return &mockEarningItemRepo{Data: map[ExploreId][]EarningItem{}}
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

func (m *mockConsumingItemRepo) Add(exploreId ExploreId, consuming []ConsumingItem) {
	m.Data[exploreId] = consuming
}

func createMockConsumingItemRepo() *mockConsumingItemRepo {
	return &mockConsumingItemRepo{Data: map[ExploreId][]ConsumingItem{}}
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

func (m *mockRequiredSkillRepo) Add(e ExploreId, skills []RequiredSkill) {
	m.Data[e] = skills
}

func createMockRequiredSkillRepo() *mockRequiredSkillRepo {
	return &mockRequiredSkillRepo{Data: map[ExploreId][]RequiredSkill{}}
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

type mockReductionStaminaSkillRepo struct {
	Data map[ExploreId][]core.SkillId
}

func (m *mockReductionStaminaSkillRepo) Get(exploreId ExploreId) ([]core.SkillId, error) {
	return m.Data[exploreId], nil
}

func (m *mockReductionStaminaSkillRepo) BatchGet(exploreIds []ExploreId) ([]BatchGetReductionStaminaSkill, error) {
	result := make([]BatchGetReductionStaminaSkill, len(exploreIds))
	for i, v := range exploreIds {
		result[i] = BatchGetReductionStaminaSkill{
			Skills:    m.Data[v],
			ExploreId: v,
		}
	}
	return result, nil
}

func createMockReductionStaminaSkillRepo() *mockReductionStaminaSkillRepo {
	return &mockReductionStaminaSkillRepo{Data: map[ExploreId][]core.SkillId{}}
}

func (m *mockReductionStaminaSkillRepo) Add(e ExploreId, skills []core.SkillId) {
	m.Data[e] = skills
}
