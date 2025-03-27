package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

// TimeArgs defines the arguments for the time tool
type TimeArgs struct {
	Format string `json:"format" jsonschema:"description=The time format to use"`
}

func main() {
	// Define command-line flags
	help := flag.Bool("help", false, "Show usage instructions")
	flag.Parse()

	// Show help message if requested
	if *help {
		fmt.Println("MCP Server with stdio transport")
		fmt.Println("\nUsage:")
		fmt.Println("  This server uses stdio transport to communicate with MCP clients.")
		fmt.Println("  To connect to a client using named pipes:")
		fmt.Println("\n  1. Create named pipes:")
		fmt.Println("     mkfifo server_in server_out")
		fmt.Println("\n  2. Start this server redirecting stdin/stdout to the pipes:")
		fmt.Println("     ./server < server_in > server_out &")
		fmt.Println("\n  3. Run the client with the pipes connected in reverse:")
		fmt.Println("     ./client < server_out > server_in")
		fmt.Println("\n  Note: The server reads from its stdin (server_in) and writes to its stdout (server_out).")
		os.Exit(0)
	}

	// Use stdio transport
	transport := stdio.NewStdioServerTransport()

	// Create a new server with the transport
	server := mcp_golang.NewServer(transport,
		mcp_golang.WithName("mcp-golang-stateless-stdio-example"),
		mcp_golang.WithVersion("0.0.1"))

	// Register a simple tool
	err := server.RegisterTool("time", "Returns the current time in the specified format", func(args TimeArgs) (*mcp_golang.ToolResponse, error) {
		format := args.Format
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(time.Now().Format(format))), nil
	})
	if err != nil {
		panic(err)
	}

	// Start the server
	fmt.Fprintln(os.Stderr, "Starting server")
	server.Serve()
	fmt.Fprintln(os.Stderr, "Stopping server")
}
