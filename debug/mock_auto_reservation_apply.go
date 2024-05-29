package debug

import (
	"fmt"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RinGo/router"
	"net/http"
)

func MockAutoReservationApply(
	autoApplyReservation reservation.AutoInsertReservationFunc,
) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := autoApplyReservation(ctx)
		if err != nil {
			http.Error(w, "error on auto apply reservation", http.StatusInternalServerError)
			fmt.Printf("error on auto apply reservation: %v\n", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
