package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec" //nolint:gosec // used by jiraExec default
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zoobzio/dotfiles/kuang/models"
)

// -- Jira param types --

// JiraIssueListParams are the parameters for listing Jira issues.
type JiraIssueListParams struct {
	Project  string `json:"project,omitempty" jsonschema:"Jira project key (e.g. PROJ)"`
	Query    string `json:"query,omitempty" jsonschema:"JQL query string"`
	Assignee string `json:"assignee,omitempty" jsonschema:"Filter by assignee username"`
	Status   string `json:"status,omitempty" jsonschema:"Filter by status name"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Max number of issues to return"`
}

// JiraIssueViewParams are the parameters for viewing a Jira issue.
type JiraIssueViewParams struct {
	Issue string `json:"issue" jsonschema:"Jira issue key (e.g. PROJ-123)"`
}

// jiraExec runs the Jira CLI and returns its output.
// Swappable for testing.
var jiraExec = func(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "jira", args...) // #nosec G204 -- intentional subprocess //nolint:gosec
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

// jiraListTemplate is a Go template that emits JSON from jira list.
const jiraListTemplate = `[{{range $i, $e := .issues}}{{if $i}},{{end}}{"key":"{{.key}}","summary":"{{.fields.summary}}","status":"{{.fields.status.name}}","assignee":"{{if .fields.assignee}}{{.fields.assignee.displayName}}{{end}}","priority":"{{if .fields.priority}}{{.fields.priority.name}}{{end}}","type":"{{.fields.issuetype.name}}"}{{end}}]`

// jiraViewTemplate is a Go template that emits JSON from jira view.
const jiraViewTemplate = `{"key":"{{.key}}","summary":"{{.fields.summary}}","description":{{toJson .fields.description}},"status":"{{.fields.status.name}}","assignee":"{{if .fields.assignee}}{{.fields.assignee.displayName}}{{end}}","reporter":"{{if .fields.reporter}}{{.fields.reporter.displayName}}{{end}}","priority":"{{if .fields.priority}}{{.fields.priority.name}}{{end}}","type":"{{.fields.issuetype.name}}","labels":{{toJson .fields.labels}},"comments":[{{range $i, $c := .fields.comment.comments}}{{if $i}},{{end}}{"author":"{{$c.author.displayName}}","body":{{toJson $c.body}},"created":"{{$c.created}}"}{{end}}]}`

// JiraIssueList lists Jira issues filtered by project, query, or assignee.
func JiraIssueList(ctx context.Context, _ *mcp.CallToolRequest, p JiraIssueListParams) (*mcp.CallToolResult, models.JiraIssueListResult, error) {
	args := []string{"list", "--template", jiraListTemplate}
	if p.Project != "" {
		args = append(args, "-p", p.Project)
	}
	if p.Query != "" {
		args = append(args, "-q", p.Query)
	}
	if p.Assignee != "" {
		args = append(args, "-a", p.Assignee)
	}
	if p.Limit > 0 {
		args = append(args, "-l", fmt.Sprintf("%d", p.Limit))
	}

	raw, err := jiraExec(ctx, args...)
	if err != nil {
		return ErrResult[models.JiraIssueListResult](err)
	}

	var items []models.JiraIssueSummary
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return ErrResult[models.JiraIssueListResult](fmt.Errorf("parsing jira output: %w", err))
	}

	// Filter by status client-side since jira CLI doesn't have a status flag
	if p.Status != "" {
		filtered := items[:0]
		for _, item := range items {
			if strings.EqualFold(item.Status, p.Status) {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	return Result(raw, models.JiraIssueListResult{Items: items})
}

// JiraIssueView retrieves full details for a specific Jira issue.
func JiraIssueView(ctx context.Context, _ *mcp.CallToolRequest, p JiraIssueViewParams) (*mcp.CallToolResult, models.JiraIssueDetail, error) {
	raw, err := jiraExec(ctx, "view", p.Issue, "--template", jiraViewTemplate)
	if err != nil {
		return ErrResult[models.JiraIssueDetail](err)
	}

	var detail models.JiraIssueDetail
	if err := json.Unmarshal([]byte(raw), &detail); err != nil {
		return ErrResult[models.JiraIssueDetail](fmt.Errorf("parsing jira output: %w", err))
	}

	return Result(raw, detail)
}
