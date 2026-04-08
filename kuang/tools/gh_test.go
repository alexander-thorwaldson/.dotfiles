package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// mockGH replaces ghExec for the duration of a test.
func mockGH(t *testing.T, fn func(ctx context.Context, args ...string) (string, error)) {
	t.Helper()
	orig := ghExec
	ghExec = fn
	t.Cleanup(func() { ghExec = orig })
}

func TestPRList(t *testing.T) {
	mockGH(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "-R", "owner/repo")
		assertContains(t, args, "-s", "open")
		assertContains(t, args, "-L", "5")
		return `[{"number":1,"title":"Fix bug","state":"OPEN","author":{"login":"alice"},"url":"https://github.com/owner/repo/pull/1"}]`, nil
	})

	result, out, err := PRList(context.Background(), nil, PRListParams{
		Repo:  "owner/repo",
		State: "open",
		Limit: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
	if len(out.Items) != 1 {
		t.Fatalf("expected 1 PR, got %d", len(out.Items))
	}
	if out.Items[0].Number != 1 {
		t.Errorf("expected PR #1, got #%d", out.Items[0].Number)
	}
	if out.Items[0].Author.Login != "alice" {
		t.Errorf("expected author alice, got %s", out.Items[0].Author.Login)
	}
}

func TestPRList_DefaultArgs(t *testing.T) {
	mockGH(t, func(_ context.Context, args ...string) (string, error) {
		for _, a := range args {
			if a == "-s" || a == "-L" {
				t.Errorf("unexpected flag %q with empty params", a)
			}
		}
		return `[]`, nil
	})

	_, out, err := PRList(context.Background(), nil, PRListParams{Repo: "owner/repo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Items) != 0 {
		t.Errorf("expected empty list, got %d items", len(out.Items))
	}
}

func TestPRView(t *testing.T) {
	mockGH(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "42")
		return `{"number":42,"title":"Add feature","state":"OPEN","body":"details","author":{"login":"bob"},"url":"https://github.com/o/r/pull/42","reviews":[],"comments":[]}`, nil
	})

	result, out, err := PRView(context.Background(), nil, PRViewParams{Repo: "o/r", Number: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
	if out.Number != 42 {
		t.Errorf("expected PR #42, got #%d", out.Number)
	}
	if out.Title != "Add feature" {
		t.Errorf("expected title %q, got %q", "Add feature", out.Title)
	}
}

func TestIssueList(t *testing.T) {
	mockGH(t, func(_ context.Context, _ ...string) (string, error) {
		return `[{"number":10,"title":"Bug report","state":"OPEN","author":{"login":"carol"},"url":"https://github.com/o/r/issues/10"}]`, nil
	})

	_, out, err := IssueList(context.Background(), nil, IssueListParams{Repo: "o/r"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Items) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(out.Items))
	}
	if out.Items[0].Title != "Bug report" {
		t.Errorf("expected title %q, got %q", "Bug report", out.Items[0].Title)
	}
}

func TestIssueView(t *testing.T) {
	mockGH(t, func(_ context.Context, _ ...string) (string, error) {
		return `{"number":10,"title":"Bug report","state":"OPEN","body":"steps to reproduce","author":{"login":"carol"},"url":"https://github.com/o/r/issues/10","comments":[]}`, nil
	})

	_, out, err := IssueView(context.Background(), nil, IssueViewParams{Repo: "o/r", Number: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Body != "steps to reproduce" {
		t.Errorf("expected body %q, got %q", "steps to reproduce", out.Body)
	}
}

func TestRepoView(t *testing.T) {
	mockGH(t, func(_ context.Context, _ ...string) (string, error) {
		return `{"name":"myrepo","description":"A repo","url":"https://github.com/o/myrepo","defaultBranchRef":{"name":"main"},"languages":[],"issues":{"totalCount":5},"pullRequests":{"totalCount":3}}`, nil
	})

	_, out, err := RepoView(context.Background(), nil, RepoViewParams{Repo: "o/myrepo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Name != "myrepo" {
		t.Errorf("expected name %q, got %q", "myrepo", out.Name)
	}
	if out.DefaultBranchRef.Name != "main" {
		t.Errorf("expected default branch %q, got %q", "main", out.DefaultBranchRef.Name)
	}
	if out.Issues.TotalCount != 5 {
		t.Errorf("expected 5 issues, got %d", out.Issues.TotalCount)
	}
}

func TestGHExecError(t *testing.T) {
	mockGH(t, func(_ context.Context, _ ...string) (string, error) {
		return "", fmt.Errorf("gh: not found")
	})

	result, _, err := PRList(context.Background(), nil, PRListParams{Repo: "o/r"})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError to be true")
	}
}

func TestGHInvalidJSON(t *testing.T) {
	mockGH(t, func(_ context.Context, _ ...string) (string, error) {
		return "not json", nil
	})

	result, _, err := PRList(context.Background(), nil, PRListParams{Repo: "o/r"})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError to be true for invalid JSON")
	}
}

func TestPRListTextContent(t *testing.T) {
	raw := `[{"number":1,"title":"Test","state":"OPEN","author":{"login":"x"},"url":"u"}]`
	mockGH(t, func(_ context.Context, _ ...string) (string, error) {
		return raw, nil
	})

	result, _, _ := PRList(context.Background(), nil, PRListParams{Repo: "o/r"})
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content, got %d", len(result.Content))
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent")
	}
	// Verify the raw JSON is passed through as text content
	var parsed []map[string]any
	if err := json.Unmarshal([]byte(tc.Text), &parsed); err != nil {
		t.Errorf("text content is not valid JSON: %v", err)
	}
}

// -- GitHub Actions tests --

func TestRunList(t *testing.T) {
	mockGH(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "-R", "o/r")
		assertContains(t, args, "-b", "main")
		assertContains(t, args, "-s", "failure")
		return `[{"number":100,"databaseId":9999,"displayTitle":"CI","headBranch":"main","event":"push","status":"completed","conclusion":"failure","workflowName":"CI","url":"u","createdAt":"2026-04-07"}]`, nil
	})

	result, out, err := RunList(context.Background(), nil, RunListParams{
		Repo:   "o/r",
		Branch: "main",
		Status: "failure",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
	if len(out.Items) != 1 {
		t.Fatalf("expected 1 run, got %d", len(out.Items))
	}
	if out.Items[0].Conclusion != "failure" {
		t.Errorf("expected conclusion failure, got %s", out.Items[0].Conclusion)
	}
}

func TestRunList_WithWorkflowAndLimit(t *testing.T) {
	mockGH(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "-w", "CI")
		assertContains(t, args, "-L", "5")
		return `[]`, nil
	})

	_, out, err := RunList(context.Background(), nil, RunListParams{
		Repo:     "o/r",
		Workflow: "CI",
		Limit:    5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Items) != 0 {
		t.Errorf("expected empty list, got %d items", len(out.Items))
	}
}

