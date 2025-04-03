package mcp

import (
	"encoding/json"
	"fmt"
	"os"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
)

// MCPServer represents the MCP server for ggplot
type MCPServer struct {
	*mcp.Server
}

// NewMCPServer creates a new MCP server with the given transport
func NewMCPServer(transport transport.Transport) (*MCPServer, error) {
	// Create a new MCP server with name and version
	server := &MCPServer{
		Server: mcp.NewServer(transport, mcp.WithName("r-server"), mcp.WithVersion("1.1.4")),
	}

	// Register the render_ggplot tool
	if err := server.RegisterTool("render_ggplot", "Render a ggplot2 visualization", RenderGGPlot); err != nil {
		return nil, fmt.Errorf("failed to register render_ggplot tool: %w", err)
	}

	// Register the execute_r_script tool
	if err := server.RegisterTool("execute_r_script", "Execute an R script and return the result", ExecuteRScriptTool); err != nil {
		return nil, fmt.Errorf("failed to register execute_r_script tool: %w", err)
	}

	return server, nil
}

// TestTool tests a tool by reading a JSON request from a file and executing it
func TestTool(filePath string) error {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the JSON
	var request struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal(data, &request); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Execute the tool
	fmt.Fprintf(os.Stderr, "Testing tool: %s\n", request.Name)
	fmt.Fprintf(os.Stderr, "Arguments: %v\n", request.Arguments)

	// For now, just print a message
	// In a real implementation, we would execute the tool and display the result
	fmt.Fprintln(os.Stderr, "Tool test not yet implemented")

	return nil
}
