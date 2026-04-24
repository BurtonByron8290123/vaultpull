package envtemplate

import (
	"os"
	"strconv"
)

const (
	envErrorOnMissing = "VAULTPULL_TEMPLATE_ERROR_ON_MISSING"
	envLeftDelim      = "VAULTPULL_TEMPLATE_LEFT_DELIM"
	envRightDelim     = "VAULTPULL_TEMPLATE_RIGHT_DELIM"
)

// FromEnv builds a Policy from environment variables, falling back to
// DefaultPolicy() for any value that is absent or invalid.
func FromEnv() Policy {
	p := DefaultPolicy()

	if raw := os.Getenv(envErrorOnMissing); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			p.ErrorOnMissing = v
		}
	}
	if v := os.Getenv(envLeftDelim); v != "" {
		p.LeftDelim = v
	}
	if v := os.Getenv(envRightDelim); v != "" {
		p.RightDelim = v
	}
	return p
}
