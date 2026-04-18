package circuitbreaker

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_CB_MAX_FAILURES  – integer, default 5
//	VAULTPULL_CB_OPEN_SECONDS  – integer seconds, default 30
func FromEnv() (Policy, error) {
	p := DefaultPolicy()

	if v := os.Getenv("VAULTPULL_CB_MAX_FAILURES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return p, fmt.Errorf("circuitbreaker: invalid VAULTPULL_CB_MAX_FAILURES: %q", v)
		}
		p.MaxFailures = n
	}

	if v := os.Getenv("VAULTPULL_CB_OPEN_SECONDS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return p, fmt.Errorf("circuitbreaker: invalid VAULTPULL_CB_OPEN_SECONDS: %q", v)
		}
		p.OpenDuration = time.Duration(n) * time.Second
	}

	return p, nil
}
