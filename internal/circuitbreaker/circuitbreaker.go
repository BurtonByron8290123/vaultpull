// Package circuitbreaker provides a simple circuit breaker for Vault requests.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is open.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Policy configures the circuit breaker.
type Policy struct {
	MaxFailures  int
	OpenDuration time.Duration
}

// DefaultPolicy returns a sensible default policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxFailures:  5,
		OpenDuration: 30 * time.Second,
	}
}

// Breaker is a circuit breaker instance.
type Breaker struct {
	mu        sync.Mutex
	policy    Policy
	failures  int
	state     State
	openUntil time.Time
}

// New creates a new Breaker with the given policy.
func New(p Policy) *Breaker {
	return &Breaker{policy: p}
}

// Allow returns nil if the request is allowed, or ErrOpen if the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateOpen:
		if time.Now().After(b.openUntil) {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	}
	return nil
}

// RecordSuccess resets the failure count and closes the circuit.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure count and may open the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.policy.MaxFailures {
		b.state = StateOpen
		b.openUntil = time.Now().Add(b.policy.OpenDuration)
	}
}

// State returns the current state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
