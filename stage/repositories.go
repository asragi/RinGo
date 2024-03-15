package stage

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type GetResourceFunc func(context.Context, core.UserId) (*GetResourceRes, error)

type GetResourceRes struct {
	UserId             core.UserId             `db:"user_id"`
	MaxStamina         core.MaxStamina         `db:"max_stamina"`
	StaminaRecoverTime core.StaminaRecoverTime `db:"stamina_recover_time"`
	Fund               core.Fund               `db:"fund"`
}

type UpdateFundFunc func(context.Context, core.UserId, core.Fund) error

type UpdateStaminaFunc func(context.Context, core.UserId, core.StaminaRecoverTime) error

type GetItemMasterRes struct {
	ItemId      core.ItemId      `db:"item_id" json:"item_id"`
	Price       core.Price       `db:"price" json:"price"`
	DisplayName core.DisplayName `db:"display_name" json:"display_name"`
	Description core.Description `db:"description" json:"description"`
	MaxStock    core.MaxStock    `db:"max_stock" json:"max_stock"`
}

type FetchItemMasterFunc func(context.Context, []core.ItemId) ([]*GetItemMasterRes, error)

type StorageData struct {
	UserId  core.UserId  `db:"user_id"`
	ItemId  core.ItemId  `db:"item_id"`
	Stock   core.Stock   `db:"stock"`
	IsKnown core.IsKnown `db:"is_known"`
}

type BatchGetStorageRes struct {
	UserId   core.UserId
	ItemData []*StorageData
}

type FetchStorageFunc func(context.Context, core.UserId, []core.ItemId) (BatchGetStorageRes, error)

type FetchAllStorageFunc func(context.Context, core.UserId) ([]*StorageData, error)

type ItemStock struct {
	ItemId     core.ItemId
	AfterStock core.Stock
	IsKnown    core.IsKnown
}

type UpdateItemStorageFunc func(context.Context, core.UserId, []*ItemStock) error

type SkillMaster struct {
	SkillId     core.SkillId     `db:"skill_id"`
	DisplayName core.DisplayName `db:"display_name"`
}

type FetchSkillMasterFunc func(context.Context, []core.SkillId) ([]*SkillMaster, error)

type UserSkillRes struct {
	UserId   core.UserId   `db:"user_id"`
	SkillId  core.SkillId  `db:"skill_id"`
	SkillExp core.SkillExp `db:"skill_exp"`
}

type BatchGetUserSkillRes struct {
	UserId core.UserId
	Skills []*UserSkillRes
}

type FetchUserSkillFunc func(context.Context, core.UserId, []core.SkillId) (BatchGetUserSkillRes, error)

type SkillGrowthData struct {
	ExploreId    ExploreId    `db:"explore_id"`
	SkillId      core.SkillId `db:"skill_id"`
	GainingPoint GainingPoint `db:"gaining_point"`
}

type FetchSkillGrowthData func(context.Context, ExploreId) ([]*SkillGrowthData, error)

type SkillGrowthPostRow struct {
	UserId   core.UserId   `db:"user_id"`
	SkillId  core.SkillId  `db:"skill_id"`
	SkillExp core.SkillExp `db:"skill_exp"`
}

type SkillGrowthPost struct {
	UserId      core.UserId `db:"user_id"`
	SkillGrowth []*SkillGrowthPostRow
}

type UpdateUserSkillExpFunc func(context.Context, SkillGrowthPost) error

type GetExploreMasterRes struct {
	ExploreId            ExploreId            `db:"explore_id"`
	DisplayName          core.DisplayName     `db:"display_name"`
	Description          core.Description     `db:"description"`
	ConsumingStamina     core.Stamina         `db:"consuming_stamina"`
	RequiredPayment      core.Price           `db:"required_payment"`
	StaminaReducibleRate StaminaReducibleRate `db:"stamina_reducible_rate"`
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

type FetchItemExploreRelationFunc func(context.Context, core.ItemId) ([]ExploreId, error)

type FetchStageExploreRelation func(context.Context, []StageId) ([]*StageExploreIdPairRow, error)

type FetchExploreMasterFunc func(context.Context, []ExploreId) ([]*GetExploreMasterRes, error)

type ExploreUserData struct {
	ExploreId ExploreId    `db:"explore_id"`
	IsKnown   core.IsKnown `db:"is_known"`
}

type GetUserExploreFunc func(context.Context, core.UserId, []ExploreId) ([]*ExploreUserData, error)

type GetActionsRes struct {
	UserId   core.UserId
	Explores []*ExploreUserData
}

type RequiredSkill struct {
	ExploreId  ExploreId    `db:"explore_id"`
	SkillId    core.SkillId `db:"skill_id"`
	RequiredLv core.SkillLv `db:"skill_lv"`
}

func (r RequiredSkill) GetId() ExploreId {
	return r.ExploreId
}

type FetchRequiredSkillsFunc func(context.Context, []ExploreId) ([]*RequiredSkill, error)

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
	StageId     StageId          `db:"stage_id"`
	DisplayName core.DisplayName `db:"display_name"`
	Description core.Description `db:"description"`
}

type GetAllStagesRes struct {
	Stages []StageMaster
}

type FetchStageMasterFunc func(context.Context, []StageId) ([]*StageMaster, error)
type FetchAllStageFunc func(context.Context) ([]*StageMaster, error)

type FetchUserStageFunc func(context.Context, core.UserId, []StageId) ([]*UserStage, error)

type UserStage struct {
	StageId StageId      `db:"stage_id"`
	IsKnown core.IsKnown `db:"is_known"`
}

type GetAllUserStagesRes struct {
	UserStage []UserStage
}

type EarningItem struct {
	ItemId      core.ItemId `db:"item_id"`
	MinCount    core.Count  `db:"min_count"`
	MaxCount    core.Count  `db:"max_count"`
	Probability EarningProb `db:"probability"`
}

type FetchEarningItemFunc func(context.Context, ExploreId) ([]*EarningItem, error)

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

type FetchConsumingItemFunc func(context.Context, []ExploreId) ([]*ConsumingItem, error)

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

type FetchReductionStaminaSkillFunc func(context.Context, []ExploreId) ([]*StaminaReductionSkillPair, error)
