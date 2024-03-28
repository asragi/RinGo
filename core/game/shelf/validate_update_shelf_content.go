package shelf

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type ValidateUpdateShelfContentFunc func(
	context.Context,
	core.UserId,
	core.ItemId,
	Index,
) error
