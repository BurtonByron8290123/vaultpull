// Package pagesize provides a configurable pagination policy for iterating
// over large sets of Vault secret keys.
//
// Usage:
//
//	p := pagesize.FromEnv()
//	if err := p.Validate(); err != nil {
//		log.Fatal(err)
//	}
//	for i := 0; i < p.Pages(len(keys)); i++ {
//		chunk := p.Slice(keys, i)
//		// fetch secrets for chunk...
//	}
//
// The page size can be tuned via the VAULTPULL_PAGE_SIZE environment variable.
// Valid values are 1–1000; the default is 100.
package pagesize
