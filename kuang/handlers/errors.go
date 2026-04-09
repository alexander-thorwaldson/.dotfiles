// Package handlers defines the HTTP endpoints for kuang's API.
package handlers

import "github.com/zoobz-io/rocco"

var (
	// ErrToolNotFound is returned when the requested tool does not exist.
	ErrToolNotFound = rocco.ErrNotFound.WithMessage("tool not found")

	// ErrUnauthorized is returned when the caller lacks permission for a tool.
	ErrUnauthorized = rocco.ErrForbidden.WithMessage("insufficient permissions")

	// ErrIceUnavailable is returned when the ice classifier is unreachable.
	ErrIceUnavailable = rocco.ErrServiceUnavailable.WithMessage("ice classifier unavailable")

	// ErrInjectionDetected is returned when ice detects prompt injection.
	ErrInjectionDetected = rocco.ErrForbidden.WithMessage("prompt injection detected")
)
