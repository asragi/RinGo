package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request)
type writeLogger func(int, error)

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func errorOnDecode(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("error on decode request: %w", err).Error(), http.StatusBadRequest)
}

func errorOnGenerateResponse(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("error on generate response: %w", err).Error(), http.StatusInternalServerError)
}

func createHandler[T any, S any](
	endpointFunc func(*T) (S, error),
	logger writeLogger,
) Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		var req T
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&req)
		if err != nil {
			errorOnDecode(w, err)
			return
		}
		res, err := endpointFunc(&req)
		if err != nil {
			errorOnGenerateResponse(w, err)
			return
		}
		resJson, err := json.Marshal(res)
		if err != nil {
			errorOnGenerateResponse(w, err)
			return
		}
		setHeader(w)
		logger(w.Write(resJson))
	}

	return h
}

func LogHttpWrite(status int, err error) {
	if err == nil {
		return
	}
	log.Printf("Write failed: %v, status: %d", err, status)
}
