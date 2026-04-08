package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FilteredHandler wraps an MCP tool handler with ice input/output scanning.
// If ice detects a prompt injection in the tool input or output, the call
// is blocked and an error result is returned.
func FilteredHandler[In, Out any](ice *ICEClient, logger *slog.Logger, name string, fn mcp.ToolHandlerFor[In, Out]) mcp.ToolHandlerFor[In, Out] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input In) (*mcp.CallToolResult, Out, error) {
		var zero Out

		// Scan input arguments as JSON
		if req != nil && req.Params != nil && req.Params.Arguments != nil {
			raw, _ := req.Params.Arguments.MarshalJSON()
			if injected, err := ice.IsInjection(string(raw)); err != nil {
				logger.Warn("ice input scan failed, allowing call", "tool", name, "err", err)
			} else if injected {
				logger.Warn("ice blocked tool call: input injection detected", "tool", name)
				r := &mcp.CallToolResult{}
				r.SetError(fmt.Errorf("request blocked: prompt injection detected in input"))
				return r, zero, nil
			}
		}

		// Execute the real handler
		result, out, err := fn(ctx, req, input)
		if err != nil {
			return result, out, err
		}

		// Scan output content
		if result != nil {
			for _, c := range result.Content {
				if tc, ok := c.(*mcp.TextContent); ok && tc.Text != "" {
					if injected, scanErr := ice.IsInjection(tc.Text); scanErr != nil {
						logger.Warn("ice output scan failed, allowing response", "tool", name, "err", scanErr)
					} else if injected {
						logger.Warn("ice blocked tool response: output injection detected", "tool", name)
						r := &mcp.CallToolResult{}
						r.SetError(fmt.Errorf("response blocked: prompt injection detected in output"))
						return r, zero, nil
					}
				}
			}
		}

		return result, out, nil
	}
}
