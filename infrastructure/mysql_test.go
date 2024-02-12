package infrastructure

import "testing"

func TestCreateListStatement(t *testing.T) {
	type testCase struct {
		args   []string
		expect string
	}

	testCases := []testCase{
		{
			[]string{"aa", "bb", "cc"},
			"(aa, bb, cc)",
		},
	}

	for _, v := range testCases {
		result := createListStatement(v.args)
		if v.expect != result {
			t.Errorf("expect: %s, got: %s", v.expect, result)
		}
	}
}
