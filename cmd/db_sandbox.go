package main

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/database"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	handleError := func(err error) {
		fmt.Printf("error: %s", err.Error())
	}
	dbSettings := &database.ConnectionSettings{
		UserName: "root",
		Password: "ringo",
		Port:     "13306",
		Protocol: "tcp",
		Host:     "127.0.0.1",
		Database: "ringo",
	}
	db, err := database.ConnectDB(dbSettings)
	if err != nil {
		handleError(err)
		return
	}
	acc := database.NewDBAccessor(db, db)
	type res struct {
		DisplayName string `db:"display_name"`
	}
	type req struct {
		SkillId core.SkillId `db:"skill_id"`
	}
	rows, err := acc.Query(
		context.Background(),
		"SELECT display_name FROM skill_masters WHERE skill_id IN (:skill_id);",
		[]*req{{SkillId: "1"}},
	)
	if err != nil {
		handleError(err)
		return
	}
	var result []*res
	for rows.Next() {
		var row res
		err = rows.StructScan(&row)
		if err != nil {
			handleError(err)
			return
		}
		result = append(result, &row)
	}
}
