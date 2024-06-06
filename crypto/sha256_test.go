package crypto

import "testing"

func TestSHA256WithKey(t *testing.T) {
	type testCase struct {
		msg    string
		key    string
		expect string
	}

	testCases := []*testCase{
		{
			msg:    "msg",
			key:    "key",
			expect: "2d93cbc1be167bcb1637a4a23cbff01a7878f0c50ee833954ea5221bb1b8c628",
		},
	}

	for i, v := range testCases {
		hashed, err := SHA256WithKey(v.key, v.msg)
		if err != nil {
			t.Fatalf("case:%d, error: %s", i, err.Error())
			return
		}
		if hashed != v.expect {
			t.Errorf("case:%d, expect: %s, got: %s", i, v.expect, hashed)
		}
	}
}
