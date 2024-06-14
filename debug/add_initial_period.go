package debug

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/database"
)

type AddInitialPeriodFunc func(ctx context.Context) error

func CreateAddInitialPeriod(execFunc database.ExecFunc) AddInitialPeriodFunc {
	return func(ctx context.Context) error {
		query := `INSERT INTO ringo.rank_period_table (rank_period) VALUES (1)`
		_, err := execFunc(ctx, query, nil)
		if err != nil {
			return fmt.Errorf("insert initial period: %w", err)
		}
		return nil
	}
}
