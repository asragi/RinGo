package main

import (
	"bytes"
	"fmt"
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
		requestBody string
		status      int
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

	testCases := []testCase{
		{
			requestBody: `{"user_id": "1", "token": "", "explore_id": "1", "exec_count": 1 }`,
			status:      http.StatusOK,
		},
		/*
			{
				requestBody: `{"user_id": "123456", "token": "", "explore_id": "1", "exec_count": 1 }`,
				status:      http.StatusBadRequest,
			},
		*/
	}

	for i, v := range testCases {
		reqBody := bytes.NewBufferString(v.requestBody)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		rec := httptest.NewRecorder()

		postActionHandler(rec, req)

		if rec.Code != v.status {
			t.Errorf("case: %d, expect :%d, got: %d", i, v.status, rec.Code)
			t.Errorf("Body is: %s", rec.Body)
		}
		fmt.Printf("Body is: %s", rec.Body)
	}
}
