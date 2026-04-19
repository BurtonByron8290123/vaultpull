package jitter

import (
	"testing"
	"time"
)

func TestDefaultPolicyIsValid(t *testing.T) {
	if err := DefaultPolicy().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRejectsZeroFactor(t *testing.T) {
	p := Policy{Factor: 0}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for factor=0")
	}
}

func TestValidateRejectsNegativeFactor(t *testing.T) {
	p := Policy{Factor: -0.1}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for negative factor")
	}
}

func TestValidateRejectsFactorAboveOne(t *testing.T) {
	p := Policy{Factor: 1.1}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for factor > 1")
	}
}

func TestApplyReturnsValueInRange(t *testing.T) {
	p := Policy{Factor: 0.5}
	base := 100 * time.Millisecond
	for i := 0; i < 100; i++ {
		got := p.Apply(base)
		if got < base {
			t.Fatalf("Apply returned %v < base %v", got, base)
		}
		if got > base+time.Duration(float64(base)*0.5) {
			t.Fatalf("Apply returned %v, exceeds max", got)
		}
	}
}

func TestApplyNReturnsCorrectCount(t *testing.T) {
	p := DefaultPolicy()
	results := p.ApplyN(50*time.Millisecond, 10)
	if len(results) != 10 {
		t.Fatalf("expected 10 results, got %d", len(results))
	}
}

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envFactor, "")
	p := FromEnv()
	if p.Factor != DefaultPolicy().Factor {
		t.Fatalf("expected default factor, got %v", p.Factor)
	}
}

func TestFromEnvReadsFactor(t *testing.T) {
	t.Setenv(envFactor, "0.35")
	p := FromEnv()
	if p.Factor != 0.35 {
		t.Fatalf("expected 0.35, got %v", p.Factor)
	}
}

func TestFromEnvIgnoresInvalidFactor(t *testing.T) {
	t.Setenv(envFactor, "not-a-number")
	p := FromEnv()
	if p.Factor != DefaultPolicy().Factor {
		t.Fatalf("expected default factor, got %v", p.Factor)
	}
}
