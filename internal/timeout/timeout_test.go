package timeout_test

import (
	"context"
	"testing"
	"time"

	"github.com/vaultpull/vaultpull/internal/timeout"
)

func TestDefaultPolicyIsValid(t *testing.T) {
	if err := timeout.DefaultPolicy().Validate(); err != nil {
		t.Fatalf("expected default policy to be valid, got: %v", err)
	}
}

func TestValidateRejectsZeroDial(t *testing.T) {
	p := timeout.DefaultPolicy()
	p.Dial = 0
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero dial timeout")
	}
}

func TestValidateRejectsZeroRead(t *testing.T) {
	p := timeout.DefaultPolicy()
	p.Read = 0
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero read timeout")
	}
}

func TestValidateRejectsZeroOverall(t *testing.T) {
	p := timeout.DefaultPolicy()
	p.Overall = 0
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero overall timeout")
	}
}

func TestApplyDeadlineIsSet(t *testing.T) {
	p := timeout.Policy{
		Dial:    1 * time.Second,
		Read:    2 * time.Second,
		Overall: 5 * time.Second,
	}
	ctx, cancel := p.Apply(context.Background())
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) > 5*time.Second {
		t.Fatalf("deadline exceeds overall timeout")
	}
}
