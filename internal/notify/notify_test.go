package notify

import (
	"bytes"
	"strings"
	"testing"
)

func newBuf(prefix string, quiet bool) (*Notifier, *bytes.Buffer) {
	var buf bytes.Buffer
	n := NewWithWriter(&buf, prefix, quiet)
	return n, &buf
}

func TestInfoWritesMessage(t *testing.T) {
	n, buf := newBuf("test", false)
	n.Info("hello world")
	if !strings.Contains(buf.String(), "hello world") {
		t.Fatalf("expected message in output, got: %q", buf.String())
	}
}

func TestInfoContainsLevelTag(t *testing.T) {
	n, buf := newBuf("", false)
	n.Info("msg")
	if !strings.Contains(buf.String(), "[INFO]") {
		t.Fatalf("expected [INFO] tag, got: %q", buf.String())
	}
}

func TestWarnContainsLevelTag(t *testing.T) {
	n, buf := newBuf("", false)
	n.Warn("something odd")
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Fatalf("expected [WARN] tag, got: %q", buf.String())
	}
}

func TestErrorContainsLevelTag(t *testing.T) {
	n, buf := newBuf("", false)
	n.Error("boom")
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Fatalf("expected [ERROR] tag, got: %q", buf.String())
	}
}

func TestInfoSuppressedInQuietMode(t *testing.T) {
	n, buf := newBuf("vp", true)
	n.Info("should be hidden")
	if buf.Len() != 0 {
		t.Fatalf("expected no output in quiet mode, got: %q", buf.String())
	}
}

func TestWarnNotSuppressedInQuietMode(t *testing.T) {
	n, buf := newBuf("vp", true)
	n.Warn("still shown")
	if buf.Len() == 0 {
		t.Fatal("expected warn to appear even in quiet mode")
	}
}

func TestPrefixAppearedInOutput(t *testing.T) {
	n, buf := newBuf("mypfx", false)
	n.Info("check prefix")
	if !strings.Contains(buf.String(), "[mypfx]") {
		t.Fatalf("expected prefix in output, got: %q", buf.String())
	}
}

func TestInfofFormatsMessage(t *testing.T) {
	n, buf := newBuf("", false)
	n.Infof("synced %d secrets", 7)
	if !strings.Contains(buf.String(), "synced 7 secrets") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}
