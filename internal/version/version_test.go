package version_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/version"
)

func TestGetReturnsDefaults(t *testing.T) {
	info := version.Get()
	if info.Version == "" {
		t.Fatal("expected non-empty Version")
	}
	if info.GoVersion == "" {
		t.Fatal("expected non-empty GoVersion")
	}
	if info.OS == "" {
		t.Fatal("expected non-empty OS")
	}
	if info.Arch == "" {
		t.Fatal("expected non-empty Arch")
	}
}

func TestGetPopulatesGoVersion(t *testing.T) {
	info := version.Get()
	if !strings.HasPrefix(info.GoVersion, "go") {
		t.Fatalf("expected GoVersion to start with 'go', got %q", info.GoVersion)
	}
}

func TestPrintContainsVersion(t *testing.T) {
	var buf bytes.Buffer
	info := version.Info{
		Version:   "1.2.3",
		Commit:    "abc1234",
		BuildDate: "2024-01-01",
		GoVersion: "go1.22.0",
		OS:        "linux",
		Arch:      "amd64",
	}
	version.Print(&buf, info)
	out := buf.String()

	for _, want := range []string{"1.2.3", "abc1234", "2024-01-01", "go1.22.0", "linux", "amd64"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPrintWritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	version.Print(&buf, version.Get())
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output from Print")
	}
}
