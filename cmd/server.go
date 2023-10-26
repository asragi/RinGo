package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asragi/RinGo/application"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/infrastructure"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type handler func(http.ResponseWriter, *http.Request)

type infrastructuresStruct struct {
	userResource   stage.UserResourceRepo
	itemMaster     stage.ItemMasterRepo
	itemStorage    stage.ItemStorageRepo
	userSkill      stage.UserSkillRepo
	stageMaster    stage.StageMasterRepo
	exploreMaster  stage.ExploreMasterRepo
	skillMaster    stage.SkillMasterRepo
	earningItem    stage.EarningItemRepo
	consumingItem  stage.ConsumingItemRepo
	requiredSkill  stage.RequiredSkillRepo
	skillGrowth    stage.SkillGrowthDataRepo
	reductionSkill stage.ReductionStaminaSkillRepo
	updateStorage  stage.ItemStorageUpdateRepo
	updateSkill    stage.SkillGrowthPostRepo
}

func createInfrastructures() (*infrastructuresStruct, error) {
	handleError := func(err error) (*infrastructuresStruct, error) {
		return nil, fmt.Errorf("error on create infrastructures: %w", err)
	}
	gwd, _ := os.Getwd()
	dataDir := gwd + "/infrastructure/data/%s.csv"
	userResource := infrastructure.CreateInMemoryUserResourceRepo()
	itemMaster, err := infrastructure.CreateInMemoryItemMasterRepo(&infrastructure.ItemMasterLoader{Path: fmt.Sprintf(dataDir, "item-master")})
	if err != nil {
		return handleError(err)
	}
	itemStorage, err := infrastructure.CreateInMemoryItemStorageRepo(&infrastructure.ItemStorageLoader{Path: fmt.Sprintf(dataDir, "item-storage")})
	if err != nil {
		return handleError(err)
	}
	userSkill, err := infrastructure.CreateInMemoryUserSkillRepo(&infrastructure.UserSkillLoader{Path: "./infrastructure/data/user-skill.csv"})
	if err != nil {
		return handleError(err)
	}
	stageMaster, err := infrastructure.CreateInMemoryStageMasterRepo(&infrastructure.StageMasterLoader{Path: fmt.Sprintf(dataDir, "stage-master")})
	if err != nil {
		return handleError(err)
	}
	skillMaster, err := infrastructure.CreateInMemorySkillMasterRepo(&infrastructure.SkillMasterLoader{Path: "./infrastructure/data/skill-master.csv"})
	if err != nil {
		return handleError(err)
	}
	exploreMaster, err := infrastructure.CreateInMemoryExploreMasterRepo(&infrastructure.ExploreMasterLoader{Path: "./infrastructure/data/explore-master.csv"})
	if err != nil {
		return handleError(err)
	}
	earningItem, err := infrastructure.CreateInMemoryEarningItemRepo(&infrastructure.EarningItemLoader{Path: fmt.Sprintf(dataDir, "earning-item")})
	if err != nil {
		return handleError(err)
	}
	consumingItem, err := infrastructure.CreateInMemoryConsumingItemRepo(&infrastructure.ConsumingItemLoader{Path: "./infrastructure/data/consuming-item.csv"})
	if err != nil {
		return handleError(err)
	}
	requiredSkill, err := infrastructure.CreateInMemoryRequiredSkillRepo(&infrastructure.RequiredSkillLoader{Path: fmt.Sprintf(dataDir, "required-skill")})
	if err != nil {
		return handleError(err)
	}
	skillGrowth, err := infrastructure.CreateInMemorySkillGrowthDataRepo(&infrastructure.SkillGrowthLoader{Path: fmt.Sprintf(dataDir, "skill-growth")})
	reductionSkill, err := infrastructure.CreateInMemoryReductionStaminaSkillRepo(&infrastructure.ReductionStaminaSkillLoader{Path: fmt.Sprintf(dataDir, "reduction-stamina-skill")})
	if err != nil {
		return handleError(err)
	}
	return &infrastructuresStruct{
		userResource:   userResource,
		itemMaster:     itemMaster,
		itemStorage:    itemStorage,
		userSkill:      userSkill,
		skillMaster:    skillMaster,
		stageMaster:    stageMaster,
		earningItem:    earningItem,
		consumingItem:  consumingItem,
		exploreMaster:  exploreMaster,
		requiredSkill:  requiredSkill,
		skillGrowth:    skillGrowth,
		reductionSkill: reductionSkill,
		updateStorage:  itemStorage,
		updateSkill:    userSkill,
	}, nil
}

