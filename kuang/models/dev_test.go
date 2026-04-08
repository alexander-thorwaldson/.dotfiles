package models

import (
	"encoding/json"
	"testing"
)

func TestDevServiceInfo_Unmarshal(t *testing.T) {
	raw := `{"repo":"app","dir":"/code/app","command":"pnpm dev","port":3000,"pid":1234,"status":"ready","url":"http://localhost:3000","startedAt":"2026-04-08T10:00:00Z"}`
	var info DevServiceInfo
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if info.Repo != "app" {
		t.Errorf("expected repo app, got %s", info.Repo)
	}
	if info.Port != 3000 {
		t.Errorf("expected port 3000, got %d", info.Port)
	}
	if info.Status != "ready" {
		t.Errorf("expected status ready, got %s", info.Status)
	}
}

func TestDevStartResult_Unmarshal(t *testing.T) {
	raw := `{"repo":"app","port":3000,"status":"ready","url":"http://localhost:3000"}`
	var result DevStartResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if result.URL != "http://localhost:3000" {
		t.Errorf("expected URL, got %q", result.URL)
	}
}

func TestDevStopResult_Marshal(t *testing.T) {
	result := DevStopResult{Repo: "app", Stopped: true}
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var rt DevStopResult
	if err := json.Unmarshal(data, &rt); err != nil {
		t.Fatalf("roundtrip error: %v", err)
	}
	if !rt.Stopped {
		t.Error("expected Stopped true after roundtrip")
	}
}

func TestDevLogResult_Unmarshal(t *testing.T) {
	raw := `{"repo":"app","lines":["line1","line2"]}`
	var result DevLogResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(result.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(result.Lines))
	}
}

func TestCmdRunResult_Unmarshal(t *testing.T) {
	raw := `{"output":"PASS\n","exitCode":0}`
	var result CmdRunResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
}

func TestCmdRunResult_NonZero(t *testing.T) {
	raw := `{"output":"FAIL\n","exitCode":1}`
	var result CmdRunResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
}
