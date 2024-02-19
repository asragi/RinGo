package utils

import "testing"

func TestGenerateUUID(t *testing.T) {
	a := GenerateUUID()
	b := GenerateUUID()
	if a == "" || b == "" {
		t.Errorf("id was not generated")
	}
	if a == b {
		t.Errorf("same id generated")
	}
}
