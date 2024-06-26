package game

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type (
	GetResourceRes struct {
		UserId             core.UserId             `db:"user_id"`
		MaxStamina         core.MaxStamina         `db:"max_stamina"`
		StaminaRecoverTime core.StaminaRecoverTime `db:"stamina_recover_time"`
		Fund               core.Fund               `db:"fund"`
	}
	GetResourceFunc func(context.Context, core.UserId) (*GetResourceRes, error)
	StaminaRes      struct {
		UserId             core.UserId             `db:"user_id"`
		MaxStamina         core.MaxStamina         `db:"max_stamina"`
		StaminaRecoverTime core.StaminaRecoverTime `db:"stamina_recover_time"`
	}
	FetchStaminaFunc func(context.Context, []core.UserId) ([]*StaminaRes, error)
	FundRes          struct {
		UserId core.UserId `db:"user_id"`
		Fund   core.Fund   `db:"fund"`
	}
	FetchFundFunc func(context.Context, []core.UserId) ([]*FundRes, error)
	UserFundPair  struct {
		UserId core.UserId `db:"user_id"`
		Fund   core.Fund   `db:"fund"`
	}
	UpdateFundFunc    func(context.Context, []*UserFundPair) error
	UpdateStaminaFunc func(context.Context, core.UserId, core.StaminaRecoverTime) error
)

func FundPairToUserId(pairs []*UserFundPair) ([]core.UserId, []core.Fund) {
	userId := make([]core.UserId, len(pairs))
	fund := make([]core.Fund, len(pairs))
	for i, v := range pairs {
		userId[i] = v.UserId
		fund[i] = v.Fund
	}
	return userId, fund
}

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

func totalItemStockToStorageData(userId core.UserId, totalItems []*totalItem) []*StorageData {
	result := make([]*StorageData, len(totalItems))
	for i, v := range totalItems {
		result[i] = &StorageData{
			UserId: userId,
			ItemId: v.ItemId,
			Stock:  v.Stock,
		}
	}
	return result
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

func ToUserItemPair(userId core.UserId, itemIds []core.ItemId) []*UserItemPair {
	result := make([]*UserItemPair, len(itemIds))
	for i, v := range itemIds {
		result[i] = &UserItemPair{
			UserId: userId,
			ItemId: v,
		}
	}
	return result
}

type FetchStorageFunc func(context.Context, []*UserItemPair) ([]*BatchGetStorageRes, error)

func FillStorageData(res []*BatchGetStorageRes, userIdPair []*UserItemPair) []*BatchGetStorageRes {
	userIds := func() []core.UserId {
		check := make(map[core.UserId]bool)
		var result []core.UserId
		for _, v := range userIdPair {
			if _, ok := check[v.UserId]; ok {
				continue
			}
			check[v.UserId] = true
			result = append(result, v.UserId)
		}
		return result
	}()
	result := make([]*BatchGetStorageRes, len(userIds))
	for i, userId := range userIds {
		itemIds := func() []core.ItemId {
			check := make(map[core.ItemId]bool)
			var result []core.ItemId
			for _, v := range userIdPair {
				if v.UserId != userId {
					continue
				}
				if _, ok := check[v.ItemId]; ok {
					continue
				}
				check[v.ItemId] = true
				result = append(result, v.ItemId)
			}
			return result
		}()
		storageData := make([]*StorageData, len(itemIds))
		for j, itemId := range itemIds {
			storage := FindStorageData(res, userId)
			if storage == nil {
				storage = &BatchGetStorageRes{
					UserId:   userId,
					ItemData: make([]*StorageData, 0),
				}
			}
			storageData[j] = FindItemStorageData(storage.ItemData, itemId)
			if storageData[j] == nil {
				storageData[j] = &StorageData{
					UserId:  userId,
					ItemId:  itemId,
					Stock:   0,
					IsKnown: false,
				}
			}
		}
		result[i] = &BatchGetStorageRes{
			UserId:   userId,
			ItemData: storageData,
		}
	}
	return result
}

type FetchAllStorageFunc func(context.Context, core.UserId) ([]*StorageData, error)

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
	ExploreId    ActionId     `db:"explore_id"`
	SkillId      core.SkillId `db:"skill_id"`
	GainingPoint GainingPoint `db:"gaining_point"`
}

type FetchSkillGrowthData func(context.Context, ActionId) ([]*SkillGrowthData, error)

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
	ExploreId  ActionId     `db:"explore_id"`
	SkillId    core.SkillId `db:"skill_id"`
	RequiredLv core.SkillLv `db:"skill_lv"`
}

type FetchRequiredSkillsFunc func(context.Context, []ActionId) ([]*RequiredSkill, error)

type RequiredSkillRow struct {
	ExploreId      ActionId
	RequiredSkills []RequiredSkill
}

type ConsumingItem struct {
	ExploreId       ActionId        `db:"explore_id"`
	ItemId          core.ItemId     `db:"item_id"`
	MaxCount        core.Count      `db:"max_count"`
	ConsumptionProb ConsumptionProb `db:"consumption_prob"`
}

type FetchConsumingItemFunc func(context.Context, []ActionId) ([]*ConsumingItem, error)

type StaminaReductionSkillPair struct {
	ExploreId ActionId     `db:"explore_id"`
	SkillId   core.SkillId `db:"skill_id"`
}

type FetchReductionStaminaSkillFunc func(context.Context, []ActionId) ([]*StaminaReductionSkillPair, error)

func ReductionStaminaSkillToMap(res []*StaminaReductionSkillPair) map[ActionId][]core.SkillId {
	result := make(map[ActionId][]core.SkillId)
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

type FetchEarningItemFunc func(context.Context, ActionId) ([]*EarningItem, error)

type GetExploreMasterRes struct {
	ExploreId            ActionId             `db:"explore_id"`
	DisplayName          core.DisplayName     `db:"display_name"`
	Description          core.Description     `db:"description"`
	ConsumingStamina     core.StaminaCost     `db:"consuming_stamina"`
	RequiredPayment      core.Cost            `db:"required_payment"`
	StaminaReducibleRate StaminaReducibleRate `db:"stamina_reducible_rate"`
}

type FetchExploreMasterFunc func(context.Context, []ActionId) ([]*GetExploreMasterRes, error)

type UserExplore struct {
	ExploreId   ActionId
	DisplayName core.DisplayName
	IsKnown     core.IsKnown
	IsPossible  core.IsPossible
}

type ExploreUserData struct {
	ExploreId ActionId     `db:"explore_id"`
	IsKnown   core.IsKnown `db:"is_known"`
}

type GetUserExploreFunc func(context.Context, core.UserId, []ActionId) ([]*ExploreUserData, error)

type GetActionsRes struct {
	UserId   core.UserId
	Explores []*ExploreUserData
}
