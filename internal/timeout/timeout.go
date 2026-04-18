package timeout

import (
	"context"
	"fmt"
	"time"
)

// Policy defines timeout behaviour for Vault operations.
type Policy struct {
	Dial    time.Duration
	Read    time.Duration
	Overall time.Duration
}

// DefaultPolicy returns sensible production defaults.
func DefaultPolicy() Policy {
	return Policy{
		Dial:    5 * time.Second,
		Read:    10 * time.Second,
		Overall: 30 * time.Second,
	}
}

// Validate returns an error if any timeout value is non-positive.
func (p Policy) Validate() error {
	if p.Dial <= 0 {
		return fmt.Errorf("timeout: dial must be positive, got %v", p.Dial)
	}
	if p.Read <= 0 {
		return fmt.Errorf("timeout: read must be positive, got %v", p.Read)
	}
	if p.Overall <= 0 {
		return fmt.Errorf("timeout: overall must be positive, got %v", p.Overall)
	}
	return nil
}

// Apply wraps ctx with the overall deadline defined in the policy.
// The returned CancelFunc must always be called by the caller.
func (p Policy) Apply(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, p.Overall)
}
