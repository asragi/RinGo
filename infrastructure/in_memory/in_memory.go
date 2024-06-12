package in_memory

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
)

func FetchSizeToActionRepoInMemory(_ context.Context, size shelf.Size) (game.ActionId, error) {
	return game.ActionId(fmt.Sprintf("size-to-%d", size)), nil
}
