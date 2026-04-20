package envfreeze_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envfreeze"
)

func newFreezer(t *testing.T, keys ...string) *envfreeze.Freezer {
	t.Helper()
	f, err := envfreeze.New(envfreeze.Policy{Keys: keys})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return f
}

func TestCheckNoFrozenKeysAlwaysPasses(t *testing.T) {
	f := newFreezer(t)
	current := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "changed"}
	if err := f.Check(current, incoming); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckFrozenKeyUnchangedPasses(t *testing.T) {
	f := newFreezer(t, "FOO")
	current := map[string]string{"FOO": "same"}
	incoming := map[string]string{"FOO": "same"}
	if err := f.Check(current, incoming); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckFrozenKeyChangedReturnsError(t *testing.T) {
	f := newFreezer(t, "SECRET_KEY")
	current := map[string]string{"SECRET_KEY": "old"}
	incoming := map[string]string{"SECRET_KEY": "new"}
	err := f.Check(current, incoming)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, envfreeze.ErrFrozenKey) {
		t.Fatalf("expected ErrFrozenKey, got %v", err)
	}
	if !strings.Contains(err.Error(), "SECRET_KEY") {
		t.Fatalf("error should mention key name, got: %v", err)
	}
}

func TestCheckFrozenKeyNotPresentInCurrentPasses(t *testing.T) {
	// Key is frozen but not yet in the current env — first write is allowed.
	f := newFreezer(t, "NEW_KEY")
	current := map[string]string{}
	incoming := map[string]string{"NEW_KEY": "value"}
	if err := f.Check(current, incoming); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckDryRunDoesNotReturnError(t *testing.T) {
	f, err := envfreeze.New(envfreeze.Policy{Keys: []string{"LOCKED"}, DryRun: true})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	current := map[string]string{"LOCKED": "old"}
	incoming := map[string]string{"LOCKED": "new"}
	if err := f.Check(current, incoming); err != nil {
		t.Fatalf("dry-run should not return error, got: %v", err)
	}
}

func TestIsFrozenCaseInsensitive(t *testing.T) {
	f := newFreezer(t, "db_password")
	if !f.IsFrozen("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be frozen")
	}
	if f.IsFrozen("OTHER") {
		t.Error("OTHER should not be frozen")
	}
}

func TestNewRejectsBlankKey(t *testing.T) {
	_, err := envfreeze.New(envfreeze.Policy{Keys: []string{"  "}})
	if err == nil {
		t.Fatal("expected error for blank key")
	}
}
