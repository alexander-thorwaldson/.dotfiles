package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zoobzio/dotfiles/kuang/models"
)

// PnpmRunParams are the parameters for running a pnpm script.
type PnpmRunParams struct {
	Dir    string `json:"dir" jsonschema:"Absolute path to the repo working directory"`
	Script string `json:"script" jsonschema:"Package.json script name to run (e.g. test, lint, build, dev)"`
}

// packageJSON is a minimal representation of a package.json file.
type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// validScript checks that the script exists in the package.json at dir.
func validScript(dir, script string) error {
	data, err := os.ReadFile(filepath.Join(dir, "package.json")) // #nosec G304 -- dir is a controlled repo path //nolint:gosec
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

// PnpmRun executes a pnpm script by name and returns its output and exit code.
// The script must exist in the package.json at the given directory.
func PnpmRun(ctx context.Context, _ *mcp.CallToolRequest, p PnpmRunParams) (*mcp.CallToolResult, models.CmdRunResult, error) {
	if p.Script == "" {
		return ErrResult[models.CmdRunResult](fmt.Errorf("script name is required"))
	}

	if err := validScript(p.Dir, p.Script); err != nil {
		return ErrResult[models.CmdRunResult](err)
	}

	cmd := exec.CommandContext(ctx, "pnpm", "run", p.Script) // #nosec G204 -- validated against package.json //nolint:gosec
	cmd.Dir = p.Dir
	out, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok { //nolint:errorlint // need concrete type for ExitCode
			exitCode = exitErr.ExitCode()
		} else {
			return ErrResult[models.CmdRunResult](fmt.Errorf("pnpm: %w", err))
		}
	}

	result := models.CmdRunResult{
		Output:   string(out),
		ExitCode: exitCode,
	}
	return Result(string(out), result)
}
