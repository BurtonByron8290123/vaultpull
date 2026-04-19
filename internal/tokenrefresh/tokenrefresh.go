// Package tokenrefresh provides automatic Vault token renewal before expiry.
package tokenrefresh

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Policy controls when and how tokens are refreshed.
type Policy struct {
	// RenewThreshold is the remaining TTL below which renewal is triggered.
	RenewThreshold time.Duration
	// MaxRetries is the number of renewal attempts before giving up.
	MaxRetries int
}

// DefaultPolicy returns a sensible default renewal policy.
func DefaultPolicy() Policy {
	return Policy{
		RenewThreshold: 5 * time.Minute,
		MaxRetries:     3,
	}
}

func (p Policy) validate() error {
	if p.RenewThreshold <= 0 {
		return fmt.Errorf("tokenrefresh: RenewThreshold must be positive")
	}
	if p.MaxRetries < 0 {
		return fmt.Errorf("tokenrefresh: MaxRetries must be non-negative")
	}
	return nil
}

// Refresher wraps a Vault client and renews the token when needed.
type Refresher struct {
	client *vaultapi.Client
	policy Policy
}

// New creates a Refresher with the given client and policy.
func New(client *vaultapi.Client, policy Policy) (*Refresher, error) {
	if err := policy.validate(); err != nil {
		return nil, err
	}
	return &Refresher{client: client, policy: policy}, nil
}

// EnsureValid checks the current token TTL and renews it if below the threshold.
func (r *Refresher) EnsureValid(ctx context.Context) error {
	secret, err := r.client.Auth().Token().LookupSelf()
	if err != nil {
		return fmt.Errorf("tokenrefresh: lookup failed: %w", err)
	}

	ttlRaw, ok := secret.Data["ttl"]
	if !ok {
		// Non-expiring token — nothing to do.
		return nil
	}

	ttlSec, ok := ttlRaw.(float64)
	if !ok {
		return fmt.Errorf("tokenrefresh: unexpected ttl type %T", ttlRaw)
	}

	ttl := time.Duration(ttlSec) * time.Second
	if ttl > r.policy.RenewThreshold {
		return nil
	}

	return r.renewWithRetry(ctx)
}

func (r *Refresher) renewWithRetry(ctx context.Context) error {
	var last error
	for i := 0; i <= r.policy.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		_, err := r.client.Auth().Token().RenewSelf(0)
		if err == nil {
			return nil
		}
		last = err
	}
	return fmt.Errorf("tokenrefresh: renewal failed after %d attempts: %w", r.policy.MaxRetries, last)
}
