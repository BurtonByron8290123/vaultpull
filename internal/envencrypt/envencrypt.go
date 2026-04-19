// Package envencrypt provides AES-GCM encryption and decryption for env file values.
package envencrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidKey is returned when the key length is not 16, 24, or 32 bytes.
var ErrInvalidKey = errors.New("envencrypt: key must be 16, 24, or 32 bytes")

// ErrDecryptFailed is returned when decryption or authentication fails.
var ErrDecryptFailed = errors.New("envencrypt: decryption failed")

// Encrypter encrypts and decrypts string values using AES-GCM.
type Encrypter struct {
	key []byte
}

// New returns an Encrypter using the supplied key.
// key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func New(key []byte) (*Encrypter, error) {
	switch len(key) {
	case 16, 24, 32:
	default:
		return nil, ErrInvalidKey
	}
	k := make([]byte, len(key))
	copy(k, key)
	return &Encrypter{key: k}, nil
}

// Encrypt encrypts plaintext and returns a base64-encoded ciphertext string.
func (e *Encrypter) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decodes and decrypts a base64-encoded ciphertext produced by Encrypt.
func (e *Encrypter) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrDecryptFailed
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", ErrDecryptFailed
	}
	plaintext, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", ErrDecryptFailed
	}
	return string(plaintext), nil
}

// EncryptMap encrypts all values in m, returning a new map.
func (e *Encrypter) EncryptMap(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		enc, err := e.Encrypt(v)
		if err != nil {
			return nil, err
		}
		out[k] = enc
	}
	return out, nil
}
