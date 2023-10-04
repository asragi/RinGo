package infrastructure

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

type itemMasterData map[core.ItemId]stage.GetItemMasterRes

type IItemMasterLoader interface {
	Load() (itemMasterData, error)
}

type InMemoryItemMasterRepo struct {
	Items itemMasterData
}

func (m *InMemoryItemMasterRepo) Get(itemId core.ItemId) (stage.GetItemMasterRes, error) {
	item := m.Items[itemId]
	return stage.GetItemMasterRes{
		ItemId:      itemId,
		Price:       item.Price,
		DisplayName: item.DisplayName,
		Description: item.Description,
		MaxStock:    item.MaxStock,
	}, nil
}

func (m *InMemoryItemMasterRepo) BatchGet(ids []core.ItemId) ([]stage.GetItemMasterRes, error) {
	result := make([]stage.GetItemMasterRes, len(ids))
	for i, v := range ids {
		result[i], _ = m.Get(v)
	}
	return result, nil
}

func CreateInMemoryItemMasterRepo(loader IItemMasterLoader) (*InMemoryItemMasterRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, err
	}
	return &InMemoryItemMasterRepo{Items: data}, err
}

type ItemStorageData map[core.UserId]map[core.ItemId]stage.ItemData

type InMemoryItemStorageRepo struct {
	Data ItemStorageData
}

type ItemStorageDataLoader interface {
	Load() ItemStorageData
}

func (m *InMemoryItemStorageRepo) Get(userId core.UserId, itemId core.ItemId, token core.AccessToken) (stage.GetItemStorageRes, error) {
	return stage.GetItemStorageRes{UserId: userId, Stock: m.GetStock(userId, itemId)}, nil
}

func (m *InMemoryItemStorageRepo) BatchGet(userId core.UserId, itemId []core.ItemId, token core.AccessToken) (stage.BatchGetStorageRes, error) {
	result := make([]stage.ItemData, len(itemId))
	for i, v := range itemId {
		itemData := stage.ItemData{
			UserId:  userId,
			ItemId:  v,
			Stock:   m.Data[userId][v].Stock,
			IsKnown: m.Data[userId][v].IsKnown,
		}
		result[i] = itemData
	}
	res := stage.BatchGetStorageRes{
		UserId:   userId,
		ItemData: result,
	}
	return res, nil
}

func (m *InMemoryItemStorageRepo) GetStock(userId core.UserId, itemId core.ItemId) core.Stock {
	return m.Data[userId][itemId].Stock
}

func CreateInMemoryItemStorageRepo(loader ItemStorageDataLoader) *InMemoryItemStorageRepo {
	itemStorageRepo := InMemoryItemStorageRepo{Data: loader.Load()}
	return &itemStorageRepo
}

type InMemoryUserExploreRepo struct {
	Data map[core.UserId]map[stage.ExploreId]stage.ExploreUserData
}

func (m *InMemoryUserExploreRepo) GetActions(userId core.UserId, exploreIds []stage.ExploreId, token core.AccessToken) (stage.GetActionsRes, error) {
	result := make([]stage.ExploreUserData, len(exploreIds))
	for i, v := range exploreIds {
		d := m.Data[userId][v]
		result[i] = d
	}
	return stage.GetActionsRes{Explores: result, UserId: userId}, nil
}

func createMockUserExploreRepo() *InMemoryUserExploreRepo {
	repo := InMemoryUserExploreRepo{Data: map[core.UserId]map[stage.ExploreId]stage.ExploreUserData{}}
	return &repo
}

type InMemoryExploreMasterRepo struct {
	Data       map[core.ItemId][]stage.GetExploreMasterRes
	StageData  map[stage.StageId][]stage.GetExploreMasterRes
	ExploreMap map[stage.ExploreId]stage.GetExploreMasterRes
}

type ExploreMasterLoader interface {
	Load() (map[core.ItemId][]stage.GetExploreMasterRes, map[stage.StageId][]stage.GetExploreMasterRes)
}

func (m *InMemoryExploreMasterRepo) GetAllExploreMaster(itemId core.ItemId) ([]stage.GetExploreMasterRes, error) {
	return m.Data[itemId], nil
}

func (m *InMemoryExploreMasterRepo) GetStageAllExploreMaster(stageIdArr []stage.StageId) (stage.BatchGetStageExploreRes, error) {
	result := []stage.StageExploreMasterRes{}
	for _, v := range stageIdArr {
		exploreMasters := m.StageData[v]
		info := stage.StageExploreMasterRes{
			StageId:  v,
			Explores: exploreMasters,
		}
		result = append(result, info)
	}
	return stage.BatchGetStageExploreRes{StageExplores: result}, nil
}

func (m *InMemoryExploreMasterRepo) Get(e stage.ExploreId) (stage.GetExploreMasterRes, error) {
	return m.ExploreMap[e], nil
}

func (m *InMemoryExploreMasterRepo) Add(e stage.ExploreId, master stage.GetExploreMasterRes) {
	m.ExploreMap[e] = master
}

