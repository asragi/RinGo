package reservation

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type ApplyAllReservationsFunc func(ctx context.Context) error

func CreateApplyAllReservations(
	fetchAllUserId core.FetchAllUserId,
	applyReservation ApplyReservationFunc,
) ApplyAllReservationsFunc {
	return func(ctx context.Context) error {
		handleError := func(err error) error {
			return fmt.Errorf("error on apply all reservations: %w", err)
		}
		userIdReq, err := fetchAllUserId(ctx)
		if err != nil {
			return handleError(err)
		}
		if len(userIdReq) == 0 {
			return nil
		}
		err = applyReservation(ctx, userIdReq)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
