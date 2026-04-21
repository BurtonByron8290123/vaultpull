// Package envstamp embeds a metadata stamp (version, timestamp, source path)
// into an env map so consumers can trace where secrets originated.
package envstamp

import (
	"fmt"
	"time"
)

const (
	KeyVersion   = "VAULTPULL_STAMP_VERSION"
	KeyTimestamp = "VAULTPULL_STAMP_TIMESTAMP"
	KeySource    = "VAULTPULL_STAMP_SOURCE"
)

// Policy controls stamp behaviour.
type Policy struct {
	// Enabled turns stamping on or off.
	Enabled bool
	// Version is an arbitrary version string embedded in the stamp.
	Version string
	// Source is the Vault path that produced the secrets.
	Source string
	// Clock is used to obtain the current time; defaults to time.Now.
	Clock func() time.Time
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		Enabled: true,
		Version: "1",
		Clock:   time.Now,
	}
}

// Stamper applies metadata stamps to env maps.
type Stamper struct {
	p Policy
}

// New creates a Stamper from the supplied Policy.
func New(p Policy) (*Stamper, error) {
	if p.Clock == nil {
		p.Clock = time.Now
	}
	return &Stamper{p: p}, nil
}

// Apply adds stamp keys to a copy of m and returns it.
// If Policy.Enabled is false the original map is returned unchanged.
func (s *Stamper) Apply(m map[string]string) map[string]string {
	if !s.p.Enabled {
		return m
	}
	out := make(map[string]string, len(m)+3)
	for k, v := range m {
		out[k] = v
	}
	out[KeyVersion] = s.p.Version
	out[KeyTimestamp] = s.p.Clock().UTC().Format(time.RFC3339)
	out[KeySource] = s.p.Source
	return out
}

// Strip removes stamp keys from a copy of m.
func Strip(m map[string]string) map[string]string {
	keys := []string{KeyVersion, KeyTimestamp, KeySource}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	for _, k := range keys {
		delete(out, k)
	}
	return out
}

// Summary returns a human-readable description of the stamp embedded in m.
func Summary(m map[string]string) string {
	v := m[KeyVersion]
	ts := m[KeyTimestamp]
	src := m[KeySource]
	if v == "" && ts == "" && src == "" {
		return "(no stamp)"
	}
	return fmt.Sprintf("version=%s timestamp=%s source=%s", v, ts, src)
}
