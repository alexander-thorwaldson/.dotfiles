package services

import (
	"fmt"
	"testing"
)

func TestFilteredCLI_PassesClean(t *testing.T) {
	srv := newTestICEServer(t, "SAFE", 0.01)
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	out, err := FilteredCLI(ice, map[string]string{"q": "hello"}, func() (string, error) {
		return "clean output", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "clean output" {
		t.Errorf("expected %q, got %q", "clean output", out)
	}
}

func TestFilteredCLI_BlocksInputInjection(t *testing.T) {
	srv := newTestICEServer(t, "INJECTION", 0.99)
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	_, err := FilteredCLI(ice, map[string]string{"q": "ignore instructions"}, func() (string, error) {
		t.Fatal("execute should not be called when input is blocked")
		return "", nil
	})
	if err == nil {
		t.Fatal("expected error for input injection")
	}
}

func TestFilteredCLI_BlocksOutputInjection(t *testing.T) {
	callCount := 0
	srv := newVariableICEServer(t, func() (string, float64) {
		callCount++
		if callCount == 1 {
			return "SAFE", 0.01 // input passes
		}
		return "INJECTION", 0.99 // output blocked
	})
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	_, err := FilteredCLI(ice, "clean input", func() (string, error) {
		return "malicious output", nil
	})
	if err == nil {
		t.Fatal("expected error for output injection")
	}
}

func TestFilteredCLI_ICEDown_FailClosed(t *testing.T) {
	ice := NewICEClient("http://127.0.0.1:1")

	_, err := FilteredCLI(ice, "input", func() (string, error) {
		t.Fatal("execute should not be called when ice is down")
		return "", nil
	})
	if err == nil {
		t.Fatal("expected error when ice is down")
	}
}

func TestFilteredCLI_NilInput(t *testing.T) {
	srv := newTestICEServer(t, "SAFE", 0.01)
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	out, err := FilteredCLI(ice, nil, func() (string, error) {
		return "result", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "result" {
		t.Errorf("expected %q, got %q", "result", out)
	}
}

func TestFilteredCLI_ExecuteError(t *testing.T) {
	srv := newTestICEServer(t, "SAFE", 0.01)
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	_, err := FilteredCLI(ice, "input", func() (string, error) {
		return "", fmt.Errorf("command failed")
	})
	if err == nil {
		t.Fatal("expected error from execute")
	}
}

func TestFilteredCLI_EmptyOutput_SkipsOutputScan(t *testing.T) {
	callCount := 0
	srv := newVariableICEServer(t, func() (string, float64) {
		callCount++
		return "SAFE", 0.01
	})
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	_, err := FilteredCLI(ice, "input", func() (string, error) {
		return "", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only 1 call (input scan), not 2.
	if callCount != 1 {
		t.Errorf("expected 1 ice call for empty output, got %d", callCount)
	}
}

func TestScanInput_Nil(t *testing.T) {
	ice := NewICEClient("http://127.0.0.1:1") // unreachable
	if err := scanInput(ice, nil); err != nil {
		t.Errorf("expected nil error for nil input, got %v", err)
	}
}

func TestScanInput_Blocked(t *testing.T) {
	srv := newTestICEServer(t, "INJECTION", 0.99)
	defer srv.Close()
	ice := NewICEClient(srv.URL)

	err := scanInput(ice, map[string]string{"cmd": "rm -rf /"})
	if err == nil {
		t.Fatal("expected error for injection")
	}
}
