package infrastructure

import (
	"fmt"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

type InMemoryUserRepo struct {
	StaminaTime core.StaminaRecoverTime
	Fund        core.Fund
}

func (m *InMemoryUserRepo) GetResource(userId core.UserId, token core.AccessToken) (stage.GetResourceRes, error) {
	return stage.GetResourceRes{
		UserId:             userId,
		MaxStamina:         6000,
		StaminaRecoverTime: m.StaminaTime,
		Fund:               m.Fund,
	}, nil
}

func (m *InMemoryUserRepo) UpdateStamina(userId core.UserId, token core.AccessToken, stamina core.StaminaRecoverTime) error {
	m.StaminaTime = stamina
	return nil
}

func CreateInMemoryUserResourceRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{}
}

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

type IItemStorageDataLoader interface {
	Load() (ItemStorageData, error)
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

func (m *InMemoryItemStorageRepo) Update(userId core.UserId, items []stage.ItemStock, token core.AccessToken) error {
	for _, v := range items {
		if _, ok := m.Data[userId]; !ok {
			m.Data[userId] = map[core.ItemId]stage.ItemData{}
		}
		m.Data[userId][v.ItemId] = stage.ItemData{
			UserId:  userId,
			ItemId:  v.ItemId,
			Stock:   v.AfterStock,
			IsKnown: true,
		}
	}
	return nil
}

func CreateInMemoryItemStorageRepo(loader IItemStorageDataLoader) (*InMemoryItemStorageRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on load item storage: %w", err)
	}
	itemStorageRepo := InMemoryItemStorageRepo{Data: data}
	return &itemStorageRepo, nil
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

type ExploreMasterData map[stage.ExploreId]stage.GetExploreMasterRes

type InMemoryExploreMasterRepo struct {
	Data ExploreMasterData
}

type IExploreMasterLoader interface {
	Load() (ExploreMasterData, error)
}

type InMemoryItemExploreRepo struct {
	Data map[core.ItemId][]stage.ExploreId
}

type InMemoryStageExploreRepo struct {
	StageData map[stage.StageId][]stage.ExploreId
}

func (m *InMemoryExploreMasterRepo) Get(e stage.ExploreId) (stage.GetExploreMasterRes, error) {
	return m.Data[e], nil
}

func (m *InMemoryExploreMasterRepo) BatchGet(e []stage.ExploreId) ([]stage.GetExploreMasterRes, error) {
	result := make([]stage.GetExploreMasterRes, len(e))
	fmt.Printf("m is nil? :%t", m == nil)
	for i, v := range e {
		result[i] = m.Data[v]
	}
	return result, nil
}

func CreateInMemoryExploreMasterRepo(loader IExploreMasterLoader) (*InMemoryExploreMasterRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on load in memory explore master: %w", err)
	}
	return &InMemoryExploreMasterRepo{
		Data: data,
	}, nil
}

type SkillMasterData map[core.SkillId]stage.SkillMaster

type InMemorySkillMasterRepo struct {
	Skills SkillMasterData
}

type ISkillMasterLoader interface {
	Load() (SkillMasterData, error)
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

func CreateInMemorySkillMasterRepo(loader ISkillMasterLoader) (*InMemorySkillMasterRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on create in memory skill master repo: %w", err)
	}
	return &InMemorySkillMasterRepo{Skills: data}, nil
}

type UserSkillData map[core.UserId]map[core.SkillId]stage.UserSkillRes

type InMemoryUserSkillRepo struct {
	Data UserSkillData
}

type IUserSkillLoader interface {
	Load() (UserSkillData, error)
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

func (m *InMemoryUserSkillRepo) Update(growth stage.SkillGrowthPost) error {
	growthData := growth.SkillGrowth
	for _, v := range growthData {
		if _, ok := m.Data[growth.UserId]; !ok {
			m.Data[growth.UserId] = map[core.SkillId]stage.UserSkillRes{}
		}
		m.Data[growth.UserId][v.SkillId] = stage.UserSkillRes{
			UserId:   growth.UserId,
			SkillId:  v.SkillId,
			SkillExp: v.SkillExp,
		}
	}
	return nil
}

func CreateInMemoryUserSkillRepo(loader IUserSkillLoader) (*InMemoryUserSkillRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on create in memory user skill data: %w", err)
	}
	return &InMemoryUserSkillRepo{Data: data}, err
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

type StageMasterData map[stage.StageId]stage.StageMaster

type InMemoryStageMasterRepo struct {
	Data StageMasterData
}

type IStageMasterLoader interface {
	Load() (StageMasterData, error)
}

func (m *InMemoryStageMasterRepo) Get(stageId stage.StageId) (stage.StageMaster, error) {
	return m.Data[stageId], nil
}

func (m *InMemoryStageMasterRepo) GetAllStages() (stage.GetAllStagesRes, error) {
	result := []stage.StageMaster{}
	for _, v := range m.Data {
		result = append(result, v)
	}
	return stage.GetAllStagesRes{Stages: result}, nil
}

func CreateInMemoryStageMasterRepo(loader IStageMasterLoader) (*InMemoryStageMasterRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on stage master repo: %w", err)
	}
	repo := InMemoryStageMasterRepo{Data: data}
	return &repo, nil
}