func CreateGetStageActionDetailHandler(
	infrastructures infrastructuresStruct,
) handler {
	calcStaminaService := stage.CreateCalcConsumingStaminaService(
		infrastructures.userSkill,
		infrastructures.exploreMaster,
		infrastructures.reductionSkill)
	getCommonService := stage.CreateCommonGetActionDetail(
		calcStaminaService.Calc,
		infrastructures.itemStorage,
		infrastructures.exploreMaster,
		infrastructures.earningItem,
		infrastructures.consumingItem,
		infrastructures.skillMaster,
		infrastructures.userSkill,
		infrastructures.requiredSkill,
	)
	getActionDetailService := stage.CreateGetStageActionDetailService(
		getCommonService.GetAction,
		infrastructures.stageMaster)
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

		resJson, err := json.Marshal(res)
		if err != nil {
			http.Error(w, fmt.Errorf("error on generate response: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resJson)
	}
	return getStageActionDetailHandler
}

func createPostHandler(
	infrastructures infrastructuresStruct,
	diContainer stage.DependencyInjectionContainer,
	random core.IRandom,
	currentTime core.ICurrentTime,
) handler {
	postActionApp := application.CreatePostActionService(
		infrastructures.userResource,
		infrastructures.exploreMaster,
		infrastructures.skillGrowth,
		infrastructures.userSkill,
		infrastructures.earningItem,
		infrastructures.consumingItem,
		infrastructures.requiredSkill,
		infrastructures.itemStorage,
		infrastructures.itemMaster,
		diContainer.ValidateAction,
		diContainer.CalcSkillGrowth,
		diContainer.CalcGrowthApply,
		diContainer.CalcEarnedItem,
		diContainer.CalcConsumedItem,
		diContainer.CalcTotalItem,
		diContainer.StaminaReduction,
		infrastructures.updateStorage.Update,
		infrastructures.updateSkill.Update,
		random,
		stage.PostAction,
		diContainer.GetPostActionArgs,
		currentTime,
	)
	postAction := endpoint.CreatePostAction(postActionApp)

	postActionHandler := func(w http.ResponseWriter, r *http.Request) {
		var req gateway.PostActionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Errorf("error on decode request: %w", err).Error(), http.StatusBadRequest)
			return
		}
		res, err := postAction.Post(&req)
		if err != nil {
			http.Error(w, fmt.Errorf("error on generate response: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		resJson, err := json.Marshal(res)
		if err != nil {
			http.Error(w, fmt.Errorf("error on generate response: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		setHeader(w)
		w.WriteHeader(http.StatusOK)
		w.Write(resJson)
	}
	return postActionHandler
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	handleError := func(err error) {
		log.Fatal(err.Error())
	}
	infrastructures, err := createInfrastructures()
	if err != nil {
		handleError(err)
		return
	}

	diContainer := stage.CreateDIContainer()
	currentTimeEmitter := core.CurrentTimeEmitter{}
	random := core.RandomEmitter{}

	postActionHandler := createPostHandler(*infrastructures, diContainer, &random, &currentTimeEmitter)
	getStageActionDetailHandler := CreateGetStageActionDetailHandler(*infrastructures)
	http.HandleFunc("/action", postActionHandler)
	http.HandleFunc("/stage", getStageActionDetailHandler)
	http.HandleFunc("/", hello)
	http.ListenAndServe(":4444", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Kisaragi!")
}
