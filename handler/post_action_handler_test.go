package handler

import (
	"github.com/asragi/RingoSuPBGo/gateway"
	"net/http"
	"net/url"
	"testing"
)

func TestCreateGetPostActionParams(t *testing.T) {
	type testCase struct {
		bodyText    string
		tokenString string
		expect      *gateway.PostActionRequest
	}

	testCases := []*testCase{
		{
			bodyText:    `{"explore_id": "some-explore-id", "exec_count": 1}`,
			tokenString: `Bearer some-token`,
			expect: &gateway.PostActionRequest{
				Token:     "some-token",
				ExploreId: "some-explore-id",
				ExecCount: 1,
			},
		},
	}
	noQueryParam := queryParameter(url.Values{})
	path := pathString("/")
	for _, v := range testCases {
		tokenHeader := http.Header{}
		tokenHeader.Add("Authorization", v.tokenString)
		header := requestHeader{tokenHeader}
		body := stringToBody(v.bodyText)
		req, err := GetPostActionParams(header, body, noQueryParam, path)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if req.Token != v.expect.Token {
			t.Errorf("expect: %s, got: %s", v.expect.Token, req.Token)
		}
		if req.ExploreId != v.expect.ExploreId {
			t.Errorf("expect: %s, got: %s", v.expect.ExploreId, req.ExploreId)
		}
		if req.ExecCount != v.expect.ExecCount {
			t.Errorf("expect: %d, got: %d", v.expect.ExecCount, req.ExecCount)
		}
	}
}
