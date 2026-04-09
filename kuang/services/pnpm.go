package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zoobzio/dotfiles/kuang/wire"
)

// Pnpm implements contracts.Pnpm with package.json validation and ice filtering.
type Pnpm struct {
	ice *ICEClient
}

// NewPnpm creates a new Pnpm service with ice filtering.
func NewPnpm(ice *ICEClient) *Pnpm { return &Pnpm{ice: ice} }

// Run executes a pnpm script after validating it exists in package.json.
func (p *Pnpm) Run(ctx context.Context, req wire.PnpmRunRequest) (*wire.PnpmRunResponse, error) {
	if req.Script == "" {
		return nil, fmt.Errorf("script name is required")
	}
	if err := validScript(req.Dir, req.Script); err != nil {
		return nil, err
	}

	// Ice scans the input (script name + dir).
	injected, err := p.ice.IsInjection(req.Script)
	if err != nil {
		return nil, fmt.Errorf("ice unavailable: %w", err)
	}
	if injected {
		return nil, fmt.Errorf("request blocked: prompt injection detected in input")
	}

	cmd := exec.CommandContext(ctx, "pnpm", "run", req.Script) // #nosec G204 -- validated against package.json //nolint:gosec
	cmd.Dir = req.Dir
	out, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok { //nolint:errorlint // need concrete type for ExitCode
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("pnpm: %w", err)
		}
	}

	// Ice scans the output.
	output := string(out)
	injected, err = p.ice.IsInjection(output)
	if err != nil {
		return nil, fmt.Errorf("ice unavailable: %w", err)
	}
	if injected {
		return nil, fmt.Errorf("response blocked: prompt injection detected in output")
	}

	return &wire.PnpmRunResponse{Output: output, ExitCode: exitCode}, nil
}

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

func validScript(dir, script string) error {
	data, err := os.ReadFile(filepath.Join(dir, "package.json")) // #nosec G304 -- controlled path //nolint:gosec
	if err != nil {
		return fmt.Errorf("reading package.json: %w", err)
	}
	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("parsing package.json: %w", err)
	}
	if _, ok := pkg.Scripts[script]; !ok {
		available := make([]string, 0, len(pkg.Scripts))
		for k := range pkg.Scripts {
			available = append(available, k)
		}
		return fmt.Errorf("script %q not found in package.json (available: %v)", script, available)
	}
	return nil
}
