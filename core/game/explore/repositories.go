package explore

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

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
	StageId   StageId       `db:"stage_id"`
	ExploreId game.ActionId `db:"explore_id"`
}

func (row StageExploreIdPairRow) GetId() StageId {
	return row.StageId
}

type FetchItemExploreRelationFunc func(context.Context, core.ItemId) ([]game.ActionId, error)

type FetchStageExploreRelation func(context.Context, []StageId) ([]*StageExploreIdPairRow, error)

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
