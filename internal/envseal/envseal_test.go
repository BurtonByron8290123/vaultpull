package envseal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/internal/envseal"
)

func newSealer(t *testing.T) *envseal.Sealer {
	t.Helper()
	p := envseal.Policy{
		EncryptionKey: make([]byte, 32),
		SigningKey:    make([]byte, 32),
	}
	s, err := envseal.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func tempFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "sealed.env")
}

func TestSealUnsealRoundTrip(t *testing.T) {
	s := newSealer(t)
	env := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}
	dst := tempFile(t)

	if err := s.Seal(env, dst); err != nil {
		t.Fatalf("Seal: %v", err)
	}
	got, err := s.Unseal(dst)
	if err != nil {
		t.Fatalf("Unseal: %v", err)
	}
	for k, v := range env {
		if got[k] != v {
			t.Errorf("key %s: want %q, got %q", k, v, got[k])
		}
	}
}

func TestUnsealDetectsTampering(t *testing.T) {
	s := newSealer(t)
	dst := tempFile(t)

	if err := s.Seal(map[string]string{"X": "1"}, dst); err != nil {
		t.Fatalf("Seal: %v", err)
	}

	// corrupt the sealed file
	data, _ := os.ReadFile(dst)
	data[len(data)-1] ^= 0xFF
	_ = os.WriteFile(dst, data, 0o600)

	_, err := s.Unseal(dst)
	if err == nil {
		t.Fatal("expected error on tampered file, got nil")
	}
}

func TestNewRejectsInvalidEncryptionKey(t *testing.T) {
	_, err := envseal.New(envseal.Policy{
		EncryptionKey: []byte{1, 2, 3}, // invalid length
		SigningKey:    make([]byte, 32),
	})
	if err == nil {
		t.Fatal("expected error for bad encryption key")
	}
}

func TestFromEnvMissingEncKeyReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_SEAL_ENC_KEY", "")
	t.Setenv("VAULTPULL_SEAL_SIGN_KEY", "")
	_, err := envseal.FromEnv()
	if err == nil {
		t.Fatal("expected error when env vars absent")
	}
}

func TestFromEnvReadsKeys(t *testing.T) {
	encHex := "0000000000000000000000000000000000000000000000000000000000000000" // 32 bytes
	signHex := "0000000000000000000000000000000000000000000000000000000000000000"
	t.Setenv("VAULTPULL_SEAL_ENC_KEY", encHex)
	t.Setenv("VAULTPULL_SEAL_SIGN_KEY", signHex)
	p, err := envseal.FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if len(p.EncryptionKey) != 32 {
		t.Errorf("enc key len: want 32, got %d", len(p.EncryptionKey))
	}
	if len(p.SigningKey) != 32 {
		t.Errorf("sign key len: want 32, got %d", len(p.SigningKey))
	}
}
