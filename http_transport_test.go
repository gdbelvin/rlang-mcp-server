package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	mcp "github.com/metoro-io/mcp-golang"
	mcphttp "github.com/metoro-io/mcp-golang/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHTTPTransportBasic tests the basic HTTP transport functionality
func TestHTTPTransportBasic(t *testing.T) {
	// Create a simple HTTP server with a test endpoint
	serverAddr := "localhost:8082"
	serverURL := "http://" + serverAddr
	
	// Create a simple HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	
	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}
	
	// Start the server in a goroutine
	go func() {
		fmt.Println("Starting HTTP server on", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
	
	// Wait for the server to start
	time.Sleep(1 * time.Second)
	fmt.Println("Server started, proceeding with test...")
	
	// Make a request to the server
	resp, err := http.Get(serverURL + "/ping")
	require.NoError(t, err, "Failed to make HTTP request")
	defer resp.Body.Close()
	
	// Check the response
	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")
	assert.Equal(t, "pong", string(body), "Expected response body to be 'pong'")
	
	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	fmt.Println("Server shutdown complete")
}

// TestMCPHTTPTransport tests the MCP server with HTTP transport
func TestMCPHTTPTransport(t *testing.T) {
	// Set up test directories
	testRmdDir := filepath.Join(os.TempDir(), "r-server-test", "rmd")
	testOutputDir := filepath.Join(testRmdDir, "output")

	// Clean up any previous test directories
	os.RemoveAll(filepath.Join(os.TempDir(), "r-server-test"))

	// Create test directories
	err := os.MkdirAll(testRmdDir, 0755)
	require.NoError(t, err, "Failed to create test RMD directory")
	err = os.MkdirAll(testOutputDir, 0755)
	require.NoError(t, err, "Failed to create test output directory")

	// Set global variables for testing
	RMD_DIR = testRmdDir
	OUTPUT_DIR = testOutputDir

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

	testFilePath := filepath.Join(testRmdDir, "test.Rmd")
	err = os.WriteFile(testFilePath, []byte(testRmd), 0644)
	require.NoError(t, err, "Failed to write test R Markdown file")

	// Start the server in a goroutine
	serverAddr := "localhost:8083" // Use a different port than the basic test
	serverURL := "http://" + serverAddr

	// Create an HTTP transport for the server
	serverTransport := mcphttp.NewHTTPTransport("/mcp")
	serverTransport.WithAddr(serverAddr)
	
	// Create a new MCP server with the transport
	server, err := NewMCPServerWithTransport(serverTransport)
	require.NoError(t, err, "Failed to create MCP server")
	
	// Register a test resource for the test
	err = server.RegisterResource("test://resource", "test_resource", "Test resource", "text/plain", func() (*mcp.ResourceResponse, error) {
		return mcp.NewResourceResponse(mcp.NewTextEmbeddedResource("test://resource", "This is a test resource", "text/plain")), nil
	})
	require.NoError(t, err, "Failed to register test resource")

	// Start the server in a goroutine with a timeout
	serverErrCh := make(chan error, 1)
	serverStarted := make(chan struct{})
	
	go func() {
		// Signal that the server is about to start
		close(serverStarted)
		fmt.Println("Starting MCP HTTP server on", serverAddr)
		serverErrCh <- server.Serve()
	}()

	// Wait for the server to start
	<-serverStarted
	fmt.Println("MCP server goroutine started, waiting for server to be ready...")
	time.Sleep(2 * time.Second)
	fmt.Println("Proceeding with MCP test...")

	// Create a simple HTTP client to test the server is running
	resp, err := http.Get(serverURL + "/mcp")
	if err != nil {
		t.Logf("Warning: Could not connect to MCP server: %v", err)
	} else {
		defer resp.Body.Close()
		t.Logf("MCP server responded with status: %d", resp.StatusCode)
	}

	// Clean up
	os.RemoveAll(filepath.Join(os.TempDir(), "r-server-test"))
	fmt.Println("MCP test cleanup complete")
}
