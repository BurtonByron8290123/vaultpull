package scope

import (
	"os"
	"strings"
)

const (
	envAllow = "VAULTPULL_SCOPE_ALLOW"
	envDeny  = "VAULTPULL_SCOPE_DENY"
)

// FromEnv builds a Policy from VAULTPULL_SCOPE_ALLOW and VAULTPULL_SCOPE_DENY
// environment variables. Each variable accepts a comma-separated list of path
// prefixes. Empty variables are silently ignored.
func FromEnv() Policy {
	return Policy{
		Allow: splitEnv(os.Getenv(envAllow)),
		Deny:  splitEnv(os.Getenv(envDeny)),
	}
}

// FromEnvWithOverrides builds a Policy like FromEnv but allows the caller to
// supply explicit allow and deny lists that take precedence over environment
// variables when non-nil.
func FromEnvWithOverrides(allow, deny []string) Policy {
	p := FromEnv()
	if allow != nil {
		p.Allow = allow
	}
	if deny != nil {
		p.Deny = deny
	}
	return p
}

func splitEnv(val string) []string {
	if strings.TrimSpace(val) == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
