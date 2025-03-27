//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/stretchr/testify/require"
)

// TimeArgs defines the arguments for the time tool
type TimeArgs struct {
	Format string `json:"format" jsonschema:"description=The time format to use"`
}

// TestExampleMCPServerClientWithStdio tests the example MCP server and client using stdio transport
func TestExampleMCPServerClientWithStdio(t *testing.T) {
	// Skip this test if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=1 to run.")
	}

	// Start the server process
	serverCmd := exec.Command("go", "run", "../examples/server/main.go")
	serverStdin, err := serverCmd.StdinPipe()
	require.NoError(t, err)
	serverStdout, err := serverCmd.StdoutPipe()
	require.NoError(t, err)

	err = serverCmd.Start()
	require.NoError(t, err)
	defer serverCmd.Process.Kill()

	// Create a stdio transport for the client that connects to the server's pipes
	transport := stdio.NewStdioServerTransportWithIO(serverStdout, serverStdin)

	// Create a new client with the transport
	client := mcp_golang.NewClient(transport)

	// Initialize the client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Initializing MCP client...")
	resp, err := client.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	fmt.Printf("Initialized client: %+v\n", resp)

	// List available tools
	tools, err := client.ListTools(ctx, nil)
	if err != nil {
		fmt.Printf("Error listing tools: %v\n", err)
		t.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Println("Available Tools:")
	for _, tool := range tools.Tools {
		desc := ""
		if tool.Description != nil {
			desc = *tool.Description
		}
		fmt.Printf("Tool: %s. Description: %s\n", tool.Name, desc)
	}

	// Call the time tool with different formats
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"Mon, 02 Jan 2006",
	}

	for _, format := range formats {
		args := map[string]interface{}{
			"format": format,
		}

		response, err := client.CallTool(ctx, "time", args)
		if err != nil {
			fmt.Printf("Error calling time tool: %v\n", err)
			t.Fatalf("Failed to call time tool: %v", err)
		}

		if len(response.Content) > 0 && response.Content[0].TextContent != nil {
			fmt.Printf("Time in format %q: %s\n", format, response.Content[0].TextContent.Text)
		}
	}
}
