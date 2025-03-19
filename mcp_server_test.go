package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	mcp "github.com/metoro-io/mcp-golang"
	mcphttp "github.com/metoro-io/mcp-golang/transport/http"
)

// NewTestMCPServer creates a new MCP server with HTTP transport for testing
func NewTestMCPServer() (*MCPServer, string, error) {
	// Let the OS choose an available port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, "", fmt.Errorf("failed to find available port: %w", err)
	}
	
	// Get the chosen port
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	
	// Create the server address
	addr := fmt.Sprintf("localhost:%d", port)
	
	// Create an HTTP transport for testing
	serverTransport := mcphttp.NewHTTPTransport("/mcp")
	serverTransport.WithAddr(addr)
	
	// Create a new MCP server with the transport
	server, err := NewMCPServerWithTransport(serverTransport)
	if err != nil {
		return nil, "", err
	}
	
	// Return the server and the full URL for the client
	return server, addr, nil
}

// TestSetup ensures the test environment is properly set up
func TestSetup(t *testing.T) {
	// Set up test directories
	testRmdDir := filepath.Join(os.TempDir(), "r-server-test", "rmd")
	testOutputDir := filepath.Join(testRmdDir, "output")

	// Clean up any previous test directories
	os.RemoveAll(filepath.Join(os.TempDir(), "r-server-test"))

	// Create test directories
	if err := os.MkdirAll(testRmdDir, 0755); err != nil {
		t.Fatalf("Failed to create test RMD directory: %v", err)
	}
	if err := os.MkdirAll(testOutputDir, 0755); err != nil {
		t.Fatalf("Failed to create test output directory: %v", err)
	}

	// Set global variables for testing
	RMD_DIR = testRmdDir
	OUTPUT_DIR = testOutputDir

	t.Logf("Test environment set up with RMD_DIR: %s, OUTPUT_DIR: %s", RMD_DIR, OUTPUT_DIR)
}

// TestCreateRmdFile tests creating an R Markdown file
func TestCreateRmdFile(t *testing.T) {
	// Create a test R Markdown file
	testRmd := `---
title: "Test R Markdown"
author: "Test"
date: "2025-03-18"
output: html_document
---

## Test

This is a test R Markdown file.

` + "```{r}\nplot(cars)\n```"

	testFilePath := filepath.Join(RMD_DIR, "test.Rmd")
	if err := os.WriteFile(testFilePath, []byte(testRmd), 0644); err != nil {
		t.Fatalf("Failed to write test R Markdown file: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Errorf("File was not created: %s", testFilePath)
	}

	// Read the file content
	content, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	// Verify content
	if string(content) != testRmd {
		t.Errorf("File content does not match expected content")
	}
}

// TestCreateRmdTool tests the create_rmd tool
func TestCreateRmdTool(t *testing.T) {
	// Create a test server
	server, addr, err := NewTestMCPServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		if err := server.Serve(); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Create a client transport that connects to the server
	clientTransport := mcphttp.NewHTTPClientTransport(fmt.Sprintf("http://%s/mcp", addr))

	// Create a client with the transport
	client := mcp.NewClient(clientTransport)

	// Create a context
	ctx := context.Background()
	
	// Initialize the client
	if _, err := client.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Test arguments
	args := map[string]interface{}{
		"filename": "test_create",
		"title":    "Test Create RMD",
		"content":  "This is a test R Markdown file.\n\n```{r}\nplot(cars)\n```",
	}

	// Call the tool using the client
	response, err := client.CallTool(ctx, "create_rmd", args)
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	// Extract the result from the response
	if len(response.Content) == 0 {
		t.Fatalf("Expected content in response")
	}

	// Get the text content
	var result string
	for _, content := range response.Content {
		if content.Type == mcp.ContentTypeText && content.TextContent != nil {
			result = content.TextContent.Text
			break
		}
	}

	if result == "" {
		t.Fatalf("No text content found in response")
	}

	// Verify result
	expectedResult := "Created R Markdown file: test_create.Rmd"
	if !strings.Contains(result, expectedResult) {
		t.Errorf("Expected result to contain '%s', got '%s'", expectedResult, result)
	}

	// Verify file was created
	filePath := filepath.Join(RMD_DIR, "test_create.Rmd")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("File was not created: %s", filePath)
	}

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	// Verify content has YAML front matter
	if !strings.Contains(string(content), "title: \"Test Create RMD\"") {
		t.Errorf("File content does not contain expected title: %s", string(content))
	}
}

