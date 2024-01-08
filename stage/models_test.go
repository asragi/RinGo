package stage

import (
	"github.com/asragi/RinGo/core"
	"testing"
)

func TestGainingPoint_Multiply(t *testing.T) {
	type testCase struct {
		point    int
		multiply int
		expect   int
	}
	testCases := []testCase{
		{
			point:    100,
			multiply: 2,
			expect:   200,
		},
		{
			point:    0,
			multiply: 2,
			expect:   0,
		},
	}
	for _, v := range testCases {
		point := GainingPoint(v.point)
		afterPoint := point.Multiply(v.multiply)
		if int(afterPoint) != v.expect {
			t.Errorf("expect: %d, got: %d", v.expect, afterPoint)
		}
	}
}

func TestGainingPoint_ApplyTo(t *testing.T) {
	type testCase struct {
		point    int
		applyExp core.SkillExp
		expect   core.SkillExp
	}
	testCases := []testCase{
		{
			point:    100,
			applyExp: 2,
			expect:   102,
		},
		{
			point:    0,
			applyExp: 2,
			expect:   2,
		},
	}
	for _, v := range testCases {
		point := GainingPoint(v.point)
		afterExp := point.ApplyTo(v.applyExp)
		if afterExp != v.expect {
			t.Errorf("expect: %d, got: %d", v.expect, afterExp)
		}
	}
}
