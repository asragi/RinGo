package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type (
	reqInterface[S any, T any] interface {
		Create(S) *T
	}
	exploreReq struct {
		ExploreId game.ExploreId `db:"explore_id"`
	}
	stageReq struct {
		StageId explore.StageId `db:"stage_id"`
	}
	itemReq struct {
		ItemId core.ItemId `db:"item_id"`
	}
	skillReq struct {
		SkillId core.SkillId `db:"skill_id"`
	}
	queryFunc database.QueryFunc
)

func (exploreReq) Create(v game.ExploreId) *exploreReq {
	return &exploreReq{ExploreId: v}
}

func (stageReq) Create(v explore.StageId) *stageReq {
	return &stageReq{StageId: v}
}

func (itemReq) Create(v core.ItemId) *itemReq {
	return &itemReq{ItemId: v}
}

func (skillReq) Create(v core.SkillId) *skillReq {
	return &skillReq{SkillId: v}
}

func CreateCheckUserExistence(queryFunc queryFunc) core.CheckDoesUserExist {
	return func(ctx context.Context, userId core.UserId) error {
		handleError := func(err error) error {
			return fmt.Errorf("check user existence: %w", err)
		}
		queryString := fmt.Sprintf(`SELECT user_id from users WHERE user_id = "%s";`, userId)
		_, err := queryFunc(ctx, queryString, nil)
		if err == nil {
			return handleError(&auth.UserAlreadyExistsError{UserId: string(userId)})
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return handleError(err)
		}
		return nil
	}
}

func CreateGetUserPassword(queryFunc queryFunc) auth.FetchHashedPassword {
	type dbResponse struct {
		HashedPassword auth.HashedPassword `db:"hashed_password"`
	}
	return func(ctx context.Context, id core.UserId) (auth.HashedPassword, error) {
		handleError := func(err error) (auth.HashedPassword, error) {
			return "", fmt.Errorf("get hashed password: %w", err)
		}
		queryString := fmt.Sprintf(`SELECT hashed_password FROM users WHERE user_id = "%s";`, id)
		rows, err := queryFunc(ctx, queryString, nil)
		if err != nil {
			return handleError(err)
		}
		var result dbResponse
		err = rows.StructScan(&result)
		if err != nil {
			return handleError(err)
		}
		return result.HashedPassword, nil
	}
}

