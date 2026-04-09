// Package main implements the kuang-proxy, a stdio MCP server that forwards
// tool calls to the kuang API over mTLS. It runs inside agent containers.
package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "kuang-proxy: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	url := os.Getenv("KUANG_URL")
	if url == "" {
		url = "https://host.docker.internal:7117"
	}
	certFile := os.Getenv("KUANG_CERT")
	keyFile := os.Getenv("KUANG_KEY")
	caFile := os.Getenv("KUANG_CA")

	if certFile == "" || keyFile == "" || caFile == "" {
		return fmt.Errorf("KUANG_CERT, KUANG_KEY, and KUANG_CA are required")
	}

	client, err := newMTLSClient(certFile, keyFile, caFile)
	if err != nil {
		return fmt.Errorf("creating mTLS client: %w", err)
	}

	// Discover tools from kuang.
	tools, err := discoverTools(client, url)
	if err != nil {
		return fmt.Errorf("discovering tools: %w", err)
	}

	// Build MCP server with discovered tools.
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "kuang-proxy",
		Version: "0.1.0",
	}, nil)

	for _, tool := range tools {
		registerTool(server, client, url, tool)
	}

	// Serve over stdio.
	return server.Run(context.Background(), &mcp.StdioTransport{})
}

// toolInfo matches wire.ToolInfo from kuang's API.
type toolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// toolListResponse matches wire.ToolListResponse.
type toolListResponse struct {
	Tools []toolInfo `json:"tools"`
}

// toolEndpoint maps a tool name to its kuang API path.
var toolEndpoint = map[string]string{
	"gh_pr_list":      "/v1/gh/prs",
	"gh_pr_view":      "/v1/gh/prs/view",
	"gh_pr_create":    "/v1/gh/prs/create",
	"gh_pr_checks":    "/v1/gh/prs/checks",
	"gh_issue_list":   "/v1/gh/issues",
	"gh_issue_view":   "/v1/gh/issues/view",
	"gh_repo_view":    "/v1/gh/repos/view",
	"gh_run_list":     "/v1/gh/runs",
	"gh_run_view":     "/v1/gh/runs/view",
	"gh_run_log":      "/v1/gh/runs/log",
	"jira_issue_list": "/v1/jira/issues",
	"jira_issue_view": "/v1/jira/issues/view",
	"dev_start":       "/v1/dev/start",
	"dev_stop":        "/v1/dev/stop",
	"dev_status":      "/v1/dev/status",
	"dev_log":         "/v1/dev/log",
	"pnpm_run":        "/v1/pnpm/run",
}

// toolMethod maps tool names to HTTP methods. GET for reads without body, POST for everything else.
var toolMethod = map[string]string{
	"dev_status": http.MethodGet,
}

func discoverTools(client *http.Client, baseURL string) ([]toolInfo, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseURL+"/v1/tools", nil) // #nosec G704 -- baseURL is from controlled env var
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	resp, err := client.Do(req) // #nosec G704
	if err != nil {
		return nil, fmt.Errorf("requesting tools: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tools endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var result toolListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding tools: %w", err)
	}
	return result.Tools, nil
}

func registerTool(server *mcp.Server, client *http.Client, baseURL string, tool toolInfo) {
	name := tool.Name
	endpoint, ok := toolEndpoint[name]
	if !ok {
		return
	}

	server.AddTool(
		&mcp.Tool{Name: name, Description: tool.Description},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return callKuang(ctx, client, baseURL, name, endpoint, req)
		},
	)
}

func callKuang(ctx context.Context, client *http.Client, baseURL, name, endpoint string, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := http.MethodPost
	if m, ok := toolMethod[name]; ok {
		method = m
	}

	var body io.Reader
	if method == http.MethodPost && req.Params != nil && req.Params.Arguments != nil {
		data, err := json.Marshal(req.Params.Arguments)
		if err != nil {
			return nil, fmt.Errorf("marshalling arguments: %w", err)
		}
		body = bytes.NewReader(data)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, baseURL+endpoint, body) // #nosec G704 -- endpoint is from controlled map
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(httpReq) // #nosec G704
	if err != nil {
		return nil, fmt.Errorf("calling kuang: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		r := &mcp.CallToolResult{}
		r.SetError(fmt.Errorf("kuang returned %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody))))
		return r, nil
	}

	// Extract the output field from the response.
	var result struct {
		Output   string `json:"output"`
		Content  string `json:"content"`
		ExitCode int    `json:"exitCode"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		// If we can't parse, return the raw response as text.
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(respBody)}},
		}, nil
	}

	// Use output field if present, otherwise use the full response.
	text := result.Output
	if text == "" && result.Content != "" {
		text = result.Content
	}
	if text == "" {
		text = string(respBody)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, nil
}

func newMTLSClient(certFile, keyFile, caFile string) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("loading client cert: %w", err)
	}

	caCert, err := os.ReadFile(caFile) // #nosec G304,G703 -- controlled path from env var
	if err != nil {
		return nil, fmt.Errorf("reading CA cert: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA cert")
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      pool,
				MinVersion:   tls.VersionTLS13,
			},
		},
	}, nil
}
