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

type GetItemMasterRes struct {
	ItemId      core.ItemId      `db:"item_id" json:"item_id"`
	Price       core.Price       `db:"price" json:"price"`
	DisplayName core.DisplayName `db:"display_name" json:"display_name"`
	Description core.Description `db:"description" json:"description"`
	MaxStock    core.MaxStock    `db:"max_stock" json:"max_stock"`
}

type FetchItemMasterFunc func([]core.ItemId) ([]GetItemMasterRes, error)

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

type FetchStorageFunc func(core.UserId, []core.ItemId, core.AccessToken) (BatchGetStorageRes, error)

type FetchAllStorageFunc func(core.UserId) ([]ItemData, error)

type ItemStock struct {
	ItemId     core.ItemId
	AfterStock core.Stock
}

type UpdateItemStorageFunc func(core.UserId, []ItemStock, core.AccessToken) error

type SkillMaster struct {
	SkillId     core.SkillId
	DisplayName core.DisplayName
}

type FetchSkillMasterFunc func([]core.SkillId) ([]SkillMaster, error)

type UserSkillRes struct {
	UserId   core.UserId
	SkillId  core.SkillId
	SkillExp core.SkillExp
}

type BatchGetUserSkillRes struct {
	UserId core.UserId
	Skills []UserSkillRes
}

type FetchUserSkillFunc func(core.UserId, []core.SkillId, core.AccessToken) (BatchGetUserSkillRes, error)

type SkillGrowthData struct {
	ExploreId    ExploreId    `db:"explore_id"`
	SkillId      core.SkillId `db:"skill_id"`
	GainingPoint GainingPoint `db:"gaining_point"`
}

type FetchSkillGrowthData func(ExploreId) ([]SkillGrowthData, error)

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
	ExploreIds []StageExploreIdPairRow
}

func (p StageExploreIdPair) CreateSelf(id StageId, data []StageExploreIdPairRow) StageExploreIdPair {
	return StageExploreIdPair{
		StageId:    id,
		ExploreIds: data,
	}
}

type StageExploreIdPairRow struct {
	StageId   StageId   `db:"stage_id"`
	ExploreId ExploreId `db:"explore_id"`
}

func (row StageExploreIdPairRow) GetId() StageId {
	return row.StageId
}

type FetchItemExploreRelationFunc func(core.ItemId) ([]ExploreId, error)

type FetchStageExploreRelation func([]StageId) ([]StageExploreIdPair, error)

type FetchExploreMasterFunc func([]ExploreId) ([]GetExploreMasterRes, error)

type ExploreUserData struct {
	ExploreId ExploreId
	IsKnown   core.IsKnown
}

type GetUserExploreFunc func(core.UserId, []ExploreId, core.AccessToken) (GetActionsRes, error)

type GetActionsRes struct {
	UserId   core.UserId
	Explores []ExploreUserData
}

type RequiredSkill struct {
	ExploreId  ExploreId    `db:"explore_id"`
	SkillId    core.SkillId `db:"skill_id"`
	RequiredLv core.SkillLv `db:"skill_lv"`
}

func (r RequiredSkill) GetId() ExploreId {
	return r.ExploreId
}

type FetchRequiredSkillsFunc func([]ExploreId) ([]RequiredSkillRow, error)

type RequiredSkillRow struct {
	ExploreId      ExploreId
	RequiredSkills []RequiredSkill
}

func (_ RequiredSkillRow) CreateSelf(id ExploreId, data []RequiredSkill) RequiredSkillRow {
	return RequiredSkillRow{
		ExploreId:      id,
		RequiredSkills: data,
	}
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

type ConsumingItem struct {
	ExploreId       ExploreId       `db:"explore_id"`
	ItemId          core.ItemId     `db:"item_id"`
	MaxCount        core.Count      `db:"max_count"`
	ConsumptionProb ConsumptionProb `db:"consumption_prob"`
}

func (c ConsumingItem) GetId() ExploreId {
	return c.ExploreId
}

type BatchGetConsumingItemRes struct {
	ExploreId      ExploreId
	ConsumingItems []ConsumingItem
}

func (_ BatchGetConsumingItemRes) CreateSelf(id ExploreId, data []ConsumingItem) BatchGetConsumingItemRes {
	return BatchGetConsumingItemRes{
		ExploreId:      id,
		ConsumingItems: data,
	}
}

type FetchConsumingItemFunc func([]ExploreId) ([]BatchGetConsumingItemRes, error)

type BatchGetReductionStaminaSkill struct {
	ExploreId ExploreId
	Skills    []StaminaReductionSkillPair
}

func (_ BatchGetReductionStaminaSkill) CreateSelf(
	id ExploreId,
	data []StaminaReductionSkillPair,
) BatchGetReductionStaminaSkill {
	return BatchGetReductionStaminaSkill{
		ExploreId: id,
		Skills:    data,
	}
}

type StaminaReductionSkillPair struct {
	ExploreId ExploreId    `db:"explore_id"`
	SkillId   core.SkillId `db:"skill_id"`
}

func (id StaminaReductionSkillPair) GetId() ExploreId {
	return id.ExploreId
}

type FetchReductionStaminaSkillFunc func([]ExploreId) ([]BatchGetReductionStaminaSkill, error)
