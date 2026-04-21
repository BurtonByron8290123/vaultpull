// Package envrotate detects secrets that are approaching or have exceeded
// their maximum allowed age, enabling rotation enforcement workflows.
//
// Usage:
//
//	p := envrotate.FromEnv()
//	r, err := envrotate.New(p)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	results, err := r.Check(map[string]time.Time{
//		"DB_PASSWORD": lastRotatedAt,
//	})
//	if errors.Is(err, envrotate.ErrRotationRequired) {
//		// handle forced rotation
//	}
//
// Environment variables:
//
//	VAULTPULL_ROTATE_MAX_AGE_DAYS   max secret age in days (default 90)
//	VAULTPULL_ROTATE_WARN_AGE_DAYS  warn threshold in days (default 75)
//	VAULTPULL_ROTATE_DRY_RUN        report only, no error (default false)
package envrotate
