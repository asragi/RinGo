package stage

import (
	"github.com/asragi/RinGo/core"
)

type fetchItemStorageTester struct {
	passedId      core.UserId
	passedItemIds []core.ItemId
	returnVal     BatchGetStorageRes
	returnErr     error
}

func (t *fetchItemStorageTester) BatchGet(
	id core.UserId,
	items []core.ItemId,
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
