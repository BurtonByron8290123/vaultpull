package envaudit

import (
	"os"
	"strconv"
)

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_AUDIT_PATH       – file path for the audit log (default: empty / disabled)
//	VAULTPULL_AUDIT_MASK       – "true"/"1" to mask values in log entries (default: true)
func FromEnv() Policy {
	p := Policy{
		AuditPath:  os.Getenv("VAULTPULL_AUDIT_PATH"),
		MaskValues: true,
	}
	if raw := os.Getenv("VAULTPULL_AUDIT_MASK"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			p.MaskValues = v
		}
	}
	return p
}
