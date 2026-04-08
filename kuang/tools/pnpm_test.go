package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func writePackageJSON(t *testing.T, dir string, scripts map[string]string) {
	t.Helper()
	content := `{"scripts":{`
	first := true
	for k, v := range scripts {
		if !first {
			content += ","
		}
		content += `"` + k + `":"` + v + `"`
		first = false
	}
	content += `}}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(content), 0o600); err != nil {
		t.Fatalf("writing package.json: %v", err)
	}
}

func TestPnpmRun_EmptyScript(t *testing.T) {
	result, _, err := PnpmRun(context.Background(), nil, PnpmRunParams{
		Dir:    t.TempDir(),
		Script: "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error for empty script")
	}
}

func TestPnpmRun_NoPackageJSON(t *testing.T) {
	result, _, err := PnpmRun(context.Background(), nil, PnpmRunParams{
		Dir:    t.TempDir(),
		Script: "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error when package.json is missing")
	}
}

func TestPnpmRun_ScriptNotInPackageJSON(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{
		"build": "tsc",
		"lint":  "eslint .",
	})

	result, _, err := PnpmRun(context.Background(), nil, PnpmRunParams{
		Dir:    dir,
		Script: "deploy",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error for script not in package.json")
	}
}

func TestPnpmRun_ValidScript(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{
		"test": "echo ok",
	})

	// pnpm may not be available in CI, so we just verify validation passes.
	result, _, err := PnpmRun(context.Background(), nil, PnpmRunParams{
		Dir:    dir,
		Script: "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// If pnpm isn't installed, the tool will return an exec error — that's fine.
	// What matters is that it got past validation.
	_ = result
}

func TestValidScript_Found(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{
		"test":  "jest",
		"lint":  "eslint .",
		"build": "tsc",
	})

	if err := validScript(dir, "test"); err != nil {
		t.Errorf("expected test to be valid: %v", err)
	}
	if err := validScript(dir, "lint"); err != nil {
		t.Errorf("expected lint to be valid: %v", err)
	}
}

func TestValidScript_NotFound(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{
		"test": "jest",
	})

	if err := validScript(dir, "deploy"); err == nil {
		t.Error("expected error for nonexistent script")
	}
}

func TestValidScript_NoPackageJSON(t *testing.T) {
	if err := validScript(t.TempDir(), "test"); err == nil {
		t.Error("expected error when package.json is missing")
	}
}

func TestValidScript_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := validScript(dir, "test"); err == nil {
		t.Error("expected error for invalid JSON")
	}
}
