package infrastructure

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

type ItemMasterLoader struct {
	Path string
}

func readCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error on opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	allDataStrings, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error on opening file: %w", err)
	}
	return allDataStrings, err
}

func (loader *ItemMasterLoader) Load() (itemMasterData, error) {
	handleError := func(err error) (itemMasterData, error) {
		return itemMasterData{}, fmt.Errorf("error on load item data: %w", err)
	}
	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}
	toItemMasterData := func(allCSVData [][]string) (itemMasterData, error) {
		handleError := func(err error) (itemMasterData, error) {
			return nil, fmt.Errorf("error on convert CSV: %w", err)
		}
		result := make(itemMasterData)
		for i, v := range allCSVData {
			if i == 0 {
				continue
			}
			price, err := strconv.Atoi(v[3])
			if err != nil {
				return handleError(err)
			}
			maxStock, err := strconv.Atoi(v[4])
			if err != nil {
				return handleError(err)
			}
			data := stage.GetItemMasterRes{
				ItemId:      core.ItemId(v[0]),
				DisplayName: core.DisplayName(v[1]),
				Description: core.Description(v[2]),
				Price:       core.Price(price),
				MaxStock:    core.MaxStock(maxStock),
			}
			result[data.ItemId] = data
		}
		return result, err
	}

	result, err := toItemMasterData(allCSVData)
	if err != nil {
		return handleError(err)
	}

	return result, err
}

type ItemStorageLoader struct {
	Path string
}

func (loader *ItemStorageLoader) Load() (ItemStorageData, error) {
	handleError := func(err error) (ItemStorageData, error) {
		return ItemStorageData{}, fmt.Errorf("error on load item storage data: %w", err)
	}

	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}

	toData := func(csvData [][]string) (ItemStorageData, error) {
		result := make(ItemStorageData)
		for i, v := range csvData {
			if i == 0 {
				continue
			}
			userId := core.UserId(v[0])
			itemId := core.ItemId(v[1])
			stockValue, err := strconv.Atoi(v[2])
			if err != nil {
				return handleError(err)
			}
			isKnownValue, err := strconv.ParseBool(v[3])
			if err != nil {
				return handleError(err)
			}
			data := stage.ItemData{
				UserId:  userId,
				ItemId:  itemId,
				Stock:   core.Stock(stockValue),
				IsKnown: core.IsKnown(isKnownValue),
			}
			result[data.UserId][data.ItemId] = data
		}
		return result, nil
	}
	data, err := toData(allCSVData)
	if err != nil {
		return handleError(err)
	}
	return data, err

}

type UserSkillLoader struct {
	Path string
}

func (loader *UserSkillLoader) Load() (UserSkillData, error) {
	handleError := func(err error) (UserSkillData, error) {
		return UserSkillData{}, fmt.Errorf("error on load user skill data: %w", err)
	}

	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}

	toData := func(csvData [][]string) (UserSkillData, error) {
		result := make(UserSkillData)
		for i, v := range csvData {
			if i == 0 {
				continue
			}
			userId := core.UserId(v[0])
			skillId := core.SkillId(v[1])
			skillExpValue, err := strconv.Atoi(v[2])
			if err != nil {
				return handleError(err)
			}
			skillExp := core.SkillExp(skillExpValue)
			data := stage.UserSkillRes{
				UserId:   userId,
				SkillId:  skillId,
				SkillExp: skillExp,
			}
			result[data.UserId][data.SkillId] = data
		}
		return result, nil
	}
	data, err := toData(allCSVData)
	if err != nil {
		return handleError(err)
	}
	return data, err
}

type StageMasterLoader struct {
	Path string
}

func (loader *StageMasterLoader) Load() (StageMasterData, error) {
	handleError := func(err error) (StageMasterData, error) {
		return StageMasterData{}, fmt.Errorf("error on load stage master data: %w", err)
	}
	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}
	toItemMasterData := func(allCSVData [][]string) (StageMasterData, error) {
		result := make(StageMasterData)
		for i, v := range allCSVData {
			if i == 0 {
				continue
			}
			data := stage.StageMaster{
				StageId:     stage.StageId(v[0]),
				DisplayName: core.DisplayName(v[1]),
				Description: core.Description(v[2]),
			}
			result[data.StageId] = data
		}
		return result, err
	}

	result, err := toItemMasterData(allCSVData)
	if err != nil {
		return handleError(err)
	}

	return result, err

}

type SkillMasterLoader struct {
	Path string
}

func (loader *SkillMasterLoader) Load() (SkillMasterData, error) {
	handleError := func(err error) (SkillMasterData, error) {
		return nil, fmt.Errorf("error on loading skill master: %w", err)
	}
	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}
	toData := func(csvData [][]string) SkillMasterData {
		result := make(SkillMasterData)
		for i, v := range csvData {
			if i == 0 {
				continue
			}
			data := stage.SkillMaster{
				SkillId:     core.SkillId(v[0]),
				DisplayName: core.DisplayName(v[1]),
			}
			result[data.SkillId] = data
		}
		return result
	}
	return toData(allCSVData), nil
}

type ExploreMasterLoader struct {
	Path string
}

