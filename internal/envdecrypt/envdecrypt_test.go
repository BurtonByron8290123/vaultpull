package envdecrypt_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"testing"

	"github.com/your-org/vaultpull/internal/envdecrypt"
)

func key32() []byte {
	k := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		panic(err)
	}
	return k
}

func encrypt(key []byte, plain string) string {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)
	ct := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return "enc:" + base64.StdEncoding.EncodeToString(ct)
}

func TestNewRejectsInvalidKey(t *testing.T) {
	_, err := envdecrypt.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestValuePassesThroughPlaintext(t *testing.T) {
	d, err := envdecrypt.New(key32())
	if err != nil {
		t.Fatal(err)
	}
	got, err := d.Value("hello")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Fatalf("want hello, got %q", got)
	}
}

func TestValueDecryptsEncPrefix(t *testing.T) {
	key := key32()
	d, _ := envdecrypt.New(key)
	enc := encrypt(key, "s3cr3t")
	got, err := d.Value(enc)
	if err != nil {
		t.Fatal(err)
	}
	if got != "s3cr3t" {
		t.Fatalf("want s3cr3t, got %q", got)
	}
}

func TestValueRejectsWrongKey(t *testing.T) {
	key := key32()
	enc := encrypt(key, "secret")
	d, _ := envdecrypt.New(key32()) // different key
	_, err := d.Value(enc)
	if err == nil {
		t.Fatal("expected decryption error with wrong key")
	}
}

func TestMapDecryptsOnlyEncValues(t *testing.T) {
	key := key32()
	d, _ := envdecrypt.New(key)
	input := map[string]string{
		"PLAIN": "visible",
		"SECRET": encrypt(key, "topsecret"),
	}
	out, err := d.Map(input)
	if err != nil {
		t.Fatal(err)
	}
	if out["PLAIN"] != "visible" {
		t.Errorf("PLAIN changed: %q", out["PLAIN"])
	}
	if out["SECRET"] != "topsecret" {
		t.Errorf("SECRET wrong: %q", out["SECRET"])
	}
}

func TestFromEnvEmptyVarReturnsNil(t *testing.T) {
	t.Setenv("VAULTPULL_DECRYPT_KEY", "")
	d, err := envdecrypt.FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if d != nil {
		t.Fatal("expected nil decrypter when env var is absent")
	}
}

func TestFromEnvLoadsValidKey(t *testing.T) {
	key := key32()
	t.Setenv("VAULTPULL_DECRYPT_KEY", hex.EncodeToString(key))
	d, err := envdecrypt.FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if d == nil {
		t.Fatal("expected non-nil decrypter")
	}
}
