package stage

import (
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

func (m *MockItemMasterRepo) Add(i core.ItemId, master MockItemMaster) {
	m.Items[i] = master
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
	itemMasterRepo.Items = items
	return &itemMasterRepo
}

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

var MockUserId = core.UserId("User")

func CreateMockItemStorageRepo() *MockItemStorageRepo {
	itemStorageRepo := MockItemStorageRepo{}
	data := make(map[core.UserId]map[core.ItemId]MockItemStorageMaster)
	itemStorageRepo.Data = data
	return &itemStorageRepo
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

func (m *MockUserExploreRepo) Add(userId core.UserId, exploreId ExploreId, exploreData ExploreUserData) {
	if _, ok := m.Data[userId]; !ok {
		m.Data[userId] = make(map[ExploreId]ExploreUserData)
	}
	m.Data[userId][exploreId] = exploreData
}

func createMockUserExploreRepo() *MockUserExploreRepo {
	repo := MockUserExploreRepo{Data: map[core.UserId]map[ExploreId]ExploreUserData{}}
	return &repo
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

func (m *MockExploreMasterRepo) Add(e ExploreId, master GetExploreMasterRes) {
	m.ExploreMap[e] = master
}

func (m *MockExploreMasterRepo) AddItem(itemId core.ItemId, e ExploreId, master GetExploreMasterRes) {
	m.Data[itemId] = append(m.Data[itemId], master)
	m.Add(e, master)
}

func (m *MockExploreMasterRepo) AddStage(stageId StageId, e ExploreId, master GetExploreMasterRes) {
	m.StageData[stageId] = append(m.StageData[stageId], master)
	m.Add(e, master)
}

func createMockExploreMasterRepo() *MockExploreMasterRepo {
	return &MockExploreMasterRepo{
		Data:       map[core.ItemId][]GetExploreMasterRes{},
		StageData:  map[StageId][]GetExploreMasterRes{},
		ExploreMap: map[ExploreId]GetExploreMasterRes{},
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

type mockUserStageRepo struct {
	Data map[core.UserId]map[StageId]UserStage
}

func (m *mockUserStageRepo) Add(userId core.UserId, stageId StageId, userData UserStage) {
	if _, ok := m.Data[userId]; !ok {
		m.Data[userId] = map[StageId]UserStage{}
	}
	m.Data[userId][stageId] = userData
}

func (m *mockUserStageRepo) GetAllUserStages(userId core.UserId, ids []StageId) (GetAllUserStagesRes, error) {
	result := []UserStage{}
	for _, v := range ids {
		result = append(result, m.Data[userId][v])
	}
	return GetAllUserStagesRes{result}, nil
}

func createMockUserStageRepo() *mockUserStageRepo {
	repo := mockUserStageRepo{Data: make(map[core.UserId]map[StageId]UserStage)}
	return &repo
}

type mockStageMasterRepo struct {
	Data map[StageId]StageMaster
}

func (m *mockStageMasterRepo) Add(id StageId, master StageMaster) {
	m.Data[id] = master
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

func createMockStageMasterRepo() *mockStageMasterRepo {
	repo := mockStageMasterRepo{Data: map[StageId]StageMaster{}}
	return &repo
}

type MockSkillGrowthDataRepo struct {
	Data map[ExploreId][]SkillGrowthData
}

func (m *MockSkillGrowthDataRepo) BatchGet(exploreId ExploreId) []SkillGrowthData {
	return m.Data[exploreId]
}

func (m *MockSkillGrowthDataRepo) Add(e ExploreId, skills []SkillGrowthData) {
	m.Data[e] = skills
}

func createMockSkillGrowthDataRepo() *MockSkillGrowthDataRepo {
	repo := MockSkillGrowthDataRepo{Data: map[ExploreId][]SkillGrowthData{}}
	return &repo
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

func createMockReductionStaminaSkillRepo() *mockReductionStaminaSkillRepo {
	return &mockReductionStaminaSkillRepo{Data: map[ExploreId][]core.SkillId{}}
}

func (m *mockReductionStaminaSkillRepo) Add(e ExploreId, skills []core.SkillId) {
	m.Data[e] = skills
}
