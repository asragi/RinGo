package core

import "testing"

func TestCalcLv(t *testing.T) {
	type testCase struct {
		input  SkillExp
		expect SkillLv
	}

	testCases := []testCase{
		{
			input:  0,
			expect: 1,
		},
		{
			input:  5,
			expect: 1,
		},
		{
			input:  10,
			expect: 2,
		},
		{
			input:  11,
			expect: 2,
		},
		{
			input:  30,
			expect: 3,
		},
		{
			input:  100000,
			expect: 100,
		},
	}

	for _, v := range testCases {
		actual := v.input.CalcLv()
		if v.expect != actual {
			t.Errorf("Expect %d, actual %d", v.expect, actual)
		}
	}
}
