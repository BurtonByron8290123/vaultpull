// Package envsign provides HMAC-based signing and verification for env file contents.
package envsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrInvalidKey      = errors.New("envsign: key must be at least 16 bytes")
	ErrSignatureMismatch = errors.New("envsign: signature mismatch")
	ErrMissingSig      = errors.New("envsign: signature file not found")
)

// Signer signs and verifies env file payloads.
type Signer struct {
	key []byte
}

// New returns a Signer using the provided HMAC key.
func New(key []byte) (*Signer, error) {
	if len(key) < 16 {
		return nil, ErrInvalidKey
	}
	return &Signer{key: key}, nil
}

// Sign computes an HMAC-SHA256 signature of data and returns the hex string.
func (s *Signer) Sign(data []byte) string {
	mac := hmac.New(sha256.New, s.key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify returns nil if sig matches the HMAC of data.
func (s *Signer) Verify(data []byte, sig string) error {
	expected := s.Sign(data)
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return ErrSignatureMismatch
	}
	return nil
}

// WriteSignature writes the HMAC signature of envFile contents to envFile + ".sig".
func (s *Signer) WriteSignature(envFile string) error {
	data, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("envsign: read file: %w", err)
	}
	sig := s.Sign(data)
	sigPath := sigFilePath(envFile)
	return os.WriteFile(sigPath, []byte(sig), 0600)
}

// VerifyFile reads envFile and its companion .sig file and verifies integrity.
func (s *Signer) VerifyFile(envFile string) error {
	sigPath := sigFilePath(envFile)
	sigBytes, err := os.ReadFile(sigPath)
	if errors.Is(err, os.ErrNotExist) {
		return ErrMissingSig
	}
	if err != nil {
		return fmt.Errorf("envsign: read sig: %w", err)
	}
	data, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("envsign: read file: %w", err)
	}
	return s.Verify(data, string(sigBytes))
}

func sigFilePath(envFile string) string {
	return filepath.Join(filepath.Dir(envFile), filepath.Base(envFile)+".sig")
}
