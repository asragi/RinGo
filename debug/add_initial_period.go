package debug

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/router"
	"net/http"
)

func CreateAddInitialPeriod(execFunc database.ExecFunc) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := addInitialPeriod(ctx, execFunc)
		if err != nil {
			http.Error(w, "error on add initial period", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func addInitialPeriod(ctx context.Context, execFunc database.ExecFunc) error {
	query := `INSERT INTO ringo.rank_period_table (rank_period) VALUES (1)`
	_, err := execFunc(ctx, query, nil)
	if err != nil {
		return fmt.Errorf("insert initial period: %w", err)
	}
	return nil
}
