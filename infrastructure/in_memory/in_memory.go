package in_memory

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
)

func FetchSizeToActionRepoInMemory(_ context.Context, size shelf.Size) (game.ExploreId, error) {
	return game.ExploreId(fmt.Sprintf("size-to-%d", size)), nil
}
