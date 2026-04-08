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

	// ICE client.
	iceEndpoint := os.Getenv("KUANG_ICE_ENDPOINT")
	if iceEndpoint == "" {
		iceEndpoint = "http://127.0.0.1:9119"
	}
	ice := handlers.NewICEClient(iceEndpoint)

	// Certs directory.
	certsDir := os.Getenv("KUANG_CERTS_DIR")
	if certsDir == "" {
		certsDir = "certs"
	}

	// Load CA pool and kuang's own key pair.
	caPool, err := handlers.LoadCAPool(certsDir + "/root_ca.crt")
	if err != nil {
		logger.Error("loading CA pool", "err", err)
		return 1
	}

	kuangCert, kuangKey, err := handlers.LoadKeyPair(certsDir+"/kuang.crt", certsDir+"/kuang.key")
	if err != nil {
		logger.Error("loading kuang key pair", "err", err)
		return 1
	}

	// Initialize sctx auth.
	auth, err := handlers.NewAuth(logger, kuangKey, kuangCert, caPool)
	if err != nil {
		logger.Error("initializing auth", "err", err)
		return 1
	}

	// MCP server.
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "kuang",
		Version: "0.1.0",
	}, &mcp.ServerOptions{
		Logger: logger,
	})

	// Dev server registry.
	reg := handlers.NewRegistry()

	// Register all tools with guard → ice → handler chain.
	ctx := context.Background()
	registerTools(ctx, logger, server, auth, ice, reg)

	// Address.
	addr := os.Getenv("KUANG_ADDR")
	if addr == "" {
		addr = ":7117"
	}

	// SSE handler wrapped with token injection middleware.
	sseHandler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return server
	}, nil)
	handler := handlers.TokenInjector(auth.Admin, logger, sseHandler)

	// TLS config for mTLS.
	tlsCfg, err := handlers.TLSConfig(certsDir+"/kuang.crt", certsDir+"/kuang.key", caPool)
	if err != nil {
		logger.Error("configuring TLS", "err", err)
		return 1
	}

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		TLSConfig:         tlsCfg,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		<-sigCtx.Done()
		if closeErr := srv.Close(); closeErr != nil {
			logger.Error("server close error", "err", closeErr)
		}
	}()

	logger.Info("kuang listening", "addr", addr, "tls", true)
	// ListenAndServeTLS with empty cert/key since TLSConfig already has them.
	if listenErr := srv.ListenAndServeTLS("", ""); listenErr != http.ErrServerClosed {
		logger.Error("server error", "err", listenErr)
		return 1
	}
	return 0
}

// addTool creates a guard for the named permission and registers the tool with
// the full middleware chain: guard → ice filter → tool handler.
func addTool[In, Out any](
	ctx context.Context,
	logger *slog.Logger,
	server *mcp.Server,
	auth *handlers.Auth,
	ice *handlers.ICEClient,
	name string,
	tool *mcp.Tool,
	handler mcp.ToolHandlerFor[In, Out],
) {
	guard, err := auth.CreateGuard(ctx, name)
	if err != nil {
		logger.Error("failed to create guard", "tool", name, "err", err)
		return
	}
	mcp.AddTool(server, tool,
		handlers.GuardedHandler(guard, logger, name,
			handlers.FilteredHandler(ice, logger, name,
				handler)))
}

func registerTools(
	ctx context.Context,
	logger *slog.Logger,
	server *mcp.Server,
	auth *handlers.Auth,
	ice *handlers.ICEClient,
	reg *handlers.Registry,
) {
	// GitHub read tools.
	addTool(ctx, logger, server, auth, ice, "gh_pr_list", &mcp.Tool{Name: "gh_pr_list", Description: "List pull requests for a GitHub repo"}, tools.PRList)
	addTool(ctx, logger, server, auth, ice, "gh_pr_view", &mcp.Tool{Name: "gh_pr_view", Description: "View a specific pull request"}, tools.PRView)
	addTool(ctx, logger, server, auth, ice, "gh_issue_list", &mcp.Tool{Name: "gh_issue_list", Description: "List issues for a GitHub repo"}, tools.IssueList)
	addTool(ctx, logger, server, auth, ice, "gh_issue_view", &mcp.Tool{Name: "gh_issue_view", Description: "View a specific issue"}, tools.IssueView)
	addTool(ctx, logger, server, auth, ice, "gh_repo_view", &mcp.Tool{Name: "gh_repo_view", Description: "View repository details"}, tools.RepoView)

	// GitHub write tools.
	addTool(ctx, logger, server, auth, ice, "gh_pr_create", &mcp.Tool{Name: "gh_pr_create", Description: "Create a pull request on GitHub"}, tools.PRCreate)

	// GitHub Actions tools.
	addTool(ctx, logger, server, auth, ice, "gh_run_list", &mcp.Tool{Name: "gh_run_list", Description: "List recent workflow runs for a GitHub repo"}, tools.RunList)
	addTool(ctx, logger, server, auth, ice, "gh_run_view", &mcp.Tool{Name: "gh_run_view", Description: "View a workflow run with jobs and steps"}, tools.RunView)
	addTool(ctx, logger, server, auth, ice, "gh_run_log", &mcp.Tool{Name: "gh_run_log", Description: "View failed step logs from a workflow run"}, tools.RunLog)
	addTool(ctx, logger, server, auth, ice, "gh_pr_checks", &mcp.Tool{Name: "gh_pr_checks", Description: "View CI check status for a pull request"}, tools.PRChecks)

	// Jira tools.
	addTool(ctx, logger, server, auth, ice, "jira_issue_list", &mcp.Tool{Name: "jira_issue_list", Description: "List Jira issues by project, JQL query, or assignee"}, tools.JiraIssueList)
	addTool(ctx, logger, server, auth, ice, "jira_issue_view", &mcp.Tool{Name: "jira_issue_view", Description: "View a Jira issue with details and comments"}, tools.JiraIssueView)

	// Dev server tools.
	addTool(ctx, logger, server, auth, ice, "dev_start", &mcp.Tool{Name: "dev_start", Description: "Start a dev server for a repo (idempotent, blocks until ready)"}, tools.NewDevStart(reg))
	addTool(ctx, logger, server, auth, ice, "dev_stop", &mcp.Tool{Name: "dev_stop", Description: "Stop a running dev server"}, tools.NewDevStop(reg))
	addTool(ctx, logger, server, auth, ice, "dev_status", &mcp.Tool{Name: "dev_status", Description: "List running dev servers and their status"}, tools.NewDevStatus(reg))
	addTool(ctx, logger, server, auth, ice, "dev_log", &mcp.Tool{Name: "dev_log", Description: "Read recent log output from a dev server"}, tools.NewDevLog(reg))

	// Package manager tools.
	addTool(ctx, logger, server, auth, ice, "pnpm_run", &mcp.Tool{Name: "pnpm_run", Description: "Run a pnpm script by name (e.g. test, lint, build)"}, tools.PnpmRun)
}
