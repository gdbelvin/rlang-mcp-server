//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	servermcp "r-server/internal/mcp"
)

// setupTestServer sets up a test MCP server on a random available port
// and returns the server, client, and a cleanup function
func setupTestServer(t *testing.T) (*servermcp.MCPServer, *mcp_golang.Client, func()) {
	// Create a listener to get a random port
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Create the server address
	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Using port %d for MCP server\n", port)

	// Create the server transport
	httpTransport := http.NewHTTPTransport("/mcp").WithAddr(serverAddr)

	// Create and configure the server
	server, err := servermcp.NewMCPServer(httpTransport)
	require.NoError(t, err)

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Starting MCP server on port %d...\n", port)
		if err := server.Serve(); err != nil {
			fmt.Printf("Error starting MCP server: %v\n", err)
		}
	}()

	// Give the server some time to start
	time.Sleep(1 * time.Second)

	// Create the client transport
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("Creating MCP client with base URL: %s\n", baseURL)
	clientTransport := http.NewHTTPClientTransport("/mcp")
	clientTransport.WithBaseURL(baseURL)

	// Create the client
	client := mcp_golang.NewClient(clientTransport)

	// Create a cleanup function
	cleanup := func() {
		// Nothing to clean up for now
		// The server will be automatically closed when the test finishes
	}

	return server, client, cleanup
}

// TestMCPClientGGPlotRendering tests rendering a ggplot visualization using the MCP client
func TestMCPClientGGPlotRendering(t *testing.T) {
	// Skip this test if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=1 to run.")
	}

	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "ggplot-client-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up the test server and client
	_, client, cleanup := setupTestServer(t)
	defer cleanup()

	// Initialize the client with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	fmt.Println("Initializing MCP client...")

	// Try to initialize the client
	initResponse, err := client.Initialize(ctx)
	if err != nil {
		fmt.Printf("Error initializing client: %v\n", err)
		t.Fatalf("Failed to initialize client: %v", err)
	}

	fmt.Printf("Client initialized successfully: %+v\n", initResponse)

	// Read the test R script
	scriptPath := filepath.Join("..", "testdata", "r_scripts", "basic_plot.R")
	scriptData, err := os.ReadFile(scriptPath)
	require.NoError(t, err)

	// Create the arguments for the render_ggplot tool
	args := map[string]interface{}{
		"code":        string(scriptData),
		"output_type": "png",
		"width":       800,
		"height":      600,
		"resolution":  96,
	}

	// Create a new context with a timeout for the tool call
	toolCtx, toolCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer toolCancel()

	// Call the render_ggplot tool
	fmt.Println("Calling render_ggplot tool via MCP client...")
	response, err := client.CallTool(toolCtx, "render_ggplot", args)
	if err != nil {
		fmt.Printf("Error calling render_ggplot tool: %v\n", err)
		t.Fatalf("Failed to call render_ggplot tool: %v", err)
	}
	require.NotNil(t, response)
	require.NotNil(t, response.Content)
	require.Len(t, response.Content, 1)

	// Verify the content type is image
	assert.Equal(t, "image", string(response.Content[0].Type))

	// Verify the image content
	require.NotNil(t, response.Content[0].ImageContent)
	assert.Equal(t, "image/png", response.Content[0].ImageContent.MimeType)
	assert.NotEmpty(t, response.Content[0].ImageContent.Data)

	// Save the image to a file for inspection
	imageFile := filepath.Join(tempDir, "output.png")
	imageData, err := json.MarshalIndent(response, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(imageFile+".json", imageData, 0644)
	require.NoError(t, err)

	fmt.Println("Successfully rendered ggplot image via MCP client")
}
