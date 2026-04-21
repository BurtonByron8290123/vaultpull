package envdecrypt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

const (
	envKeyHex = "VAULTPULL_DECRYPT_KEY"
)

// FromEnv constructs a Decrypter from the hex-encoded key stored in the
// VAULTPULL_DECRYPT_KEY environment variable. If the variable is empty a nil
// Decrypter and nil error are returned, signalling that decryption is
// disabled.
func FromEnv() (*Decrypter, error) {
	hex_val := os.Getenv(envKeyHex)
	if hex_val == "" {
		return nil, nil
	}
	key, err := hex.DecodeString(hex_val)
	if err != nil {
		return nil, fmt.Errorf("envdecrypt: %s is not valid hex: %w", envKeyHex, err)
	}
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("envdecrypt: key must be 16, 24, or 32 bytes (32, 48, or 64 hex chars)")
	}
	return New(key)
}
