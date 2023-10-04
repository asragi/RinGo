package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/infrastructure"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type handler func(http.ResponseWriter, *http.Request)

func main() {
	handleError := func(err error) {
		log.Fatal(err)
	}
	dataDir := "./infrastructure/data/%s.csv"
	_, err := infrastructure.CreateInMemoryItemMasterRepo(&infrastructure.ItemMasterLoader{Path: fmt.Sprintf(dataDir, "item-master")})
	if err != nil {
		handleError(err)
	}
	itemStorage, err := infrastructure.CreateInMemoryItemStorageRepo(&infrastructure.ItemStorageLoader{Path: fmt.Sprintf(dataDir, "item-storage")})
	if err != nil {
		handleError(err)
	}
	userSkill, err := infrastructure.CreateInMemoryUserSkillRepo(&infrastructure.UserSkillLoader{Path: "./infrastructure/data/user-skill.csv"})
	if err != nil {
		handleError(err)
	}
	stageMaster, err := infrastructure.CreateInMemoryStageMasterRepo(&infrastructure.StageMasterLoader{Path: fmt.Sprintf(dataDir, "stage-master")})
	if err != nil {
		handleError(err)
	}
	skillMaster, err := infrastructure.CreateInMemorySkillMasterRepo(&infrastructure.SkillMasterLoader{Path: "./infrastructure/data/skill-master.csv"})
	if err != nil {
		handleError(err)
	}
	exploreMaster, err := infrastructure.CreateInMemoryExploreMasterRepo(&infrastructure.ExploreMasterLoader{Path: "./infrastructure/data/explore-master.csv"})
	if err != nil {
		handleError(err)
	}
	earningItem, err := infrastructure.CreateInMemoryEarningItemRepo(&infrastructure.EarningItemLoader{Path: fmt.Sprintf(dataDir, "earning-item")})
	if err != nil {
		handleError(err)
	}
	consumingItem, err := infrastructure.CreateInMemoryConsumingItemRepo(&infrastructure.ConsumingItemLoader{Path: "./infrastructure/data/consuming-item.csv"})
	if err != nil {
		handleError(err)
	}
	requiredSkill, err := infrastructure.CreateInMemoryRequiredSkillRepo(&infrastructure.RequiredSkillLoader{Path: fmt.Sprintf(dataDir, "required-skill")})
	if err != nil {
		handleError(err)
	}
	reductionSkill, err := infrastructure.CreateInMemoryReductionStaminaSkillRepo(&infrastructure.ReductionStaminaSkillLoader{Path: fmt.Sprintf(dataDir, "reduction-stamina-skill")})
	if err != nil {
		handleError(err)
	}
	calcStaminaService := stage.CreateCalcConsumingStaminaService(userSkill, exploreMaster, reductionSkill)
	getCommonService := stage.CreateCommonGetActionDetail(
		calcStaminaService.Calc,
		itemStorage,
		exploreMaster,
		earningItem,
		consumingItem,
		skillMaster,
		userSkill,
		requiredSkill,
	)
	getActionDetailService := stage.CreateGetStageActionDetailService(
		getCommonService.GetAction,
		stageMaster)
	getStageActionDetail := endpoint.CreateGetStageActionDetail(getActionDetailService.GetAction)

	getStageActionDetailHandler := func(w http.ResponseWriter, r *http.Request) {
		var req gateway.GetStageActionDetailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Errorf("error on decode request: %w", err).Error(), http.StatusBadRequest)
			return
		}
		res, err := getStageActionDetail(&req)
		if err != nil {
			http.Error(w, fmt.Errorf("error on generate response: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		/*
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(&res); err != nil {
				http.Error(w, fmt.Errorf("error on generate response: %w", err).Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(w, buf.String())
		*/
		resJson, err := json.Marshal(res)
		if err != nil {
			http.Error(w, fmt.Errorf("error on generate response: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resJson)
	}

	http.HandleFunc("/stage", getStageActionDetailHandler)
	http.HandleFunc("/", hello)
	http.ListenAndServe(":4444", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Kisaragi!")
}
