package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/circuitbreaker"
)

func fastPolicy() circuitbreaker.Policy {
	return circuitbreaker.Policy{MaxFailures: 3, OpenDuration: 50 * time.Millisecond}
}

func TestAllowWhenClosed(t *testing.T) {
	b := circuitbreaker.New(fastPolicy())
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestOpensAfterMaxFailures(t *testing.T) {
	b := circuitbreaker.New(fastPolicy())
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected StateOpen")
	}
}

func TestHalfOpenAfterDuration(t *testing.T) {
	b := circuitbreaker.New(fastPolicy())
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if b.State() != circuitbreaker.StateHalfOpen {
		t.Fatal("expected StateHalfOpen")
	}
}

func TestRecordSuccessCloses(t *testing.T) {
	b := circuitbreaker.New(fastPolicy())
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	_ = b.Allow() // transition to half-open
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected StateClosed after success")
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after close, got %v", err)
	}
}

func TestBelowThresholdStaysClosed(t *testing.T) {
	b := circuitbreaker.New(fastPolicy())
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected StateClosed below threshold")
	}
}
