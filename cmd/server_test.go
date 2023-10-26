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
	type testCase struct {
	}
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

	reqBody := bytes.NewBufferString(`{"user_id": "1", "token": "", "explore_id": "1", "exec_count": 1 }`)
	req := httptest.NewRequest(http.MethodPost, "/", reqBody)
	rec := httptest.NewRecorder()

	postActionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status is :%d", rec.Code)
		t.Errorf("Body is: %s", rec.Body)
	}
}
