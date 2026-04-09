package wire

import (
	"encoding/json"
	"testing"
)

func TestToolListResponse_Clone(t *testing.T) {
	orig := ToolListResponse{Tools: []ToolInfo{{Name: "a", Description: "b"}}}
	clone := orig.Clone()
	clone.Tools[0].Name = "changed"
	if orig.Tools[0].Name != "a" {
		t.Error("clone mutated original")
	}
}

func TestCLIOutput_Marshal(t *testing.T) {
	out := CLIOutput{Output: `{"key":"val"}`}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var rt CLIOutput
	if err := json.Unmarshal(data, &rt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if rt.Output != out.Output {
		t.Errorf("roundtrip mismatch: %q != %q", rt.Output, out.Output)
	}
}

func TestDevStartResponse_Clone(t *testing.T) {
	orig := DevStartResponse{Repo: "app", Port: 3000, Status: "ready", URL: "http://localhost:3000"}
	clone := orig.Clone()
	if clone != orig {
		t.Error("clone mismatch")
	}
}

func TestDevStatusResponse_Clone(t *testing.T) {
	orig := DevStatusResponse{Services: []DevServiceInfo{{Repo: "a", Port: 3000}}}
	clone := orig.Clone()
	clone.Services[0].Repo = "changed"
	if orig.Services[0].Repo != "a" {
		t.Error("clone mutated original")
	}
}

func TestDevLogResponse_Clone(t *testing.T) {
	orig := DevLogResponse{Repo: "app", Lines: []string{"line1", "line2"}}
	clone := orig.Clone()
	clone.Lines[0] = "changed"
	if orig.Lines[0] != "line1" {
		t.Error("clone mutated original")
	}
}

func TestPnpmRunResponse_Clone(t *testing.T) {
	orig := PnpmRunResponse{Output: "PASS", ExitCode: 0}
	clone := orig.Clone()
	if clone != orig {
		t.Error("clone mismatch")
	}
}
