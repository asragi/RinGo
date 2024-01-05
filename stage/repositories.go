package stage

import (
	"github.com/asragi/RinGo/core"
)

type GetResourceFunc func(core.UserId, core.AccessToken) (GetResourceRes, error)

type GetResourceRes struct {
	UserId             core.UserId
	MaxStamina         core.MaxStamina
	StaminaRecoverTime core.StaminaRecoverTime
	Fund               core.Fund
}

type UpdateFundFunc func(core.UserId, core.Fund) error

type UpdateStaminaFunc func(core.UserId, core.StaminaRecoverTime) error

// Deprecated: use GetResourceFunc or UpdateStaminaFunc
type UserResourceRepo interface {
	GetResource(core.UserId, core.AccessToken) (GetResourceRes, error)
	UpdateStamina(core.UserId, core.AccessToken, core.StaminaRecoverTime) error
}

type GetItemMasterRes struct {
	ItemId      core.ItemId
	Price       core.Price
	DisplayName core.DisplayName
	Description core.Description
	MaxStock    core.MaxStock
}

// Deprecated: use BatchGetItemMasterFunc
type GetItemMasterFunc func(core.ItemId) (GetItemMasterRes, error)
type BatchGetItemMasterFunc func([]core.ItemId) ([]GetItemMasterRes, error)

type ItemMasterRepo interface {
	Get(core.ItemId) (GetItemMasterRes, error)
	BatchGet([]core.ItemId) ([]GetItemMasterRes, error)
}

type GetItemStorageRes struct {
	UserId core.UserId
	Stock  core.Stock
}

type ItemData struct {
	UserId  core.UserId
	ItemId  core.ItemId
	Stock   core.Stock
	IsKnown core.IsKnown
}

type BatchGetStorageRes struct {
	UserId   core.UserId
	ItemData []ItemData
}

// Deprecated: user BatchGetStorageFunc
type GetItemStorageFunc func(core.UserId, core.ItemId, core.AccessToken) (GetItemStorageRes, error)

type BatchGetStorageFunc func(core.UserId, []core.ItemId, core.AccessToken) (BatchGetStorageRes, error)

type GetAllStorageFunc func(core.UserId) ([]ItemData, error)

type ItemStorageRepo interface {
	Get(core.UserId, core.ItemId, core.AccessToken) (GetItemStorageRes, error)
	BatchGet(core.UserId, []core.ItemId, core.AccessToken) (BatchGetStorageRes, error)
}

type ItemStock struct {
	ItemId     core.ItemId
	AfterStock core.Stock
}

type UpdateItemStorageFunc func(core.UserId, []ItemStock, core.AccessToken) error

// Skill
type SkillMaster struct {
	SkillId     core.SkillId
	DisplayName core.DisplayName
}

type BatchGetSkillMasterRes struct {
	Skills []SkillMaster
}

type FetchSkillMasterFunc func([]core.SkillId) ([]SkillMaster, error)

// Deprecated: use FetchSkillMasterFunc
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

type BatchGetUserSkillFunc func(core.UserId, []core.SkillId, core.AccessToken) (BatchGetUserSkillRes, error)

// Deprecated: use BatchGetUserSkillFunc
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

type SkillGrowthPostRow struct {
	SkillId  core.SkillId
	SkillExp core.SkillExp
}

type SkillGrowthPost struct {
	UserId      core.UserId
	AccessToken core.AccessToken
	SkillGrowth []SkillGrowthPostRow
}

type SkillGrowthPostFunc func(SkillGrowthPost) error

// Explore
type GetExploreMasterRes struct {
	ExploreId            ExploreId
	DisplayName          core.DisplayName
	Description          core.Description
	ConsumingStamina     core.Stamina
	RequiredPayment      core.Price
	StaminaReducibleRate StaminaReducibleRate
}

type StageExploreIdPair struct {
	StageId    StageId
	ExploreIds []ExploreId
}

type GetItemExploreRelationFunc func(core.ItemId) ([]ExploreId, error)

type FetchStageExploreRelation func([]StageId) ([]StageExploreIdPair, error)

type FetchExploreMasterFunc func([]ExploreId) ([]GetExploreMasterRes, error)

// Deprecated: use fetchExploreMasterFunc
type ExploreMasterRepo interface {
	Get(ExploreId) (GetExploreMasterRes, error)
	BatchGet([]ExploreId) ([]GetExploreMasterRes, error)
}

type ExploreUserData struct {
	ExploreId ExploreId
	IsKnown   core.IsKnown
}

type GetActionsFunc func(core.UserId, []ExploreId, core.AccessToken) (GetActionsRes, error)

type GetActionsRes struct {
	UserId   core.UserId
	Explores []ExploreUserData
}

type RequiredSkill struct {
	SkillId    core.SkillId
	RequiredLv core.SkillLv
}

type FetchRequiredSkillsFunc func([]ExploreId) ([]RequiredSkillRow, error)

type RequiredSkillRow struct {
	ExploreId      ExploreId
	RequiredSkills []RequiredSkill
}

// Deprecated: use FetchRequiredSkillsFunc
type RequiredSkillRepo interface {
	Get(ExploreId) ([]RequiredSkill, error)
	BatchGet([]ExploreId) ([]RequiredSkillRow, error)
}

type StageMaster struct {
	StageId     StageId
	DisplayName core.DisplayName
	Description core.Description
}

type GetAllStagesRes struct {
	Stages []StageMaster
}

type FetchStageMasterFunc func(StageId) (StageMaster, error)
type FetchAllStageFunc func() (GetAllStagesRes, error)

type FetchUserStageFunc func(core.UserId, []StageId) (GetAllUserStagesRes, error)

type UserStage struct {
	StageId StageId
	IsKnown core.IsKnown
}

type GetAllUserStagesRes struct {
	UserStage []UserStage
}

type EarningItem struct {
	ItemId   core.ItemId
	MinCount core.Count
	MaxCount core.Count
}

type FetchEarningItemFunc func(ExploreId) ([]EarningItem, error)

// Deprecated: Replace with FetchEarningItemFunc
type EarningItemRepo interface {
	BatchGet(ExploreId) []EarningItem
}

type ConsumingItem struct {
	ItemId          core.ItemId
	MaxCount        core.Count
	ConsumptionProb ConsumptionProb
}

type BatchGetConsumingItemRes struct {
	ExploreId      ExploreId
	ConsumingItems []ConsumingItem
}

type GetConsumingItemFunc func([]ExploreId) ([]BatchGetConsumingItemRes, error)

// Deprecated: use GetConsumingItemFunc
type ConsumingItemRepo interface {
	BatchGet(ExploreId) ([]ConsumingItem, error)
	AllGet([]ExploreId) ([]BatchGetConsumingItemRes, error)
}

type BatchGetReductionStaminaSkill struct {
	ExploreId ExploreId
	Skills    []core.SkillId
}

type FetchReductionStaminaSkillFunc func([]ExploreId) ([]BatchGetReductionStaminaSkill, error)
