package stage

import (
	"testing"
)

func check(t *testing.T, expect string, actual string) {
	if expect != actual {
		t.Errorf("want %s, actual %s", expect, actual)
	}
}

func checkBool(t *testing.T, title string, expect bool, actual bool) {
	if expect != actual {
		t.Errorf("%s: want %t, actual %t", title, expect, actual)
	}
}

func checkInt(t *testing.T, title string, expect int, actual int) {
	if expect != actual {
		t.Errorf("%s: want %d, actual %d", title, expect, actual)
	}
}
