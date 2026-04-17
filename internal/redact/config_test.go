package redact_test

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/redact"
)

func TestFromConfigIncludesDefaults(t *testing.T) {
	r := redact.FromConfig(redact.Config{})
	for _, k := range redact.DefaultSensitiveKeys {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive by default", k)
		}
	}
}

func TestFromConfigMergesExtraKeys(t *testing.T) {
	r := redact.FromConfig(redact.Config{ExtraKeys: []string{"MY_CUSTOM_SECRET"}})
	if !r.IsSensitive("MY_CUSTOM_SECRET") {
		t.Error("expected MY_CUSTOM_SECRET to be sensitive")
	}
}

func TestFromConfigCustomMask(t *testing.T) {
	r := redact.FromConfig(redact.Config{Mask: "---"})
	got := r.Value("PASSWORD", "hunter2")
	if got != "---" {
		t.Errorf("expected ---, got %q", got)
	}
}

func TestFromEnvAddsExtraKeys(t *testing.T) {
	t.Setenv("VAULTPULL_REDACT_KEYS", "MY_KEY, ANOTHER_KEY")
	r := redact.FromEnv()
	if !r.IsSensitive("MY_KEY") {
		t.Error("expected MY_KEY from env to be sensitive")
	}
	if !r.IsSensitive("ANOTHER_KEY") {
		t.Error("expected ANOTHER_KEY from env to be sensitive")
	}
}

func TestFromEnvEmptyVarUsesDefaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_REDACT_KEYS")
	r := redact.FromEnv()
	if !r.IsSensitive("TOKEN") {
		t.Error("expected TOKEN to be sensitive by default")
	}
}
