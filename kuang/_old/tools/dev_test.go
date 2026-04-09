package tools

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/zoobzio/dotfiles/kuang/handlers"
)

func freePort(t *testing.T) int {
	t.Helper()
	var lc net.ListenConfig
	l, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("unexpected address type")
	}
	port := addr.Port
	_ = l.Close()
	return port
}

func TestDevStart_And_Stop(t *testing.T) {
	reg := handlers.NewRegistry()
	port := freePort(t)
	start := NewDevStart(reg)
	stop := NewDevStop(reg)

	result, out, err := start(context.Background(), nil, DevStartParams{
		Repo:    "test-repo",
		Dir:     t.TempDir(),
		Command: fmt.Sprintf("python3 -c 'import socket,time;s=socket.socket();s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1);s.bind((\"127.0.0.1\",%d));s.listen(128);time.sleep(60)'", port),
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
	if out.Status != "ready" {
		t.Errorf("expected status ready, got %s", out.Status)
	}
	if out.Port != port {
		t.Errorf("expected port %d, got %d", port, out.Port)
	}

	// Stop.
	stopResult, stopOut, err := stop(context.Background(), nil, DevStopParams{Repo: "test-repo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stopResult.IsError {
		t.Fatal("expected success")
	}
	if !stopOut.Stopped {
		t.Error("expected Stopped to be true")
	}
}

func TestDevStart_DefaultCommand(t *testing.T) {
	reg := handlers.NewRegistry()
	port := freePort(t)
	start := NewDevStart(reg)

	// pnpm dev won't exist, so this will fail — but we're testing the default is applied.
	_, _, _ = start(context.Background(), nil, DevStartParams{
		Repo:    "test-repo",
		Dir:     t.TempDir(),
		Port:    port,
		Timeout: 1,
	})

	// Check that the registry entry was created with the default command.
	e, ok := reg.Get("test-repo")
	if ok && e.Command != "pnpm dev" {
		t.Errorf("expected default command 'pnpm dev', got %q", e.Command)
	}

	reg.Stop("test-repo")
}

func TestDevStart_Idempotent(t *testing.T) {
	reg := handlers.NewRegistry()
	port := freePort(t)
	start := NewDevStart(reg)

	_, out1, err := start(context.Background(), nil, DevStartParams{
		Repo:    "repo",
		Dir:     t.TempDir(),
		Command: fmt.Sprintf("python3 -c 'import socket,time;s=socket.socket();s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1);s.bind((\"127.0.0.1\",%d));s.listen(128);time.sleep(60)'", port),
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("first start error: %v", err)
	}

	// Second call should be a no-op.
	_, out2, err := start(context.Background(), nil, DevStartParams{
		Repo:    "repo",
		Dir:     t.TempDir(),
		Command: "different-cmd",
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("second start error: %v", err)
	}
	if out2.Port != out1.Port {
		t.Error("expected idempotent return")
	}

	reg.Stop("repo")
}

func TestDevStop_Nonexistent(t *testing.T) {
	reg := handlers.NewRegistry()
	stop := NewDevStop(reg)

	result, out, err := stop(context.Background(), nil, DevStopParams{Repo: "nope"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success even for non-existent")
	}
	if out.Stopped {
		t.Error("expected Stopped to be false")
	}
}

func TestDevStatus_Empty(t *testing.T) {
	reg := handlers.NewRegistry()
	status := NewDevStatus(reg)

	_, out, err := status(context.Background(), nil, DevStatusParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(out.Services))
	}
}

func TestDevStatus_WithRunning(t *testing.T) {
	reg := handlers.NewRegistry()
	port := freePort(t)
	start := NewDevStart(reg)
	status := NewDevStatus(reg)

	_, _, err := start(context.Background(), nil, DevStartParams{
		Repo:    "repo",
		Dir:     t.TempDir(),
		Command: fmt.Sprintf("python3 -c 'import socket,time;s=socket.socket();s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1);s.bind((\"127.0.0.1\",%d));s.listen(128);time.sleep(60)'", port),
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("start error: %v", err)
	}

	_, out, err := status(context.Background(), nil, DevStatusParams{})
	if err != nil {
		t.Fatalf("status error: %v", err)
	}
	if len(out.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(out.Services))
	}
	if out.Services[0].Repo != "repo" {
		t.Errorf("expected repo, got %s", out.Services[0].Repo)
	}

	reg.Stop("repo")
}

func TestDevLog_NoServer(t *testing.T) {
	reg := handlers.NewRegistry()
	log := NewDevLog(reg)

	result, _, err := log(context.Background(), nil, DevLogParams{Repo: "nope"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error for non-existent server")
	}
}

func TestDevLog_WithOutput(t *testing.T) {
	reg := handlers.NewRegistry()
	port := freePort(t)
	start := NewDevStart(reg)
	logFn := NewDevLog(reg)

	_, _, err := start(context.Background(), nil, DevStartParams{
		Repo:    "repo",
		Dir:     t.TempDir(),
		Command: fmt.Sprintf("python3 -c 'import socket,time,sys;sys.stdout.write(\"hello\\n\");sys.stdout.flush();s=socket.socket();s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1);s.bind((\"127.0.0.1\",%d));s.listen(128);time.sleep(60)'", port),
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("start error: %v", err)
	}

	// Give the output pipe a moment to flush.
	time.Sleep(100 * time.Millisecond)

	_, out, err := logFn(context.Background(), nil, DevLogParams{Repo: "repo", Lines: 10})
	if err != nil {
		t.Fatalf("log error: %v", err)
	}
	if len(out.Lines) == 0 {
		t.Error("expected at least one log line")
	}

	reg.Stop("repo")
}
