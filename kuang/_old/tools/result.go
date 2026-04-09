package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Result returns a successful tool result with text content and typed output.
func Result[Out any](text string, out Out) (*mcp.CallToolResult, Out, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, out, nil
}

// ErrResult returns a failed tool result with the zero value of Out.
func ErrResult[Out any](err error) (*mcp.CallToolResult, Out, error) {
	r := &mcp.CallToolResult{}
	r.SetError(err)
	var zero Out
	return r, zero, nil
}
