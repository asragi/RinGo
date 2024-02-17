package utils

import (
	"encoding/json"
	"fmt"
)

type JsonToStructFunc[S any] func(json string) (*S, error)

func JsonToStruct[S any](jsonString string) (*S, error) {
	var targetStruct S
	err := json.Unmarshal([]byte(jsonString), &targetStruct)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	return &targetStruct, nil
}

type StructToJsonFunc[S any] func(*S) (*string, error)

func StructToJson[S any](obj *S) (*string, error) {
	jsonData, err := json.Marshal(*obj)
	if err != nil {
		return nil, fmt.Errorf("struct to json: %w", err)
	}
	jsonString := string(jsonData)
	return &jsonString, nil
}
