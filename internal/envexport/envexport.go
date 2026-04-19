// Package envexport writes secret maps to various output formats.
package envexport

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format identifies the output serialisation format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatExport Format = "export"
)

// Exporter serialises a secret map to an io.Writer.
type Exporter struct {
	format Format
}

// New returns an Exporter for the given format.
// An error is returned for unrecognised formats.
func New(format Format) (*Exporter, error) {
	switch format {
	case FormatDotenv, FormatJSON, FormatExport:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("envexport: unknown format %q", format)
	}
}

// Write serialises secrets to w in the configured format.
func (e *Exporter) Write(w io.Writer, secrets map[string]string) error {
	switch e.format {
	case FormatJSON:
		return writeJSON(w, secrets)
	case FormatExport:
		return writeDotenv(w, secrets, true)
	default:
		return writeDotenv(w, secrets, false)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func writeDotenv(w io.Writer, secrets map[string]string, withExport bool) error {
	prefix := ""
	if withExport {
		prefix = "export "
	}
	for _, k := range sortedKeys(secrets) {
		v := secrets[k]
		if strings.ContainsAny(v, " \t\n") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(w, "%s%s=%s\n", prefix, k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, secrets map[string]string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	ordered := make(map[string]string, len(secrets))
	for k, v := range secrets {
		ordered[k] = v
	}
	return enc.Encode(ordered)
}
