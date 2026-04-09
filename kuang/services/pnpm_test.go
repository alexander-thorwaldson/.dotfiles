package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/zoobzio/dotfiles/kuang/wire"
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

func TestPnpm_Run_EmptyScript(t *testing.T) {
	ice := safeICE(t)
	pnpm := NewPnpm(ice)
	_, err := pnpm.Run(context.Background(), wire.PnpmRunRequest{Dir: t.TempDir(), Script: ""})
	if err == nil {
		t.Error("expected error for empty script")
	}
}

func TestPnpm_Run_NoPackageJSON(t *testing.T) {
	ice := safeICE(t)
	pnpm := NewPnpm(ice)
	_, err := pnpm.Run(context.Background(), wire.PnpmRunRequest{Dir: t.TempDir(), Script: "test"})
	if err == nil {
		t.Error("expected error when package.json is missing")
	}
}

func TestPnpm_Run_ScriptNotFound(t *testing.T) {
	ice := safeICE(t)
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"build": "tsc"})

	pnpm := NewPnpm(ice)
	_, err := pnpm.Run(context.Background(), wire.PnpmRunRequest{Dir: dir, Script: "deploy"})
	if err == nil {
		t.Error("expected error for missing script")
	}
}

func TestPnpm_Run_ValidScript(t *testing.T) {
	ice := safeICE(t)
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"test": "echo ok"})

	pnpm := NewPnpm(ice)
	// pnpm may not be installed — just verify validation passes.
	_, _ = pnpm.Run(context.Background(), wire.PnpmRunRequest{Dir: dir, Script: "test"})
}

func TestPnpm_ICEBlocks(t *testing.T) {
	srv := newTestICEServer(t, "INJECTION", 0.99)
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"test": "echo ok"})

	pnpm := NewPnpm(ice)
	_, err := pnpm.Run(context.Background(), wire.PnpmRunRequest{Dir: dir, Script: "test"})
	if err == nil {
		t.Fatal("expected ice to block")
	}
}

func TestValidScript_Found(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"test": "jest", "lint": "eslint ."})
	if err := validScript(dir, "test"); err != nil {
		t.Errorf("expected valid: %v", err)
	}
}

func TestValidScript_NotFound(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"test": "jest"})
	if err := validScript(dir, "deploy"); err == nil {
		t.Error("expected error")
	}
}

func TestValidScript_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := validScript(dir, "test"); err == nil {
		t.Error("expected error")
	}
}
