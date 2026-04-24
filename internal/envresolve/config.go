package envresolve

import (
	"os"
	"strings"
)

const (
	envAllowFallback  = "VAULTPULL_RESOLVE_ENV_FALLBACK"
	envErrorOnMissing = "VAULTPULL_RESOLVE_ERROR_ON_MISSING"
)

// FromEnv constructs a Policy from environment variables.
//
//	VAULTPULL_RESOLVE_ENV_FALLBACK    – "true"/"false" (default true)
//	VAULTPULL_RESOLVE_ERROR_ON_MISSING – "true"/"false" (default false)
func FromEnv() Policy {
	p := DefaultPolicy()
	if v := strings.TrimSpace(os.Getenv(envAllowFallback)); v != "" {
		p.AllowEnvFallback = v == "true"
	}
	if v := strings.TrimSpace(os.Getenv(envErrorOnMissing)); v != "" {
		p.ErrorOnMissing = v == "true"
	}
	return p
}
