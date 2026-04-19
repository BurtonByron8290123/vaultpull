package envsign_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envsign"
)

func newSigner(t *testing.T) *envsign.Signer {
	t.Helper()
	s, err := envsign.New([]byte("supersecretkey16"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNewRejectsShortKey(t *testing.T) {
	_, err := envsign.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestSignAndVerifyRoundTrip(t *testing.T) {
	s := newSigner(t)
	data := []byte("FOO=bar\nBAZ=qux\n")
	sig := s.Sign(data)
	if err := s.Verify(data, sig); err != nil {
		t.Fatalf("Verify: %v", err)
	}
}

func TestVerifyDetectsTampering(t *testing.T) {
	s := newSigner(t)
	data := []byte("FOO=bar\n")
	sig := s.Sign(data)
	if err := s.Verify([]byte("FOO=evil\n"), sig); err == nil {
		t.Fatal("expected mismatch error")
	}
}

func TestWriteAndVerifyFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("KEY=val\n"), 0600); err != nil {
		t.Fatal(err)
	}
	s := newSigner(t)
	if err := s.WriteSignature(envFile); err != nil {
		t.Fatalf("WriteSignature: %v", err)
	}
	if err := s.VerifyFile(envFile); err != nil {
		t.Fatalf("VerifyFile: %v", err)
	}
}

func TestVerifyFileMissingSig(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte("KEY=val\n"), 0600)
	s := newSigner(t)
	if err := s.VerifyFile(envFile); err != envsign.ErrMissingSig {
		t.Fatalf("expected ErrMissingSig, got %v", err)
	}
}

func TestVerifyFileDetectsTampering(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte("KEY=original\n"), 0600)
	s := newSigner(t)
	_ = s.WriteSignature(envFile)
	// tamper
	_ = os.WriteFile(envFile, []byte("KEY=tampered\n"), 0600)
	if err := s.VerifyFile(envFile); err == nil {
		t.Fatal("expected signature mismatch")
	}
}

func TestSigFilePermissions(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte("K=v\n"), 0600)
	s := newSigner(t)
	_ = s.WriteSignature(envFile)
	info, err := os.Stat(envFile + ".sig")
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
