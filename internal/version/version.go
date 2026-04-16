package version

import (
	"fmt"
	"io"
	"runtime"
)

// Build information, set via -ldflags at build time.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Info holds structured version information.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get returns the current build Info.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// Print writes a human-readable version string to w.
func Print(w io.Writer, info Info) {
	fmt.Fprintf(w, "vaultpull %s\n", info.Version)
	fmt.Fprintf(w, "  commit:     %s\n", info.Commit)
	fmt.Fprintf(w, "  built:      %s\n", info.BuildDate)
	fmt.Fprintf(w, "  go version: %s\n", info.GoVersion)
	fmt.Fprintf(w, "  os/arch:    %s/%s\n", info.OS, info.Arch)
}
