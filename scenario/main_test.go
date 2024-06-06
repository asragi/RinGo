package scenario

import (
	"context"
	"testing"
)

func TestE2E(t *testing.T) {
	c := newClient("localhost:50051")
	ctx := context.Background()
	err := signUp(ctx, c)
	if err != nil {
		t.Errorf("sign up: %v", err)
	}
}
