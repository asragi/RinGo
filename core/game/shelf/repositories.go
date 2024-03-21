package shelf

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type (
	FetchSizeToActionRepoFunc func(context.Context, Size) (game.ExploreId, error)
	UpdateShelfSizeRepoFunc   func(
		ctx context.Context,
		userId core.UserId,
		shelfSize Size,
	) error
)
