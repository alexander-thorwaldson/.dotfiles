package services

import (
	"context"
	"testing"

	"github.com/zoobzio/dotfiles/kuang/wire"
)

func TestJira_IssueList(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, name string, _ ...string) (string, error) {
		if name != "jira" {
			t.Errorf("expected jira, got %s", name)
		}
		return `[{"key":"PROJ-1","summary":"Bug","status":"Open"}]`, nil
	})

	jira := NewJira(ice)
	out, err := jira.IssueList(context.Background(), wire.JiraIssueListRequest{Project: "PROJ"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestJira_IssueView(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		return `{"key":"PROJ-42","summary":"Fix auth"}`, nil
	})

	jira := NewJira(ice)
	out, err := jira.IssueView(context.Background(), wire.JiraIssueViewRequest{Issue: "PROJ-42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestJira_ICEBlocks(t *testing.T) {
	srv := newTestICEServer(t, "INJECTION", 0.99)
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		t.Fatal("CLI should not be called when ice blocks")
		return "", nil
	})

	jira := NewJira(ice)
	_, err := jira.IssueList(context.Background(), wire.JiraIssueListRequest{Project: "PROJ"})
	if err == nil {
		t.Fatal("expected ice to block")
	}
}
