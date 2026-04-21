package envaudit

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ErrNoLog is returned when the audit log file does not exist.
var ErrNoLog = errors.New("envaudit: log file not found")

// ReadAll reads all entries from the given audit log file.
// Returns ErrNoLog if the file does not exist.
func ReadAll(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoLog
		}
		return nil, fmt.Errorf("envaudit: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("envaudit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envaudit: scan: %w", err)
	}
	return entries, nil
}