// TestReadResource tests the ReadResource method
func TestReadResource(t *testing.T) {
	// Create a test server
	server, addr, err := NewTestMCPServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		if err := server.Serve(); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Create a test R Markdown file
	testRmd := `---
title: "Test R Markdown"
author: "Test"
date: "2025-03-18"
output: html_document
---

## Test

This is a test R Markdown file.

` + "```{r}\nplot(cars)\n```"

	testFilePath := filepath.Join(RMD_DIR, "test_read.Rmd")
	if err := os.WriteFile(testFilePath, []byte(testRmd), 0644); err != nil {
		t.Fatalf("Failed to write test R Markdown file: %v", err)
	}

	// Create a client transport that connects to the server
	clientTransport := mcphttp.NewHTTPClientTransport(fmt.Sprintf("http://%s/mcp", addr))

	// Create a client with the transport
	client := mcp.NewClient(clientTransport)

	// Create a context
	ctx := context.Background()
	
	// Initialize the client
	if _, err := client.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Read the resource
	response, err := client.ReadResource(ctx, "rmd:///test_read.Rmd")
	if err != nil {
		t.Fatalf("ReadResource failed: %v", err)
	}

	// Extract the content from the response
	if len(response.Contents) == 0 {
		t.Fatalf("Expected contents in response")
	}

	// Get the text content
	var content string
	var mimeType string
	for _, res := range response.Contents {
		if res.TextResourceContents != nil {
			content = res.TextResourceContents.Text
			if res.TextResourceContents.MimeType != nil {
				mimeType = *res.TextResourceContents.MimeType
			}
			break
		}
	}

	if content == "" {
		t.Fatalf("No text content found in response")
	}

	// Verify content and MIME type
	if content != testRmd {
		t.Errorf("Expected content to match test R Markdown, got '%s'", content)
	}
	if mimeType != "text/markdown" {
		t.Errorf("Expected MIME type 'text/markdown', got '%s'", mimeType)
	}
}

// TestRMarkdownToHTML tests the workflow of creating an RMarkdown document and simulating rendering to HTML
func TestRMarkdownToHTML(t *testing.T) {
	// Create a test server
	server, addr, err := NewTestMCPServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		if err := server.Serve(); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Create a client transport that connects to the server
	clientTransport := mcphttp.NewHTTPClientTransport(fmt.Sprintf("http://%s/mcp", addr))

	// Create a client with the transport
	client := mcp.NewClient(clientTransport)

	// Create a context
	ctx := context.Background()
	
	// Initialize the client
	if _, err := client.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Step 1: Create an R Markdown document
	createArgs := map[string]interface{}{
		"filename": "test_workflow",
		"title":    "Test Workflow",
		"content": `
## R Markdown Test

This is a test R Markdown document for the MCP server.

` + "```{r}\n# Generate some data\nset.seed(123)\nx <- 1:10\ny <- x + rnorm(10)\nplot(x, y)\n```",
	}

	response, err := client.CallTool(ctx, "create_rmd", createArgs)
	if err != nil {
		t.Fatalf("Failed to create R Markdown document: %v", err)
	}

	// Verify we got a response
	if len(response.Content) == 0 {
		t.Fatalf("Expected content in response")
	}

	// Verify the file was created
	rmdPath := filepath.Join(RMD_DIR, "test_workflow.Rmd")
	if _, err := os.Stat(rmdPath); os.IsNotExist(err) {
		t.Fatalf("R Markdown file was not created: %s", rmdPath)
	}

	// Step 2: Simulate rendering by creating a mock HTML file
	mockHtml := `<!DOCTYPE html>
<html>
<head>
  <title>Test Workflow</title>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <style>
    body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
    h1, h2 { color: #2c3e50; }
    pre { background-color: #f8f8f8; padding: 10px; border-radius: 5px; }
    img { max-width: 100%; height: auto; }
  </style>
</head>
<body>
  <h1>Test Workflow</h1>
  <h2>R Markdown Test</h2>
  <p>This is a test R Markdown document for the MCP server.</p>
  <pre><code>
# Generate some data
set.seed(123)
x <- 1:10
y <- x + rnorm(10)
plot(x, y)
  </code></pre>
  <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFeAJ5jfZixgAAAABJRU5ErkJggg==" alt="Plot of x vs y">
</body>
</html>`

	htmlPath := filepath.Join(OUTPUT_DIR, "test_workflow.html")
	if err := os.WriteFile(htmlPath, []byte(mockHtml), 0644); err != nil {
		t.Fatalf("Failed to write mock HTML file: %v", err)
	}

	// Step 3: Read the HTML resource
	resourceResponse, err := client.ReadResource(ctx, "rmd-output:///test_workflow.html")
	if err != nil {
		t.Fatalf("Failed to read HTML resource: %v", err)
	}

	// Extract the content from the response
	if len(resourceResponse.Contents) == 0 {
		t.Fatalf("Expected contents in response")
	}

	// Get the text content
	var content string
	var mimeType string
	for _, res := range resourceResponse.Contents {
		if res.TextResourceContents != nil {
			content = res.TextResourceContents.Text
			if res.TextResourceContents.MimeType != nil {
				mimeType = *res.TextResourceContents.MimeType
			}
			break
		}
	}

	if content == "" {
		t.Fatalf("No text content found in response")
	}

	// Verify the content and MIME type
	if !strings.Contains(content, "<title>Test Workflow</title>") {
		t.Errorf("HTML content does not contain expected title")
	}
	if mimeType != "text/html" {
		t.Errorf("Expected MIME type 'text/html', got '%s'", mimeType)
	}
}

// TestCleanup cleans up the test environment
func TestCleanup(t *testing.T) {
	// Clean up test directories
	os.RemoveAll(filepath.Join(os.TempDir(), "r-server-test"))
}
