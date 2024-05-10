package utils

import "testing"

func TestSet(t *testing.T) {
	type innerStructForTestSet struct {
		id   int
		body string
	}

	type testCase struct {
		data []*innerStructForTestSet
	}

	testCases := []testCase{
		{
			data: []*innerStructForTestSet{
				{
					id:   1,
					body: "test1",
				},
				{
					id:   2,
					body: "test2",
				},
			},
		},
	}

	for _, v := range testCases {
		set := NewSet[*innerStructForTestSet](v.data)
		if set.Length() != len(v.data) {
			t.Errorf("expected: %d, got: %d", len(v.data), set.Length())
		}
		for i := 0; i < set.Length(); i++ {
			if set.Get(i) != v.data[i] {
				t.Errorf("expected: %v, got: %v", v.data[i], set.Get(i))
			}
		}
		for _, w := range v.data {
			foundData := set.Find(func(d *innerStructForTestSet) bool { return d.id == w.id })
			if foundData != w {
				t.Errorf("expected: %v, got: %v", w, foundData)
			}
		}
		m := SetToMap(set, func(d *innerStructForTestSet) int { return d.id })
		for _, w := range v.data {
			if m[w.id] != w {
				t.Errorf("expected: %v, got: %v", w, m[w.id])
			}
		}
	}
}
