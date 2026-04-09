package services

import (
	"encoding/json"
	"fmt"
)

// scanInput checks the input for prompt injection. Returns an error if blocked.
func scanInput(ice *ICEClient, input any) error {
	if input == nil {
		return nil
	}
	raw, _ := json.Marshal(input)
	injected, err := ice.IsInjection(string(raw))
	if err != nil {
		return fmt.Errorf("ice unavailable: %w", err)
	}
	if injected {
		return fmt.Errorf("request blocked: prompt injection detected in input")
	}
	return nil
}

// FilteredCLI wraps a CLI execution with ice input/output scanning.
// If ice detects prompt injection, it returns an error.
// If ice is unreachable, it blocks the call (fail-closed).
func FilteredCLI(ice *ICEClient, input any, execute func() (string, error)) (string, error) {
	// Scan input.
	if input != nil {
		raw, _ := json.Marshal(input)
		injected, err := ice.IsInjection(string(raw))
		if err != nil {
			return "", fmt.Errorf("ice unavailable: %w", err)
		}
		if injected {
			return "", fmt.Errorf("request blocked: prompt injection detected in input")
		}
	}

	// Execute.
	output, err := execute()
	if err != nil {
		return "", err
	}

	// Scan output.
	if output != "" {
		injected, err := ice.IsInjection(output)
		if err != nil {
			return "", fmt.Errorf("ice unavailable: %w", err)
		}
		if injected {
			return "", fmt.Errorf("response blocked: prompt injection detected in output")
		}
	}

	return output, nil
}
