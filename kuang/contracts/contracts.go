// Package contracts defines the service interfaces consumed by API handlers.
package contracts

import (
	"context"

	"github.com/zoobzio/dotfiles/kuang/wire"
)

// GH provides GitHub operations.
type GH interface {
	PRList(ctx context.Context, req wire.PRListRequest) (string, error)
	PRView(ctx context.Context, req wire.PRViewRequest) (string, error)
	PRCreate(ctx context.Context, req wire.PRCreateRequest) (string, error)
	IssueList(ctx context.Context, req wire.IssueListRequest) (string, error)
	IssueView(ctx context.Context, req wire.IssueViewRequest) (string, error)
	RepoView(ctx context.Context, req wire.RepoViewRequest) (string, error)
	RunList(ctx context.Context, req wire.RunListRequest) (string, error)
	RunView(ctx context.Context, req wire.RunViewRequest) (string, error)
	RunLog(ctx context.Context, req wire.RunLogRequest) (string, error)
	PRChecks(ctx context.Context, req wire.PRChecksRequest) (string, error)
}

// Jira provides Jira operations.
type Jira interface {
	IssueList(ctx context.Context, req wire.JiraIssueListRequest) (string, error)
	IssueView(ctx context.Context, req wire.JiraIssueViewRequest) (string, error)
}

// Dev provides dev server lifecycle operations.
type Dev interface {
	Start(ctx context.Context, req wire.DevStartRequest) (*wire.DevStartResponse, error)
	Stop(ctx context.Context, req wire.DevStopRequest) (*wire.DevStopResponse, error)
	Status(ctx context.Context) (*wire.DevStatusResponse, error)
	Log(ctx context.Context, req wire.DevLogRequest) (*wire.DevLogResponse, error)
}

// Pnpm provides package script execution.
type Pnpm interface {
	Run(ctx context.Context, req wire.PnpmRunRequest) (*wire.PnpmRunResponse, error)
}
