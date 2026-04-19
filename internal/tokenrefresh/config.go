package tokenrefresh

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	envThresholdSec = "VAULTPULL_TOKEN_RENEW_THRESHOLD_SEC"
	envMaxRetries   = "VAULTPULL_TOKEN_RENEW_MAX_RETRIES"
)

// FromEnv builds a Policy from environment variables, falling back to defaults.
func FromEnv() (Policy, error) {
	p := DefaultPolicy()

	if v := os.Getenv(envThresholdSec); v != "" {
		sec, err := strconv.Atoi(v)
		if err != nil || sec <= 0 {
			return p, fmt.Errorf("tokenrefresh: invalid %s: %q", envThresholdSec, v)
		}
		p.RenewThreshold = time.Duration(sec) * time.Second
	}

	if v := os.Getenv(envMaxRetries); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return p, fmt.Errorf("tokenrefresh: invalid %s: %q", envMaxRetries, v)
		}
		p.MaxRetries = n
	}

	return p, nil
}
