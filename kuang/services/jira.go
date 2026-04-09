package services

import (
	"context"
	"fmt"

	"github.com/zoobzio/dotfiles/kuang/wire"
)

// jiraListTemplate emits JSON from jira list.
const jiraListTemplate = `[{{range $i, $e := .issues}}{{if $i}},{{end}}{"key":"{{.key}}","summary":"{{.fields.summary}}","status":"{{.fields.status.name}}","assignee":"{{if .fields.assignee}}{{.fields.assignee.displayName}}{{end}}","priority":"{{if .fields.priority}}{{.fields.priority.name}}{{end}}","type":"{{.fields.issuetype.name}}"}{{end}}]`

// jiraViewTemplate emits JSON from jira view.
const jiraViewTemplate = `{"key":"{{.key}}","summary":"{{.fields.summary}}","description":{{toJson .fields.description}},"status":"{{.fields.status.name}}","assignee":"{{if .fields.assignee}}{{.fields.assignee.displayName}}{{end}}","reporter":"{{if .fields.reporter}}{{.fields.reporter.displayName}}{{end}}","priority":"{{if .fields.priority}}{{.fields.priority.name}}{{end}}","type":"{{.fields.issuetype.name}}","labels":{{toJson .fields.labels}},"comments":[{{range $i, $c := .fields.comment.comments}}{{if $i}},{{end}}{"author":"{{$c.author.displayName}}","body":{{toJson $c.body}},"created":"{{$c.created}}"}{{end}}]}`

// Jira implements contracts.Jira by shelling out to the jira CLI with ice filtering.
type Jira struct {
	ice *ICEClient
}

// NewJira creates a new Jira service with ice filtering.
func NewJira(ice *ICEClient) *Jira { return &Jira{ice: ice} }

// IssueList lists Jira issues.
func (j *Jira) IssueList(ctx context.Context, req wire.JiraIssueListRequest) (string, error) {
	return FilteredCLI(j.ice, req, func() (string, error) {
		args := []string{"list", "--template", jiraListTemplate}
		if req.Project != "" {
			args = append(args, "-p", req.Project)
		}
		if req.Query != "" {
			args = append(args, "-q", req.Query)
		}
		if req.Assignee != "" {
			args = append(args, "-a", req.Assignee)
		}
		if req.Limit > 0 {
			args = append(args, "-l", fmt.Sprintf("%d", req.Limit))
		}
		return cli(ctx, "jira", args...)
	})
}

// IssueView retrieves details for a Jira issue.
func (j *Jira) IssueView(ctx context.Context, req wire.JiraIssueViewRequest) (string, error) {
	return FilteredCLI(j.ice, req, func() (string, error) {
		return cli(ctx, "jira", "view", req.Issue, "--template", jiraViewTemplate)
	})
}
