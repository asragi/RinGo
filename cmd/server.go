package main

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/initialize"
	"github.com/asragi/RinGo/server"
	"log"
)

func main() {
	handleError := func(err error) {
		log.Fatal(err.Error())
	}
	// TODO: secretKey should be stored in a secure place
	secretKey := auth.SecretHashKey("secret")
	constants := &initialize.Constants{
		InitialFund:        core.Fund(100000),
		InitialMaxStamina:  core.MaxStamina(6000),
		InitialPopularity:  shelf.ShopPopularity(0),
		UserIdChallengeNum: 3,
	}

	db, err := CreateDB()
	if err != nil {
		handleError(err)
		return
	}

	endpoints := initialize.CreateEndpoints(secretKey, constants, db.Exec, db.Query)
	serve, stopDB, err := server.SetUpServer(4444, endpoints)
	if err != nil {
		handleError(err)
		return
	}
	defer stopDB()

	err = serve()
	if err != nil {
		log.Printf("Http Server Error: %v", err)
		return
	}
}
