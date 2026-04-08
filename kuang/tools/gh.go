// Package tools implements MCP tool handlers that shell out to CLI tools.
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

// PRListParams are the parameters for listing pull requests.
type PRListParams struct {
	Repo  string `json:"repo" jsonschema:"owner/repo to list PRs for"`
	State string `json:"state,omitempty" jsonschema:"Filter by state: open, closed, merged, all (default: open)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Max number of PRs to return (default: 30)"`
}

// PRViewParams are the parameters for viewing a pull request.
type PRViewParams struct {
	Repo   string `json:"repo" jsonschema:"owner/repo"`
	Number int    `json:"number" jsonschema:"PR number"`
}

// IssueListParams are the parameters for listing issues.
type IssueListParams struct {
	Repo  string `json:"repo" jsonschema:"owner/repo to list issues for"`
	State string `json:"state,omitempty" jsonschema:"Filter by state: open, closed, all (default: open)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Max number of issues to return (default: 30)"`
}

// IssueViewParams are the parameters for viewing an issue.
type IssueViewParams struct {
	Repo   string `json:"repo" jsonschema:"owner/repo"`
	Number int    `json:"number" jsonschema:"Issue number"`
}

// RepoViewParams are the parameters for viewing a repository.
type RepoViewParams struct {
	Repo string `json:"repo" jsonschema:"owner/repo to view"`
}

// ghExec runs the GitHub CLI and returns its output.
// Swappable for testing.
var ghExec = func(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "gh", args...) // #nosec G204 -- intentional subprocess //nolint:gosec
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
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

// PRList lists pull requests for a GitHub repository.
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

// PRView retrieves full details for a specific pull request.
func PRView(ctx context.Context, _ *mcp.CallToolRequest, p PRViewParams) (*mcp.CallToolResult, models.PRDetail, error) {
	var out models.PRDetail
	raw, err := ghJSON(ctx, &out, "pr", "view", "-R", p.Repo, fmt.Sprintf("%d", p.Number), "--json", "number,title,state,body,author,url,reviews,comments")
	if err != nil {
		return ErrResult[models.PRDetail](err)
	}
	return Result(raw, out)
}

// IssueList lists issues for a GitHub repository.
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

// IssueView retrieves full details for a specific issue.
func IssueView(ctx context.Context, _ *mcp.CallToolRequest, p IssueViewParams) (*mcp.CallToolResult, models.IssueDetail, error) {
	var out models.IssueDetail
	raw, err := ghJSON(ctx, &out, "issue", "view", "-R", p.Repo, fmt.Sprintf("%d", p.Number), "--json", "number,title,state,body,author,url,comments")
	if err != nil {
		return ErrResult[models.IssueDetail](err)
	}
	return Result(raw, out)
}

// RepoView retrieves full details for a repository.
func RepoView(ctx context.Context, _ *mcp.CallToolRequest, p RepoViewParams) (*mcp.CallToolResult, models.RepoDetail, error) {
	var out models.RepoDetail
	raw, err := ghJSON(ctx, &out, "repo", "view", p.Repo, "--json", "name,description,url,defaultBranchRef,languages,issues,pullRequests")
	if err != nil {
		return ErrResult[models.RepoDetail](err)
	}
	return Result(raw, out)
}

// -- PR Write Operations --

// PRCreateParams are the parameters for creating a pull request.
type PRCreateParams struct {
	Repo     string   `json:"repo" jsonschema:"owner/repo to create the PR in"`
	Title    string   `json:"title" jsonschema:"Title for the pull request"`
	Body     string   `json:"body,omitempty" jsonschema:"Body/description for the pull request"`
	Base     string   `json:"base,omitempty" jsonschema:"Base branch to merge into (default: repo default branch)"`
	Head     string   `json:"head" jsonschema:"Branch that contains commits for the PR"`
	Labels   []string `json:"labels,omitempty" jsonschema:"Labels to add by name"`
	Reviewer []string `json:"reviewers,omitempty" jsonschema:"Request reviews from people or teams by handle"`
	Draft    bool     `json:"draft,omitempty" jsonschema:"Mark pull request as a draft"`
}

