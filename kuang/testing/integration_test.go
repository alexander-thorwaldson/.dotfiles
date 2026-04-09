package testing

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zoobz-io/sctx"
	"github.com/zoobz-io/sum"
	"github.com/zoobzio/dotfiles/kuang/contracts"
	"github.com/zoobzio/dotfiles/kuang/handlers"
	"github.com/zoobzio/dotfiles/kuang/services"
	"github.com/zoobzio/dotfiles/kuang/wire"
)

var (
	testCA         *TestCA
	kuangBaseURL   string
	engineerClient *http.Client
	architectClient *http.Client
	noCertClient   *http.Client
	badRoleClient  *http.Client
)

func TestMain(m *testing.M) {
	sctx.ResetAdminForTesting()

	// Generate test PKI.
	t := &testing.T{}
	testCA = NewTestCA(t)
	serverCert := testCA.Issue(t, "kuang", nil)
	engineerCert := testCA.Issue(t, "case", []string{"engineer"})
	architectCert := testCA.Issue(t, "dixie", []string{"architect"})
	badRoleCert := testCA.Issue(t, "hacker", []string{"superadmin"})

	// sctx admin.
	roleExpand := func(role string) []string {
		return services.ExpandRole(services.Role(role))
	}
	admin, err := sctx.NewAdminService[services.AgentMeta](serverCert.Key, testCA.Pool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating admin: %v\n", err)
		os.Exit(1)
	}
	if err := admin.SetPolicy(services.AgentPolicy(roleExpand)); err != nil {
		fmt.Fprintf(os.Stderr, "setting policy: %v\n", err)
		os.Exit(1)
	}
	if _, err := admin.GenerateTrusted(context.Background(), serverCert.Cert); err != nil {
		fmt.Fprintf(os.Stderr, "authenticating kuang: %v\n", err)
		os.Exit(1)
	}

	// Mock ICE — unreachable but won't be called for non-CLI tests.
	ice := services.NewICEClient("http://127.0.0.1:1")

	// Registry.
	k := sum.Start()
	sum.Register[contracts.GH](k, services.NewGH(ice))
	sum.Register[contracts.Jira](k, services.NewJira(ice))
	sum.Register[contracts.Dev](k, services.NewDev(ice))
	sum.Register[contracts.Pnpm](k, services.NewPnpm(ice))
	sum.Freeze(k)

	// Engine.
	svc := sum.New()
	svc.Engine().
		WithTLSConfig(&tls.Config{
			Certificates: []tls.Certificate{serverCert.TLS},
			ClientCAs:    testCA.Pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			MinVersion:   tls.VersionTLS13,
		}).
		WithAuthenticator(services.Authenticator(admin, roleExpand))
	svc.Handle(handlers.All()...)

	// Find a free port.
	lc := net.ListenConfig{}
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		os.Exit(1)
	}
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		fmt.Fprintf(os.Stderr, "unexpected address type: %T\n", ln.Addr())
		os.Exit(1)
	}
	port := tcpAddr.Port
	addr := tcpAddr.String()
	_ = ln.Close()

	go func() {
		_ = svc.Run("127.0.0.1", port)
	}()

	// Wait for server.
	kuangBaseURL = fmt.Sprintf("https://%s", addr)
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		dialer := &tls.Dialer{
			NetDialer: &net.Dialer{Timeout: 200 * time.Millisecond},
			Config: &tls.Config{
				Certificates: []tls.Certificate{serverCert.TLS},
				RootCAs:      testCA.Pool,
				MinVersion:   tls.VersionTLS13,
			},
		}
		conn, dialErr := dialer.DialContext(context.Background(), "tcp", addr)
		if dialErr == nil {
			_ = conn.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Build clients.
	engineerClient = MTLSClient(testCA, engineerCert)
	architectClient = MTLSClient(testCA, architectCert)
	badRoleClient = MTLSClient(testCA, badRoleCert)
	noCertClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    testCA.Pool,
				MinVersion: tls.VersionTLS13,
			},
		},
	}

	os.Exit(m.Run())
}

func doGet(t *testing.T, client *http.Client, url string) (int, []byte) {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

func doPost(t *testing.T, client *http.Client, url string, payload any) (int, []byte) {
	t.Helper()
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshalling payload: %v", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

func toolNames(t *testing.T, body []byte) []string {
	t.Helper()
	var resp wire.ToolListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decoding tools: %v", err)
	}
	names := make([]string, len(resp.Tools))
	for i, tool := range resp.Tools {
		names[i] = tool.Name
	}
	return names
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// -- Discovery tests --

func TestDiscovery_EngineerSeesWriteTools(t *testing.T) {
	status, body := doGet(t, engineerClient, kuangBaseURL+"/v1/tools")
	if status != 200 {
		t.Fatalf("expected 200, got %d: %s", status, body)
	}
	names := toolNames(t, body)
	if !contains(names, "gh_pr_create") {
		t.Error("engineer should see gh_pr_create")
	}
	if !contains(names, "pnpm_run") {
		t.Error("engineer should see pnpm_run")
	}
	if !contains(names, "dev_start") {
		t.Error("engineer should see dev_start")
	}
}

func TestDiscovery_ArchitectCannotSeeWriteTools(t *testing.T) {
	status, body := doGet(t, architectClient, kuangBaseURL+"/v1/tools")
	if status != 200 {
		t.Fatalf("expected 200, got %d: %s", status, body)
	}
	names := toolNames(t, body)
	if contains(names, "gh_pr_create") {
		t.Error("architect should NOT see gh_pr_create")
	}
	if contains(names, "pnpm_run") {
		t.Error("architect should NOT see pnpm_run")
	}
	if contains(names, "dev_start") {
		t.Error("architect should NOT see dev_start")
	}
	if !contains(names, "gh_pr_view") {
		t.Error("architect should see gh_pr_view")
	}
	if !contains(names, "jira_issue_list") {
		t.Error("architect should see jira_issue_list")
	}
}

// -- Authorization tests --

func TestAuthorization_ArchitectCannotCallWriteTool(t *testing.T) {
	status, _ := doPost(t, architectClient, kuangBaseURL+"/v1/gh/prs/create", wire.PRCreateRequest{
		Repo: "o/r", Title: "Test", Head: "branch",
	})
	if status != 403 {
		t.Errorf("expected 403 for architect calling gh_pr_create, got %d", status)
	}
}

func TestAuthorization_EngineerCanCallReadTool(t *testing.T) {
	status, _ := doGet(t, engineerClient, kuangBaseURL+"/v1/dev/status")
	if status != 200 {
		t.Errorf("expected 200 for engineer calling dev_status, got %d", status)
	}
}

// -- mTLS tests --

func TestNoClientCert_Rejected(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, kuangBaseURL+"/v1/tools", nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}
	resp, err := noCertClient.Do(req)
	if err == nil {
		_ = resp.Body.Close()
		t.Error("expected TLS error with no client cert")
	}
}

// -- Unknown role tests --

func TestUnknownRole_Rejected(t *testing.T) {
	status, _ := doGet(t, badRoleClient, kuangBaseURL+"/v1/tools")
	if status == 200 {
		t.Error("expected non-200 for unknown role")
	}
}
