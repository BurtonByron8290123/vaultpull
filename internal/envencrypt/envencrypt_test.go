package envencrypt_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envencrypt"
)

func key32() []byte { return []byte("12345678901234567890123456789012") }

func TestNewRejectsInvalidKeyLength(t *testing.T) {
	_, err := envencrypt.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestNewAcceptsValidKeyLengths(t *testing.T) {
	for _, n := range []int{16, 24, 32} {
		_, err := envencrypt.New(make([]byte, n))
		if err != nil {
			t.Fatalf("unexpected error for key length %d: %v", n, err)
		}
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	e, _ := envencrypt.New(key32())
	plain := "super-secret-value"
	enc, err := e.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	got, err := e.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if got != plain {
		t.Fatalf("want %q got %q", plain, got)
	}
}

func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	e, _ := envencrypt.New(key32())
	a, _ := e.Encrypt("value")
	b, _ := e.Encrypt("value")
	if a == b {
		t.Fatal("expected different ciphertexts due to random nonce")
	}
}

func TestDecryptFailsOnTamperedData(t *testing.T) {
	e, _ := envencrypt.New(key32())
	enc, _ := e.Encrypt("hello")
	tampered := enc[:len(enc)-4] + "XXXX"
	_, err := e.Decrypt(tampered)
	if err == nil {
		t.Fatal("expected decryption error on tampered data")
	}
}

func TestDecryptFailsOnGarbage(t *testing.T) {
	e, _ := envencrypt.New(key32())
	_, err := e.Decrypt("not-base64!!!")
	if err == nil {
		t.Fatal("expected error on invalid base64")
	}
}

func TestEncryptMapRoundTrip(t *testing.T) {
	e, _ := envencrypt.New(key32())
	m := map[string]string{"DB_PASS": "hunter2", "API_KEY": "abc123"}
	enc, err := e.EncryptMap(m)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	for k, v := range enc {
		if strings.Contains(v, m[k]) {
			t.Errorf("key %s: ciphertext should not contain plaintext", k)
		}
		dec, err := e.Decrypt(v)
		if err != nil {
			t.Fatalf("Decrypt key %s: %v", k, err)
		}
		if dec != m[k] {
			t.Errorf("key %s: want %q got %q", k, m[k], dec)
		}
	}
}
