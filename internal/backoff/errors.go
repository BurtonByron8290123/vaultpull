package backoff

import "errors"

// Sentinel errors returned by Policy.Validate.
var (
	ErrInvalidInitial    = errors.New("backoff: initial interval must be > 0")
	ErrInvalidMultiplier = errors.New("backoff: multiplier must be >= 1")
	ErrInvalidMaxInterval = errors.New("backoff: max interval must be >= initial interval")
	ErrInvalidJitter     = errors.New("backoff: jitter must be in [0, 1]")
)
