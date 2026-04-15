package envwriter

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath string
	overwrite bool
}

// New creates a new Writer for the given file path.
func New(filePath string, overwrite bool) *Writer {
	return &Writer{
		filePath: filePath,
		overwrite: overwrite,
	}
}

// Write writes the provided key-value secrets map to the .env file.
// If overwrite is false, existing keys are preserved.
func (w *Writer) Write(secrets map[string]string) error {
	existing := map[string]string{}

	if !w.overwrite {
		var err error
		existing, err = readEnvFile(w.filePath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("reading existing env file: %w", err)
		}
	}

	for k, v := range secrets {
		existing[k] = v
	}

	return writeEnvFile(w.filePath, existing)
}

// readEnvFile parses a .env file into a key-value map.
func readEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result, nil
}

// writeEnvFile writes a key-value map to a .env file in sorted key order.
func writeEnvFile(path string, data map[string]string) error {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, data[k])
	}

	return os.WriteFile(path, []byte(sb.String()), 0600)
}
