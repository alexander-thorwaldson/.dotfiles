package tools

import (
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestResult(t *testing.T) {
	type Out struct {
		Name string `json:"name"`
	}

	result, out, err := Result("hello", Out{Name: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Name != "test" {
		t.Errorf("expected out.Name = %q, got %q", "test", out.Name)
	}
	if result.IsError {
		t.Error("expected IsError to be false")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}
	if tc.Text != "hello" {
		t.Errorf("expected text %q, got %q", "hello", tc.Text)
	}
}

func TestErrResult(t *testing.T) {
	type Out struct {
		Value int `json:"value"`
	}

	result, out, err := ErrResult[Out](errors.New("something broke"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Value != 0 {
		t.Errorf("expected zero value, got %+v", out)
	}
	if !result.IsError {
		t.Error("expected IsError to be true")
	}
	if result.GetError() == nil {
		t.Error("expected GetError to return non-nil")
	}
}