func CreateInsertNewUser(
	dbExec database.DBExecFunc,
	initialFund core.Fund,
	initialMaxStamina core.MaxStamina,
	getTime core.GetCurrentTimeFunc,
) auth.InsertNewUser {
	return func(ctx context.Context, id core.UserId, userName core.UserName, password auth.HashedPassword) error {
		handleError := func(err error) error {
			return fmt.Errorf("insert new user: %w", err)
		}
		queryText := `INSERT INTO users (user_id, name, fund, max_stamina, stamina_recover_time, hashed_password) VALUES (:user_id, :name, :fund, :max_stamina, :stamina_recover_time, :hashed_password);`

		type UserToDB struct {
			UserId             core.UserId         `db:"user_id"`
			Name               core.UserName       `db:"name"`
			Fund               core.Fund           `db:"fund"`
			MaxStamina         core.MaxStamina     `db:"max_stamina"`
			StaminaRecoverTime time.Time           `db:"stamina_recover_time"`
			HashedPassword     auth.HashedPassword `db:"hashed_password"`
		}

		createUserData := UserToDB{
			UserId:             id,
			Name:               userName,
			Fund:               initialFund,
			MaxStamina:         initialMaxStamina,
			StaminaRecoverTime: getTime(),
			HashedPassword:     password,
		}

		_, err := dbExec(ctx, queryText, []*UserToDB{&createUserData})
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}

func CreateGetResourceMySQL(q queryFunc) game.GetResourceFunc {
	type responseStruct struct {
		UserId             core.UserId     `db:"user_id"`
		MaxStamina         core.MaxStamina `db:"max_stamina"`
		StaminaRecoverTime time.Time       `db:"stamina_recover_time"`
		Fund               core.Fund       `db:"fund"`
	}

	return func(ctx context.Context, userId core.UserId) (*game.GetResourceRes, error) {
		handleError := func(err error) (*game.GetResourceRes, error) {
			return nil, fmt.Errorf("get resource from mysql: %w", err)
		}
		rows, err := q(
			ctx,
			fmt.Sprintf(
				`SELECT user_id, max_stamina, stamina_recover_time, fund FROM users WHERE user_id = "%s";`,
				userId,
			),
			nil,
		)
		if err != nil {
			return handleError(err)
		}
		var result responseStruct
		err = rows.StructScan(&result)
		if err != nil {
			return handleError(err)
		}
		return &game.GetResourceRes{
			UserId:             result.UserId,
			MaxStamina:         result.MaxStamina,
			StaminaRecoverTime: core.StaminaRecoverTime(result.StaminaRecoverTime),
			Fund:               result.Fund,
		}, err
	}
}

func CreateUpdateStamina(execDb database.DBExecFunc) game.UpdateStaminaFunc {
	type updateStaminaReq struct {
		stamina core.StaminaRecoverTime `db:"stamina_recover_time"`
	}
	query := func(userId core.UserId) string {
		return fmt.Sprintf(`UPDATE users SET stamina_recover_time = ? WHERE user_id = "%s";`, userId)
	}
	return func(ctx context.Context, userId core.UserId, recoverTime core.StaminaRecoverTime) error {
		return CreateExec[updateStaminaReq](
			execDb,
			"update stamina: %w",
			query(userId),
		)(ctx, []*updateStaminaReq{{stamina: recoverTime}})
	}
}

func CreateGetItemMasterMySQL(q queryFunc) game.FetchItemMasterFunc {
	return CreateGetQueryFromReq[core.ItemId, itemReq, game.GetItemMasterRes](
		q,
		"get item master from mysql: %w",
		"SELECT item_id, price, display_name, description, max_stock from item_masters WHERE item_id IN (:item_id);",
	)
}

func CreateGetStageMaster(q queryFunc) explore.FetchStageMasterFunc {
	return CreateGetQueryFromReq[explore.StageId, stageReq, explore.StageMaster](
		q,
		"get stage master: %w",
		"SELECT stage_id, display_name, description from stage_masters WHERE stage_id IN (:stage_id);",
	)
}

func CreateGetAllStageMaster(q queryFunc) explore.FetchAllStageFunc {
	f := func(ctx context.Context) ([]*explore.StageMaster, error) {
		handleError := func(err error) ([]*explore.StageMaster, error) {
			return nil, fmt.Errorf("get all stage master from mysql: %w", err)
		}
		query := "SELECT stage_id, display_name, description from stage_masters;"
		rows, err := q(ctx, query, nil)
		if err != nil {
			return handleError(err)
		}
		var result []*explore.StageMaster
		for rows.Next() {
			var res explore.StageMaster
			err = rows.Scan(&res.StageId, &res.DisplayName, &res.Description)
			if err != nil {
				return handleError(err)
			}
			result = append(result, &res)
		}
		return result, nil
	}

	return f
}

func CreateGetQueryFromReq[S any, SReq reqInterface[S, SReq], T any](
	q queryFunc,
	errorMessageFormat string,
	queryText string,
) func(context.Context, []S) ([]*T, error) {
	f := CreateGetQuery[SReq, T](q, errorMessageFormat, queryText)
	return func(ctx context.Context, s []S) ([]*T, error) {
		req := func(s []S) []*SReq {
			result := make([]*SReq, len(s))
			for i, v := range s {
				var tmp SReq
				result[i] = tmp.Create(v)
			}
			return result
		}(s)
		return f(ctx, req)
	}
}

func CreateGetExploreMasterMySQL(q queryFunc) game.FetchExploreMasterFunc {
	f := CreateGetQuery[exploreReq, game.GetExploreMasterRes](
		q,
		"get explore master from mysql: %w",
		"SELECT explore_id, display_name, description, consuming_stamina, required_payment, stamina_reducible_rate from explore_masters WHERE explore_id IN (:explore_id);",
	)

	return func(ctx context.Context, ids []game.ExploreId) ([]*game.GetExploreMasterRes, error) {
		req := func(ids []game.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateGetSkillMaster(q queryFunc) game.FetchSkillMasterFunc {
	f := CreateGetQuery[skillReq, game.SkillMaster](
		q,
		"get skill master from mysql: %w",
		"SELECT skill_id, display_name from skill_masters WHERE skill_id IN (:skill_id);",
	)
	return func(ctx context.Context, ids []core.SkillId) ([]*game.SkillMaster, error) {
		req := func(ids []core.SkillId) []*skillReq {
			result := make([]*skillReq, len(ids))
			for i, v := range ids {
				result[i] = &skillReq{SkillId: v}
			}
			return result
		}(ids)
		res, err := f(ctx, req)
		return res, err
	}
}

func CreateGetEarningItem(q queryFunc) game.FetchEarningItemFunc {
	f := CreateGetQuery[exploreReq, game.EarningItem](
		q,
		"get earning item data from mysql: %w",
		"SELECT item_id, min_count, max_count, probability from earning_items WHERE explore_id IN (:explore_id);",
	)

	return func(ctx context.Context, id game.ExploreId) ([]*game.EarningItem, error) {
		req := &exploreReq{ExploreId: id}
		return f(ctx, []*exploreReq{req})
	}
}

func CreateGetConsumingItem(q queryFunc) game.FetchConsumingItemFunc {
	f := CreateGetQuery[exploreReq, game.ConsumingItem](
		q,
		"get consuming item data from mysql: %w",
		"SELECT explore_id, item_id, max_count, consumption_prob from consuming_items",
	)

	return func(ctx context.Context, ids []game.ExploreId) ([]*game.ConsumingItem, error) {
		req := func(ids []game.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateGetRequiredSkills(q queryFunc) game.FetchRequiredSkillsFunc {
	f := CreateGetQuery[exploreReq, game.RequiredSkill](
		q,
		"get required skill from mysql :%w",
		"SELECT explore_id, skill_id, skill_lv from required_skills",
	)
	return func(ctx context.Context, ids []game.ExploreId) ([]*game.RequiredSkill, error) {
		req := func(ids []game.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateGetSkillGrowth(q queryFunc) game.FetchSkillGrowthData {
	f := CreateGetQuery[exploreReq, game.SkillGrowthData](
		q,
		"get skill growth from mysql: %w",
		`SELECT explore_id, skill_id, gaining_point FROM skill_growth_data WHERE explore_id IN (:explore_id);`,
	)

	return func(ctx context.Context, id game.ExploreId) ([]*game.SkillGrowthData, error) {
		req := &exploreReq{ExploreId: id}
		return f(ctx, []*exploreReq{req})
	}
}

func CreateGetReductionSkill(q queryFunc) game.FetchReductionStaminaSkillFunc {
	f := CreateGetQuery[exploreReq, game.StaminaReductionSkillPair](
		q,
		"get stamina reduction skill from mysql: %w",
		`SELECT explore_id, skill_id FROM stamina_reduction_skills WHERE explore_id IN (:explore_id);`,
	)

	return func(ctx context.Context, ids []game.ExploreId) ([]*game.StaminaReductionSkillPair, error) {
		req := func(ids []game.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateStageExploreRelation(q queryFunc) explore.FetchStageExploreRelation {
	f := CreateGetQuery[stageReq, explore.StageExploreIdPairRow](
		q,
		"get stage explore relation from mysql: %w",
		"SELECT explore_id, stage_id FROM stage_explore_relations WHERE stage_id IN (:stage_id);",
	)

	return func(ctx context.Context, ids []explore.StageId) ([]*explore.StageExploreIdPairRow, error) {
		req := func(ids []explore.StageId) []*stageReq {
			result := make([]*stageReq, len(ids))
			for i, v := range ids {
				result[i] = &stageReq{StageId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateItemExploreRelation(q queryFunc) explore.FetchItemExploreRelationFunc {
	type fetchExploreIdRes struct {
		ExploreId game.ExploreId `db:"explore_id"`
	}
	f := CreateGetQuery[itemReq, fetchExploreIdRes](
		q,
		"get item explore relation from mysql: %w",
		"SELECT explore_id FROM item_explore_relations WHERE item_id IN (:item_id);",
	)

	return func(ctx context.Context, id core.ItemId) ([]game.ExploreId, error) {
		req := &itemReq{ItemId: id}
		res, err := f(ctx, []*itemReq{req})
		if err != nil {
			return nil, err
		}
		result := make([]game.ExploreId, len(res))
		for i, v := range res {
			result[i] = v.ExploreId
		}
		return result, nil
	}
}

func CreateGetUserExplore(q queryFunc) game.GetUserExploreFunc {
	f := CreateUserQuery[exploreReq, game.ExploreUserData](
		q,
		"get user explore data: %w",
		createQueryFromUserId(`SELECT explore_id, is_known FROM user_explore_data WHERE user_id = "%s" AND explore_id IN (:explore_id);`),
	)

	return func(ctx context.Context, id core.UserId, ids []game.ExploreId) ([]*game.ExploreUserData, error) {
		req := func(ids []game.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, id, req)
	}
}

func CreateGetUserStageData(queryFunc queryFunc) explore.FetchUserStageFunc {
	f := CreateUserQuery[stageReq, explore.UserStage](
		queryFunc,
		"get user stage data: %w",
		createQueryFromUserId(`SELECT stage_id, is_known FROM user_stage_data WHERE user_id = '%s' AND stage_id IN (:stage_id);`),
	)

	return func(ctx context.Context, id core.UserId, ids []explore.StageId) ([]*explore.UserStage, error) {
		req := func(ids []explore.StageId) []*stageReq {
			result := make([]*stageReq, len(ids))
			for i, v := range ids {
				result[i] = &stageReq{StageId: v}
			}
			return result
		}(ids)
		return f(ctx, id, req)
	}
}

func CreateUpdateFund(dbExec database.DBExecFunc) game.UpdateFundFunc {
	query := func(userId core.UserId) string {
		return fmt.Sprintf(`UPDATE users SET fund = ? WHERE user_id = "%s";`, userId)
	}
	return func(ctx context.Context, userId core.UserId, fund core.Fund) error {
		return CreateExec[core.Fund](
			dbExec,
			"insert user fund: %w",
			query(userId),
		)(ctx, []*core.Fund{&fund})
	}
}

func CreateGetStorage(queryFunc queryFunc) game.FetchStorageFunc {
	type ItemDataRes struct {
		UserId  core.UserId `db:"user_id"`
		ItemId  core.ItemId `db:"item_id"`
		Stock   core.Stock  `db:"stock"`
		IsKnown bool        `db:"is_known"`
	}
	type itemIdReq struct {
		ItemId core.ItemId `db:"item_id"`
	}
	g := CreateUserQuery[itemIdReq, ItemDataRes](
		queryFunc,
		"get user storage: %w",
		createQueryFromUserId(`SELECT user_id, item_id, stock, is_known FROM item_storages WHERE user_id = "%s" AND item_id IN (:item_id);`),
	)
	return func(ctx context.Context, userId core.UserId, itemId []core.ItemId) (game.BatchGetStorageRes, error) {
		if len(itemId) <= 0 {
			return game.BatchGetStorageRes{}, nil
		}
		req := &itemIdReq{ItemId: itemId[0]}
		res, err := g(ctx, userId, []*itemIdReq{req})
		result := func() []*game.StorageData {
			r := make([]*game.StorageData, len(res))
			for i, v := range res {
				r[i] = &game.StorageData{
					UserId:  v.UserId,
					ItemId:  v.ItemId,
					Stock:   v.Stock,
					IsKnown: core.IsKnown(v.IsKnown),
				}
			}
			return r
		}()
		return game.BatchGetStorageRes{
			UserId:   userId,
			ItemData: result,
		}, err
	}
}

func CreateGetAllStorage(queryFunc queryFunc) game.FetchAllStorageFunc {
	return func(ctx context.Context, userId core.UserId) ([]*game.StorageData, error) {
		handleError := func(err error) ([]*game.StorageData, error) {
			return nil, fmt.Errorf("get all storage from mysql: %w", err)
		}
		query := fmt.Sprintf(
			`SELECT user_id, item_id, stock, is_known from item_storages WHERE user_id = "%s";`,
			userId,
		)
		rows, err := queryFunc(ctx, query, nil)
		if err != nil {
			return handleError(err)
		}
		var result []*game.StorageData
		for rows.Next() {
			var res game.StorageData
			err = rows.Scan(&res.UserId, &res.ItemId, &res.Stock, &res.IsKnown)
			if err != nil {
				return handleError(err)
			}
			result = append(result, &res)
		}
		if result == nil || len(result) == 0 {
			return []*game.StorageData{}, sql.ErrNoRows
		}
		return result, nil
	}
}

func CreateUpdateItemStorage(dbExec database.DBExecFunc) game.UpdateItemStorageFunc {
	return func(ctx context.Context, userId core.UserId, stocks []*game.ItemStock) error {
		type userItemStock struct {
			UserId  core.UserId  `db:"user_id"`
			ItemId  core.ItemId  `db:"item_id"`
			Stock   core.Stock   `db:"stock"`
			IsKnown core.IsKnown `db:"is_known"`
		}

		stockData := func(stocks []*game.ItemStock) []*userItemStock {
			result := make([]*userItemStock, len(stocks))
			for i, v := range stocks {
				result[i] = &userItemStock{
					UserId:  userId,
					ItemId:  v.ItemId,
					Stock:   v.AfterStock,
					IsKnown: v.IsKnown,
				}
			}
			return result
		}(stocks)

		query := createQueryFromUserId(
			`INSERT INTO item_storages (user_id, item_id, stock, is_known) VALUES (:user_id, :item_id, :stock, :is_known) ON DUPLICATE KEY UPDATE stock =VALUES(stock);`,
		)(userId)

		return CreateExec[userItemStock](
			dbExec,
			"update item storage: %w",
			query,
		)(ctx, stockData)
	}
}

func CreateGetUserSkill(dbExec queryFunc) game.FetchUserSkillFunc {
	type skillReq struct {
		skillId string `db:"skill_id"`
	}
	queryFromUserId := createQueryFromUserId(
		`SELECT user_id, skill_id, skill_exp FROM user_skills WHERE user_id = "%s" AND skill_id IN (:skill_id);`,
	)
	g := CreateUserQuery[skillReq, game.UserSkillRes](
		dbExec,
		"get user skill data :%w",
		queryFromUserId,
	)
	f := func(ctx context.Context, userId core.UserId, skillIds []core.SkillId) (game.BatchGetUserSkillRes, error) {
		skillReqStructs := func(ids []core.SkillId) []*skillReq {
			result := make([]*skillReq, len(ids))
			for i, v := range ids {
				result[i] = &skillReq{skillId: v.ToString()}
			}
			return result
		}(skillIds)
		res, err := g(ctx, userId, skillReqStructs)
		if err != nil {
			return game.BatchGetUserSkillRes{}, err
		}
		return game.BatchGetUserSkillRes{
			Skills: res,
			UserId: userId,
		}, nil
	}
	return f
}

func CreateUpdateUserSkill(dbExec database.DBExecFunc) game.UpdateUserSkillExpFunc {
	g := CreateExec[game.SkillGrowthPostRow]
	f := func(ctx context.Context, growthData game.SkillGrowthPost) error {
		query := `INSERT INTO user_skills (user_id, skill_id, skill_exp) VALUES (:user_id, :skill_id, :skill_exp) ON DUPLICATE KEY UPDATE skill_exp =VALUES(skill_exp);`

		return g(
			dbExec,
			"update skill growth: %w",
			query,
		)(
			ctx,
			growthData.SkillGrowth,
		)
	}

	return f
}

func CreateGetQuery[S any, T any](
	queryFunc queryFunc,
	errorMessageFormat string,
	queryText string,
) func(context.Context, []*S) ([]*T, error) {
	f := func(ctx context.Context, ids []*S) ([]*T, error) {
		handleError := func(err error) ([]*T, error) {
			return nil, fmt.Errorf(errorMessageFormat, err)
		}
		if len(ids) <= 0 {
			return nil, nil
		}
		rows, err := queryFunc(ctx, queryText, ids)
		if err != nil {
			return handleError(err)
		}
		var result []*T
		for rows.Next() {
			var row T
			err = rows.StructScan(&row)
			if err != nil {
				return handleError(err)
			}
			result = append(result, &row)
		}
		return result, nil
	}
	return f
}

func createQueryFromUserId(queryText string) func(core.UserId) string {
	return func(userId core.UserId) string {
		return fmt.Sprintf(queryText, userId)
	}
}

func CreateUserQuery[S any, T any](
	queryFunc queryFunc,
	errorMessageFormat string,
	queryTextFromUserId func(core.UserId) string,
) func(context.Context, core.UserId, []*S) ([]*T, error) {
	f := func(ctx context.Context, userId core.UserId, ids []*S) ([]*T, error) {
		handleError := func(err error) ([]*T, error) {
			return nil, fmt.Errorf(errorMessageFormat, err)
		}
		if len(ids) <= 0 {
			return nil, nil
		}
		queryText := queryTextFromUserId(userId)
		rows, err := queryFunc(ctx, queryText, ids)
		if err != nil {
			return handleError(err)
		}
		var result []*T
		for rows.Next() {
			var row T
			err = rows.StructScan(&row)
			if err != nil {
				return handleError(err)
			}
			result = append(result, &row)
		}
		return result, nil
	}
	return f
}

func CreateExec[S any](
	dbExec database.DBExecFunc,
	errorMessageFormat string,
	query string,
) func(context.Context, []*S) error {
	return func(ctx context.Context, data []*S) error {
		handleError := func(err error) error {
			return fmt.Errorf(errorMessageFormat, err)
		}
		_, err := dbExec(ctx, query, data)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
