// Package pagesize controls paginated secret fetching from Vault.
package pagesize

import (
	"errors"
	"os"
	"strconv"
)

const (
	DefaultPageSize = 100
	MinPageSize     = 1
	MaxPageSize     = 1000
)

// Policy holds pagination configuration.
type Policy struct {
	PageSize int
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{PageSize: DefaultPageSize}
}

// Validate returns an error if the policy is misconfigured.
func (p Policy) Validate() error {
	if p.PageSize < MinPageSize {
		return errors.New("pagesize: page size must be at least 1")
	}
	if p.PageSize > MaxPageSize {
		return errors.New("pagesize: page size must not exceed 1000")
	}
	return nil
}

// Pages returns the number of pages required to cover total items.
func (p Policy) Pages(total int) int {
	if total <= 0 {
		return 0
	}
	return (total + p.PageSize - 1) / p.PageSize
}

// Slice returns the sub-slice for the given zero-based page index.
func (p Policy) Slice(items []string, page int) []string {
	start := page * p.PageSize
	if start >= len(items) {
		return nil
	}
	end := start + p.PageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

// FromEnv reads VAULTPULL_PAGE_SIZE from the environment.
func FromEnv() Policy {
	p := DefaultPolicy()
	if raw := os.Getenv("VAULTPULL_PAGE_SIZE"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			p.PageSize = v
		}
	}
	return p
}
