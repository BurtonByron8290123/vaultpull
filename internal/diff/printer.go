package diff

import (
	"fmt"
	"io"
	"sort"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorGray   = "\033[90m"
)

// PrintOptions controls output behaviour of the diff printer.
type PrintOptions struct {
	ShowUnchanged bool
	NoColor       bool
}

// Print writes a formatted diff to w.
func Print(w io.Writer, result *Result, opts PrintOptions) {
	changes := make([]Change, len(result.Changes))
	copy(changes, result.Changes)

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	for _, c := range changes {
		printChange(w, c, opts)
	}

	fmt.Fprintf(w, "\n%s\n", result.Summary())
}

// printChange writes a single Change entry to w using the given options.
func printChange(w io.Writer, c Change, opts PrintOptions) {
	color := changeColor(c.Type)
	prefix := changePrefix(c.Type)

	switch c.Type {
	case Added:
		fmt.Fprintf(w, "%s%s %s = %s%s\n", colorize(color, opts.NoColor), prefix, c.Key, c.NewVal, colorize(colorReset, opts.NoColor))
	case Updated:
		fmt.Fprintf(w, "%s%s %s: %s → %s%s\n", colorize(color, opts.NoColor), prefix, c.Key, c.OldVal, c.NewVal, colorize(colorReset, opts.NoColor))
	case Removed:
		fmt.Fprintf(w, "%s%s %s (was %s)%s\n", colorize(color, opts.NoColor), prefix, c.Key, c.OldVal, colorize(colorReset, opts.NoColor))
	case Unchanged:
		if opts.ShowUnchanged {
			fmt.Fprintf(w, "%s%s %s%s\n", colorize(color, opts.NoColor), prefix, c.Key, colorize(colorReset, opts.NoColor))
		}
	}
}

// changeColor returns the ANSI color code associated with a ChangeType.
func changeColor(t ChangeType) string {
	switch t {
	case Added:
		return colorGreen
	case Updated:
		return colorYellow
	case Removed:
		return colorRed
	default:
		return colorGray
	}
}

// changePrefix returns the single-character prefix symbol for a ChangeType.
func changePrefix(t ChangeType) string {
	switch t {
	case Added:
		return "+"
	case Updated:
		return "~"
	case Removed:
		return "-"
	default:
		return " "
	}
}

func colorize(code string, noColor bool) string {
	if noColor {
		return ""
	}
	return code
}
