package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	type mockReq struct {
		Message string `json:"message"`
	}
	type mockRes struct {
		message string
	}
	type testCase struct {
		expectStatus int
		request      string
		mockRes      mockRes
		mockError    error
	}
	testCases := []testCase{
		{
			expectStatus: http.StatusBadRequest,
			request:      `{"user_id": "invalid request", "stage_id": "1", "explore_id": "1", "token": "1"}`,
			mockRes: mockRes{
				message: "test message",
			},
		},
		{
			expectStatus: http.StatusOK,
			request:      `{"message": "valid request"}`,
			mockRes: mockRes{
				message: "test message",
			},
		},
	}

	for _, v := range testCases {
		mockEndpoint := func(req *mockReq) (*mockRes, error) {
			return &v.mockRes, v.mockError
		}
		mockLogger := func(i int, err error) {}

		reqBody := bytes.NewBufferString(v.request)
		req := httptest.NewRequest(http.MethodGet, "/", reqBody)
		rec := httptest.NewRecorder()
		handler := createHandler(mockEndpoint, mockLogger)
		handler(rec, req)
		if rec.Code != v.expectStatus {
			t.Errorf("expect: %d, got: %d", v.expectStatus, rec.Code)
		}
		fmt.Printf("Body is: %s\n", rec.Body)
	}
}
