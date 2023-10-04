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
