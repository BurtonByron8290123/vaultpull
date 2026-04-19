// Package preflight runs sanity checks before a pull operation begins.
package preflight

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Check represents a single preflight check.
type Check struct {
	Name string
	Fn   func() error
}

// Runner executes a set of preflight checks.
type Runner struct {
	checks []Check
}

// New returns a Runner with the provided checks.
func New(checks ...Check) *Runner {
	return &Runner{checks: checks}
}

// Default returns a Runner with the standard set of checks.
func Default(vaultAddr, token, outputPath string) *Runner {
	return New(
		CheckVaultAddr(vaultAddr),
		CheckToken(token),
		CheckOutputWritable(outputPath),
	)
}

// Run executes all checks and returns a combined error if any fail.
func (r *Runner) Run() error {
	var errs []string
	for _, c := range r.checks {
		if err := c.Fn(); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", c.Name, err))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// CheckVaultAddr verifies the VAULT_ADDR is non-empty.
func CheckVaultAddr(addr string) Check {
	return Check{
		Name: "vault_addr",
		Fn: func() error {
			if strings.TrimSpace(addr) == "" {
				return errors.New("VAULT_ADDR is not set")
			}
			return nil
		},
	}
}

// CheckToken verifies the Vault token is non-empty.
func CheckToken(token string) Check {
	return Check{
		Name: "vault_token",
		Fn: func() error {
			if strings.TrimSpace(token) == "" {
				return errors.New("VAULT_TOKEN is not set")
			}
			return nil
		},
	}
}

// CheckOutputWritable verifies the output path's directory is writable.
func CheckOutputWritable(path string) Check {
	return Check{
		Name: "output_writable",
		Fn: func() error {
			if path == "" {
				return errors.New("output path is empty")
			}
			dir := dirOf(path)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("cannot create output directory %q: %w", dir, err)
			}
			tmp, err := os.CreateTemp(dir, ".vaultpull-probe-*")
			if err != nil {
				return fmt.Errorf("directory %q is not writable: %w", dir, err)
			}
			tmp.Close()
			os.Remove(tmp.Name())
			return nil
		},
	}
}

func dirOf(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}
