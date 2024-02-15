package stage

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
)

var (
	exploreMasterRepo  = createMockExploreMasterRepo()
	userSkillRepo      = createMockUserSkillRepo()
	reductionSkillRepo = createMockReductionStaminaSkillRepo()
)

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

func (m *MockExploreMasterRepo) Add(e ExploreId, master GetExploreMasterRes) {
	m.Data[e] = master
}

func createMockExploreMasterRepo() *MockExploreMasterRepo {
	return &MockExploreMasterRepo{
		Data: map[ExploreId]GetExploreMasterRes{},
	}
}

type MockUserSkillRepo struct {
	Data map[core.UserId]map[core.SkillId]UserSkillRes
}

func (m *MockUserSkillRepo) BatchGet(
	userId core.UserId,
	skillIds []core.SkillId,
	_ auth.AccessToken,
) (BatchGetUserSkillRes, error) {
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

type mockReductionStaminaSkillRepo struct {
	Data map[ExploreId][]core.SkillId
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

type fetchItemStorageTester struct {
	passedId      core.UserId
	passedItemIds []core.ItemId
	returnVal     BatchGetStorageRes
	returnErr     error
}

func (t *fetchItemStorageTester) BatchGet(
	id core.UserId,
	items []core.ItemId,
	_ auth.AccessToken,
) (BatchGetStorageRes, error) {
	t.passedId = id
	t.passedItemIds = items
	return t.returnVal, t.returnErr
}

type fetchExploreTester struct {
	passedArgs []ExploreId
	returnVal  []GetExploreMasterRes
	returnErr  error
}

func (t *fetchExploreTester) BatchGet(ids []ExploreId) ([]GetExploreMasterRes, error) {
	t.passedArgs = ids
	return t.returnVal, t.returnErr
}

type fetchEarningItemTester struct {
	passedId  ExploreId
	returnVal []EarningItem
	returnErr error
}

func (t *fetchEarningItemTester) Get(id ExploreId) ([]EarningItem, error) {
	t.passedId = id
	return t.returnVal, t.returnErr
}

type fetchConsumingTester struct {
	passedArgs []ExploreId
	returnVal  []BatchGetConsumingItemRes
	returnErr  error
}

func (t *fetchConsumingTester) BatchGet(ids []ExploreId) ([]BatchGetConsumingItemRes, error) {
	t.passedArgs = ids
	return t.returnVal, t.returnErr
}

type fetchSkillMasterTester struct {
	passedArgs []core.SkillId
	returnVal  []SkillMaster
	returnErr  error
}

func (t *fetchSkillMasterTester) BatchGet(args []core.SkillId) ([]SkillMaster, error) {
	t.passedArgs = args
	return t.returnVal, t.returnErr
}

type fetchUserSkillTester struct {
	passedId     core.UserId
	passedSkills []core.SkillId
	returnValue  BatchGetUserSkillRes
	returnErr    error
}

func (t *fetchUserSkillTester) BatchGet(
	id core.UserId,
	skills []core.SkillId,
	_ auth.AccessToken,
) (
	BatchGetUserSkillRes,
	error,
) {
	t.passedId = id
	t.passedSkills = skills
	return t.returnValue, t.returnErr
}

type fetchRequiredSkillTester struct {
	passedArgs []ExploreId
	returnVal  []RequiredSkillRow
	returnErr  error
}

func (t *fetchRequiredSkillTester) BatchGet(ids []ExploreId) ([]RequiredSkillRow, error) {
	t.passedArgs = ids
	return t.returnVal, t.returnErr
}
