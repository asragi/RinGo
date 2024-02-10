package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Handler func(http.ResponseWriter, *http.Request)
type writeLogger func(int, error)

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

type ReturnResponseOnErrorFunc func(http.ResponseWriter, error)

func ErrorOnDecode(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("error on decode request: %w", err).Error(), http.StatusBadRequest)
}

func ErrorOnGenerateResponse(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("error on generate response: %w", err).Error(), http.StatusInternalServerError)
}

func ErrorOnInternalError(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("internal server error: %w", err).Error(), http.StatusInternalServerError)
}

func ErrorOnMethodNotAllowed(w http.ResponseWriter, err error) {

}

func ErrorOnPageNotFound(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Errorf("not found: %w", err).Error(), http.StatusNotFound)
}

type RequestBody io.ReadCloser
type QueryParameter url.Values

func (q *QueryParameter) GetFirstQuery(name string) (string, error) {
	arr := (*q)[name]
	if len(arr) <= 0 {
		return "", NoQueryProvidedError{Message: name}
	}
	return arr[0], nil
}

type PathString string

func DecodeBody[T any](body io.ReadCloser) (*T, error) {
	var req T
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)

	return &req, err
}

func createHandlerWithParameter[T any, S any](
	endpointFunc func(*T) (S, error),
	selectParam func(RequestBody, QueryParameter, PathString) (*T, error),
	logger writeLogger,
) Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		passedUrl := r.URL
		query := passedUrl.Query()
		req, err := selectParam(r.Body, QueryParameter(query), PathString(passedUrl.Path))
		if err != nil {
			ErrorOnDecode(w, err)
			return
		}
		res, err := endpointFunc(req)
		if err != nil {
			ErrorOnGenerateResponse(w, err)
			return
		}
		resJson, err := json.Marshal(res)
		if err != nil {
			ErrorOnGenerateResponse(w, err)
			return
		}
		setHeader(w)
		logger(w.Write(resJson))
	}

	return h
}

// Deprecated: use createHandlerWithParameter
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
			ErrorOnDecode(w, err)
			return
		}
		res, err := endpointFunc(&req)
		if err != nil {
			ErrorOnGenerateResponse(w, err)
			return
		}
		resJson, err := json.Marshal(res)
		if err != nil {
			ErrorOnGenerateResponse(w, err)
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
