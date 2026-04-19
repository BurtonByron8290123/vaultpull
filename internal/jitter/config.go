package jitter

import (
	"os"
	"strconv"
)

const (
	envFactor = "VAULTPULL_JITTER_FACTOR"
)

// FromEnv builds a Policy from environment variables, falling back to
// DefaultPolicy() for any value that is absent or invalid.
//
//	VAULTPULL_JITTER_FACTOR  float in (0,1]  default 0.20
func FromEnv() Policy {
	p := DefaultPolicy()

	if v := os.Getenv(envFactor); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			tmp := Policy{Factor: f}
			if tmp.Validate() == nil {
				p.Factor = f
			}
		}
	}

	return p
}
