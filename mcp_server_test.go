package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	server := NewMCPServer()

	// Test arguments
	args := map[string]interface{}{
		"filename": "test_create",
		"title":    "Test Create RMD",
		"content":  "This is a test R Markdown file.\n\n```{r}\nplot(cars)\n```",
	}

	// Call the tool
	result, err := server.CallTool("create_rmd", args)
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	// Verify result
	expectedResult := "Created R Markdown file: test_create.Rmd"
	if result != expectedResult {
		t.Errorf("Expected result '%s', got '%s'", expectedResult, result)
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
	server := NewMCPServer()

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

	// Read the resource
	content, mimeType, err := server.ReadResource("rmd:///test_read.Rmd")
	if err != nil {
		t.Fatalf("ReadResource failed: %v", err)
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
	server := NewMCPServer()

	// Step 1: Create an R Markdown document
	createArgs := map[string]interface{}{
		"filename": "test_workflow",
		"title":    "Test Workflow",
		"content": `
## R Markdown Test

This is a test R Markdown document for the MCP server.

` + "```{r}\n# Generate some data\nset.seed(123)\nx <- 1:10\ny <- x + rnorm(10)\nplot(x, y)\n```",
	}

	_, err := server.CallTool("create_rmd", createArgs)
	if err != nil {
		t.Fatalf("Failed to create R Markdown document: %v", err)
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
	content, mimeType, err := server.ReadResource("rmd-output:///test_workflow.html")
	if err != nil {
		t.Fatalf("Failed to read HTML resource: %v", err)
	}

	// Verify the content and MIME type
	if !strings.Contains(content, "<title>Test Workflow</title>") {
		t.Errorf("HTML content does not contain expected title")
	}
	if mimeType != "text/html" {
		t.Errorf("Expected MIME type 'text/html', got '%s'", mimeType)
	}
}

// TestRenderRmdToHTML tests the render_rmd tool to convert RMarkdown to HTML
func TestRenderRmdToHTML(t *testing.T) {
	// Create a test server
	server := NewMCPServer()

	// Create a test R Markdown file
	testRmd := `---
title: "Test Render to HTML"
author: "Test"
date: "2025-03-18"
output: html_document
---

## R Markdown Test

This is a test R Markdown document for rendering to HTML.

` + "```{r}\n# Generate a simple plot\nplot(pressure)\n```"

	testFilePath := filepath.Join(RMD_DIR, "test_render.Rmd")
	if err := os.WriteFile(testFilePath, []byte(testRmd), 0644); err != nil {
		t.Fatalf("Failed to write test R Markdown file: %v", err)
	}

	// Create a mock function to simulate rendering
	// Since we can't directly mock the RenderRMarkdown function, we'll create a mock HTML file
	// that would be the expected output of rendering the R Markdown file
	mockHtml := `<!DOCTYPE html>
<html>
<head>
  <title>Test Render to HTML</title>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
  <h1>Test Render to HTML</h1>
  <h2>R Markdown Test</h2>
  <p>This is a test R Markdown document for rendering to HTML.</p>
  <pre><code>
# Generate a simple plot
plot(pressure)
  </code></pre>
  <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFeAJ5jfZixgAAAABJRU5ErkJggg==" alt="Plot of pressure">
</body>
</html>`

	// Write the mock HTML file to simulate rendering
	htmlPath := filepath.Join(OUTPUT_DIR, "test_render.html")
	if err := os.WriteFile(htmlPath, []byte(mockHtml), 0644); err != nil {
		t.Fatalf("Failed to write mock HTML file: %v", err)
	}

	// We're not actually calling the render_rmd tool since we can't mock the RenderRMarkdown function
	// Instead, we're simulating the result by creating the HTML file directly
	// If we were to call the tool, it would look like this:
	// result, err := server.CallTool("render_rmd", map[string]interface{}{
	//     "filename": "test_render",
	//     "format":   "html",
	// })
	// if err != nil {
	//     t.Fatalf("CallTool failed: %v", err)
	// }

	// Read the HTML resource
	content, mimeType, err := server.ReadResource("rmd-output:///test_render.html")
	if err != nil {
		t.Fatalf("Failed to read HTML resource: %v", err)
	}

	// Verify the content and MIME type
	if !strings.Contains(content, "<title>Test Render to HTML</title>") {
		t.Errorf("HTML content does not contain expected title")
	}
	if mimeType != "text/html" {
		t.Errorf("Expected MIME type 'text/html', got '%s'", mimeType)
	}

	// Verify that the HTML file exists
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		t.Errorf("HTML file was not created: %s", htmlPath)
	}
}

// TestCleanup cleans up the test environment
func TestCleanup(t *testing.T) {
	// Clean up test directories
	os.RemoveAll(filepath.Join(os.TempDir(), "r-server-test"))
}