type SkillGrowthRepoData map[stage.ExploreId][]stage.SkillGrowthData

type ISkillGrowthLoader interface {
	Load() (SkillGrowthRepoData, error)
}

type InMemorySkillGrowthDataRepo struct {
	Data SkillGrowthRepoData
}

func (m *InMemorySkillGrowthDataRepo) BatchGet(exploreId stage.ExploreId) []stage.SkillGrowthData {
	return m.Data[exploreId]
}

func (m *InMemorySkillGrowthDataRepo) Add(e stage.ExploreId, skills []stage.SkillGrowthData) {
	m.Data[e] = skills
}

func CreateInMemorySkillGrowthDataRepo(loader ISkillGrowthLoader) (*InMemorySkillGrowthDataRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return &InMemorySkillGrowthDataRepo{}, err
	}
	repo := InMemorySkillGrowthDataRepo{Data: data}
	return &repo, nil
}

type EarningItemData map[stage.ExploreId][]stage.EarningItem

type InMemoryEarningItemRepo struct {
	Data EarningItemData
}

type IEarningItemLoader interface {
	Load() (EarningItemData, error)
}

func (m *InMemoryEarningItemRepo) BatchGet(exploreId stage.ExploreId) []stage.EarningItem {
	return m.Data[exploreId]
}

func CreateInMemoryEarningItemRepo(loader IEarningItemLoader) (*InMemoryEarningItemRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on create in memory earning item repo: %w", err)
	}
	return &InMemoryEarningItemRepo{Data: data}, nil
}

type ConsumingItemData map[stage.ExploreId][]stage.ConsumingItem

type InMemoryConsumingItemRepo struct {
	Data ConsumingItemData
}

func (m *InMemoryConsumingItemRepo) BatchGet(exploreId stage.ExploreId) ([]stage.ConsumingItem, error) {
	data, ok := m.Data[exploreId]
	if !ok {
		return []stage.ConsumingItem{}, nil
	}
	return data, nil
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

func CreateInMemoryConsumingItemRepo(loader IConsumingItemLoader) (*InMemoryConsumingItemRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on creating in memory consuming item repo: %w", err)
	}
	return &InMemoryConsumingItemRepo{Data: data}, nil
}

type RequiredSkillData map[stage.ExploreId][]stage.RequiredSkill

type IRequiredSkillLoader interface {
	Load() (RequiredSkillData, error)
}

type InMemoryRequiredSkillRepo struct {
	Data RequiredSkillData
}

func (m *InMemoryRequiredSkillRepo) Get(exploreId stage.ExploreId) ([]stage.RequiredSkill, error) {
	if _, ok := m.Data[exploreId]; !ok {
		return []stage.RequiredSkill{}, nil
	}
	return m.Data[exploreId], nil
}

func (m *InMemoryRequiredSkillRepo) BatchGet(ids []stage.ExploreId) ([]stage.RequiredSkillRow, error) {
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

func CreateInMemoryRequiredSkillRepo(loader IRequiredSkillLoader) (*InMemoryRequiredSkillRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on create required skill repo: %w", err)
	}
	return &InMemoryRequiredSkillRepo{Data: data}, nil
}

type InMemoryItemStorageUpdateRepo struct {
	Data map[core.UserId][]stage.ItemStock
}

func (m *InMemoryItemStorageUpdateRepo) Update(userId core.UserId, items []stage.ItemStock, _ core.AccessToken) error {
	m.Data = make(map[core.UserId][]stage.ItemStock)
	m.Data[userId] = items
	return nil
}

func (m *InMemoryItemStorageUpdateRepo) Get(userId core.UserId) []stage.ItemStock {
	return m.Data[userId]
}

func CreateInMemoryItemStorageUpdateRepo() *InMemoryItemStorageUpdateRepo {
	return &InMemoryItemStorageUpdateRepo{}
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

type ReductionStaminaSkillData map[stage.ExploreId][]core.SkillId

type InMemoryReductionStaminaSkillRepo struct {
	Data ReductionStaminaSkillData
}

type IReductionStaminaSkillLoader interface {
	Load() (ReductionStaminaSkillData, error)
}

func (m *InMemoryReductionStaminaSkillRepo) Get(exploreId stage.ExploreId) ([]core.SkillId, error) {
	return m.Data[exploreId], nil
}

func (m *InMemoryReductionStaminaSkillRepo) BatchGet(exploreIds []stage.ExploreId) ([]stage.BatchGetReductionStaminaSkill, error) {
	result := make([]stage.BatchGetReductionStaminaSkill, len(exploreIds))
	for i, v := range exploreIds {
		result[i] = stage.BatchGetReductionStaminaSkill{
			ExploreId: v,
			Skills:    m.Data[v],
		}
	}
	return result, nil
}

func CreateInMemoryReductionStaminaSkillRepo(loader IReductionStaminaSkillLoader) (*InMemoryReductionStaminaSkillRepo, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error on create reduction skill repo: %w", err)
	}
	return &InMemoryReductionStaminaSkillRepo{Data: data}, nil
}

func (m *InMemoryReductionStaminaSkillRepo) Add(e stage.ExploreId, skills []core.SkillId) {
	m.Data[e] = skills
}
