package envstamp_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envstamp"
)

func fixedClock(ts time.Time) func() time.Time {
	return func() time.Time { return ts }
}

func TestApplyAddsStampKeys(t *testing.T) {
	p := envstamp.DefaultPolicy()
	p.Version = "42"
	p.Source = "secret/myapp"
	p.Clock = fixedClock(time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC))

	s, err := envstamp.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	out := s.Apply(map[string]string{"FOO": "bar"})

	if out[envstamp.KeyVersion] != "42" {
		t.Errorf("version: got %q", out[envstamp.KeyVersion])
	}
	if out[envstamp.KeySource] != "secret/myapp" {
		t.Errorf("source: got %q", out[envstamp.KeySource])
	}
	if !strings.HasPrefix(out[envstamp.KeyTimestamp], "2024-01-15") {
		t.Errorf("timestamp: got %q", out[envstamp.KeyTimestamp])
	}
	if out["FOO"] != "bar" {
		t.Errorf("original key lost")
	}
}

func TestApplyDisabledReturnsOriginal(t *testing.T) {
	p := envstamp.DefaultPolicy()
	p.Enabled = false

	s, _ := envstamp.New(p)
	in := map[string]string{"A": "1"}
	out := s.Apply(in)

	if _, ok := out[envstamp.KeyVersion]; ok {
		t.Error("stamp key present when disabled")
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	p := envstamp.DefaultPolicy()
	s, _ := envstamp.New(p)

	in := map[string]string{"X": "y"}
	_ = s.Apply(in)

	if _, ok := in[envstamp.KeyVersion]; ok {
		t.Error("input map was mutated")
	}
}

func TestStripRemovesStampKeys(t *testing.T) {
	m := map[string]string{
		"APP_KEY":              "value",
		envstamp.KeyVersion:   "1",
		envstamp.KeyTimestamp: "2024-01-01T00:00:00Z",
		envstamp.KeySource:    "secret/app",
	}

	out := envstamp.Strip(m)

	if _, ok := out[envstamp.KeyVersion]; ok {
		t.Error("KeyVersion still present after Strip")
	}
	if out["APP_KEY"] != "value" {
		t.Error("non-stamp key removed by Strip")
	}
}

func TestSummaryNoStamp(t *testing.T) {
	s := envstamp.Summary(map[string]string{"FOO": "bar"})
	if s != "(no stamp)" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummaryContainsFields(t *testing.T) {
	m := map[string]string{
		envstamp.KeyVersion:   "3",
		envstamp.KeyTimestamp: "2024-06-01T10:00:00Z",
		envstamp.KeySource:    "secret/prod",
	}
	s := envstamp.Summary(m)
	for _, want := range []string{"version=3", "secret/prod", "2024-06-01"} {
		if !strings.Contains(s, want) {
			t.Errorf("summary %q missing %q", s, want)
		}
	}
}
