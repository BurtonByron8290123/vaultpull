package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// MarshalOptions controls how key-value pairs are serialized.
type MarshalOptions struct {
	// SortKeys ensures deterministic output order.
	SortKeys bool
	// QuoteValues wraps all values in double quotes.
	QuoteValues bool
}

// Marshal serializes a map of key-value pairs into .env file bytes.
// Keys containing characters outside [A-Za-z0-9_] are skipped.
func Marshal(vars map[string]string, opts MarshalOptions) ([]byte, error) {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		if !isValidKey(k) {
			return nil, fmt.Errorf("env: invalid key %q", k)
		}
		keys = append(keys, k)
	}

	if opts.SortKeys {
		sort.Strings(keys)
	}

	var sb strings.Builder
	for _, k := range keys {
		v := vars[k]
		if opts.QuoteValues {
			v = quoteValue(v)
		} else if needsQuoting(v) {
			v = quoteValue(v)
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}

	return []byte(sb.String()), nil
}

// WriteFile marshals vars and writes them to path with the given file mode.
func WriteFile(path string, vars map[string]string, opts MarshalOptions, perm os.FileMode) error {
	data, err := Marshal(vars, opts)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

func isValidKey(k string) bool {
	if len(k) == 0 {
		return false
	}
	for _, r := range k {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

func needsQuoting(v string) bool {
	return strings.ContainsAny(v, " \t\n\r#$'\"\\")
}

func quoteValue(v string) string {
	v = strings.ReplaceAll(v, `\`, `\\`)
	v = strings.ReplaceAll(v, `"`, `\"`)
	return `"` + v + `"`
}
