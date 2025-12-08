package dev

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"careco/backend/config"
	"careco/backend/infra/firestore"
	"careco/backend/o11y"

	sdk "cloud.google.com/go/firestore"
)

func ProvideFirestoreDatabaseID(_ *config.Environment) (firestore.DatabaseID, error) {
	return sdk.DefaultDatabaseID, nil
}

func ProvideGoogleProjectID(_ *config.Environment) (firestore.ProjectID, error) {
	return "dummy", nil
}

func ProvideDeploymentEnvironmentName(_ *config.Environment) (o11y.DeploymentEnvironmentName, error) {
	return "dev", nil
}

func ProvideServiceVersion(ctx context.Context, _ *config.Environment) (o11y.ServiceVersion, error) {
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
