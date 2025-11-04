package providers

import (
	"bytes"
	"careco/backend/o11y"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func ProvideServiceVersionFromGitRevision(ctx context.Context) (o11y.ServiceVersion, error) {
	c := exec.CommandContext(ctx, "git", "describe", "--always", "--tags", "--dirty", "--abbrev=0")
	c.WaitDelay = time.Second * 1
	c.Cancel = func() error { return c.Process.Signal(os.Interrupt) }
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		return "", fmt.Errorf("command failed: stderr=%s: %w", stderr, err)
	}
	return o11y.ServiceVersion(bytes.TrimSpace(stdout.Bytes())), nil
}
