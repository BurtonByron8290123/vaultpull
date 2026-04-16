package env

// MergeStrategy controls how existing keys are handled during a merge.
type MergeStrategy int

const (
	// StrategyOverwrite replaces existing keys with incoming values.
	StrategyOverwrite MergeStrategy = iota
	// StrategyPreserve keeps existing keys unchanged.
	StrategyPreserve
)

// Merge combines base and incoming maps according to the given strategy.
// Keys present only in incoming are always added.
// Keys present only in base are always preserved.
func Merge(base, incoming map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range incoming {
		if _, exists := result[k]; !exists || strategy == StrategyOverwrite {
			result[k] = v
		}
	}
	return result
}

// Subtract returns a copy of base without keys present in remove.
func Subtract(base map[string]string, remove map[string]string) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		if _, found := remove[k]; !found {
			result[k] = v
		}
	}
	return result
}
