package handlers

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
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

func TestLogBuffer_Write_And_Lines(t *testing.T) {
	buf := NewLogBuffer(3)
	buf.Write("a")
	buf.Write("b")
	buf.Write("c")
	buf.Write("d")

	lines := buf.Lines(0)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "b" || lines[1] != "c" || lines[2] != "d" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestLogBuffer_Lines_Partial(t *testing.T) {
	buf := NewLogBuffer(10)
	buf.Write("x")
	buf.Write("y")

	lines := buf.Lines(1)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != "y" {
		t.Errorf("expected y, got %s", lines[0])
	}
}

func TestLogBuffer_Lines_MoreThanAvailable(t *testing.T) {
	buf := NewLogBuffer(10)
	buf.Write("only")

	lines := buf.Lines(100)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestRegistry_Start_And_Stop(t *testing.T) {
	reg := NewRegistry()
	port := freePort(t)

	// Start a simple TCP listener as our "dev server".
	entry, err := reg.Start("test-repo", t.TempDir(), fmt.Sprintf("python3 -c 'import socket,time;s=socket.socket();s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1);s.bind((\"127.0.0.1\",%d));s.listen(128);time.sleep(60)'", port), port, 5*time.Second)
	if err != nil {
		t.Fatalf("start error: %v", err)
	}
	if entry.Status != "ready" {
		t.Errorf("expected status ready, got %s", entry.Status)
	}
	if entry.PID == 0 {
		t.Error("expected non-zero PID")
	}

	// Verify it's in the registry.
	e, ok := reg.Get("test-repo")
	if !ok {
		t.Fatal("expected entry in registry")
	}
	if e.Port != port {
		t.Errorf("expected port %d, got %d", port, e.Port)
	}

	// Stop it.
	reg.Stop("test-repo")
	_, ok = reg.Get("test-repo")
	if ok {
		t.Error("expected entry to be removed after stop")
	}
}

func TestRegistry_Start_Idempotent(t *testing.T) {
	reg := NewRegistry()
	port := freePort(t)

	entry1, err := reg.Start("repo", t.TempDir(), fmt.Sprintf("python3 -c 'import socket,time;s=socket.socket();s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1);s.bind((\"127.0.0.1\",%d));s.listen(128);time.sleep(60)'", port), port, 5*time.Second)
	if err != nil {
		t.Fatalf("first start error: %v", err)
	}

	// Second start should return the same entry.
	entry2, err := reg.Start("repo", t.TempDir(), "pnpm dev", port, 5*time.Second)
	if err != nil {
		t.Fatalf("second start error: %v", err)
	}
	if entry2.PID != entry1.PID {
		t.Errorf("expected same PID %d, got %d (not idempotent)", entry1.PID, entry2.PID)
	}

	reg.Stop("repo")
}

func TestRegistry_Start_PortInUse(t *testing.T) {
	// Occupy a port externally.
	var lc net.ListenConfig
	l, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = l.Close() }()
	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("unexpected address type")
	}
	port := addr.Port

	reg := NewRegistry()
	_, err = reg.Start("repo", t.TempDir(), "pnpm dev", port, 2*time.Second)
	if err == nil {
		t.Fatal("expected error for port already in use")
	}
}

func TestRegistry_Stop_Nonexistent(_ *testing.T) {
	reg := NewRegistry()
	reg.Stop("nope") // should not panic
}

func TestRegistry_All_Empty(t *testing.T) {
	reg := NewRegistry()
	entries := reg.All()
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestIsPortOpen_Closed(t *testing.T) {
	port := freePort(t)
	if isPortOpen(port) {
		t.Error("expected port to be closed")
	}
}

func TestIsPortOpen_Open(t *testing.T) {
	var lc net.ListenConfig
	l, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = l.Close() }()
	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("unexpected address type")
	}
	port := addr.Port
	if !isPortOpen(port) {
		t.Error("expected port to be open")
	}
}

func TestWaitForPort_AlreadyOpen(t *testing.T) {
	var lc net.ListenConfig
	l, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = l.Close() }()
	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("unexpected address type")
	}
	port := addr.Port

	if err := waitForPort(port, 1*time.Second); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestWaitForPort_Timeout(t *testing.T) {
	port := freePort(t)
	err := waitForPort(port, 500*time.Millisecond)
	if err == nil {
		t.Error("expected timeout error")
	}
}
