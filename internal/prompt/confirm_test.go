package prompt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/vaultpull/internal/prompt"
)

func newConfirmer(input string) (*prompt.Confirmer, *bytes.Buffer) {
	out := &bytes.Buffer{}
	c := prompt.NewWithIO(strings.NewReader(input), out)
	return c, out
}

func TestAskYesReturnsTrue(t *testing.T) {
	for _, ans := range []string{"y\n", "yes\n", "Y\n", "YES\n"} {
		c, _ := newConfirmer(ans)
		ok, err := c.Ask("Continue?", false)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", ans, err)
		}
		if !ok {
			t.Errorf("expected true for answer %q", ans)
		}
	}
}

func TestAskNoReturnsFalse(t *testing.T) {
	for _, ans := range []string{"n\n", "no\n", "N\n", "NO\n"} {
		c, _ := newConfirmer(ans)
		ok, err := c.Ask("Continue?", true)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", ans, err)
		}
		if ok {
			t.Errorf("expected false for answer %q", ans)
		}
	}
}

func TestAskEmptyUsesDefault(t *testing.T) {
	c, _ := newConfirmer("\n")
	ok, err := c.Ask("Continue?", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected default true when input is empty")
	}

	c2, _ := newConfirmer("\n")
	ok2, err := c2.Ask("Continue?", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok2 {
		t.Error("expected default false when input is empty")
	}
}

func TestAskEOFUsesDefault(t *testing.T) {
	c, _ := newConfirmer("") // immediate EOF
	ok, err := c.Ask("Continue?", true)
	if err != nil {
		t.Fatalf("unexpected error on EOF: %v", err)
	}
	if !ok {
		t.Error("expected defaultYes=true on EOF")
	}
}

func TestAskUnrecognisedAnswerReturnsError(t *testing.T) {
	c, _ := newConfirmer("maybe\n")
	_, err := c.Ask("Continue?", false)
	if err == nil {
		t.Error("expected error for unrecognised answer")
	}
}

func TestAskWritesPromptToOutput(t *testing.T) {
	c, out := newConfirmer("y\n")
	_, _ = c.Ask("Deploy now?", false)
	if !strings.Contains(out.String(), "Deploy now?") {
		t.Errorf("prompt not written to output: %q", out.String())
	}
}

func TestAskHintShowsDefaultYes(t *testing.T) {
	c, out := newConfirmer("y\n")
	_, _ = c.Ask("Deploy?", true)
	if !strings.Contains(out.String(), "[Y/n]") {
		t.Errorf("expected [Y/n] hint, got: %q", out.String())
	}
}
