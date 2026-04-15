// Package hooks provides pre/post pull lifecycle hook execution.
package hooks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Config holds the hook configuration.
type Config struct {
	// PrePull is a shell command to run before pulling secrets.
	PrePull string `mapstructure:"pre_pull"`
	// PostPull is a shell command to run after pulling secrets.
	PostPull string `mapstructure:"post_pull"`
	// Timeout is the maximum duration for a hook to run.
	Timeout time.Duration `mapstructure:"timeout"`
}

// Runner executes lifecycle hooks.
type Runner struct {
	cfg Config
	env []string
}

// New creates a new Runner with the given config.
// If cfg.Timeout is zero, a default of 30 seconds is used.
func New(cfg Config) *Runner {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &Runner{cfg: cfg, env: os.Environ()}
}

// RunPrePull executes the pre-pull hook if configured.
func (r *Runner) RunPrePull(ctx context.Context) error {
	if r.cfg.PrePull == "" {
		return nil
	}
	return r.run(ctx, "pre_pull", r.cfg.PrePull)
}

// RunPostPull executes the post-pull hook if configured.
func (r *Runner) RunPostPull(ctx context.Context) error {
	if r.cfg.PostPull == "" {
		return nil
	}
	return r.run(ctx, "post_pull", r.cfg.PostPull)
}

func (r *Runner) run(ctx context.Context, name, command string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("hooks: %s command is empty", name)
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Env = r.env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("hooks: %s timed out after %s", name, r.cfg.Timeout)
		}
		return fmt.Errorf("hooks: %s failed: %w", name, err)
	}
	return nil
}
