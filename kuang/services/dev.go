package services

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/zoobzio/dotfiles/kuang/wire"
)

// Dev implements contracts.Dev with an in-memory process registry and ice filtering.
type Dev struct {
	ice      *ICEClient
	services map[string]*serviceEntry
	mu       sync.Mutex
}

type serviceEntry struct {
	startedAt time.Time
	cmd       *exec.Cmd
	logBuf    *logBuffer
	cancel    context.CancelFunc
	repo      string
	dir       string
	command   string
	status    string
	port      int
	pid       int
}

// NewDev creates a new Dev service with ice filtering.
func NewDev(ice *ICEClient) *Dev {
	return &Dev{ice: ice, services: make(map[string]*serviceEntry)}
}

// Start launches a dev server. Idempotent — returns existing if running.
func (d *Dev) Start(_ context.Context, req wire.DevStartRequest) (*wire.DevStartResponse, error) {
	if err := scanInput(d.ice, req); err != nil {
		return nil, err
	}

	cmd := req.Command
	if cmd == "" {
		cmd = "pnpm dev"
	}
	timeout := 30 * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	d.mu.Lock()

	// Idempotent check.
	if e, ok := d.services[req.Repo]; ok {
		d.mu.Unlock()
		if isPortOpen(e.port) {
			return &wire.DevStartResponse{Repo: e.repo, Port: e.port, Status: e.status, URL: fmt.Sprintf("http://localhost:%d", e.port)}, nil
		}
		d.stop(req.Repo)
		d.mu.Lock()
	}

	if isPortOpen(req.Port) {
		d.mu.Unlock()
		return nil, fmt.Errorf("port %d is already in use", req.Port)
	}

	ctx, cancel := context.WithCancel(context.Background())
	proc := exec.CommandContext(ctx, "sh", "-c", cmd) // #nosec G204 //nolint:gosec
	proc.Dir = req.Dir
	proc.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", req.Port))

	buf := newLogBuffer(1000)
	stdout, err := proc.StdoutPipe()
	if err != nil {
		cancel()
		d.mu.Unlock()
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := proc.StderrPipe()
	if err != nil {
		cancel()
		d.mu.Unlock()
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := proc.Start(); err != nil {
		cancel()
		d.mu.Unlock()
		return nil, fmt.Errorf("starting process: %w", err)
	}

	entry := &serviceEntry{
		repo: req.Repo, dir: req.Dir, command: cmd,
		port: req.Port, pid: proc.Process.Pid,
		cmd: proc, cancel: cancel, logBuf: buf,
		status: "starting", startedAt: time.Now(),
	}
	d.services[req.Repo] = entry
	d.mu.Unlock()

	go pipeToLog(stdout, buf)
	go pipeToLog(stderr, buf)
	go func() {
		_ = proc.Wait()
		d.mu.Lock()
		if cur, ok := d.services[req.Repo]; ok && cur == entry {
			cur.status = "stopped"
		}
		d.mu.Unlock()
	}()

	if err := waitForPort(req.Port, timeout); err != nil {
		d.mu.Lock()
		entry.status = "failed"
		d.mu.Unlock()
		return nil, fmt.Errorf("dev server did not become ready: %w", err)
	}

	d.mu.Lock()
	entry.status = "ready"
	d.mu.Unlock()

	return &wire.DevStartResponse{
		Repo: req.Repo, Port: req.Port, Status: "ready",
		URL: fmt.Sprintf("http://localhost:%d", req.Port),
	}, nil
}

// Stop terminates a dev server.
func (d *Dev) Stop(_ context.Context, req wire.DevStopRequest) (*wire.DevStopResponse, error) {
	if err := scanInput(d.ice, req); err != nil {
		return nil, err
	}

	d.mu.Lock()
	_, exists := d.services[req.Repo]
	d.mu.Unlock()
	if exists {
		d.stop(req.Repo)
	}
	return &wire.DevStopResponse{Repo: req.Repo, Stopped: exists}, nil
}

func (d *Dev) stop(repo string) {
	d.mu.Lock()
	e, ok := d.services[repo]
	if ok {
		delete(d.services, repo)
	}
	d.mu.Unlock()
	if ok && e.cancel != nil {
		e.cancel()
		if e.cmd != nil && e.cmd.Process != nil {
			_ = e.cmd.Process.Kill()
		}
	}
}

// Status lists all running dev servers.
func (d *Dev) Status(_ context.Context) (*wire.DevStatusResponse, error) { //nolint:unparam // ctx reserved for future use
	d.mu.Lock()
	defer d.mu.Unlock()
	svcs := make([]wire.DevServiceInfo, 0, len(d.services))
	for _, e := range d.services {
		svcs = append(svcs, wire.DevServiceInfo{
			Repo: e.repo, Dir: e.dir, Command: e.command,
			Port: e.port, PID: e.pid, Status: e.status,
			URL:       fmt.Sprintf("http://localhost:%d", e.port),
			StartedAt: e.startedAt.Format(time.RFC3339),
		})
	}
	return &wire.DevStatusResponse{Services: svcs}, nil
}

// Log returns recent log lines for a dev server.
func (d *Dev) Log(_ context.Context, req wire.DevLogRequest) (*wire.DevLogResponse, error) {
	if err := scanInput(d.ice, req); err != nil {
		return nil, err
	}

	d.mu.Lock()
	e, ok := d.services[req.Repo]
	d.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("no dev server running for %s", req.Repo)
	}
	n := req.Lines
	if n <= 0 {
		n = 50
	}
	return &wire.DevLogResponse{Repo: req.Repo, Lines: e.logBuf.lines(n)}, nil
}

// -- internal helpers --

func isPortOpen(port int) bool {
	d := net.Dialer{Timeout: 200 * time.Millisecond}
	conn, err := d.DialContext(context.Background(), "tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func waitForPort(port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if isPortOpen(port) {
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("port %d not ready after %s", port, timeout)
}

func pipeToLog(r io.Reader, buf *logBuffer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		buf.write(scanner.Text())
	}
}

type logBuffer struct {
	data     []string
	mu       sync.Mutex
	capacity int
}

func newLogBuffer(capacity int) *logBuffer {
	return &logBuffer{data: make([]string, 0, capacity), capacity: capacity}
}

func (b *logBuffer) write(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.data) >= b.capacity {
		b.data = b.data[1:]
	}
	b.data = append(b.data, line)
}

func (b *logBuffer) lines(n int) []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if n <= 0 || n > len(b.data) {
		n = len(b.data)
	}
	start := len(b.data) - n
	out := make([]string, n)
	copy(out, b.data[start:])
	return out
}
