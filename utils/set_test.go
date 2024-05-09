package utils

import "testing"

type innerStructForTestSet struct {
	id   int
	body string
}

func (i *innerStructForTestSet) Id() int {
	return i.id
}

type testCase struct {
	data []*innerStructForTestSet
}

func TestSet(t *testing.T) {
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
		set := NewSet[int, *innerStructForTestSet](v.data)
		if set.Length() != len(v.data) {
			t.Errorf("expected: %d, got: %d", len(v.data), set.Length())
		}
		for i := 0; i < set.Length(); i++ {
			if set.Get(i) != v.data[i] {
				t.Errorf("expected: %v, got: %v", v.data[i], set.Get(i))
			}
		}
		for _, w := range v.data {
			if set.Find(w.id) != w {
				t.Errorf("expected: %v, got: %v", w, set.Find(w.id))
			}
		}
		m := set.ToMap()
		for _, w := range v.data {
			if m[w.id] != w {
				t.Errorf("expected: %v, got: %v", w, m[w.id])
			}
		}
	}
}
