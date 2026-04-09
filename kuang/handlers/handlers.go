package handlers

import (
	"github.com/zoobz-io/rocco"
	"github.com/zoobzio/dotfiles/kuang/wire"
)

var discover = rocco.GET[rocco.NoBody, wire.ToolListResponse]("/v1/tools", func(r *rocco.Request[rocco.NoBody]) (wire.ToolListResponse, error) {
	// Return all registered endpoints as tool descriptions.
	// Role-based filtering will be applied by middleware.
	tools := []wire.ToolInfo{
		{Name: "gh_pr_list", Description: "List pull requests for a GitHub repo"},
		{Name: "gh_pr_view", Description: "View a specific pull request"},
		{Name: "gh_pr_create", Description: "Create a pull request on GitHub"},
		{Name: "gh_issue_list", Description: "List issues for a GitHub repo"},
		{Name: "gh_issue_view", Description: "View a specific issue"},
		{Name: "gh_repo_view", Description: "View repository details"},
		{Name: "gh_run_list", Description: "List recent workflow runs"},
		{Name: "gh_run_view", Description: "View a workflow run with jobs and steps"},
		{Name: "gh_run_log", Description: "View failed step logs from a workflow run"},
		{Name: "gh_pr_checks", Description: "View CI check status for a PR"},
		{Name: "jira_issue_list", Description: "List Jira issues"},
		{Name: "jira_issue_view", Description: "View a Jira issue"},
		{Name: "dev_start", Description: "Start a dev server"},
		{Name: "dev_stop", Description: "Stop a dev server"},
		{Name: "dev_status", Description: "List running dev servers"},
		{Name: "dev_log", Description: "Read dev server logs"},
		{Name: "pnpm_run", Description: "Run a pnpm script"},
	}
	return wire.ToolListResponse{Tools: tools}, nil
}).
	WithSummary("List available tools").
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
