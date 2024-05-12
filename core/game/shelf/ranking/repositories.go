package ranking

import (
	"context"
	"github.com/asragi/RinGo/core"
	"time"
)

type UserScorePair struct {
	UserId     core.UserId `db:"user_id"`
	TotalScore TotalScore  `db:"total_score"`
}
type FetchUserScore func(ctx context.Context, userId []core.UserId, currentTime time.Time) ([]*UserScorePair, error)
type UpsertScoreFunc func(context.Context, []*UserScorePair, time.Time) error
