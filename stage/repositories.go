package stage

import (
	"github.com/asragi/RinGo/core"
)

type GetItemMasterRes struct {
	ItemId      core.ItemId
	Price       core.Price
	DisplayName core.DisplayName
	Description core.Description
	MaxStock    core.MaxStock
}

type ItemMasterRepo interface {
	Get(core.ItemId) (GetItemMasterRes, error)
}

type GetItemStorageRes struct {
	UserId core.UserId
	Stock  core.Stock
}

type ItemData struct {
	UserId core.UserId
	ItemId core.ItemId
	Stock  core.Stock
}

type BatchGetStorageRes struct {
	UserId   core.UserId
	ItemData []ItemData
}

type ItemStorageRepo interface {
	Get(core.UserId, core.ItemId, core.AccessToken) (GetItemStorageRes, error)
	BatchGet(core.UserId, []core.ItemId, core.AccessToken) (BatchGetStorageRes, error)
}

type RequiredItemData struct {
	ExploreId       ExploreId
	ItemId          core.ItemId
	ConsumptionProb ConsumptionProb
}

type RequiredItemRepo interface {
	Get(ExploreId) []RequiredItemData
}

// Skill
type SkillMaster struct {
	SkillId     core.SkillId
	DisplayName core.DisplayName
}

type BatchGetSkillMasterRes struct {
	Skills []SkillMaster
}

type SkillMasterRepo interface {
	BatchGet([]core.SkillId) (BatchGetSkillMasterRes, error)
}

type UserSkillRes struct {
	UserId   core.UserId
	SkillId  core.SkillId
	SkillExp core.SkillExp
}

type BatchGetUserSkillRes struct {
	UserId core.UserId
	Skills []UserSkillRes
}

type UserSkillRepo interface {
	BatchGet(core.UserId, []core.SkillId, core.AccessToken) (BatchGetUserSkillRes, error)
}

// Skill Growth
type SkillGrowthData struct {
	ExploreId    ExploreId
	SkillId      core.SkillId
	GainingPoint GainingPoint
}

type SkillGrowthDataRepo interface {
	BatchGet(ExploreId) []SkillGrowthData
}

// Explore
type GetAllExploreMasterRes struct {
	ExploreId   ExploreId
	DisplayName core.DisplayName
	Description core.Description
}

type StageExploreMasterRes struct {
	StageId  StageId
	Explores []GetAllExploreMasterRes
}

type BatchGetStageExploreRes struct {
	StageExplores []StageExploreMasterRes
}

type ExploreMasterRepo interface {
	GetAllExploreMaster(core.ItemId) ([]GetAllExploreMasterRes, error)
	GetStageAllExploreMaster([]StageId) (BatchGetStageExploreRes, error)
}

type ExploreUserData struct {
	ExploreId ExploreId
	IsKnown   core.IsKnown
}

type GetActionsRes struct {
	UserId   core.UserId
	Explores []ExploreUserData
}

type UserExploreRepo interface {
	GetActions(core.UserId, []ExploreId, core.AccessToken) (GetActionsRes, error)
}

type Condition struct {
	ConditionId          ConditionId
	ConditionType        ConditionType
	ConditionTargetId    ConditionTargetId
	ConditionTargetValue ConditionTargetValue
}

type ExploreConditions struct {
	ExploreId  ExploreId
	Conditions []Condition
}

type GetAllConditionsRes struct {
	Explores []ExploreConditions
}

type ConditionRepo interface {
	GetAllConditions([]ExploreId) (GetAllConditionsRes, error)
}

type StageMaster struct {
	StageId     StageId
	DisplayName core.DisplayName
	Description core.Description
}

type GetAllStagesRes struct {
	Stages []StageMaster
}

type StageMasterRepo interface {
	GetAllStages() (GetAllStagesRes, error)
}

type UserStage struct {
	StageId StageId
	IsKnown core.IsKnown
}

type GetAllUserStagesRes struct {
	UserStage []UserStage
}

type UserStageRepo interface {
	GetAllUserStages(core.UserId, []StageId) (GetAllUserStagesRes, error)
}

type EarningItem struct {
	ItemId   core.ItemId
	MinCount core.Count
	MaxCount core.Count
}

type EarningItemRepo interface {
	BatchGet(ExploreId) []EarningItem
}

type ConsumingItem struct {
	ItemId          core.ItemId
	MaxCount        core.Count
	ConsumptionProb ConsumptionProb
}

type ConsumingItemRepo interface {
	BatchGet(ExploreId) []ConsumingItem
}
