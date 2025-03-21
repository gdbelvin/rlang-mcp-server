package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"r-server/internal/mcp"
	"github.com/metoro-io/mcp-golang/transport/http"
)

func main() {
	// Parse command-line flags
	port := flag.Int("port", 22011, "Port number for the MCP server")
	testTool := flag.String("test-tool", "", "Path to a JSON file containing a tool request for testing")
	flag.Parse()

	// If test-tool flag is provided, run the tool test
	if *testTool != "" {
		if err := mcp.TestTool(*testTool); err != nil {
			fmt.Printf("Error testing tool: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Create output directories if they don't exist
	outputDir := filepath.Join("output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Create HTTP transport
	httpTransport := http.NewHTTPTransport("/mcp").WithAddr(fmt.Sprintf(":%d", *port))

	// Create and configure the server
	server, err := mcp.NewMCPServer(httpTransport)
	if err != nil {
		fmt.Printf("Error creating MCP server: %v\n", err)
		os.Exit(1)
	}

	// Start the server
	fmt.Printf("Starting MCP server on port %d\n", *port)
	if err := server.Serve(); err != nil {
		fmt.Printf("Error starting MCP server: %v\n", err)
		os.Exit(1)
	}
}
