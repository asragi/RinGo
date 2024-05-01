package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type WriteLogger func(int, error)

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

type requestBody io.ReadCloser

type queryParameter url.Values

func (q *queryParameter) GetFirstQuery(name string) (string, error) {
	arr := (*q)[name]
	if len(arr) <= 0 {
		return "", NoQueryProvidedError{Message: name}
	}
	return arr[0], nil
}

type pathString string
type requestHeader struct {
	header http.Header
}

func (h *requestHeader) Get(key string) string {
	return h.header.Get(key)
}

func DecodeBody[T any](body io.ReadCloser) (*T, error) {
	var req T
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)

	return &req, err
}

type selectParam[T any] func(requestHeader, requestBody, queryParameter, pathString) (*T, error)

func createHandlerWithParameter[T any, S any](
	endpointFunc func(context.Context, *T) (*S, error),
	createContext utils.CreateContextFunc,
	selectParam selectParam[T],
	logger WriteLogger,
) router.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		passedUrl := r.URL
		query := passedUrl.Query()
		req, err := selectParam(requestHeader{r.Header}, r.Body, queryParameter(query), pathString(passedUrl.Path))
		if err != nil {
			ErrorOnDecode(w, err)
			return
		}
		ctx := createContext()
		res, err := endpointFunc(ctx, req)
		if err != nil {
			ErrorOnGenerateResponse(w, err)
			return
		}
		resJson, err := json.Marshal(*res)
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
		log.Printf("Write success: %d", status)
		return
	}
	log.Printf("Write failed: %v, status: %d", err, status)
}

func (h *requestHeader) getTokenFromHeader() (string, error) {
	tokenHeader := h.Get("Authorization")
	if tokenHeader == "" {
		return "", fmt.Errorf("no token provided")
	}
	headerData := strings.Split(tokenHeader, " ")
	if len(headerData) != 2 {
		return "", fmt.Errorf("invalid token format")
	}
	if headerData[0] != "Bearer" {
		return "", fmt.Errorf("invalid token type")
	}
	return headerData[1], nil
}
