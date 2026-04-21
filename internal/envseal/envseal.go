// Package envseal provides functionality to seal (encrypt + sign) and unseal
// .env files, combining encryption and HMAC signing into a single operation.
package envseal

import (
	"errors"
	"fmt"

	"github.com/vaultpull/internal/envencrypt"
	"github.com/vaultpull/internal/envsign"
)

// ErrSealedFileTampered is returned when the sealed file fails signature verification.
var ErrSealedFileTampered = errors.New("envseal: sealed file has been tampered")

// Policy holds the keys required for sealing and unsealing.
type Policy struct {
	EncryptionKey []byte // 16, 24, or 32 bytes for AES-GCM
	SigningKey    []byte // at least 32 bytes for HMAC-SHA256
}

// Sealer seals and unseals env maps.
type Sealer struct {
	enc *envencrypt.Encrypter
	sig *envsign.Signer
}

// New creates a Sealer from the given Policy.
func New(p Policy) (*Sealer, error) {
	enc, err := envencrypt.New(p.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("envseal: encryption init: %w", err)
	}
	sig, err := envsign.New(p.SigningKey)
	if err != nil {
		return nil, fmt.Errorf("envseal: signing init: %w", err)
	}
	return &Sealer{enc: enc, sig: sig}, nil
}

// Seal encrypts and signs the env map, writing the sealed blob to dst.
func (s *Sealer) Seal(env map[string]string, dst string) error {
	if err := s.enc.EncryptToFile(env, dst); err != nil {
		return fmt.Errorf("envseal: encrypt: %w", err)
	}
	if err := s.sig.SignFile(dst); err != nil {
		return fmt.Errorf("envseal: sign: %w", err)
	}
	return nil
}

// Unseal verifies the signature and decrypts the sealed blob at src.
func (s *Sealer) Unseal(src string) (map[string]string, error) {
	if err := s.sig.VerifyFile(src); err != nil {
		return nil, ErrSealedFileTampered
	}
	env, err := s.enc.DecryptFromFile(src)
	if err != nil {
		return nil, fmt.Errorf("envseal: decrypt: %w", err)
	}
	return env, nil
}
