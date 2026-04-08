// Package handlers provides MCP tool handler wrappers including prompt
// injection filtering via the ice classifier.
package handlers

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
)

// ServiceEntry tracks a running dev server process.
type ServiceEntry struct {
	StartedAt time.Time
	Cmd       *exec.Cmd
	LogBuf    *LogBuffer
	Cancel    context.CancelFunc
	Repo      string
	Dir       string
	Command   string
	Status    string
	Port      int
	PID       int
}

// Registry tracks running dev server processes keyed by repo name.
type Registry struct {
	services map[string]*ServiceEntry
	mu       sync.Mutex
}

// NewRegistry creates an empty service registry.
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*ServiceEntry),
	}
}

// Get returns the service entry for a repo, if any.
func (r *Registry) Get(repo string) (*ServiceEntry, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.services[repo]
	return e, ok
}

// All returns all service entries.
func (r *Registry) All() []*ServiceEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	entries := make([]*ServiceEntry, 0, len(r.services))
	for _, e := range r.services {
		entries = append(entries, e)
	}
	return entries
}

// Start launches a dev server for a repo. It is idempotent — if a server
// is already running on the specified port, it returns the existing entry.
// It blocks until the port is accepting connections or the timeout expires.
func (r *Registry) Start(repo, dir, command string, port int, timeout time.Duration) (*ServiceEntry, error) {
	r.mu.Lock()

	// Idempotent: already tracked and port is open.
	if e, ok := r.services[repo]; ok {
		r.mu.Unlock()
		if isPortOpen(e.Port) {
			return e, nil
		}
		// Process died — clean up and re-start.
		r.Stop(repo)
		r.mu.Lock()
	}

	// Check if something else is already on the port.
	if isPortOpen(port) {
		r.mu.Unlock()
		return nil, fmt.Errorf("port %d is already in use", port)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "sh", "-c", command) // #nosec G204 -- intentional subprocess //nolint:gosec
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))

	logBuf := NewLogBuffer(1000)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		r.mu.Unlock()
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		r.mu.Unlock()
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		r.mu.Unlock()
		return nil, fmt.Errorf("starting process: %w", err)
	}

	entry := &ServiceEntry{
		Repo:      repo,
		Dir:       dir,
		Command:   command,
		Port:      port,
		PID:       cmd.Process.Pid,
		Cmd:       cmd,
		Cancel:    cancel,
		LogBuf:    logBuf,
		Status:    "starting",
		StartedAt: time.Now(),
	}
	r.services[repo] = entry
	r.mu.Unlock()

	// Stream stdout/stderr into the log buffer.
	go pipeToLog(stdout, logBuf)
	go pipeToLog(stderr, logBuf)

	// Wait for process exit in background to update status.
	go func() {
		_ = cmd.Wait()
		r.mu.Lock()
		if cur, ok := r.services[repo]; ok && cur == entry {
			cur.Status = "stopped"
		}
		r.mu.Unlock()
	}()

	// Block until port is ready or timeout.
	if err := waitForPort(port, timeout); err != nil {
		r.mu.Lock()
		entry.Status = "failed"
		r.mu.Unlock()
		return entry, fmt.Errorf("dev server did not become ready: %w", err)
	}

	r.mu.Lock()
	entry.Status = "ready"
	r.mu.Unlock()
	return entry, nil
}

// Stop terminates the dev server for a repo and removes it from the registry.
func (r *Registry) Stop(repo string) {
	r.mu.Lock()
	e, ok := r.services[repo]
	if ok {
		delete(r.services, repo)
	}
	r.mu.Unlock()

	if ok && e.Cancel != nil {
		e.Cancel()
		if e.Cmd != nil && e.Cmd.Process != nil {
			_ = e.Cmd.Process.Kill()
		}
	}
}

// isPortOpen checks if a TCP port is accepting connections on localhost.
func isPortOpen(port int) bool {
	d := net.Dialer{Timeout: 200 * time.Millisecond}
	conn, err := d.DialContext(context.Background(), "tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// waitForPort blocks until the port is accepting connections or timeout.
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

// pipeToLog reads lines from r and writes them to the log buffer.
func pipeToLog(r io.Reader, buf *LogBuffer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		buf.Write(scanner.Text())
	}
}

// LogBuffer is a thread-safe ring buffer of log lines.
type LogBuffer struct {
	lines    []string
	mu       sync.Mutex
	capacity int
}

// NewLogBuffer creates a log buffer with the given capacity.
func NewLogBuffer(capacity int) *LogBuffer {
	return &LogBuffer{
		lines:    make([]string, 0, capacity),
		capacity: capacity,
	}
}

// Write appends a line to the buffer, evicting the oldest if full.
func (b *LogBuffer) Write(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.lines) >= b.capacity {
		b.lines = b.lines[1:]
	}
	b.lines = append(b.lines, line)
}

// Lines returns the last n lines from the buffer.
func (b *LogBuffer) Lines(n int) []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if n <= 0 || n > len(b.lines) {
		n = len(b.lines)
	}
	start := len(b.lines) - n
	out := make([]string, n)
	copy(out, b.lines[start:])
	return out
}
