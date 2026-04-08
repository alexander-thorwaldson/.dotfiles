package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec" //nolint:gosec // used by ghExec default
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zoobzio/dotfiles/kuang/models"
)

// -- GitHub param types --

type PRListParams struct {
	Repo  string `json:"repo" jsonschema:"owner/repo to list PRs for"`
	State string `json:"state,omitempty" jsonschema:"Filter by state: open, closed, merged, all (default: open)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Max number of PRs to return (default: 30)"`
}

type PRViewParams struct {
	Repo   string `json:"repo" jsonschema:"owner/repo"`
	Number int    `json:"number" jsonschema:"PR number"`
}

type IssueListParams struct {
	Repo  string `json:"repo" jsonschema:"owner/repo to list issues for"`
	State string `json:"state,omitempty" jsonschema:"Filter by state: open, closed, all (default: open)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Max number of issues to return (default: 30)"`
}

type IssueViewParams struct {
	Repo   string `json:"repo" jsonschema:"owner/repo"`
	Number int    `json:"number" jsonschema:"Issue number"`
}

type RepoViewParams struct {
	Repo string `json:"repo" jsonschema:"owner/repo to view"`
}

// ghExec runs the GitHub CLI and returns its output.
// Swappable for testing.
var ghExec = func(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "gh", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

// ghJSON runs the GitHub CLI, parses the JSON output into v, and returns the raw text.
func ghJSON[Out any](ctx context.Context, v *Out, args ...string) (string, error) {
	raw, err := ghExec(ctx, args...)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal([]byte(raw), v); err != nil {
		return raw, fmt.Errorf("parsing gh output: %w", err)
	}
	return raw, nil
}

func PRList(ctx context.Context, _ *mcp.CallToolRequest, p PRListParams) (*mcp.CallToolResult, models.PRListResult, error) {
	args := []string{"pr", "list", "-R", p.Repo, "--json", "number,title,state,author,url"}
	if p.State != "" {
		args = append(args, "-s", p.State)
	}
	if p.Limit > 0 {
		args = append(args, "-L", fmt.Sprintf("%d", p.Limit))
	}
	var items []models.PRSummary
	raw, err := ghJSON(ctx, &items, args...)
	if err != nil {
		return ErrResult[models.PRListResult](err)
	}
	return Result(raw, models.PRListResult{Items: items})
}

func PRView(ctx context.Context, _ *mcp.CallToolRequest, p PRViewParams) (*mcp.CallToolResult, models.PRDetail, error) {
	var out models.PRDetail
	raw, err := ghJSON(ctx, &out, "pr", "view", "-R", p.Repo, fmt.Sprintf("%d", p.Number), "--json", "number,title,state,body,author,url,reviews,comments")
	if err != nil {
		return ErrResult[models.PRDetail](err)
	}
	return Result(raw, out)
}

func IssueList(ctx context.Context, _ *mcp.CallToolRequest, p IssueListParams) (*mcp.CallToolResult, models.IssueListResult, error) {
	args := []string{"issue", "list", "-R", p.Repo, "--json", "number,title,state,author,url"}
	if p.State != "" {
		args = append(args, "-s", p.State)
	}
	if p.Limit > 0 {
		args = append(args, "-L", fmt.Sprintf("%d", p.Limit))
	}
	var items []models.IssueSummary
	raw, err := ghJSON(ctx, &items, args...)
	if err != nil {
		return ErrResult[models.IssueListResult](err)
	}
	return Result(raw, models.IssueListResult{Items: items})
}

func IssueView(ctx context.Context, _ *mcp.CallToolRequest, p IssueViewParams) (*mcp.CallToolResult, models.IssueDetail, error) {
	var out models.IssueDetail
	raw, err := ghJSON(ctx, &out, "issue", "view", "-R", p.Repo, fmt.Sprintf("%d", p.Number), "--json", "number,title,state,body,author,url,comments")
	if err != nil {
		return ErrResult[models.IssueDetail](err)
	}
	return Result(raw, out)
}

func RepoView(ctx context.Context, _ *mcp.CallToolRequest, p RepoViewParams) (*mcp.CallToolResult, models.RepoDetail, error) {
	var out models.RepoDetail
	raw, err := ghJSON(ctx, &out, "repo", "view", p.Repo, "--json", "name,description,url,defaultBranchRef,languages,issues,pullRequests")
	if err != nil {
		return ErrResult[models.RepoDetail](err)
	}
	return Result(raw, out)
}
