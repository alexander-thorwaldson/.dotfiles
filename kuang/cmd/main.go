package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zoobzio/dotfiles/kuang/handlers"
	"github.com/zoobzio/dotfiles/kuang/tools"
)

func main() {
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

	addr := os.Getenv("KUANG_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	handler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return server
	}, nil)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srv := &http.Server{Addr: addr, Handler: handler}

	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	logger.Info("kuang listening", "addr", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
