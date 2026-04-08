// Package main is the entry point for the kuang MCP server.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zoobzio/dotfiles/kuang/handlers"
	"github.com/zoobzio/dotfiles/kuang/tools"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	iceEndpoint := os.Getenv("KUANG_ICE_ENDPOINT")
	if iceEndpoint == "" {
		iceEndpoint = "http://127.0.0.1:9119"
	}
	ice := handlers.NewICEClient(iceEndpoint)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "kuang",
		Version: "0.1.0",
	}, &mcp.ServerOptions{
		Logger: logger,
	})

	// GitHub tools — all wrapped with ice filtering
	mcp.AddTool(server, &mcp.Tool{Name: "gh_pr_list", Description: "List pull requests for a GitHub repo"}, handlers.FilteredHandler(ice, logger, "gh_pr_list", tools.PRList))
	mcp.AddTool(server, &mcp.Tool{Name: "gh_pr_view", Description: "View a specific pull request"}, handlers.FilteredHandler(ice, logger, "gh_pr_view", tools.PRView))
	mcp.AddTool(server, &mcp.Tool{Name: "gh_issue_list", Description: "List issues for a GitHub repo"}, handlers.FilteredHandler(ice, logger, "gh_issue_list", tools.IssueList))
	mcp.AddTool(server, &mcp.Tool{Name: "gh_issue_view", Description: "View a specific issue"}, handlers.FilteredHandler(ice, logger, "gh_issue_view", tools.IssueView))
	mcp.AddTool(server, &mcp.Tool{Name: "gh_repo_view", Description: "View repository details"}, handlers.FilteredHandler(ice, logger, "gh_repo_view", tools.RepoView))

	// GitHub PR write tools
	mcp.AddTool(server, &mcp.Tool{Name: "gh_pr_create", Description: "Create a pull request on GitHub"}, handlers.FilteredHandler(ice, logger, "gh_pr_create", tools.PRCreate))

	// GitHub Actions tools
	mcp.AddTool(server, &mcp.Tool{Name: "gh_run_list", Description: "List recent workflow runs for a GitHub repo"}, handlers.FilteredHandler(ice, logger, "gh_run_list", tools.RunList))
	mcp.AddTool(server, &mcp.Tool{Name: "gh_run_view", Description: "View a workflow run with jobs and steps"}, handlers.FilteredHandler(ice, logger, "gh_run_view", tools.RunView))
	mcp.AddTool(server, &mcp.Tool{Name: "gh_run_log", Description: "View failed step logs from a workflow run"}, handlers.FilteredHandler(ice, logger, "gh_run_log", tools.RunLog))
	mcp.AddTool(server, &mcp.Tool{Name: "gh_pr_checks", Description: "View CI check status for a pull request"}, handlers.FilteredHandler(ice, logger, "gh_pr_checks", tools.PRChecks))

	// Jira tools
	mcp.AddTool(server, &mcp.Tool{Name: "jira_issue_list", Description: "List Jira issues by project, JQL query, or assignee"}, handlers.FilteredHandler(ice, logger, "jira_issue_list", tools.JiraIssueList))
	mcp.AddTool(server, &mcp.Tool{Name: "jira_issue_view", Description: "View a Jira issue with details and comments"}, handlers.FilteredHandler(ice, logger, "jira_issue_view", tools.JiraIssueView))

	addr := os.Getenv("KUANG_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	handler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return server
	}, nil)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Close(); err != nil {
			logger.Error("server close error", "err", err)
		}
	}()

	logger.Info("kuang listening", "addr", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server error", "err", err)
		return 1
	}
	return 0
}
