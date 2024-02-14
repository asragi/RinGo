package infrastructure

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type ConnectionSettings struct {
	UserName string
	Password string
	Port     string
	Protocol string
	Host     string
	Database string
}

func ConnectDB(settings *ConnectionSettings) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		settings.UserName,
		settings.Password,
		settings.Protocol,
		settings.Host,
		settings.Port,
		settings.Database,
	)
	db, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createListStatement(keywords []string) string {
	result := "("
	for i, v := range keywords {
		if len(keywords) == (i + 1) {
			result = fmt.Sprintf("%s%s)", result, v)
			break
		}
		result = fmt.Sprintf("%s%s, ", result, v)
	}
	return result
}

func idToStringArray[T core.Stringer](ids []T) []string {
	result := make([]string, len(ids))
	for i, v := range ids {
		result[i] = v.ToString()
	}
	return result
}

type ConnectDBFunc func() (*sqlx.DB, error)

func CreateGetResourceMySQL(connect ConnectDBFunc) stage.GetResourceFunc {
	return CreateGetUserDataQueryRow[stage.GetResourceRes](
		connect,
		"get resource: %w",
		"SELECT user_id, max_stamina, stamina_recover_time, fund FROM users",
	)
}

func CreateUpdateStamina(connect ConnectDBFunc) stage.UpdateStaminaFunc {
	return CreateUpdateUserData[core.StaminaRecoverTime](
		connect,
		"update stamimna: %w",
		`UPDATE users SET stamina_recover_time = %s WHERE user_id = "%s";`,
		func(stamina core.StaminaRecoverTime) string {
			t := time.Time(stamina)
			return t.Format("2006-01-02 15:04:05")
		},
	)
}

func CreateGetItemMasterMySQL(connect ConnectDBFunc) stage.FetchItemMasterFunc {
	return CreateBatchGetQuery[core.ItemId, stage.GetItemMasterRes](
		connect,
		"get item master from mysql: %w",
		"SELECT item_id, price, display_name, description, max_stock from item_masters",
		"item_id",
	)
}

func CreateGetStageMaster(connect ConnectDBFunc) stage.FetchStageMasterFunc {
	f := CreateBatchGetQuery[stage.StageId, stage.StageMaster](
		connect,
		"get stage master from mysql: %w",
		"SELECT stage_id, display_name, description from stage_masters",
		"stage_id",
	)
	return func(s stage.StageId) (stage.StageMaster, error) {
		res, err := f([]stage.StageId{s})
		return res[0], err
	}
}

