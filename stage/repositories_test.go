package stage

import (
	"time"

	"github.com/asragi/RinGo/core"
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

var MockConditions map[ExploreId][]Condition = map[ExploreId][]Condition{
	mockExploreIds[0]: {
		{
			ConditionId:          "enough-stick",
			ConditionType:        ConditionTypeItem,
			ConditionTargetId:    ConditionTargetId(MockItems[2].ItemId),
			ConditionTargetValue: ConditionTargetValue(10),
		},
		{
			ConditionId:          "enough-apple-lv",
			ConditionType:        ConditionTypeSkill,
			ConditionTargetId:    ConditionTargetId(MockSkillMaster[0].SkillId),
			ConditionTargetValue: ConditionTargetValue(3),
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

var mockExploreMaster = map[core.ItemId][]GetAllExploreMasterRes{
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
}

var mockStageIds = []StageId{
	StageId("forest"),
	StageId("volcano"),
}

var mockStageExploreMaster = map[StageId][]GetAllExploreMasterRes{
	mockStageIds[0]: {
		{
			ExploreId:   mockStageExploreIds[0],
			DisplayName: "りんごを拾いに行く",
			Description: "木くずや石も拾えるかも",
		},
	},
}

type MockExploreMasterRepo struct {
	Data      map[core.ItemId][]GetAllExploreMasterRes
	StageData map[StageId][]GetAllExploreMasterRes
}

func (m *MockExploreMasterRepo) GetAllExploreMaster(itemId core.ItemId) ([]GetAllExploreMasterRes, error) {
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

func createMockExploreMasterRepo() *MockExploreMasterRepo {
	repo := MockExploreMasterRepo{}
	repo.Data = mockExploreMaster
	repo.StageData = mockStageExploreMaster
	return &repo
}

var MockSkillMaster = []SkillMaster{
	{
		SkillId:     "apple",
		DisplayName: "りんご愛好家",
	},
	{
		SkillId:     "fire",
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
			UserId:  MockUserId,
			SkillId: MockSkillMaster[0].SkillId,
			SkillLv: 3,
		},
		MockSkillMaster[1].SkillId: {
			UserId:  MockUserId,
			SkillId: MockSkillMaster[0].SkillId,
			SkillLv: 1,
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
			UserId:  userId,
			SkillId: v,
			SkillLv: list[v].SkillLv,
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
