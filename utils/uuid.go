package utils

import "github.com/google/uuid"

func GenerateUUID() string {
	uuidObj, _ := uuid.NewUUID()
	return uuidObj.String()
}
