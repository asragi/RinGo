package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asragi/RinGo/core"
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

		userId := core.UserId("1")
		itemId := core.ItemId("1")
		itemStorage := infrastructures.itemStorage
		beforeItemRes, _ := itemStorage.Get(
			userId,
			itemId,
			"",
		)
		beforeItemNum := beforeItemRes.Stock
		postActionHandler(rec, req)

		if rec.Code != v.status {
			t.Errorf("case: %d, expect :%d, got: %d", i, v.status, rec.Code)
			t.Errorf("Body is: %s", rec.Body)
		}
		afterItemNum, _ := itemStorage.Get(
			userId,
			itemId,
			"",
		)
		if beforeItemNum == afterItemNum.Stock {
			t.Errorf("")
		}
		fmt.Printf("Body is: %s\n", rec.Body)
		fmt.Printf("Num: %d -> %d\n", beforeItemNum, afterItemNum.Stock)
	}
}

func TestGetStageActionDetail(t *testing.T) {
	infrastructures, err := createInfrastructures()
	if err != nil {
		t.Fatalf("error on test post action: %s", err.Error())
	}
	type expect struct {
		StatusCode int
	}
	type testCase struct {
		expect  expect
		request string
	}

	testCases := []testCase{
		{
			expect: expect{
				StatusCode: http.StatusOK,
			},
			request: `{"user_id": "1", "stage_id": "1", "explore_id": "1", "token": "1"}`,
		},
	}

	for _, v := range testCases {
		reqBody := bytes.NewBufferString(v.request)
		expect := v.expect
		req := httptest.NewRequest(http.MethodGet, "/", reqBody)
		rec := httptest.NewRecorder()
		handler := CreateGetStageActionDetailHandler(*infrastructures)
		handler(rec, req)
		if rec.Code != expect.StatusCode {
			t.Errorf("expect: %d, got: %d", expect.StatusCode, rec.Code)
		}
		fmt.Printf("Body is: %s\n", rec.Body)
	}
}