func CreateGetAllStageMaster(connect ConnectDBFunc) stage.FetchAllStageFunc {
	f := func() (stage.GetAllStagesRes, error) {
		handleError := func(err error) (stage.GetAllStagesRes, error) {
			return stage.GetAllStagesRes{}, fmt.Errorf("get all stage master from mysql: %w", err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := "SELECT stage_id, display_name, description from stage_masters;"
		rows, err := db.Query(query)
		if err != nil {
			return handleError(err)
		}
		var result []stage.StageMaster
		for rows.Next() {
			var res stage.StageMaster
			err = rows.Scan(&res.StageId, &res.DisplayName, &res.Description)
			if err != nil {
				return handleError(err)
			}
			result = append(result, res)
		}
		return stage.GetAllStagesRes{Stages: result}, nil

	}

	return f
}

func CreateGetExploreMasterMySQL(connect ConnectDBFunc) stage.FetchExploreMasterFunc {
	return CreateBatchGetQuery[stage.ExploreId, stage.GetExploreMasterRes](
		connect,
		"get explore master from mysql: %w",
		"SELECT explore_id, display_name, description, consuming_stamina, required_payment, stamina_reducible_rate from explore_masters",
		"explore_id",
	)
}

func CreateGetSkillMaster(connect ConnectDBFunc) stage.FetchSkillMasterFunc {
	return CreateBatchGetQuery[core.SkillId, stage.SkillMaster](
		connect,
		"get explore master from mysql: %w",
		"SELECT skill_id, display_name from skill_masters",
		"explore_id",
	)
}

func CreateGetEarningItem(connect ConnectDBFunc) stage.FetchEarningItemFunc {
	f := CreateBatchGetQuery[stage.ExploreId, stage.EarningItem](
		connect,
		"get earning item data from mysql: %w",
		"SELECT item_id, min_count, max_count, probability from earning_items",
		"explore_id",
	)

	return func(id stage.ExploreId) ([]stage.EarningItem, error) {
		return f([]stage.ExploreId{id})
	}
}

func CreateGetConsumingItem(connect ConnectDBFunc) stage.FetchConsumingItemFunc {
	f := CreateBatchGetMultiQuery[stage.ExploreId, stage.ConsumingItem, stage.BatchGetConsumingItemRes](
		connect,
		"get consuming item data from mysql: %w",
		"SELECT explore_id, item_id, max_count, consumption_prob from consuming_items",
		"explore_id",
	)

	return f
}

func CreateGetRequiredSkills(connect ConnectDBFunc) stage.FetchRequiredSkillsFunc {
	f := CreateBatchGetMultiQuery[stage.ExploreId, stage.RequiredSkill, stage.RequiredSkillRow](
		connect,
		"get required skill from mysql :%w",
		"SELECT explore_id, skill_id, skill_lv from required_skills",
		"explore_id",
	)
	return f
}

func CreateGetSkillGrowth(connect ConnectDBFunc) stage.FetchSkillGrowthData {
	f := CreateGetQuery[stage.ExploreId, stage.SkillGrowthData](
		connect,
		"get skill growth from mysql: %w",
		"SELECT explore_id, skill_id, gaining_point FROM skill_growth_data",
		"explore_id",
	)

	return f
}

func CreateGetReductionSkill(connect ConnectDBFunc) stage.FetchReductionStaminaSkillFunc {
	f := CreateBatchGetMultiQuery[stage.ExploreId, stage.StaminaReductionSkillPair, stage.BatchGetReductionStaminaSkill](
		connect,
		"get skill growth from mysql: %w",
		"SELECT explore_id, skill_id, gaining_point FROM skill_growth_data",
		"explore_id",
	)

	return f
}

func CreateStageExploreRelation(connect ConnectDBFunc) stage.FetchStageExploreRelation {
	f := CreateBatchGetMultiQuery[stage.StageId, stage.StageExploreIdPairRow, stage.StageExploreIdPair](
		connect,
		"get stage explore relation from mysql: %w",
		"SELECT explore_id, stage_id FROM stage_explore_relations",
		"stage_id",
	)

	return f
}

func CreateItemExploreRelation(connect ConnectDBFunc) stage.FetchItemExploreRelationFunc {
	f := CreateGetQuery[core.ItemId, stage.ExploreId](
		connect,
		"get item explore relation from mysql: %w",
		"SELECT item_id, explore_id FROM item_explore_relations",
		"item_id",
	)

	return f
}

func CreateGetUserExplore(connect ConnectDBFunc) stage.GetUserExploreFunc {
	f := CreateBatchGetUserDataQuery[stage.ExploreId, stage.ExploreUserData](
		connect,
		"get user explore data: %w",
		"SELECT explore_id, is_known FROM user_explore_data",
		"explore_id",
	)

	return f
}

func CreateGetUserStageData(connect ConnectDBFunc) stage.FetchUserStageFunc {
	f := CreateBatchGetUserDataQuery[stage.StageId, stage.UserStage](
		connect,
		"get user stage data: %w",
		"SELECT stage_id, is_known FROM user_stage_data",
		"stage_id",
	)

	return f
}

func CreateGetQuery[S core.Stringer, T any](
	connect ConnectDBFunc,
	errorMessageFormat string,
	queryBase string,
	columnName string,
) func(S) ([]T, error) {
	f := func(id S) ([]T, error) {
		handleError := func(err error) ([]T, error) {
			return nil, fmt.Errorf(errorMessageFormat, err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := queryBase
		queryString := fmt.Sprintf("%s WHERE %s = %s;", query, columnName, id)
		rows, err := db.Queryx(queryString)
		var result []T
		for rows.Next() {
			var row T
			err = rows.StructScan(&row)
			if err != nil {
				return handleError(err)
			}
			result = append(result, row)
		}
		if err != nil {
			return handleError(err)
		}
		return result, nil
	}

	return f
}

// CreateBatchGetQuery returns function that receives N args and returns N values
func CreateBatchGetQuery[S core.Stringer, T any](
	connect ConnectDBFunc,
	errorMessageFormat string,
	queryBase string,
	columnName string,
) func([]S) ([]T, error) {
	f := func(ids []S) ([]T, error) {
		handleError := func(err error) ([]T, error) {
			return nil, fmt.Errorf(errorMessageFormat, err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := queryBase
		listStatement := createListStatement(idToStringArray(ids))
		queryString := fmt.Sprintf("%s WHERE %s IN %s;", query, columnName, listStatement)
		rows, err := db.Queryx(queryString)
		var result []T
		for rows.Next() {
			var row T
			err = rows.StructScan(&row)
			if err != nil {
				return handleError(err)
			}
			result = append(result, row)
		}
		if err != nil {
			return handleError(err)
		}
		return result, nil
	}

	return f
}

// CreateBatchGetMultiQuery returns function that receives N args and returns N values which have M values as array.
func CreateBatchGetMultiQuery[S core.Stringer, T core.ProvideId[S], U core.MultiResponseReceiver[S, T, U]](
	connect ConnectDBFunc,
	errorMessageFormat string,
	queryBase string,
	columnName string,
) func([]S) ([]U, error) {
	f := func(ids []S) ([]U, error) {
		handleError := func(err error) ([]U, error) {
			return nil, fmt.Errorf(errorMessageFormat, err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := queryBase
		listStatement := createListStatement(idToStringArray(ids))
		queryString := fmt.Sprintf("%s WHERE %s IN %s;", query, columnName, listStatement)
		rows, err := db.Queryx(queryString)
		if err != nil {
			return handleError(err)
		}
		var sqlResponse []T
		for rows.Next() {
			var row T
			err = rows.StructScan(&row)
			if err != nil {
				return handleError(err)
			}
			sqlResponse = append(sqlResponse, row)
		}
		mapping := make(map[string][]T)
		for _, v := range sqlResponse {
			id := v.GetId()
			idStr := id.ToString()
			if _, ok := mapping[idStr]; !ok {
				mapping[idStr] = []T{}
			}
			mapping[idStr] = append(mapping[idStr], v)
		}
		result := make([]U, len(ids))
		for i, v := range ids {
			arr, ok := mapping[v.ToString()]
			if !ok {
				arr = []T{}
			}
			result[i] = result[i].CreateSelf(v, arr)
		}

		return result, nil
	}

	return f
}

// CreateGetUserDataQueryRow returns function that receives N args and returns N values
func CreateGetUserDataQueryRow[T any](
	connect ConnectDBFunc,
	errorMessageFormat string,
	query string,
) func(core.UserId) (T, error) {
	f := func(userId core.UserId) (T, error) {
		handleError := func(err error) (T, error) {
			var empty T
			return empty, fmt.Errorf(errorMessageFormat, err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		queryString := fmt.Sprintf(
			`%s WHERE user_id = "%s";`,
			query,
			userId,
		)
		row := db.QueryRowx(queryString)
		var result T
		err = row.StructScan(&result)
		if err != nil {
			return handleError(err)
		}
		return result, nil
	}

	return f
}

// CreateBatchGetUserDataQuery returns function that receives N args and returns N values
func CreateBatchGetUserDataQuery[S core.Stringer, T any](
	connect ConnectDBFunc,
	errorMessageFormat string,
	queryBase string,
	columnName string,
) func(core.UserId, []S) ([]T, error) {
	f := func(userId core.UserId, ids []S) ([]T, error) {
		handleError := func(err error) ([]T, error) {
			return nil, fmt.Errorf(errorMessageFormat, err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := queryBase
		listStatement := createListStatement(idToStringArray(ids))
		queryString := fmt.Sprintf(
			`%s WHERE user_id = "%s" AND %s IN %s;`,
			query,
			userId,
			columnName,
			listStatement,
		)
		rows, err := db.Queryx(queryString)
		var result []T
		for rows.Next() {
			var row T
			err = rows.StructScan(&row)
			if err != nil {
				return handleError(err)
			}
			result = append(result, row)
		}
		if err != nil {
			return handleError(err)
		}
		return result, nil
	}

	return f
}

func CreateUpdateFund(connect ConnectDBFunc) stage.UpdateFundFunc {
	return CreateUpdateUserData[core.Fund](
		connect,
		"insert user fund: %w",
		`UPDATE users SET fund = %s WHERE user_id = "%s";`,
		func(f core.Fund) string { return strconv.Itoa(int(f)) },
	)
}

func CreateGetStorage(connect ConnectDBFunc) stage.FetchStorageFunc {
	g := CreateBatchGetUserDataQuery[core.ItemId, stage.ItemData](
		connect,
		"get user storage: %w",
		"SELECT user_id, item_id, stock, is_known FROM item_storages",
		"item_id",
	)
	return func(userId core.UserId, itemId []core.ItemId) (stage.BatchGetStorageRes, error) {
		res, err := g(userId, itemId)
		return stage.BatchGetStorageRes{
			UserId:   userId,
			ItemData: res,
		}, err
	}
}

func CreateGetAllStorage(connect ConnectDBFunc) stage.FetchAllStorageFunc {
	f := func(userId core.UserId) ([]stage.ItemData, error) {
		handleError := func(err error) ([]stage.ItemData, error) {
			return nil, fmt.Errorf("get all storage from mysql: %w", err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := fmt.Sprintf(
			`SELECT user_id, item_id, stock, is_known from item_storages WHERE user_id = "%s";`,
			userId,
		)
		rows, err := db.Query(query)
		if err != nil {
			return handleError(err)
		}
		var result []stage.ItemData
		for rows.Next() {
			var res stage.ItemData
			err = rows.Scan(&res.UserId, &res.ItemId, &res.Stock, &res.IsKnown)
			if err != nil {
				return handleError(err)
			}
			result = append(result, res)
		}
		return result, nil
	}

	return f
}

func CreateUpdateItemStorage(connect ConnectDBFunc) stage.UpdateItemStorageFunc {
	return CreateBulkUpdateUserData[stage.ItemStock](
		connect,
		"update user item storage",
		`INSERT INTO item_storages (user_id, item_id, stock) VALUES (:user_id, :item_id, :stock_exp) ON DUPLICATE KEY UPDATE stock =VALUES(stock);`,
	)
}

func CreateGetUserSkill(connect ConnectDBFunc) stage.FetchUserSkillFunc {
	g := CreateBatchGetUserDataQuery[core.SkillId, stage.UserSkillRes](
		connect,
		"get user skill data :%w",
		"SELECT user_id, skill_id, skill_exp FROM user_skills",
		"skill_id",
	)
	f := func(userId core.UserId, skillIds []core.SkillId) (stage.BatchGetUserSkillRes, error) {
		res, err := g(userId, skillIds)
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

func CreateUpdateUserSkill(connect ConnectDBFunc) stage.UpdateUserSkillExpFunc {
	g := CreateBulkUpdateUserData[stage.SkillGrowthPostRow]
	f := func(growthData stage.SkillGrowthPost) error {
		return g(
			connect,
			"update skill growth: %w",
			`INSERT INTO user_skills (user_id, skill_id, skill_exp) VALUES (:user_id, :skill_id, :skill_exp) ON DUPLICATE KEY UPDATE skill_exp =VALUES(skill_exp);`,
		)(
			growthData.UserId,
			growthData.SkillGrowth,
		)
	}

	return f
}

func CreateUpdateUserData[S any](
	connect ConnectDBFunc,
	errorMessageFormat string,
	queryFormat string,
	dataToString func(S) string,
) func(core.UserId, S) error {
	f := func(userId core.UserId, data S) error {
		handleError := func(err error) error {
			return fmt.Errorf(errorMessageFormat, err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		queryString := fmt.Sprintf(queryFormat, dataToString(data), userId)
		_, err = db.Exec(queryString)
		if err != nil {
			return handleError(err)
		}
		return nil
	}

	return f
}

func CreateBulkUpdateUserData[S any](
	connect ConnectDBFunc,
	errorMessageFormat string,
	query string,
) func(core.UserId, []S) error {
	f := func(userId core.UserId, data []S) error {
		handleError := func(err error) error {
			return fmt.Errorf(errorMessageFormat, err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}

		_, err = db.NamedExec(query, data)
		if err != nil {
			return handleError(err)
		}
		return nil
	}

	return f
}
