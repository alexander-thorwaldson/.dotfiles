package tools

import (
	"context"
	"encoding/json"
	"fmt"
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
