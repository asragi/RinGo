package router

import (
	"fmt"
	"net/http"
)

func ErrorOnDecode(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("error on decode request: %w", err).Error(), http.StatusBadRequest)
}

func ErrorOnInternalError(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("internal server error: %w", err).Error(), http.StatusInternalServerError)
}

func ErrorOnPageNotFound(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("not found: %w", err).Error(), http.StatusNotFound)
}
