// Package services implements the business logic behind kuang's API contracts.
package services

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// cliExec runs a command and returns its output. Swappable for testing.
var cliExec = func(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...) // #nosec G204 -- intentional subprocess //nolint:gosec
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

// cli is the package-level shorthand for cliExec.
func cli(ctx context.Context, name string, args ...string) (string, error) {
	return cliExec(ctx, name, args...)
}
