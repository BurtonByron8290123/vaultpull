package envsign

import (
	"encoding/hex"
	"fmt"
	"os"
)

const (
	envKeyHex = "VAULTPULL_SIGN_KEY"
)

// FromEnv constructs a Signer from the VAULTPULL_SIGN_KEY environment variable.
// The variable must be a hex-encoded key of at least 16 bytes (32 hex chars).
func FromEnv() (*Signer, error) {
	raw := os.Getenv(envKeyHex)
	if raw == "" {
		return nil, fmt.Errorf("envsign: %s is not set", envKeyHex)
	}
	key, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("envsign: decode key: %w", err)
	}
	return New(key)
}
