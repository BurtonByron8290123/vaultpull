package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair parsed from a .env file.
type Entry struct {
	Key   string
	Value string
	Raw   string // original line as it appeared in the file
}

// Parse reads a .env file from the given path and returns its entries.
// Lines starting with '#' and empty lines are skipped.
// Values may optionally be quoted with single or double quotes.
func Parse(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("env: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		entry, ok := parseLine(line)
		if !ok {
			continue
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scan %s: %w", path, err)
	}
	return entries, nil
}

// ToMap converts a slice of Entry into a key→value map.
// Later entries win on duplicate keys.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// parseLine parses a single line into an Entry.
// Returns (Entry{}, false) for blank lines and comments.
func parseLine(line string) (Entry, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return Entry{}, false
	}
	// Strip optional export prefix
	trimmed = strings.TrimPrefix(trimmed, "export ")

	idx := strings.IndexByte(trimmed, '=')
	if idx < 0 {
		return Entry{}, false
	}
	key := strings.TrimSpace(trimmed[:idx])
	val := strings.TrimSpace(trimmed[idx+1:])
	val = unquote(val)
	return Entry{Key: key, Value: val, Raw: line}, true
}

// unquote strips matching surrounding quotes from a value string.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
