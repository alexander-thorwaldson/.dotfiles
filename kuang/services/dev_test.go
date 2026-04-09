package services

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/zoobzio/dotfiles/kuang/wire"
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

func pythonListener(port int) string {
	return fmt.Sprintf("python3 -c 'import socket,time;s=socket.socket();s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1);s.bind((\"127.0.0.1\",%d));s.listen(128);time.sleep(60)'", port)
}

func TestDev_StartAndStop(t *testing.T) {
	ice := safeICE(t)
	dev := NewDev(ice)
	port := freePort(t)

	resp, err := dev.Start(context.Background(), wire.DevStartRequest{
		Repo:    "test-repo",
		Dir:     t.TempDir(),
		Command: pythonListener(port),
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("start error: %v", err)
	}
	if resp.Status != "ready" {
		t.Errorf("expected ready, got %s", resp.Status)
	}

	stopResp, err := dev.Stop(context.Background(), wire.DevStopRequest{Repo: "test-repo"})
	if err != nil {
		t.Fatalf("stop error: %v", err)
	}
	if !stopResp.Stopped {
		t.Error("expected Stopped true")
	}
}

func TestDev_Start_Idempotent(t *testing.T) {
	ice := safeICE(t)
	dev := NewDev(ice)
	port := freePort(t)

	resp1, err := dev.Start(context.Background(), wire.DevStartRequest{
		Repo:    "repo",
		Dir:     t.TempDir(),
		Command: pythonListener(port),
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("first start: %v", err)
	}

	resp2, err := dev.Start(context.Background(), wire.DevStartRequest{
		Repo:    "repo",
		Dir:     t.TempDir(),
		Command: "different",
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("second start: %v", err)
	}
	if resp2.Port != resp1.Port {
		t.Error("expected idempotent return")
	}

	dev.stop("repo")
}

func TestDev_Stop_Nonexistent(t *testing.T) {
	ice := safeICE(t)
	dev := NewDev(ice)

	resp, err := dev.Stop(context.Background(), wire.DevStopRequest{Repo: "nope"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Stopped {
		t.Error("expected Stopped false")
	}
}

func TestDev_Status_Empty(t *testing.T) {
	ice := safeICE(t)
	dev := NewDev(ice)

	resp, err := dev.Status(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(resp.Services))
	}
}

func TestDev_Log_NoServer(t *testing.T) {
	ice := safeICE(t)
	dev := NewDev(ice)

	_, err := dev.Log(context.Background(), wire.DevLogRequest{Repo: "nope"})
	if err == nil {
		t.Error("expected error for non-existent server")
	}
}

func TestDev_Log_WithOutput(t *testing.T) {
	ice := safeICE(t)
	dev := NewDev(ice)
	port := freePort(t)

	_, err := dev.Start(context.Background(), wire.DevStartRequest{
		Repo:    "repo",
		Dir:     t.TempDir(),
		Command: fmt.Sprintf("python3 -c 'import socket,time,sys;sys.stdout.write(\"hello\\n\");sys.stdout.flush();s=socket.socket();s.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR,1);s.bind((\"127.0.0.1\",%d));s.listen(128);time.sleep(60)'", port),
		Port:    port,
		Timeout: 5,
	})
	if err != nil {
		t.Fatalf("start error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	resp, err := dev.Log(context.Background(), wire.DevLogRequest{Repo: "repo", Lines: 10})
	if err != nil {
		t.Fatalf("log error: %v", err)
	}
	if len(resp.Lines) == 0 {
		t.Error("expected at least one log line")
	}

	dev.stop("repo")
}

func TestDev_PortInUse(t *testing.T) {
	ice := safeICE(t)
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

	dev := NewDev(ice)
	_, err = dev.Start(context.Background(), wire.DevStartRequest{
		Repo: "repo", Dir: t.TempDir(), Command: "pnpm dev", Port: addr.Port, Timeout: 2,
	})
	if err == nil {
		t.Fatal("expected error for port in use")
	}
}
