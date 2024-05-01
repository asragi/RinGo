package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/asragi/RinGo/test"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func stringToBody(bodyText string) requestBody {
	return io.NopCloser(strings.NewReader(bodyText))
}

func TestCreateHandlerWithParameter(t *testing.T) {
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
			expectStatus: http.StatusOK,
			request:      `{"message": "valid request"}`,
			mockRes: mockRes{
				message: "test message",
			},
		},
	}

	for _, v := range testCases {
		mockEndpoint := func(_ context.Context, req *mockReq) (*mockRes, error) {
			return &v.mockRes, v.mockError
		}
		mockLogger := func(i int, err error) {}
		handleParam := func(
			header requestHeader,
			body requestBody,
			queryParameter queryParameter,
			pathString pathString,
		) (*mockReq, error) {
			return &mockReq{}, nil
		}

		reqBody := bytes.NewBufferString(v.request)
		req := httptest.NewRequest(http.MethodGet, "/", reqBody)
		rec := httptest.NewRecorder()
		handler := createHandlerWithParameter(mockEndpoint, test.MockCreateContext, handleParam, mockLogger)
		handler(rec, req)
		if rec.Code != v.expectStatus {
			t.Errorf("expect: %d, got: %d", v.expectStatus, rec.Code)
			fmt.Printf("Body is: %s\n", rec.Body)
		}
	}
}
