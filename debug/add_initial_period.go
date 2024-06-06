package debug

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/database"
)

// addInitialPeriod required
func addInitialPeriod(ctx context.Context, execFunc database.ExecFunc) error {
	query := `INSERT INTO ringo.rank_period_table (rank_period) VALUES (1)`
	_, err := execFunc(ctx, query, nil)
	if err != nil {
		return fmt.Errorf("insert initial period: %w", err)
	}
	return nil
}
