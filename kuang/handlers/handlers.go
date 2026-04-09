package handlers

import (
	"github.com/zoobz-io/rocco"
	"github.com/zoobzio/dotfiles/kuang/wire"
)

// toolDef pairs a tool name/description with its required scope.
type toolDef struct {
	Name        string
	Description string
	Scope       string
}

// allTools is the canonical list of tools and their required scopes.
var allTools = []toolDef{
	{"gh_pr_list", "List pull requests for a GitHub repo", "gh_pr_list"},
	{"gh_pr_view", "View a specific pull request", "gh_pr_view"},
	{"gh_pr_create", "Create a pull request on GitHub", "gh_pr_create"},
	{"gh_issue_list", "List issues for a GitHub repo", "gh_issue_list"},
	{"gh_issue_view", "View a specific issue", "gh_issue_view"},
	{"gh_repo_view", "View repository details", "gh_repo_view"},
	{"gh_run_list", "List recent workflow runs", "gh_run_list"},
	{"gh_run_view", "View a workflow run with jobs and steps", "gh_run_view"},
	{"gh_run_log", "View failed step logs from a workflow run", "gh_run_log"},
	{"gh_pr_checks", "View CI check status for a PR", "gh_pr_checks"},
	{"jira_issue_list", "List Jira issues", "jira_issue_list"},
	{"jira_issue_view", "View a Jira issue", "jira_issue_view"},
	{"dev_start", "Start a dev server", "dev_start"},
	{"dev_stop", "Stop a dev server", "dev_stop"},
	{"dev_status", "List running dev servers", "dev_status"},
	{"dev_log", "Read dev server logs", "dev_log"},
	{"pnpm_run", "Run a pnpm script", "pnpm_run"},
}

var discover = rocco.GET[rocco.NoBody, wire.ToolListResponse]("/v1/tools", func(r *rocco.Request[rocco.NoBody]) (wire.ToolListResponse, error) {
	var tools []wire.ToolInfo
	for _, td := range allTools {
		if r.Identity != nil && r.Identity.HasScope(td.Scope) {
			tools = append(tools, wire.ToolInfo{Name: td.Name, Description: td.Description})
		}
	}
	return wire.ToolListResponse{Tools: tools}, nil
}).
	WithSummary("List available tools filtered by caller's role").
	WithTags("tools").
	WithAuthentication()

// All returns every API endpoint.
func All() []rocco.Endpoint {
	return []rocco.Endpoint{
		// Discovery.
		discover,
		// GitHub.
		prList, prView, prCreate,
		issueList, issueView,
		repoView,
		runList, runView, runLog,
		prChecks,
		// Jira.
		jiraIssueList, jiraIssueView,
		// Dev servers.
		devStart, devStop, devStatus, devLog,
		// Package manager.
		pnpmRun,
	}
}
