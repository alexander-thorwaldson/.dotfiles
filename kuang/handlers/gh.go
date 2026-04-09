package handlers

import (
	"github.com/zoobz-io/rocco"
	"github.com/zoobz-io/sum"
	"github.com/zoobzio/dotfiles/kuang/contracts"
	"github.com/zoobzio/dotfiles/kuang/wire"
)

var prList = rocco.POST[wire.PRListRequest, wire.CLIOutput]("/v1/gh/prs", func(r *rocco.Request[wire.PRListRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.PRList(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("List pull requests").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_pr_list")

var prView = rocco.POST[wire.PRViewRequest, wire.CLIOutput]("/v1/gh/prs/view", func(r *rocco.Request[wire.PRViewRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.PRView(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("View a pull request").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_pr_view")

var prCreate = rocco.POST[wire.PRCreateRequest, wire.CLIOutput]("/v1/gh/prs/create", func(r *rocco.Request[wire.PRCreateRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.PRCreate(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("Create a pull request").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_pr_create")

var issueList = rocco.POST[wire.IssueListRequest, wire.CLIOutput]("/v1/gh/issues", func(r *rocco.Request[wire.IssueListRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.IssueList(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("List issues").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_issue_list")

var issueView = rocco.POST[wire.IssueViewRequest, wire.CLIOutput]("/v1/gh/issues/view", func(r *rocco.Request[wire.IssueViewRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.IssueView(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("View an issue").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_issue_view")

var repoView = rocco.POST[wire.RepoViewRequest, wire.CLIOutput]("/v1/gh/repos/view", func(r *rocco.Request[wire.RepoViewRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.RepoView(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("View repository details").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_repo_view")

var runList = rocco.POST[wire.RunListRequest, wire.CLIOutput]("/v1/gh/runs", func(r *rocco.Request[wire.RunListRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.RunList(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("List workflow runs").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_run_list")

var runView = rocco.POST[wire.RunViewRequest, wire.CLIOutput]("/v1/gh/runs/view", func(r *rocco.Request[wire.RunViewRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.RunView(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("View a workflow run").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_run_view")

var runLog = rocco.POST[wire.RunLogRequest, wire.CLIOutput]("/v1/gh/runs/log", func(r *rocco.Request[wire.RunLogRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.RunLog(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("View failed step logs").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_run_log")

var prChecks = rocco.POST[wire.PRChecksRequest, wire.CLIOutput]("/v1/gh/prs/checks", func(r *rocco.Request[wire.PRChecksRequest]) (wire.CLIOutput, error) {
	gh := sum.MustUse[contracts.GH](r)
	out, err := gh.PRChecks(r, r.Body)
	if err != nil {
		return wire.CLIOutput{}, err
	}
	return wire.CLIOutput{Output: out}, nil
}).
	WithSummary("View CI checks for a PR").
	WithTags("gh").
	WithAuthentication().
	WithScopes("gh_pr_checks")