// PRCreate creates a pull request on GitHub.
func PRCreate(ctx context.Context, _ *mcp.CallToolRequest, p PRCreateParams) (*mcp.CallToolResult, models.PRCreateResult, error) {
	args := []string{"pr", "create", "-R", p.Repo, "--title", p.Title, "--head", p.Head}
	if p.Body != "" {
		args = append(args, "--body", p.Body)
	} else {
		args = append(args, "--body", "")
	}
	if p.Base != "" {
		args = append(args, "--base", p.Base)
	}
	if p.Draft {
		args = append(args, "--draft")
	}
	for _, l := range p.Labels {
		args = append(args, "--label", l)
	}
	for _, r := range p.Reviewer {
		args = append(args, "--reviewer", r)
	}
	raw, err := ghExec(ctx, args...)
	if err != nil {
		return ErrResult[models.PRCreateResult](err)
	}
	url := strings.TrimSpace(raw)
	return Result(url, models.PRCreateResult{URL: url})
}

// -- GitHub Actions --

// RunListParams are the parameters for listing workflow runs.
type RunListParams struct {
	Repo     string `json:"repo" jsonschema:"owner/repo to list runs for"`
	Branch   string `json:"branch,omitempty" jsonschema:"Filter by branch name"`
	Status   string `json:"status,omitempty" jsonschema:"Filter by status: queued, completed, in_progress, failure, success, etc."`
	Workflow string `json:"workflow,omitempty" jsonschema:"Filter by workflow name"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Max number of runs to return (default: 20)"`
}

// RunViewParams are the parameters for viewing a workflow run.
type RunViewParams struct {
	Repo  string `json:"repo" jsonschema:"owner/repo"`
	RunID int    `json:"run_id" jsonschema:"Workflow run ID"`
}

// RunLogParams are the parameters for viewing failed logs from a workflow run.
type RunLogParams struct {
	Repo  string `json:"repo" jsonschema:"owner/repo"`
	RunID int    `json:"run_id" jsonschema:"Workflow run ID"`
}

// PRChecksParams are the parameters for viewing checks on a pull request.
type PRChecksParams struct {
	Repo   string `json:"repo" jsonschema:"owner/repo"`
	Number int    `json:"number" jsonschema:"PR number"`
}

// RunList lists recent workflow runs for a repository.
func RunList(ctx context.Context, _ *mcp.CallToolRequest, p RunListParams) (*mcp.CallToolResult, models.RunListResult, error) {
	args := []string{"run", "list", "-R", p.Repo, "--json", "number,databaseId,displayTitle,headBranch,event,status,conclusion,workflowName,url,createdAt"}
	if p.Branch != "" {
		args = append(args, "-b", p.Branch)
	}
	if p.Status != "" {
		args = append(args, "-s", p.Status)
	}
	if p.Workflow != "" {
		args = append(args, "-w", p.Workflow)
	}
	if p.Limit > 0 {
		args = append(args, "-L", fmt.Sprintf("%d", p.Limit))
	}
	var items []models.RunSummary
	raw, err := ghJSON(ctx, &items, args...)
	if err != nil {
		return ErrResult[models.RunListResult](err)
	}
	return Result(raw, models.RunListResult{Items: items})
}

// RunView retrieves full details for a specific workflow run including jobs and steps.
func RunView(ctx context.Context, _ *mcp.CallToolRequest, p RunViewParams) (*mcp.CallToolResult, models.RunDetail, error) {
	var out models.RunDetail
	raw, err := ghJSON(ctx, &out, "run", "view", "-R", p.Repo, fmt.Sprintf("%d", p.RunID), "--json", "number,attempt,displayTitle,headBranch,headSha,event,status,conclusion,workflowName,url,createdAt,jobs")
	if err != nil {
		return ErrResult[models.RunDetail](err)
	}
	return Result(raw, out)
}

// RunLog retrieves failed step logs from a workflow run.
func RunLog(ctx context.Context, _ *mcp.CallToolRequest, p RunLogParams) (*mcp.CallToolResult, models.RunDetail, error) {
	raw, err := ghExec(ctx, "run", "view", "-R", p.Repo, fmt.Sprintf("%d", p.RunID), "--log-failed")
	if err != nil {
		return ErrResult[models.RunDetail](err)
	}
	// Log output is plain text, not JSON — return as text with zero-value struct.
	return Result[models.RunDetail](raw, models.RunDetail{})
}

// PRChecks retrieves CI check status for a pull request.
func PRChecks(ctx context.Context, _ *mcp.CallToolRequest, p PRChecksParams) (*mcp.CallToolResult, models.PRChecksResult, error) {
	var items []models.PRCheck
	raw, err := ghJSON(ctx, &items, "pr", "checks", "-R", p.Repo, fmt.Sprintf("%d", p.Number), "--json", "name,state,bucket,description,workflow,link,event,startedAt,completedAt")
	if err != nil {
		return ErrResult[models.PRChecksResult](err)
	}
	return Result(raw, models.PRChecksResult{Items: items})
}
