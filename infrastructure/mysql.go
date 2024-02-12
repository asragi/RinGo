package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	_ "github.com/go-sql-driver/mysql"
)

type ConnectionSettings struct {
	UserName string
	Password string
	Port     string
	Protocol string
	Host     string
	Database string
}

func ConnectDB(settings *ConnectionSettings) (*sql.DB, error) {
	dataSource := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		settings.UserName,
		settings.Password,
		settings.Protocol,
		settings.Host,
		settings.Port,
		settings.Database,
	)
	db, err := sql.Open("mysql", dataSource)
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

type ConnectDBFunc func() (*sql.DB, error)

func CreateGetItemMasterMySQL(connect ConnectDBFunc) stage.BatchGetItemMasterFunc {
	f := func(itemIds []core.ItemId) ([]stage.GetItemMasterRes, error) {
		handleError := func(err error) ([]stage.GetItemMasterRes, error) {
			return nil, fmt.Errorf("get item master from mysql: %w", err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := "SELECT item_id, price, display_name, description, max_stock from item_masters"
		listStatement := createListStatement(idToStringArray(itemIds))
		queryString := fmt.Sprintf("%s WHERE item_id IN %s;", query, listStatement)
		rows, err := db.Query(queryString)
		if err != nil {
			return handleError(err)
		}
		var result []stage.GetItemMasterRes
		for rows.Next() {
			u := &stage.GetItemMasterRes{}
			if err := rows.Scan(&u.ItemId, &u.Price, &u.DisplayName, &u.Description, &u.MaxStock); err != nil {
				return handleError(err)
			}
			result = append(result, *u)
		}
		if err = rows.Err(); err != nil {
			return handleError(err)
		}
		return result, nil
	}
	return f
}

func CreateGetStageMaster(connect ConnectDBFunc) stage.FetchStageMasterFunc {
	f := func(stageId stage.StageId) (stage.StageMaster, error) {
		handleError := func(err error) (stage.StageMaster, error) {
			return stage.StageMaster{}, fmt.Errorf("get stage master from mysql: %w", err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := "SELECT stage_id, display_name, description from stage_masters"
		var res stage.StageMaster
		queryString := fmt.Sprintf("%s WHERE stage_id = %s;", query, stageId)
		row := db.QueryRow(queryString)
		err = row.Scan(&res.StageId, &res.DisplayName, &res.Description)
		if err != nil {
			return handleError(err)
		}
		return res, nil
	}
	return f
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
	f := func(exploreIds []stage.ExploreId) ([]stage.GetExploreMasterRes, error) {
		handleError := func(err error) ([]stage.GetExploreMasterRes, error) {
			return nil, fmt.Errorf("get explore master from mysql: %w", err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := "SELECT explore_id, display_name, description, consuming_stamina, required_payment, stamina_reducible_rate from explore_masters"
		listStatement := createListStatement(idToStringArray(exploreIds))
		queryString := fmt.Sprintf("%s WHERE explore_id IN %s;", query, listStatement)
		rows, err := db.Query(queryString)
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()
		var result []stage.GetExploreMasterRes
		for rows.Next() {
			u := &stage.GetExploreMasterRes{}
			if err := rows.Scan(
				&u.ExploreId,
				&u.DisplayName,
				&u.Description,
				&u.ConsumingStamina,
				&u.RequiredPayment,
				&u.StaminaReducibleRate,
			); err != nil {
				return handleError(err)
			}
			result = append(result, *u)
		}
		if err = rows.Err(); err != nil {
			return handleError(err)
		}
		return result, nil
	}
	return f
}

func CreateGetSkillMaster(connect ConnectDBFunc) stage.FetchSkillMasterFunc {
	f := func(skillIds []core.SkillId) ([]stage.SkillMaster, error) {
		handleError := func(err error) ([]stage.SkillMaster, error) {
			return nil, fmt.Errorf("get explore master from mysql: %w", err)
		}
		db, err := connect()
		if err != nil {
			return handleError(err)
		}
		query := "SELECT skill_id, display_name from skill_masters"
		listStatement := createListStatement(idToStringArray(skillIds))
		queryString := fmt.Sprintf("%s WHERE explore_id IN %s;", query, listStatement)
		rows, err := db.Query(queryString)
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()
		var result []stage.SkillMaster
		for rows.Next() {
			u := &stage.SkillMaster{}
			if err := rows.Scan(
				&u.SkillId,
				&u.DisplayName,
			); err != nil {
				return handleError(err)
			}
			result = append(result, *u)
		}
		if err = rows.Err(); err != nil {
			return handleError(err)
		}
		return result, nil
	}
	return f
}

/*
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
		rows, err := db.Query(queryString)
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()
		var result []T
		for rows.Next() {
			var u *T
			if err := rows.Scan(
				&u.SkillId,
				&u.DisplayName,
			); err != nil {
				return handleError(err)
			}
			result = append(result, *u)
		}
		if err = rows.Err(); err != nil {
			return handleError(err)
		}
		return result, nil
	}

	return f
}
*/