func (m *InMemoryExploreMasterRepo) AddItem(itemId core.ItemId, e stage.ExploreId, master stage.GetExploreMasterRes) {
	m.Data[itemId] = append(m.Data[itemId], master)
	m.Add(e, master)
}

func (m *InMemoryExploreMasterRepo) AddStage(stageId stage.StageId, e stage.ExploreId, master stage.GetExploreMasterRes) {
	m.StageData[stageId] = append(m.StageData[stageId], master)
	m.Add(e, master)
}

func createInMemoryExploreMasterRepo() *InMemoryExploreMasterRepo {
	return &InMemoryExploreMasterRepo{
		Data:       map[core.ItemId][]stage.GetExploreMasterRes{},
		StageData:  map[stage.StageId][]stage.GetExploreMasterRes{},
		ExploreMap: map[stage.ExploreId]stage.GetExploreMasterRes{},
	}
}

type SkillMasterData map[core.SkillId]stage.SkillMaster

type InMemorySkillMasterRepo struct {
	Skills SkillMasterData
}

type SkillMasterLoader interface {
	Load() SkillMasterData
}

func (m *InMemorySkillMasterRepo) BatchGet(skills []core.SkillId) (stage.BatchGetSkillMasterRes, error) {
	result := make([]stage.SkillMaster, len(skills))
	for i, id := range skills {
		result[i] = m.Skills[id]
	}
	return stage.BatchGetSkillMasterRes{
		Skills: result,
	}, nil
}

func createInMemorySkillMasterRepo(loader SkillMasterLoader) *InMemorySkillMasterRepo {
	return &InMemorySkillMasterRepo{Skills: loader.Load()}
}

type InMemoryUserSkillRepo struct {
	Data map[core.UserId]map[core.SkillId]stage.UserSkillRes
}

func (m *InMemoryUserSkillRepo) BatchGet(userId core.UserId, skillIds []core.SkillId, token core.AccessToken) (stage.BatchGetUserSkillRes, error) {
	list := m.Data[userId]
	result := make([]stage.UserSkillRes, len(skillIds))
	for i, v := range skillIds {
		result[i] = stage.UserSkillRes{
			UserId:   userId,
			SkillId:  v,
			SkillExp: list[v].SkillExp,
		}
	}
	return stage.BatchGetUserSkillRes{
		UserId: userId,
		Skills: result,
	}, nil
}

func (m *InMemoryUserSkillRepo) Add(userId core.UserId, skills []stage.UserSkillRes) {
	if _, ok := m.Data[userId]; !ok {
		m.Data[userId] = map[core.SkillId]stage.UserSkillRes{}
	}
	for _, v := range skills {
		m.Data[userId][v.SkillId] = v
	}
}

func createMockUserSkillRepo() *InMemoryUserSkillRepo {
	return &InMemoryUserSkillRepo{Data: map[core.UserId]map[core.SkillId]stage.UserSkillRes{}}
}

type mockUserStageRepo struct {
	Data map[core.UserId]map[stage.StageId]stage.UserStage
}

func (m *mockUserStageRepo) Add(userId core.UserId, stageId stage.StageId, userData stage.UserStage) {
	if _, ok := m.Data[userId]; !ok {
		m.Data[userId] = map[stage.StageId]stage.UserStage{}
	}
	m.Data[userId][stageId] = userData
}

func (m *mockUserStageRepo) GetAllUserStages(userId core.UserId, ids []stage.StageId) (stage.GetAllUserStagesRes, error) {
	result := []stage.UserStage{}
	for _, v := range ids {
		result = append(result, m.Data[userId][v])
	}
	return stage.GetAllUserStagesRes{UserStage: result}, nil
}

func createMockUserStageRepo() *mockUserStageRepo {
	repo := mockUserStageRepo{Data: make(map[core.UserId]map[stage.StageId]stage.UserStage)}
	return &repo
}

type mockStageMasterRepo struct {
	Data map[stage.StageId]stage.StageMaster
}

func (m *mockStageMasterRepo) Add(id stage.StageId, master stage.StageMaster) {
	m.Data[id] = master
}

func (m *mockStageMasterRepo) Get(stageId stage.StageId) (stage.StageMaster, error) {
	return m.Data[stageId], nil
}

func (m *mockStageMasterRepo) GetAllStages() (stage.GetAllStagesRes, error) {
	result := []stage.StageMaster{}
	for _, v := range m.Data {
		result = append(result, v)
	}
	return stage.GetAllStagesRes{Stages: result}, nil
}

func createMockStageMasterRepo() *mockStageMasterRepo {
	repo := mockStageMasterRepo{Data: map[stage.StageId]stage.StageMaster{}}
	return &repo
}

type MockSkillGrowthDataRepo struct {
	Data map[stage.ExploreId][]stage.SkillGrowthData
}

func (m *MockSkillGrowthDataRepo) BatchGet(exploreId stage.ExploreId) []stage.SkillGrowthData {
	return m.Data[exploreId]
}

