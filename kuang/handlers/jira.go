package handlers

import (
	"github.com/zoobz-io/rocco"
	"github.com/zoobz-io/sum"
	"github.com/zoobzio/dotfiles/kuang/contracts"
	"github.com/zoobzio/dotfiles/kuang/wire"
)

var jiraIssueList = rocco.POST[wire.JiraIssueListRequest, wire.CLIOutput]("/v1/jira/issues", func(r *rocco.Request[wire.JiraIssueListRequest]) (wire.CLIOutput, error) {
	jira := sum.MustUse[contracts.Jira](r)
	out, err := jira.IssueList(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("List Jira issues").
	WithTags("jira").
	WithAuthentication().
	WithScopes("jira_issue_list")

var jiraIssueView = rocco.POST[wire.JiraIssueViewRequest, wire.CLIOutput]("/v1/jira/issues/view", func(r *rocco.Request[wire.JiraIssueViewRequest]) (wire.CLIOutput, error) {
	jira := sum.MustUse[contracts.Jira](r)
	out, err := jira.IssueView(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("View a Jira issue").
	WithTags("jira").
	WithAuthentication().
	WithScopes("jira_issue_view")
