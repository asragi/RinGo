package router

import "testing"

func TestNewPathMatchPattern(t *testing.T) {
	type testCase struct {
		samplePath  string
		validPath   []string
		invalidPath []string
	}
	testCases := []testCase{
		{
			samplePath: "/users/{userId}/items/{itemId}",
			validPath: []string{
				"/users/test-path-sample1/items/test-item-id",
			},
			invalidPath: []string{
				"/users/test-path/items",
				"/users/test-path/items/some-id/some",
				"/actions/test-path/items/some-id",
			},
		},
	}

	for _, v := range testCases {
		samplePath, err := NewSamplePath(v.samplePath)
		if err != nil {
			t.Fatalf("error: %+v", err)
		}
		checker, err := NewPathMatchPattern(samplePath)
		if err != nil {
			t.Fatalf("error: %+v", err)
		}
		for _, w := range v.validPath {
			path, err := NewPathData(w)
			if err != nil {
				t.Fatalf("error: %+v", err)
			}
			isValid := checker.Match(path)
			if !isValid {
				t.Errorf("expect: true, got: false")
			}
		}
		for _, w := range v.invalidPath {
			path, err := NewPathData(w)
			if err != nil {
				t.Fatalf("error: %+v", err)
			}
			isValid := checker.Match(path)
			if isValid {
				t.Errorf("expect: false, got: true")
			}
		}
	}
}
