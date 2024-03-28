package shelf

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type ValidateUpdateShelfSizeFunc func(
	context.Context,
	core.UserId,
	Size,
) error
