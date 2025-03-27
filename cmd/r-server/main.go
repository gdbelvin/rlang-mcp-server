package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"r-server/internal/mcp"

	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	// Parse command-line flags
	testTool := flag.String("test-tool", "", "Path to a JSON file containing a tool request for testing")
	flag.Parse()

	// If test-tool flag is provided, run the tool test
	if *testTool != "" {
	if err := mcp.TestTool(*testTool); err != nil {
		fmt.Fprintf(os.Stderr, "Error testing tool: %v\n", err)
		os.Exit(1)
	}
		return
	}

	// Create output directories if they don't exist
	outputDir := filepath.Join("output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Create stdio transport
	stdioTransport := stdio.NewStdioServerTransport()

	// Create and configure the server
	server, err := mcp.NewMCPServer(stdioTransport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating MCP server: %v\n", err)
		os.Exit(1)
	}

	// Start the server
	fmt.Fprintf(os.Stderr, "Starting MCP server with stdio transport\n")
	if err := server.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting MCP server: %v\n", err)
		os.Exit(1)
	}
	// Never exit?
	done := make(chan struct{})
	<-done

	fmt.Fprintf(os.Stderr, "Exiting\n")
}
