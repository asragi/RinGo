package core

import "testing"

func TestUserIdIsInvalidError(t *testing.T) {
	e := UserIdIsInvalidError{userId: "text"}
	expect := "id is invalid: text"

	if e.Error() != expect {
		t.Errorf("expect: %s, got: %s", expect, e.Error())
	}
}

func TestTokenIsInvalidError(t *testing.T) {
	e := TokenIsInvalidError{token: "token"}
	expect := "token is invalid: token"

	if e.Error() != expect {
		t.Errorf("expect: %s, got: %s", expect, e.Error())
	}
}
