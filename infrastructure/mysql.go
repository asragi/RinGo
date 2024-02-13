package infrastructure

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
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

func CreateUpdateFund(connect ConnectDBFunc) stage.UpdateFundFunc {
	return CreateUpdateUserData[core.Fund](
		connect,
		"insert user fund: %w",
		`UPDATE users SET fund = %s WHERE user_id = "%s";`,
		func(f core.Fund) string { return strconv.Itoa(int(f)) },
	)
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
