// Package wire defines the request and response types for kuang's API.
package wire

// ToolInfo describes an available tool.
type ToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToolListResponse is returned by the tool discovery endpoint.
type ToolListResponse struct {
	Tools []ToolInfo `json:"tools"`
}

// Clone returns a deep copy of ToolListResponse.
func (r ToolListResponse) Clone() ToolListResponse {
	tools := make([]ToolInfo, len(r.Tools))
	copy(tools, r.Tools)
	return ToolListResponse{Tools: tools}
}

// CLIOutput is the common response for tool endpoints that shell out to a CLI.
type CLIOutput struct {
	Output string `json:"output"`
}

// Clone returns a copy of CLIOutput.
func (r CLIOutput) Clone() CLIOutput { return r }
