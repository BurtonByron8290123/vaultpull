package expire

import "errors"

// Sentinel errors returned by the expire package.
var (
	ErrInvalidTTL = errors.New("expire: TTL must be greater than zero")
	ErrNilClock   = errors.New("expire: ClockFunc must not be nil")
)
