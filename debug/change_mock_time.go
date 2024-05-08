package debug

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/handler"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"net/http"
	"time"
)

type changeTimeInterface interface {
	SetTimer(core.GetCurrentTimeFunc)
}

func ChangeMockTimeHandler(timer changeTimeInterface) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Time string `json:"time"`
		}

		res, err := handler.DecodeBody[request](r.Body)
		if err != nil {
			http.Error(w, "error on decode request", http.StatusBadRequest)
			return
		}
		utcTime, err := utils.StringToTime(res.Time)
		if err != nil {
			http.Error(w, "error on parse time", http.StatusBadRequest)
			return
		}
		timer.SetTimer(
			func() time.Time {
				return utcTime
			},
		)
		w.WriteHeader(http.StatusNoContent)
	}
}
