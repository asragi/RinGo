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

func itemIdToStringArray[T core.Stringer](ids []T) []string {
	result := make([]string, len(ids))
	for i, v := range ids {
		result[i] = v.ToString()
	}
	return result
}

func CreateGetItemMasterMySQL(db *sql.DB) func([]core.ItemId) ([]stage.GetItemMasterRes, error) {
	f := func(itemIds []core.ItemId) ([]stage.GetItemMasterRes, error) {
		handleError := func(err error) ([]stage.GetItemMasterRes, error) {
			return nil, fmt.Errorf("get item master from mysql: %w", err)
		}
		query := "SELECT item_id, price, display_name, description, max_stock from item_masters"
		if len(itemIds) == 1 {
			itemId := itemIds[0]
			var res stage.GetItemMasterRes
			queryString := fmt.Sprintf("%s WHERE item_id = %s;", query, itemId)
			row := db.QueryRow(queryString)
			err := row.Scan(&res.ItemId, &res.Price, &res.DisplayName, &res.Description, &res.MaxStock)
			if err != nil {
				return handleError(err)
			}
			return []stage.GetItemMasterRes{res}, nil
		}
		listStatement := createListStatement(itemIdToStringArray(itemIds))
		queryString := fmt.Sprintf("%s WHERE item_id IN %s;", query, listStatement)
		rows, err := db.Query(queryString)
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()
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
