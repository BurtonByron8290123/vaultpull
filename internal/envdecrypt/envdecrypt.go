// Package envdecrypt provides on-the-fly decryption of encrypted .env values
// using AES-GCM. Values prefixed with "enc:" are treated as base64-encoded
// ciphertext and decrypted before being written to the environment map.
package envdecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const encPrefix = "enc:"

// ErrNotEncrypted is returned when a value does not carry the enc: prefix.
var ErrNotEncrypted = errors.New("envdecrypt: value is not encrypted")

// Decrypter decrypts individual values or entire env maps.
type Decrypter struct {
	gcm cipher.AEAD
}

// New creates a Decrypter using the supplied key (16, 24, or 32 bytes).
func New(key []byte) (*Decrypter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("envdecrypt: invalid key: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("envdecrypt: gcm init: %w", err)
	}
	return &Decrypter{gcm: gcm}, nil
}

// Value decrypts a single value. If the value does not start with "enc:"
// it is returned unchanged.
func (d *Decrypter) Value(v string) (string, error) {
	if !strings.HasPrefix(v, encPrefix) {
		return v, nil
	}
	ciphertext, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(v, encPrefix))
	if err != nil {
		return "", fmt.Errorf("envdecrypt: base64 decode: %w", err)
	}
	ns := d.gcm.NonceSize()
	if len(ciphertext) < ns {
		return "", errors.New("envdecrypt: ciphertext too short")
	}
	nonce, ct := ciphertext[:ns], ciphertext[ns:]
	plain, err := d.gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("envdecrypt: decrypt: %w", err)
	}
	return string(plain), nil
}

// Map decrypts all values in m that carry the enc: prefix, returning a new
// map. Non-encrypted values are copied as-is.
func (d *Decrypter) Map(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		dec, err := d.Value(v)
		if err != nil {
			return nil, fmt.Errorf("envdecrypt: key %q: %w", k, err)
		}
		out[k] = dec
	}
	return out, nil
}
