package envseal

import (
	"encoding/hex"
	"fmt"
	"os"
)

const (
	envEncKey  = "VAULTPULL_SEAL_ENC_KEY"
	envSignKey = "VAULTPULL_SEAL_SIGN_KEY"
)

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_SEAL_ENC_KEY  — hex-encoded AES key (32, 48, or 64 hex chars → 16/24/32 bytes)
//	VAULTPULL_SEAL_SIGN_KEY — hex-encoded HMAC key (≥64 hex chars → ≥32 bytes)
func FromEnv() (Policy, error) {
	encHex := os.Getenv(envEncKey)
	if encHex == "" {
		return Policy{}, fmt.Errorf("envseal: %s is not set", envEncKey)
	}
	encKey, err := hex.DecodeString(encHex)
	if err != nil {
		return Policy{}, fmt.Errorf("envseal: invalid %s: %w", envEncKey, err)
	}

	signHex := os.Getenv(envSignKey)
	if signHex == "" {
		return Policy{}, fmt.Errorf("envseal: %s is not set", envSignKey)
	}
	signKey, err := hex.DecodeString(signHex)
	if err != nil {
		return Policy{}, fmt.Errorf("envseal: invalid %s: %w", envSignKey, err)
	}

	return Policy{EncryptionKey: encKey, SigningKey: signKey}, nil
}
