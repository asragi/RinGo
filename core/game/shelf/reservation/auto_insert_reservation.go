package reservation

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type AutoInsertReservationFunc func(context.Context) error

func CreateAutoInsertReservation(
	fetchAllUserId core.FetchAllUserId,
	insertReservation BatchInsertReservationFunc,
) AutoInsertReservationFunc {
	return func(ctx context.Context) error {
		handleError := func(err error) error {
			return fmt.Errorf("error on auto insert reservation: %w", err)
		}
		allUserId, err := fetchAllUserId(ctx)
		if err != nil {
			return handleError(err)
		}
		if len(allUserId) == 0 {
			return nil
		}
		err = insertReservation(ctx, allUserId)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
