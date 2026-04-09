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
	raw, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshalling input for ice scan: %w", err)
	}
	injected, iceErr := ice.IsInjection(string(raw))
	if iceErr != nil {
		return fmt.Errorf("ice unavailable: %w", iceErr)
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
	if err := scanInput(ice, input); err != nil {
		return "", err
	}

	// Execute.
	output, err := execute()
	if err != nil {
		return "", err
	}

	// Scan output.
	if output != "" {
		injected, iceErr := ice.IsInjection(output)
		if iceErr != nil {
			return "", fmt.Errorf("ice unavailable: %w", iceErr)
		}
		if injected {
			return "", fmt.Errorf("response blocked: prompt injection detected in output")
		}
	}

	return output, nil
}
