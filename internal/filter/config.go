package filter

import (
	"bufio"
	"os"
	"strings"
)

// PatternConfig holds filter patterns that can be loaded from config or a file.
type PatternConfig struct {
	Patterns []string `mapstructure:"patterns" yaml:"patterns"`
	File     string   `mapstructure:"pattern_file" yaml:"pattern_file"`
}

// Resolve returns a Filter built from the config's patterns and optional
// pattern file. Patterns from the file are appended after inline patterns.
func (pc PatternConfig) Resolve() (*Filter, error) {
	patterns := make([]string, len(pc.Patterns))
	copy(patterns, pc.Patterns)

	if pc.File != "" {
		filePatterns, err := loadPatternsFromFile(pc.File)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, filePatterns...)
	}

	return New(patterns), nil
}

// loadPatternsFromFile reads one pattern per line from a file.
// Lines starting with '#' and blank lines are ignored.
func loadPatternsFromFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, scanner.Err()
}
