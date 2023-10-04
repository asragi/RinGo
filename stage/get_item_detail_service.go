package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type GetUserItemDetailReq struct {
	UserId      core.UserId
	ItemId      core.ItemId
	AccessToken core.AccessToken
}

type getUserItemDetailRes struct {
	UserId       core.UserId
	ItemId       core.ItemId
	Price        core.Price
	DisplayName  core.DisplayName
	Description  core.Description
	MaxStock     core.MaxStock
	Stock        core.Stock
	UserExplores []userExplore
}

type itemService struct {
	GetUserItemDetail func(GetUserItemDetailReq) (getUserItemDetailRes, error)
}

func CreateGetItemDetailService(
	makeUserExploreArray makeUserExploreArrayFunc,
	itemMasterRepo ItemMasterRepo,
	itemStorageRepo ItemStorageRepo,
	exploreMasterRepo ExploreMasterRepo,
	itemExploreRelationRepo ItemExploreRelationRepo,
) itemService {
	getAllAction := func(req GetUserItemDetailReq) ([]userExplore, error) {
		handleError := func(err error) ([]userExplore, error) {
			return []userExplore{}, fmt.Errorf("error on getAllAction: %w", err)
		}
		itemExploreIds, err := itemExploreRelationRepo.Get(req.ItemId)
		if err != nil {
			return handleError(err)
		}
		explores, err := exploreMasterRepo.BatchGet(itemExploreIds)
		if err != nil {
			return handleError(err)
		}
		exploreIds := make([]ExploreId, len(explores))
		for i, v := range explores {
			exploreIds[i] = v.ExploreId
		}
		exploreMap := make(map[ExploreId]GetExploreMasterRes)
		for _, v := range explores {
			exploreMap[v.ExploreId] = v
		}

		result, err := makeUserExploreArray(
			req.UserId,
			req.AccessToken,
			exploreIds,
			exploreMap,
		)

		if err != nil {
			return handleError(err)
		}

		return result, nil
	}

	getUserItemDetail := func(req GetUserItemDetailReq) (getUserItemDetailRes, error) {
		handleError := func(err error) (getUserItemDetailRes, error) {
			return getUserItemDetailRes{}, fmt.Errorf("error on getUserItemDetail: %w", err)
		}
		masterRes, err := itemMasterRepo.Get(req.ItemId)
		if err != nil {
			return handleError(err)
		}
		storageRes, err := itemStorageRepo.Get(req.UserId, req.ItemId, req.AccessToken)
		if err != nil {
			return handleError(err)
		}
		explores, err := getAllAction(req)
		if err != nil {
			return handleError(err)
		}
		return getUserItemDetailRes{
			UserId:       storageRes.UserId,
			ItemId:       masterRes.ItemId,
			Price:        masterRes.Price,
			DisplayName:  masterRes.DisplayName,
			Description:  masterRes.Description,
			MaxStock:     masterRes.MaxStock,
			Stock:        storageRes.Stock,
			UserExplores: explores,
		}, nil
	}

	return itemService{
		GetUserItemDetail: getUserItemDetail,
	}
}
