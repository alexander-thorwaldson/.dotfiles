// Package models defines the data structures returned by tool commands.
package models

// DevServiceInfo describes a running dev server.
type DevServiceInfo struct {
	Repo      string `json:"repo"`
	Dir       string `json:"dir"`
	Command   string `json:"command"`
	Status    string `json:"status"`
	URL       string `json:"url"`
	StartedAt string `json:"startedAt"`
	Port      int    `json:"port"`
	PID       int    `json:"pid"`
}

// DevStartResult is returned when a dev server is started.
type DevStartResult struct {
	Repo   string `json:"repo"`
	Status string `json:"status"`
	URL    string `json:"url"`
	Port   int    `json:"port"`
}

// DevStatusResult wraps a list of running dev services.
type DevStatusResult struct {
	Services []DevServiceInfo `json:"services"`
}

// DevStopResult is returned when a dev server is stopped.
type DevStopResult struct {
	Repo    string `json:"repo"`
	Stopped bool   `json:"stopped"`
}

// DevLogResult holds recent log output from a dev server.
type DevLogResult struct {
	Repo  string   `json:"repo"`
	Lines []string `json:"lines"`
}

// CmdRunResult holds output from a one-shot command (test, lint, etc.).
type CmdRunResult struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exitCode"`
}