func (loader *ExploreMasterLoader) Load() (ExploreMasterData, error) {
	handleError := func(err error) (ExploreMasterData, error) {
		return nil, fmt.Errorf("error on load explore master: %w", err)
	}
	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}
	toData := func(csvData [][]string) (ExploreMasterData, error) {
		result := make(ExploreMasterData)
		for i, v := range csvData {
			if i == 0 {
				continue
			}
			// TODO: implement all field
			exploreId := stage.ExploreId(v[0])
			data := stage.GetExploreMasterRes{
				ExploreId: exploreId,
			}
			result[data.ExploreId] = data
		}
		return result, nil
	}
	data, err := toData(allCSVData)
	if err != nil {
		return handleError(err)
	}

	return data, nil
}

type EarningItemLoader struct {
	Path string
}

func (loader EarningItemLoader) Load() (EarningItemData, error) {
	handleError := func(err error) (EarningItemData, error) {
		return nil, fmt.Errorf("error on load earning item master: %w", err)
	}
	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}
	toData := func(csvData [][]string) (EarningItemData, error) {
		handleError := func(err error) (EarningItemData, error) {
			return nil, fmt.Errorf("error on convert value type: %w", err)
		}
		result := make(EarningItemData)
		for i, v := range csvData {
			if i == 0 {
				continue
			}
			exploreId := stage.ExploreId(v[0])
			maxCountValue, err := strconv.Atoi(v[3])
			if err != nil {
				return handleError(err)
			}
			minCountValue, err := strconv.Atoi(v[2])
			if err != nil {
				return handleError(err)
			}
			data := stage.EarningItem{
				ItemId:   core.ItemId(v[1]),
				MaxCount: core.Count(maxCountValue),
				MinCount: core.Count(minCountValue),
			}
			if _, ok := result[exploreId]; !ok {
				result[exploreId] = []stage.EarningItem{}
			}
			result[exploreId] = append(result[exploreId], data)
		}
		return result, nil
	}
	data, err := toData(allCSVData)
	if err != nil {
		return handleError(err)
	}

	return data, nil

}

type ConsumingItemLoader struct {
	Path string
}

type IConsumingItemLoader interface {
	Load() (ConsumingItemData, error)
}

func (loader *ConsumingItemLoader) Load() (ConsumingItemData, error) {
	handleError := func(err error) (ConsumingItemData, error) {
		return nil, fmt.Errorf("error on load consuming item master: %w", err)
	}
	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}
	toData := func(csvData [][]string) (ConsumingItemData, error) {
		handleError := func(err error) (ConsumingItemData, error) {
			return nil, fmt.Errorf("error on convert value type: %w", err)
		}
		result := make(ConsumingItemData)
		for i, v := range csvData {
			if i == 0 {
				continue
			}
			exploreId := stage.ExploreId(v[0])
			maxCountValue, err := strconv.Atoi(v[2])
			if err != nil {
				return handleError(err)
			}
			consumeProbValue, err := strconv.Atoi(v[3])
			if err != nil {
				return handleError(err)
			}
			data := stage.ConsumingItem{
				ItemId:          core.ItemId(v[1]),
				MaxCount:        core.Count(maxCountValue),
				ConsumptionProb: stage.ConsumptionProb(consumeProbValue),
			}
			if _, ok := result[exploreId]; !ok {
				result[exploreId] = []stage.ConsumingItem{}
			}
			result[exploreId] = append(result[exploreId], data)
		}
		return result, nil
	}
	data, err := toData(allCSVData)
	if err != nil {
		return handleError(err)
	}

	return data, nil

}

type RequiredSkillLoader struct {
	Path string
}

func (loader *RequiredSkillLoader) Load() (RequiredSkillData, error) {
	handleError := func(err error) (RequiredSkillData, error) {
		return nil, fmt.Errorf("error on load reduction skill master: %w", err)
	}
	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}
	toData := func(csvData [][]string) (RequiredSkillData, error) {
		result := make(RequiredSkillData)
		for i, v := range csvData {
			if i == 0 {
				continue
			}
			exploreId := stage.ExploreId(v[1])
			skillId := core.SkillId(v[0])
			requiredLvValue, err := strconv.Atoi(v[2])
			if err != nil {
				return handleError(err)
			}
			data := stage.RequiredSkill{
				SkillId:    skillId,
				RequiredLv: core.SkillLv(requiredLvValue),
			}
			if _, ok := result[exploreId]; !ok {
				result[exploreId] = []stage.RequiredSkill{}
			}
			result[exploreId] = append(result[exploreId], data)
		}
		return result, nil
	}
	data, err := toData(allCSVData)
	if err != nil {
		return handleError(err)
	}

	return data, nil

}

type ReductionStaminaSkillLoader struct {
	Path string
}

func (loader *ReductionStaminaSkillLoader) Load() (ReductionStaminaSkillData, error) {
	handleError := func(err error) (ReductionStaminaSkillData, error) {
		return nil, fmt.Errorf("error on load reduction skill master: %w", err)
	}
	allCSVData, err := readCSV(loader.Path)
	if err != nil {
		return handleError(err)
	}
	toData := func(csvData [][]string) (ReductionStaminaSkillData, error) {
		result := make(ReductionStaminaSkillData)
		for i, v := range csvData {
			if i == 0 {
				continue
			}
			exploreId := stage.ExploreId(v[0])
			skillId := core.SkillId(v[1])
			if _, ok := result[exploreId]; !ok {
				result[exploreId] = []core.SkillId{}
			}
			result[exploreId] = append(result[exploreId], skillId)
		}
		return result, nil
	}
	data, err := toData(allCSVData)
	if err != nil {
		return handleError(err)
	}

	return data, nil

}
