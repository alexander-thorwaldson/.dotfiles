package tools

import (
	"context"
	"fmt"
	"testing"
)

func mockJira(t *testing.T, fn func(ctx context.Context, args ...string) (string, error)) {
	t.Helper()
	orig := jiraExec
	jiraExec = fn
	t.Cleanup(func() { jiraExec = orig })
}

func TestJiraIssueList(t *testing.T) {
	mockJira(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "-p", "PROJ")
		return `[{"key":"PROJ-1","summary":"Fix login","status":"In Progress","assignee":"alice","priority":"High","type":"Bug"}]`, nil
	})

	result, out, err := JiraIssueList(context.Background(), nil, JiraIssueListParams{
		Project: "PROJ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
	if len(out.Items) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(out.Items))
	}
	if out.Items[0].Key != "PROJ-1" {
		t.Errorf("expected key PROJ-1, got %s", out.Items[0].Key)
	}
	if out.Items[0].Status != "In Progress" {
		t.Errorf("expected status In Progress, got %s", out.Items[0].Status)
	}
}

func TestJiraIssueList_WithJQL(t *testing.T) {
	mockJira(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "-q", "type = Bug")
		return `[]`, nil
	})

	_, out, err := JiraIssueList(context.Background(), nil, JiraIssueListParams{
		Query: "type = Bug",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Items) != 0 {
		t.Errorf("expected empty list, got %d items", len(out.Items))
	}
}

func TestJiraIssueList_StatusFilter(t *testing.T) {
	mockJira(t, func(_ context.Context, _ ...string) (string, error) {
		return `[
			{"key":"PROJ-1","summary":"A","status":"Done","assignee":"","priority":"Low","type":"Task"},
			{"key":"PROJ-2","summary":"B","status":"In Progress","assignee":"","priority":"High","type":"Bug"},
			{"key":"PROJ-3","summary":"C","status":"Done","assignee":"","priority":"Medium","type":"Story"}
		]`, nil
	})

	_, out, err := JiraIssueList(context.Background(), nil, JiraIssueListParams{
		Project: "PROJ",
		Status:  "Done",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Items) != 2 {
		t.Fatalf("expected 2 Done issues, got %d", len(out.Items))
	}
	for _, item := range out.Items {
		if item.Status != "Done" {
			t.Errorf("expected status Done, got %s", item.Status)
		}
	}
}

func TestJiraIssueList_WithLimit(t *testing.T) {
	mockJira(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "-l", "5")
		return `[]`, nil
	})

	_, _, err := JiraIssueList(context.Background(), nil, JiraIssueListParams{
		Project: "PROJ",
		Limit:   5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJiraIssueView(t *testing.T) {
	mockJira(t, func(_ context.Context, args ...string) (string, error) {
		assertContains(t, args, "PROJ-42")
		return `{"key":"PROJ-42","summary":"Fix auth","description":"Auth is broken","status":"Open","assignee":"bob","reporter":"carol","priority":"Critical","type":"Bug","labels":["security"],"comments":[{"author":"dave","body":"Looking into it","created":"2026-04-01"}]}`, nil
	})

	result, out, err := JiraIssueView(context.Background(), nil, JiraIssueViewParams{
		Issue: "PROJ-42",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
	if out.Key != "PROJ-42" {
		t.Errorf("expected key PROJ-42, got %s", out.Key)
	}
	if out.Description != "Auth is broken" {
		t.Errorf("expected description %q, got %q", "Auth is broken", out.Description)
	}
	if len(out.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(out.Comments))
	}
	if out.Comments[0].Author != "dave" {
		t.Errorf("expected comment author dave, got %s", out.Comments[0].Author)
	}
	if len(out.Labels) != 1 || out.Labels[0] != "security" {
		t.Errorf("unexpected labels: %v", out.Labels)
	}
}

func TestJiraIssueList_Error(t *testing.T) {
	mockJira(t, func(_ context.Context, _ ...string) (string, error) {
		return "", fmt.Errorf("jira: not configured")
	})

	result, _, err := JiraIssueList(context.Background(), nil, JiraIssueListParams{Project: "PROJ"})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError to be true")
	}
}

func TestJiraIssueView_InvalidJSON(t *testing.T) {
	mockJira(t, func(_ context.Context, _ ...string) (string, error) {
		return "not json", nil
	})

	result, _, err := JiraIssueView(context.Background(), nil, JiraIssueViewParams{Issue: "PROJ-1"})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError for invalid JSON")
	}
}
