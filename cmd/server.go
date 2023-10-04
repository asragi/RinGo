package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asragi/RinGo/infrastructure"
)

type handler func(http.ResponseWriter, *http.Request)

func main() {
	handleError := func(err error) {
		log.Fatal(err)
	}
	itemMaster, err := infrastructure.CreateInMemoryItemMasterRepo(&infrastructure.ItemMasterLoader{Path: "./infrastructure/data/item-master.csv"})
	if err != nil {
		handleError(err)
	}
	fmt.Println(itemMaster)
	/*
		getCommonService := stage.CreateCommonGetActionDetail()
		getActionDetailService := stage.CreateGetStageActionDetailService(getCommonService, stageMasterRepo)
		getStageActionDetail := endpoint.CreateGetStageActionDetail(getActionDetailService.GetAction)

		getStageActionDetailHandler := func(w http.ResponseWriter, r *http.Request) {
			_, _ = getStageActionDetail()
			fmt.Println(w, "getStageAction")
		}

		http.HandleFunc("items", getStageActionDetailHandler)
	*/
	http.HandleFunc("/", hello)
	http.ListenAndServe(":4444", nil)
}

func getItem(w http.ResponseWriter, r *http.Request) {

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Kisaragi!")
}
