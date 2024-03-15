package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/stage"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type (
	reqInterface[S any, T any] interface {
		Create(S) *T
	}
	exploreReq struct {
		ExploreId stage.ExploreId `db:"explore_id"`
	}
	stageReq struct {
		StageId stage.StageId `db:"stage_id"`
	}
	itemReq struct {
		ItemId core.ItemId `db:"item_id"`
	}
	skillReq struct {
		SkillId core.SkillId `db:"skill_id"`
	}
	queryFunc database.QueryFunc
)

func (exploreReq) Create(v stage.ExploreId) *exploreReq {
	return &exploreReq{ExploreId: v}
}

func (stageReq) Create(v stage.StageId) *stageReq {
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

func CreateGetResourceMySQL(q queryFunc) stage.GetResourceFunc {
	type responseStruct struct {
		UserId             core.UserId     `db:"user_id"`
		MaxStamina         core.MaxStamina `db:"max_stamina"`
		StaminaRecoverTime time.Time       `db:"stamina_recover_time"`
		Fund               core.Fund       `db:"fund"`
	}

	return func(ctx context.Context, userId core.UserId) (*stage.GetResourceRes, error) {
		handleError := func(err error) (*stage.GetResourceRes, error) {
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
		return &stage.GetResourceRes{
			UserId:             result.UserId,
			MaxStamina:         result.MaxStamina,
			StaminaRecoverTime: core.StaminaRecoverTime(result.StaminaRecoverTime),
			Fund:               result.Fund,
		}, err
	}
}

func CreateUpdateStamina(execDb database.DBExecFunc) stage.UpdateStaminaFunc {
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

func CreateGetItemMasterMySQL(q queryFunc) stage.FetchItemMasterFunc {
	return CreateGetQueryFromReq[core.ItemId, itemReq, stage.GetItemMasterRes](
		q,
		"get item master from mysql: %w",
		"SELECT item_id, price, display_name, description, max_stock from item_masters WHERE item_id IN (:item_id);",
	)
}

func CreateGetStageMaster(q queryFunc) stage.FetchStageMasterFunc {
	return CreateGetQueryFromReq[stage.StageId, stageReq, stage.StageMaster](
		q,
		"get stage master: %w",
		"SELECT stage_id, display_name, description from stage_masters WHERE stage_id IN (:stage_id);",
	)
}

func CreateGetAllStageMaster(q queryFunc) stage.FetchAllStageFunc {
	f := func(ctx context.Context) ([]*stage.StageMaster, error) {
		handleError := func(err error) ([]*stage.StageMaster, error) {
			return nil, fmt.Errorf("get all stage master from mysql: %w", err)
		}
		query := "SELECT stage_id, display_name, description from stage_masters;"
		rows, err := q(ctx, query, nil)
		if err != nil {
			return handleError(err)
		}
		var result []*stage.StageMaster
		for rows.Next() {
			var res stage.StageMaster
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

func CreateGetExploreMasterMySQL(q queryFunc) stage.FetchExploreMasterFunc {
	f := CreateGetQuery[exploreReq, stage.GetExploreMasterRes](
		q,
		"get explore master from mysql: %w",
		"SELECT explore_id, display_name, description, consuming_stamina, required_payment, stamina_reducible_rate from explore_masters WHERE explore_id IN (:explore_id);",
	)

	return func(ctx context.Context, ids []stage.ExploreId) ([]*stage.GetExploreMasterRes, error) {
		req := func(ids []stage.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateGetSkillMaster(q queryFunc) stage.FetchSkillMasterFunc {
	f := CreateGetQuery[skillReq, stage.SkillMaster](
		q,
		"get skill master from mysql: %w",
		"SELECT skill_id, display_name from skill_masters WHERE skill_id IN (:skill_id);",
	)
	return func(ctx context.Context, ids []core.SkillId) ([]*stage.SkillMaster, error) {
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

func CreateGetEarningItem(q queryFunc) stage.FetchEarningItemFunc {
	f := CreateGetQuery[exploreReq, stage.EarningItem](
		q,
		"get earning item data from mysql: %w",
		"SELECT item_id, min_count, max_count, probability from earning_items WHERE explore_id IN (:explore_id);",
	)

	return func(ctx context.Context, id stage.ExploreId) ([]*stage.EarningItem, error) {
		req := &exploreReq{ExploreId: id}
		return f(ctx, []*exploreReq{req})
	}
}

func CreateGetConsumingItem(q queryFunc) stage.FetchConsumingItemFunc {
	f := CreateGetQuery[exploreReq, stage.ConsumingItem](
		q,
		"get consuming item data from mysql: %w",
		"SELECT explore_id, item_id, max_count, consumption_prob from consuming_items",
	)

	return func(ctx context.Context, ids []stage.ExploreId) ([]*stage.ConsumingItem, error) {
		req := func(ids []stage.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateGetRequiredSkills(q queryFunc) stage.FetchRequiredSkillsFunc {
	f := CreateGetQuery[exploreReq, stage.RequiredSkill](
		q,
		"get required skill from mysql :%w",
		"SELECT explore_id, skill_id, skill_lv from required_skills",
	)
	return func(ctx context.Context, ids []stage.ExploreId) ([]*stage.RequiredSkill, error) {
		req := func(ids []stage.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateGetSkillGrowth(q queryFunc) stage.FetchSkillGrowthData {
	f := CreateGetQuery[exploreReq, stage.SkillGrowthData](
		q,
		"get skill growth from mysql: %w",
		`SELECT explore_id, skill_id, gaining_point FROM skill_growth_data WHERE explore_id IN (:explore_id);`,
	)

	return func(ctx context.Context, id stage.ExploreId) ([]*stage.SkillGrowthData, error) {
		req := &exploreReq{ExploreId: id}
		return f(ctx, []*exploreReq{req})
	}
}

func CreateGetReductionSkill(q queryFunc) stage.FetchReductionStaminaSkillFunc {
	f := CreateGetQuery[exploreReq, stage.StaminaReductionSkillPair](
		q,
		"get stamina reduction skill from mysql: %w",
		`SELECT explore_id, skill_id FROM stamina_reduction_skills WHERE explore_id IN (:explore_id);`,
	)

	return func(ctx context.Context, ids []stage.ExploreId) ([]*stage.StaminaReductionSkillPair, error) {
		req := func(ids []stage.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateStageExploreRelation(q queryFunc) stage.FetchStageExploreRelation {
	f := CreateGetQuery[stageReq, stage.StageExploreIdPairRow](
		q,
		"get stage explore relation from mysql: %w",
		"SELECT explore_id, stage_id FROM stage_explore_relations WHERE stage_id IN (:stage_id);",
	)

	return func(ctx context.Context, ids []stage.StageId) ([]*stage.StageExploreIdPairRow, error) {
		req := func(ids []stage.StageId) []*stageReq {
			result := make([]*stageReq, len(ids))
			for i, v := range ids {
				result[i] = &stageReq{StageId: v}
			}
			return result
		}(ids)
		return f(ctx, req)
	}
}

func CreateItemExploreRelation(q queryFunc) stage.FetchItemExploreRelationFunc {
	type fetchExploreIdRes struct {
		ExploreId stage.ExploreId `db:"explore_id"`
	}
	f := CreateGetQuery[itemReq, fetchExploreIdRes](
		q,
		"get item explore relation from mysql: %w",
		"SELECT explore_id FROM item_explore_relations WHERE item_id IN (:item_id);",
	)

	return func(ctx context.Context, id core.ItemId) ([]stage.ExploreId, error) {
		req := &itemReq{ItemId: id}
		res, err := f(ctx, []*itemReq{req})
		if err != nil {
			return nil, err
		}
		result := make([]stage.ExploreId, len(res))
		for i, v := range res {
			result[i] = v.ExploreId
		}
		return result, nil
	}
}

func CreateGetUserExplore(q queryFunc) stage.GetUserExploreFunc {
	f := CreateUserQuery[exploreReq, stage.ExploreUserData](
		q,
		"get user explore data: %w",
		createQueryFromUserId(`SELECT explore_id, is_known FROM user_explore_data WHERE user_id = "%s" AND explore_id IN (:explore_id);`),
	)

	return func(ctx context.Context, id core.UserId, ids []stage.ExploreId) ([]*stage.ExploreUserData, error) {
		req := func(ids []stage.ExploreId) []*exploreReq {
			result := make([]*exploreReq, len(ids))
			for i, v := range ids {
				result[i] = &exploreReq{ExploreId: v}
			}
			return result
		}(ids)
		return f(ctx, id, req)
	}
}

func CreateGetUserStageData(queryFunc queryFunc) stage.FetchUserStageFunc {
	f := CreateUserQuery[stageReq, stage.UserStage](
		queryFunc,
		"get user stage data: %w",
		createQueryFromUserId(`SELECT stage_id, is_known FROM user_stage_data WHERE user_id = '%s' AND stage_id IN (:stage_id);`),
	)

	return func(ctx context.Context, id core.UserId, ids []stage.StageId) ([]*stage.UserStage, error) {
		req := func(ids []stage.StageId) []*stageReq {
			result := make([]*stageReq, len(ids))
			for i, v := range ids {
				result[i] = &stageReq{StageId: v}
			}
			return result
		}(ids)
		return f(ctx, id, req)
	}
}

func CreateUpdateFund(dbExec database.DBExecFunc) stage.UpdateFundFunc {
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

func CreateGetStorage(queryFunc queryFunc) stage.FetchStorageFunc {
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
	return func(ctx context.Context, userId core.UserId, itemId []core.ItemId) (stage.BatchGetStorageRes, error) {
		if len(itemId) <= 0 {
			return stage.BatchGetStorageRes{}, nil
		}
		req := &itemIdReq{ItemId: itemId[0]}
		res, err := g(ctx, userId, []*itemIdReq{req})
		result := func() []*stage.StorageData {
			r := make([]*stage.StorageData, len(res))
			for i, v := range res {
				r[i] = &stage.StorageData{
					UserId:  v.UserId,
					ItemId:  v.ItemId,
					Stock:   v.Stock,
					IsKnown: core.IsKnown(v.IsKnown),
				}
			}
			return r
		}()
		return stage.BatchGetStorageRes{
			UserId:   userId,
			ItemData: result,
		}, err
	}
}

func CreateGetAllStorage(queryFunc queryFunc) stage.FetchAllStorageFunc {
	return func(ctx context.Context, userId core.UserId) ([]*stage.StorageData, error) {
		handleError := func(err error) ([]*stage.StorageData, error) {
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
		var result []*stage.StorageData
		for rows.Next() {
			var res stage.StorageData
			err = rows.Scan(&res.UserId, &res.ItemId, &res.Stock, &res.IsKnown)
			if err != nil {
				return handleError(err)
			}
			result = append(result, &res)
		}
		if result == nil || len(result) == 0 {
			return []*stage.StorageData{}, sql.ErrNoRows
		}
		return result, nil
	}
}

func CreateUpdateItemStorage(dbExec database.DBExecFunc) stage.UpdateItemStorageFunc {
	return func(ctx context.Context, userId core.UserId, stocks []*stage.ItemStock) error {
		type userItemStock struct {
			UserId  core.UserId  `db:"user_id"`
			ItemId  core.ItemId  `db:"item_id"`
			Stock   core.Stock   `db:"stock"`
			IsKnown core.IsKnown `db:"is_known"`
		}

		stockData := func(stocks []*stage.ItemStock) []*userItemStock {
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

func CreateGetUserSkill(dbExec queryFunc) stage.FetchUserSkillFunc {
	type skillReq struct {
		skillId string `db:"skill_id"`
	}
	queryFromUserId := createQueryFromUserId(
		`SELECT user_id, skill_id, skill_exp FROM user_skills WHERE user_id = "%s" AND skill_id IN (:skill_id);`,
	)
	g := CreateUserQuery[skillReq, stage.UserSkillRes](
		dbExec,
		"get user skill data :%w",
		queryFromUserId,
	)
	f := func(ctx context.Context, userId core.UserId, skillIds []core.SkillId) (stage.BatchGetUserSkillRes, error) {
		skillReqStructs := func(ids []core.SkillId) []*skillReq {
			result := make([]*skillReq, len(ids))
			for i, v := range ids {
				result[i] = &skillReq{skillId: v.ToString()}
			}
			return result
		}(skillIds)
		res, err := g(ctx, userId, skillReqStructs)
		if err != nil {
			return stage.BatchGetUserSkillRes{}, err
		}
		return stage.BatchGetUserSkillRes{
			Skills: res,
			UserId: userId,
		}, nil
	}
	return f
}

func CreateUpdateUserSkill(dbExec database.DBExecFunc) stage.UpdateUserSkillExpFunc {
	g := CreateExec[stage.SkillGrowthPostRow]
	f := func(ctx context.Context, growthData stage.SkillGrowthPost) error {
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