func (m *MockSkillGrowthDataRepo) Add(e stage.ExploreId, skills []stage.SkillGrowthData) {
	m.Data[e] = skills
}

func createMockSkillGrowthDataRepo() *MockSkillGrowthDataRepo {
	repo := MockSkillGrowthDataRepo{Data: map[stage.ExploreId][]stage.SkillGrowthData{}}
	return &repo
}

type mockEarningItemRepo struct {
	Data map[stage.ExploreId][]stage.EarningItem
}

func (m *mockEarningItemRepo) BatchGet(exploreId stage.ExploreId) []stage.EarningItem {
	return m.Data[exploreId]
}

func (m *mockEarningItemRepo) Add(e stage.ExploreId, items []stage.EarningItem) {
	m.Data[e] = items
}

func createMockEarningItemRepo() *mockEarningItemRepo {
	return &mockEarningItemRepo{Data: map[stage.ExploreId][]stage.EarningItem{}}
}

type InMemoryConsumingItemRepo struct {
	Data map[stage.ExploreId][]stage.ConsumingItem
}

func (m *InMemoryConsumingItemRepo) BatchGet(exploreId stage.ExploreId) ([]stage.ConsumingItem, error) {
	return m.Data[exploreId], nil
}

func (m *InMemoryConsumingItemRepo) AllGet(exploreId []stage.ExploreId) ([]stage.BatchGetConsumingItemRes, error) {
	result := make([]stage.BatchGetConsumingItemRes, len(exploreId))
	for i, v := range exploreId {
		result[i] = stage.BatchGetConsumingItemRes{
			ExploreId:      v,
			ConsumingItems: m.Data[v],
		}
	}
	return result, nil
}

func (m *InMemoryConsumingItemRepo) Add(exploreId stage.ExploreId, consuming []stage.ConsumingItem) {
	m.Data[exploreId] = consuming
}

func createMockConsumingItemRepo() *InMemoryConsumingItemRepo {
	return &InMemoryConsumingItemRepo{Data: map[stage.ExploreId][]stage.ConsumingItem{}}
}

type mockRequiredSkillRepo struct {
	Data map[stage.ExploreId][]stage.RequiredSkill
}

func (m *mockRequiredSkillRepo) Get(exploreId stage.ExploreId) ([]stage.RequiredSkill, error) {
	if _, ok := m.Data[exploreId]; !ok {
		return []stage.RequiredSkill{}, nil
	}
	return m.Data[exploreId], nil
}

func (m *mockRequiredSkillRepo) BatchGet(ids []stage.ExploreId) ([]stage.RequiredSkillRow, error) {
	result := make([]stage.RequiredSkillRow, len(ids))
	for i, v := range ids {
		row := stage.RequiredSkillRow{
			ExploreId:      v,
			RequiredSkills: m.Data[v],
		}
		result[i] = row
	}
	return result, nil
}

func (m *mockRequiredSkillRepo) Add(e stage.ExploreId, skills []stage.RequiredSkill) {
	m.Data[e] = skills
}

func createMockRequiredSkillRepo() *mockRequiredSkillRepo {
	return &mockRequiredSkillRepo{Data: map[stage.ExploreId][]stage.RequiredSkill{}}
}

type mockItemStorageUpdateRepo struct {
	Data map[core.UserId][]stage.ItemStock
}

func (m *mockItemStorageUpdateRepo) Update(userId core.UserId, items []stage.ItemStock, _ core.AccessToken) error {
	m.Data = make(map[core.UserId][]stage.ItemStock)
	m.Data[userId] = items
	return nil
}

func (m *mockItemStorageUpdateRepo) Get(userId core.UserId) []stage.ItemStock {
	return m.Data[userId]
}

func createMockItemStorageUpdateRepo() *mockItemStorageUpdateRepo {
	return &mockItemStorageUpdateRepo{}
}

type mockSkillUpdateRepo struct {
	Data map[core.UserId][]stage.SkillGrowthPostRow
}

func (m *mockSkillUpdateRepo) Update(req stage.SkillGrowthPost) error {
	m.Data = make(map[core.UserId][]stage.SkillGrowthPostRow)
	m.Data[req.UserId] = req.SkillGrowth
	return nil
}

func (m *mockSkillUpdateRepo) Get(userId core.UserId) []stage.SkillGrowthPostRow {
	return m.Data[userId]
}

func createMockSkillUpdateRepo() *mockSkillUpdateRepo {
	return &mockSkillUpdateRepo{}
}

type mockReductionStaminaSkillRepo struct {
	Data map[stage.ExploreId][]core.SkillId
}

func (m *mockReductionStaminaSkillRepo) Get(exploreId stage.ExploreId) ([]core.SkillId, error) {
	return m.Data[exploreId], nil
}

func createMockReductionStaminaSkillRepo() *mockReductionStaminaSkillRepo {
	return &mockReductionStaminaSkillRepo{Data: map[stage.ExploreId][]core.SkillId{}}
}

func (m *mockReductionStaminaSkillRepo) Add(e stage.ExploreId, skills []core.SkillId) {
	m.Data[e] = skills
}
