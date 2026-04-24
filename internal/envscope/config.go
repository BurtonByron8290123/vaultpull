package envscope

import (
	"os"
	"strings"
)

const (
	envAllow = "VAULTPULL_SCOPE_ALLOW"
	envDeny  = "VAULTPULL_SCOPE_DENY"
)

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_SCOPE_ALLOW – comma-separated list of allowed key prefixes
//	VAULTPULL_SCOPE_DENY  – comma-separated list of denied key prefixes
func FromEnv() Policy {
	return Policy{
		Allow: splitCSV(os.Getenv(envAllow)),
		Deny:  splitCSV(os.Getenv(envDeny)),
	}
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
