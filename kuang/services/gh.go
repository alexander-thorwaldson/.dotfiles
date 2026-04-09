package services

import (
	"context"
	"fmt"

	"github.com/zoobzio/dotfiles/kuang/wire"
)

// GH implements contracts.GH by shelling out to the gh CLI with ice filtering.
type GH struct {
	ice *ICEClient
}

// NewGH creates a new GitHub service with ice filtering.
func NewGH(ice *ICEClient) *GH { return &GH{ice: ice} }

// PRList lists pull requests for a repository.
func (g *GH) PRList(ctx context.Context, req wire.PRListRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		args := []string{"pr", "list", "-R", req.Repo, "--json", "number,title,state,author,url"}
		if req.State != "" {
			args = append(args, "-s", req.State)
		}
		if req.Limit > 0 {
			args = append(args, "-L", fmt.Sprintf("%d", req.Limit))
		}
		return cli(ctx, "gh", args...)
	})
}

// PRView retrieves details for a specific pull request.
func (g *GH) PRView(ctx context.Context, req wire.PRViewRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		return cli(ctx, "gh", "pr", "view", "-R", req.Repo, fmt.Sprintf("%d", req.Number), "--json", "number,title,state,body,author,url,reviews,comments")
	})
}

// PRCreate creates a pull request.
func (g *GH) PRCreate(ctx context.Context, req wire.PRCreateRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		args := []string{"pr", "create", "-R", req.Repo, "--title", req.Title, "--head", req.Head}
		if req.Body != "" {
			args = append(args, "--body", req.Body)
		} else {
			args = append(args, "--body", "")
		}
		if req.Base != "" {
			args = append(args, "--base", req.Base)
		}
		if req.Draft {
			args = append(args, "--draft")
		}
		for _, l := range req.Labels {
			args = append(args, "--label", l)
		}
		for _, r := range req.Reviewer {
			args = append(args, "--reviewer", r)
		}
		return cli(ctx, "gh", args...)
	})
}

// IssueList lists issues for a repository.
func (g *GH) IssueList(ctx context.Context, req wire.IssueListRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		args := []string{"issue", "list", "-R", req.Repo, "--json", "number,title,state,author,url"}
		if req.State != "" {
			args = append(args, "-s", req.State)
		}
		if req.Limit > 0 {
			args = append(args, "-L", fmt.Sprintf("%d", req.Limit))
		}
		return cli(ctx, "gh", args...)
	})
}

// IssueView retrieves details for a specific issue.
func (g *GH) IssueView(ctx context.Context, req wire.IssueViewRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		return cli(ctx, "gh", "issue", "view", "-R", req.Repo, fmt.Sprintf("%d", req.Number), "--json", "number,title,state,body,author,url,comments")
	})
}

// RepoView retrieves details for a repository.
func (g *GH) RepoView(ctx context.Context, req wire.RepoViewRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		return cli(ctx, "gh", "repo", "view", req.Repo, "--json", "name,description,url,defaultBranchRef,languages,issues,pullRequests")
	})
}

// RunList lists recent workflow runs.
func (g *GH) RunList(ctx context.Context, req wire.RunListRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		args := []string{"run", "list", "-R", req.Repo, "--json", "number,databaseId,displayTitle,headBranch,event,status,conclusion,workflowName,url,createdAt"}
		if req.Branch != "" {
			args = append(args, "-b", req.Branch)
		}
		if req.Status != "" {
			args = append(args, "-s", req.Status)
		}
		if req.Workflow != "" {
			args = append(args, "-w", req.Workflow)
		}
		if req.Limit > 0 {
			args = append(args, "-L", fmt.Sprintf("%d", req.Limit))
		}
		return cli(ctx, "gh", args...)
	})
}

// RunView retrieves details for a workflow run.
func (g *GH) RunView(ctx context.Context, req wire.RunViewRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		return cli(ctx, "gh", "run", "view", "-R", req.Repo, fmt.Sprintf("%d", req.RunID), "--json", "number,attempt,displayTitle,headBranch,headSha,event,status,conclusion,workflowName,url,createdAt,jobs")
	})
}

// RunLog retrieves failed step logs from a workflow run.
func (g *GH) RunLog(ctx context.Context, req wire.RunLogRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		return cli(ctx, "gh", "run", "view", "-R", req.Repo, fmt.Sprintf("%d", req.RunID), "--log-failed")
	})
}

// PRChecks retrieves CI check status for a pull request.
func (g *GH) PRChecks(ctx context.Context, req wire.PRChecksRequest) (string, error) {
	return FilteredCLI(g.ice, req, func() (string, error) {
		return cli(ctx, "gh", "pr", "checks", "-R", req.Repo, fmt.Sprintf("%d", req.Number), "--json", "name,state,bucket,description,workflow,link,event,startedAt,completedAt")
	})
}
