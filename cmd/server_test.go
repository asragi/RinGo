package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RinGo/test"
)

func TestCreateInfrastructures(t *testing.T) {
	_, err := createInfrastructures()
	if err != nil {
		t.Errorf("error occurred:%s", err.Error())
	}
}

func TestPostActionHttp(t *testing.T) {
	infrastructures, err := createInfrastructures()
	if err != nil {
		t.Fatalf("error on test post action: %s", err.Error())
	}
	diContainer := stage.CreateDIContainer()
	timer := test.MockTimer{}
	random := test.TestRandom{Value: 0.5}
	postActionHandler := createPostHandler(
		*infrastructures,
		diContainer,
		&random,
		&timer,
	)

	reqBody := bytes.NewBufferString("")
	req := httptest.NewRequest(http.MethodGet, "/action", reqBody)
	rec := httptest.NewRecorder()

	postActionHandler(rec, req)
}
