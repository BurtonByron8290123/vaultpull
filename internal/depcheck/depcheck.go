// Package depcheck verifies that required environment variables are present
// after secrets are pulled from Vault, reporting any that are missing.
package depcheck

import "fmt"

// Result holds the outcome of a dependency check.
type Result struct {
	Missing []string
	Present []string
}

// Policy configures which keys are required.
type Policy struct {
	Required []string
}

// Checker verifies required keys exist in a secret map.
type Checker struct {
	policy Policy
}

// New returns a Checker for the given policy.
func New(p Policy) *Checker {
	return &Checker{policy: p}
}

// Check inspects secrets and returns a Result.
func (c *Checker) Check(secrets map[string]string) Result {
	var res Result
	for _, key := range c.policy.Required {
		if _, ok := secrets[key]; ok {
			res.Present = append(res.Present, key)
		} else {
			res.Missing = append(res.Missing, key)
		}
	}
	return res
}

// Err returns a non-nil error if any required keys are missing.
func (r Result) Err() error {
	if len(r.Missing) == 0 {
		return nil
	}
	return fmt.Errorf("depcheck: missing required keys: %v", r.Missing)
}

// OK returns true when no keys are missing.
func (r Result) OK() bool { return len(r.Missing) == 0 }
