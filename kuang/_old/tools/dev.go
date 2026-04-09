package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zoobzio/dotfiles/kuang/handlers"
	"github.com/zoobzio/dotfiles/kuang/models"
)

// DefaultDevTimeout is how long to wait for a dev server port to become ready.
const DefaultDevTimeout = 30 * time.Second

// DevStartParams are the parameters for starting a dev server.
type DevStartParams struct {
	Repo    string `json:"repo" jsonschema:"Repository name used as the service key"`
	Dir     string `json:"dir" jsonschema:"Absolute path to the repo working directory"`
	Command string `json:"command,omitempty" jsonschema:"Command to run (default: pnpm dev)"`
	Port    int    `json:"port" jsonschema:"Port the dev server listens on"`
	Timeout int    `json:"timeout,omitempty" jsonschema:"Seconds to wait for port readiness (default: 30)"`
}

// DevStopParams are the parameters for stopping a dev server.
type DevStopParams struct {
	Repo string `json:"repo" jsonschema:"Repository name to stop"`
}

// DevStatusParams are the parameters for listing dev server status.
type DevStatusParams struct {
	Repo string `json:"repo,omitempty" jsonschema:"Filter by repository name (omit for all)"`
}

// DevLogParams are the parameters for reading dev server logs.
type DevLogParams struct {
	Repo  string `json:"repo" jsonschema:"Repository name to read logs from"`
	Lines int    `json:"lines,omitempty" jsonschema:"Number of recent lines to return (default: 50)"`
}

// NewDevStart creates a DevStart handler bound to the given registry.
func NewDevStart(reg *handlers.Registry) mcp.ToolHandlerFor[DevStartParams, models.DevStartResult] {
	return func(_ context.Context, _ *mcp.CallToolRequest, p DevStartParams) (*mcp.CallToolResult, models.DevStartResult, error) {
		cmd := p.Command
		if cmd == "" {
			cmd = "pnpm dev"
		}
		timeout := DefaultDevTimeout
		if p.Timeout > 0 {
			timeout = time.Duration(p.Timeout) * time.Second
		}

		entry, err := reg.Start(p.Repo, p.Dir, cmd, p.Port, timeout)
		if err != nil {
			return ErrResult[models.DevStartResult](err)
		}

		out := models.DevStartResult{
			Repo:   entry.Repo,
			Port:   entry.Port,
			Status: entry.Status,
			URL:    fmt.Sprintf("http://localhost:%d", entry.Port),
		}
		return Result(fmt.Sprintf("dev server for %s is %s at http://localhost:%d", entry.Repo, entry.Status, entry.Port), out)
	}
}

// NewDevStop creates a DevStop handler bound to the given registry.
func NewDevStop(reg *handlers.Registry) mcp.ToolHandlerFor[DevStopParams, models.DevStopResult] {
	return func(_ context.Context, _ *mcp.CallToolRequest, p DevStopParams) (*mcp.CallToolResult, models.DevStopResult, error) {
		_, exists := reg.Get(p.Repo)
		reg.Stop(p.Repo)

		out := models.DevStopResult{
			Repo:    p.Repo,
			Stopped: exists,
		}
		msg := fmt.Sprintf("dev server for %s stopped", p.Repo)
		if !exists {
			msg = fmt.Sprintf("no dev server running for %s", p.Repo)
		}
		return Result(msg, out)
	}
}

// NewDevStatus creates a DevStatus handler bound to the given registry.
func NewDevStatus(reg *handlers.Registry) mcp.ToolHandlerFor[DevStatusParams, models.DevStatusResult] {
	return func(_ context.Context, _ *mcp.CallToolRequest, p DevStatusParams) (*mcp.CallToolResult, models.DevStatusResult, error) {
		var entries []*handlers.ServiceEntry
		if p.Repo != "" {
			if e, ok := reg.Get(p.Repo); ok {
				entries = []*handlers.ServiceEntry{e}
			}
		} else {
			entries = reg.All()
		}

		services := make([]models.DevServiceInfo, 0, len(entries))
		for _, e := range entries {
			services = append(services, models.DevServiceInfo{
				Repo:      e.Repo,
				Dir:       e.Dir,
				Command:   e.Command,
				Port:      e.Port,
				PID:       e.PID,
				Status:    e.Status,
				URL:       fmt.Sprintf("http://localhost:%d", e.Port),
				StartedAt: e.StartedAt.Format(time.RFC3339),
			})
		}

		out := models.DevStatusResult{Services: services}
		lines := make([]string, 0, len(services))
		for _, s := range services {
			lines = append(lines, fmt.Sprintf("%s: %s on :%d (%s)", s.Repo, s.Status, s.Port, s.Command))
		}
		if len(lines) == 0 {
			return Result("no dev servers running", out)
		}
		return Result(strings.Join(lines, "\n"), out)
	}
}

// NewDevLog creates a DevLog handler bound to the given registry.
func NewDevLog(reg *handlers.Registry) mcp.ToolHandlerFor[DevLogParams, models.DevLogResult] {
	return func(_ context.Context, _ *mcp.CallToolRequest, p DevLogParams) (*mcp.CallToolResult, models.DevLogResult, error) {
		e, ok := reg.Get(p.Repo)
		if !ok {
			return ErrResult[models.DevLogResult](fmt.Errorf("no dev server running for %s", p.Repo))
		}

		n := p.Lines
		if n <= 0 {
			n = 50
		}
		lines := e.LogBuf.Lines(n)
		out := models.DevLogResult{
			Repo:  p.Repo,
			Lines: lines,
		}
		return Result(strings.Join(lines, "\n"), out)
	}
}