func TestRunView(t *testing.T) {
	mockGH(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "9999")
		return `{"number":100,"attempt":1,"displayTitle":"CI","headBranch":"main","headSha":"abc","event":"push","status":"completed","conclusion":"success","workflowName":"CI","url":"u","createdAt":"2026-04-07","jobs":[{"name":"build","status":"completed","conclusion":"success","steps":[]}]}`, nil
	})

	result, out, err := RunView(context.Background(), nil, RunViewParams{Repo: "o/r", RunID: 9999})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
	if len(out.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(out.Jobs))
	}
	if out.Jobs[0].Name != "build" {
		t.Errorf("expected job name build, got %s", out.Jobs[0].Name)
	}
}

func TestRunLog(t *testing.T) {
	mockGH(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "--log-failed")
		return "build\tStep 2: FAIL npm test\nError: tests failed", nil
	})

	result, _, err := RunLog(context.Background(), nil, RunLogParams{Repo: "o/r", RunID: 9999})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent")
	}
	if !strings.Contains(tc.Text, "FAIL") {
		t.Errorf("expected log to contain FAIL, got %q", tc.Text)
	}
}

func TestPRChecks(t *testing.T) {
	mockGH(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "42")
		return `[{"name":"CI","state":"SUCCESS","bucket":"pass","description":"ok","workflow":"CI","link":"u","event":"pull_request","startedAt":"t1","completedAt":"t2"}]`, nil
	})

	result, out, err := PRChecks(context.Background(), nil, PRChecksParams{Repo: "o/r", Number: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
	if len(out.Items) != 1 {
		t.Fatalf("expected 1 check, got %d", len(out.Items))
	}
	if out.Items[0].Bucket != "pass" {
		t.Errorf("expected bucket pass, got %s", out.Items[0].Bucket)
	}
}

// assertContains checks that args contains the expected key/value pair.
func assertContains(t *testing.T, args []string, vals ...string) {
	t.Helper()
	for _, v := range vals {
		found := false
		for _, a := range args {
			if a == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected args to contain %q, got %v", v, args)
		}
	}
}
