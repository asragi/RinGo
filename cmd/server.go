package main

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/handler"
	"github.com/asragi/RinGo/server"
	"log"
)

func main() {
	handleError := func(err error) {
		log.Fatal(err.Error())
	}
	constants := &server.Constants{
		InitialFund:       core.Fund(100000),
		InitialMaxStamina: core.MaxStamina(6000),
		InitialPopularity: game.ShopPopularity(0),
	}

	writeLogger := handler.LogHttpWrite

	err, closeDB, serve := server.InitializeServer(constants, writeLogger)
	defer func() {
		if e := closeDB(); e != nil {
			handleError(e)
		}
	}()
	if err != nil {
		handleError(err)
		return
	}

	err = serve()
	if err != nil {
		log.Printf("Http Server Error: %v", err)
		return
	}
}
