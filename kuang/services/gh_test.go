package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/zoobzio/dotfiles/kuang/wire"
)

func mockCLI(t *testing.T, fn func(ctx context.Context, name string, args ...string) (string, error)) {
	t.Helper()
	orig := cliExec
	cliExec = fn
	t.Cleanup(func() { cliExec = orig })
}

func safeICE(t *testing.T) *ICEClient {
	t.Helper()
	srv := newTestICEServer(t, "SAFE", 0.01)
	t.Cleanup(srv.Close)
	return NewICEClient(srv.URL)
}

func TestGH_PRList(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, name string, _ ...string) (string, error) {
		if name != "gh" {
			return "", fmt.Errorf("expected gh, got %s", name)
		}
		return `[{"number":1,"title":"Fix","state":"OPEN"}]`, nil
	})

	gh := NewGH(ice)
	out, err := gh.PRList(context.Background(), wire.PRListRequest{Repo: "o/r", State: "open", Limit: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestGH_PRView(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		return `{"number":42,"title":"Feature"}`, nil
	})

	gh := NewGH(ice)
	out, err := gh.PRView(context.Background(), wire.PRViewRequest{Repo: "o/r", Number: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestGH_PRCreate(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		return "https://github.com/o/r/pull/99\n", nil
	})

	gh := NewGH(ice)
	out, err := gh.PRCreate(context.Background(), wire.PRCreateRequest{Repo: "o/r", Title: "New", Head: "feat"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestGH_IssueList(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		return `[{"number":10,"title":"Bug"}]`, nil
	})

	gh := NewGH(ice)
	out, err := gh.IssueList(context.Background(), wire.IssueListRequest{Repo: "o/r"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestGH_RepoView(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		return `{"name":"repo"}`, nil
	})

	gh := NewGH(ice)
	out, err := gh.RepoView(context.Background(), wire.RepoViewRequest{Repo: "o/r"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestGH_RunList(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		return `[{"number":100}]`, nil
	})

	gh := NewGH(ice)
	out, err := gh.RunList(context.Background(), wire.RunListRequest{Repo: "o/r", Branch: "main", Status: "failure"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestGH_CLIError(t *testing.T) {
	ice := safeICE(t)
	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		return "", fmt.Errorf("gh: not found")
	})

	gh := NewGH(ice)
	_, err := gh.PRList(context.Background(), wire.PRListRequest{Repo: "o/r"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGH_ICEBlocks(t *testing.T) {
	srv := newTestICEServer(t, "INJECTION", 0.99)
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	mockCLI(t, func(_ context.Context, _ string, _ ...string) (string, error) {
		t.Fatal("CLI should not be called when ice blocks input")
		return "", nil
	})

	gh := NewGH(ice)
	_, err := gh.PRList(context.Background(), wire.PRListRequest{Repo: "o/r"})
	if err == nil {
		t.Fatal("expected ice to block")
	}
}
