package timeout

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// FromEnv builds a Policy from environment variables, falling back to
// DefaultPolicy values for any variable that is absent or empty.
//
//	VAULTPULL_TIMEOUT_DIAL    – dial timeout  (seconds)
//	VAULTPULL_TIMEOUT_READ    – read timeout  (seconds)
//	VAULTPULL_TIMEOUT_OVERALL – overall timeout (seconds)
func FromEnv() (Policy, error) {
	p := DefaultPolicy()

	if v := os.Getenv("VAULTPULL_TIMEOUT_DIAL"); v != "" {
		d, err := parseSec("VAULTPULL_TIMEOUT_DIAL", v)
		if err != nil {
			return p, err
		}
		p.Dial = d
	}

	if v := os.Getenv("VAULTPULL_TIMEOUT_READ"); v != "" {
		d, err := parseSec("VAULTPULL_TIMEOUT_READ", v)
		if err != nil {
			return p, err
		}
		p.Read = d
	}

	if v := os.Getenv("VAULTPULL_TIMEOUT_OVERALL"); v != "" {
		d, err := parseSec("VAULTPULL_TIMEOUT_OVERALL", v)
		if err != nil {
			return p, err
		}
		p.Overall = d
	}

	return p, p.Validate()
}

func parseSec(name, v string) (time.Duration, error) {
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, fmt.Errorf("timeout: %s is not a valid number: %w", name, err)
	}
	return time.Duration(n * float64(time.Second)), nil
}
