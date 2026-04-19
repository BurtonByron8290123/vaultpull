package envlock_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultpull/internal/envlock"
)

func tempEnvFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, ".env")
}

func TestAcquireCreatesLockFile(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, time.Second)
	if err := l.Acquire(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Release() //nolint:errcheck

	if _, err := os.Stat(l.LockPath()); err != nil {
		t.Fatalf("lock file not found: %v", err)
	}
}

func TestReleaseRemovesLockFile(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, time.Second)
	_ = l.Acquire()

	if err := l.Release(); err != nil {
		t.Fatalf("release error: %v", err)
	}
	if _, err := os.Stat(l.LockPath()); !os.IsNotExist(err) {
		t.Fatal("expected lock file to be removed")
	}
}

func TestReleaseNoopWhenLockAbsent(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, time.Second)
	if err := l.Release(); err != nil {
		t.Fatalf("expected no error releasing non-existent lock, got: %v", err)
	}
}

func TestAcquireTimesOutWhenLocked(t *testing.T) {
	path := tempEnvFile(t)
	owner := envlock.New(path, time.Second)
	_ = owner.Acquire()
	defer owner.Release() //nolint:errcheck

	contender := envlock.New(path, 100*time.Millisecond)
	err := contender.Acquire()
	if err == nil {
		_ = contender.Release()
		t.Fatal("expected timeout error")
	}
	if err != envlock.ErrLockTimeout {
		t.Fatalf("expected ErrLockTimeout, got: %v", err)
	}
}

func TestLockFilePermissions(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, time.Second)
	_ = l.Acquire()
	defer l.Release() //nolint:errcheck

	info, err := os.Stat(l.LockPath())
	if err != nil {
		t.Fatalf("stat lock file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Fatalf("expected 0600, got %04o", perm)
	}
}
