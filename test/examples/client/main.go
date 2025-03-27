package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	// Define command-line flags
	help := flag.Bool("help", false, "Show usage instructions")
	flag.Parse()

	// Show help message if requested
	if *help {
		fmt.Println("MCP Client with stdio transport")
		fmt.Println("\nUsage:")
		fmt.Println("  This client uses stdio transport to communicate with an MCP server.")
		fmt.Println("  To connect to a server using named pipes:")
		fmt.Println("\n  1. Create named pipes:")
		fmt.Println("     mkfifo server_in server_out")
		fmt.Println("\n  2. Start the server redirecting stdin/stdout to the pipes:")
		fmt.Println("     ./server < server_in > server_out &")
		fmt.Println("\n  3. Run this client with the pipes connected in reverse:")
		fmt.Println("     ./client < server_out > server_in")
		fmt.Println("\n  Note: The client reads from the server's output and writes to the server's input.")
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Starting client\n")
	// Create a stdio transport using stdin/stdout
	// When connected with named pipes, this will communicate with the server
	transport := stdio.NewStdioServerTransport()

	// Create a new client with the transport
	client := mcp_golang.NewClient(transport)

	// Initialize the client
	if resp, err := client.Initialize(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize client: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stderr, "Initialized client: %v\n", spew.Sdump(resp))
	}

	// List available tools
	tools, err := client.ListTools(context.Background(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list tools: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Available Tools:")
	for _, tool := range tools.Tools {
		desc := ""
		if tool.Description != nil {
			desc = *tool.Description
		}
		fmt.Fprintf(os.Stderr, "Tool: %s. Description: %s\n", tool.Name, desc)
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

		response, err := client.CallTool(context.Background(), "time", args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to call time tool: %v\n", err)
			continue
		}

		if len(response.Content) > 0 && response.Content[0].TextContent != nil {
			fmt.Fprintf(os.Stderr, "Time in format %q: %s\n", format, response.Content[0].TextContent.Text)
		}
	}
}
