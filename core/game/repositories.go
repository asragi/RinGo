package game

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type (
	// Deprecated: use GetFundFunc, or GetStaminaFunc
	GetResourceFunc func(context.Context, core.UserId) (*GetResourceRes, error)
	StaminaRes      struct {
		UserId             core.UserId             `db:"user_id"`
		MaxStamina         core.MaxStamina         `db:"max_stamina"`
		StaminaRecoverTime core.StaminaRecoverTime `db:"stamina_recover_time"`
	}
	GetStaminaFunc func(context.Context, []core.UserId) ([]*StaminaRes, error)
	FundRes        struct {
		UserId core.UserId `db:"user_id"`
		Fund   core.Fund   `db:"fund"`
	}
	GetFundFunc    func(context.Context, []core.UserId) ([]*FundRes, error)
	GetResourceRes struct {
		UserId             core.UserId             `db:"user_id"`
		MaxStamina         core.MaxStamina         `db:"max_stamina"`
		StaminaRecoverTime core.StaminaRecoverTime `db:"stamina_recover_time"`
		Fund               core.Fund               `db:"fund"`
	}
	// Deprecated: use UpdateFundFunc
	UpdateFundFuncDeprecated func(context.Context, core.UserId, core.Fund) error
	UserFundPair             struct {
		UserId core.UserId `db:"user_id"`
		Fund   core.Fund   `db:"fund"`
	}
	UpdateFundFunc    func(context.Context, []*UserFundPair) error
	UpdateStaminaFunc func(context.Context, core.UserId, core.StaminaRecoverTime) error
)

type GetItemMasterRes struct {
	ItemId      core.ItemId      `db:"item_id" json:"item_id"`
	Price       core.Price       `db:"price" json:"price"`
	DisplayName core.DisplayName `db:"display_name" json:"display_name"`
	Description core.Description `db:"description" json:"description"`
	MaxStock    core.MaxStock    `db:"max_stock" json:"max_stock"`
}

func ItemMasterResToMap(res []*GetItemMasterRes) map[core.ItemId]*GetItemMasterRes {
	result := make(map[core.ItemId]*GetItemMasterRes)
	for _, v := range res {
		result[v.ItemId] = v
	}
	return result
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

func SpreadGetStorageRes(res []*BatchGetStorageRes) []*StorageData {
	result := make([]*StorageData, 0)
	for _, v := range res {
		result = append(result, v.ItemData...)
	}
	return result
}

func FindItemStorageData(data []*StorageData, itemId core.ItemId) *StorageData {
	for _, v := range data {
		if v.ItemId == itemId {
			return v
		}
	}
	return nil
}

func FindStorageData(res []*BatchGetStorageRes, userId core.UserId) *BatchGetStorageRes {
	for _, v := range res {
		if v.UserId == userId {
			return v
		}
	}
	return nil
}

func StorageDataToMap(res []*BatchGetStorageRes) map[core.UserId]map[core.ItemId]*StorageData {
	result := make(map[core.UserId]map[core.ItemId]*StorageData)
	for _, v := range res {
		result[v.UserId] = make(map[core.ItemId]*StorageData)
		for _, item := range v.ItemData {
			result[v.UserId][item.ItemId] = item
		}
	}
	return result
}

type UserItemPair struct {
	UserId core.UserId `db:"user_id"`
	ItemId core.ItemId `db:"item_id"`
}

type FetchStorageFunc func(context.Context, []*UserItemPair) ([]*BatchGetStorageRes, error)

// Deprecated: use FetchStorageFunc
type FetchStorageFuncDeprecated func(context.Context, core.UserId, []core.ItemId) (BatchGetStorageRes, error)

type FetchAllStorageFunc func(context.Context, core.UserId) ([]*StorageData, error)

type TotalItemStock struct {
	ItemId     core.ItemId
	AfterStock core.Stock
	IsKnown    core.IsKnown
}

func totalItemToItemStock(totalItems []*totalItem) []*TotalItemStock {
	result := make([]*TotalItemStock, len(totalItems))
	for i, v := range totalItems {
		result[i] = &TotalItemStock{
			ItemId:     v.ItemId,
			AfterStock: v.Stock,
			IsKnown:    true,
		}
	}
	return result
}

// Deprecated: use UpdateItemStorageFunc
type UpdateItemStorageFuncDeprecated func(context.Context, core.UserId, []*TotalItemStock) error
type UpdateItemStorageFunc func(context.Context, []*StorageData) error

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

type RequiredSkill struct {
	ExploreId  ExploreId    `db:"explore_id"`
	SkillId    core.SkillId `db:"skill_id"`
	RequiredLv core.SkillLv `db:"skill_lv"`
}

type FetchRequiredSkillsFunc func(context.Context, []ExploreId) ([]*RequiredSkill, error)

type RequiredSkillRow struct {
	ExploreId      ExploreId
	RequiredSkills []RequiredSkill
}

type ConsumingItem struct {
	ExploreId       ExploreId       `db:"explore_id"`
	ItemId          core.ItemId     `db:"item_id"`
	MaxCount        core.Count      `db:"max_count"`
	ConsumptionProb ConsumptionProb `db:"consumption_prob"`
}

type BatchGetConsumingItemRes struct {
	ExploreId      ExploreId
	ConsumingItems []ConsumingItem
}

type FetchConsumingItemFunc func(context.Context, []ExploreId) ([]*ConsumingItem, error)

type BatchGetReductionStaminaSkill struct {
	ExploreId ExploreId
	Skills    []StaminaReductionSkillPair
}

type StaminaReductionSkillPair struct {
	ExploreId ExploreId    `db:"explore_id"`
	SkillId   core.SkillId `db:"skill_id"`
}

type FetchReductionStaminaSkillFunc func(context.Context, []ExploreId) ([]*StaminaReductionSkillPair, error)

func ReductionStaminaSkillToMap(res []*StaminaReductionSkillPair) map[ExploreId][]core.SkillId {
	result := make(map[ExploreId][]core.SkillId)
	for _, v := range res {
		result[v.ExploreId] = append(result[v.ExploreId], v.SkillId)
	}
	return result
}

type EarningItem struct {
	ItemId      core.ItemId `db:"item_id"`
	MinCount    core.Count  `db:"min_count"`
	MaxCount    core.Count  `db:"max_count"`
	Probability EarningProb `db:"probability"`
}

type FetchEarningItemFunc func(context.Context, ExploreId) ([]*EarningItem, error)

type GetExploreMasterRes struct {
	ExploreId            ExploreId            `db:"explore_id"`
	DisplayName          core.DisplayName     `db:"display_name"`
	Description          core.Description     `db:"description"`
	ConsumingStamina     core.StaminaCost     `db:"consuming_stamina"`
	RequiredPayment      core.Cost            `db:"required_payment"`
	StaminaReducibleRate StaminaReducibleRate `db:"stamina_reducible_rate"`
}

type FetchExploreMasterFunc func(context.Context, []ExploreId) ([]*GetExploreMasterRes, error)

type UserExplore struct {
	ExploreId   ExploreId
	DisplayName core.DisplayName
	IsKnown     core.IsKnown
	IsPossible  core.IsPossible
}

type ExploreUserData struct {
	ExploreId ExploreId    `db:"explore_id"`
	IsKnown   core.IsKnown `db:"is_known"`
}

type GetUserExploreFunc func(context.Context, core.UserId, []ExploreId) ([]*ExploreUserData, error)

type GetActionsRes struct {
	UserId   core.UserId
	Explores []*ExploreUserData
}
