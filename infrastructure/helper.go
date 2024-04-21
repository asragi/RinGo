package infrastructure

import (
	"fmt"
	"github.com/asragi/RinGo/core"
)

func UserIdsToString(userIds []core.UserId) []string {
	result := make([]string, len(userIds))
	for i, v := range userIds {
		result[i] = fmt.Sprintf(`"%s"`, v)
	}
	return result
}
