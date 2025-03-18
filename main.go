package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Base directory for R Markdown files
var RMD_DIR string
var OUTPUT_DIR string

// BuildDockerImage builds the Docker image for R Markdown rendering
func BuildDockerImage() error {
	fmt.Println("Building Docker image for R Markdown rendering...")
	cmd := exec.Command("docker", "build", "-t", "r-server-rmd", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build Docker image: %w", err)
	}
	fmt.Println("Docker image built successfully")
	return nil
}

// RenderRMarkdown renders an R Markdown file using Docker
func RenderRMarkdown(filename, outputFormat string, useDockerCompose bool) (string, error) {
	var outputFile string
	var err error
	
	if useDockerCompose {
		fmt.Printf("Rendering %s using docker-compose\n", filename)
		outputFile, err = RenderWithDockerCompose(RMD_DIR, filename, outputFormat)
	} else {
		fmt.Printf("Rendering %s using Dockerode\n", filename)
		outputFile, err = RenderWithDockerode(RMD_DIR, filename, outputFormat)
	}
	
	if err != nil {
		return "", fmt.Errorf("failed to render R Markdown file: %w", err)
	}
	
	return outputFile, nil
}

// StartMCPServer initializes and starts the MCP server
func StartMCPServer() error {
	fmt.Println("Starting MCP server...")
	
	// Create a new MCP server instance
	server := NewMCPServer()
	
	// Note: This is a simplified implementation for demonstration purposes.
	// In a real implementation with a Go MCP SDK, we would:
	// 1. Create an MCP server instance using the MCP SDK
	// 2. Set up request handlers for resources and tools
	// 3. Connect the server to a transport (e.g., stdio)
	//
	// For example:
	// mcpServer := mcp.NewServer(mcp.ServerConfig{
	//     Name:    "R-Server",
	//     Version: "0.1.0",
	// })
	// mcpServer.SetRequestHandler(mcp.ListResourcesRequestSchema, handleListResources)
	// mcpServer.SetRequestHandler(mcp.ReadResourceRequestSchema, handleReadResource)
	// mcpServer.SetRequestHandler(mcp.ListToolsRequestSchema, handleListTools)
	// mcpServer.SetRequestHandler(mcp.CallToolRequestSchema, handleCallTool)
	// transport := mcp.NewStdioTransport()
	// mcpServer.Connect(transport)
	
	// For demonstration purposes, let's list the available resources and tools
	resources, err := server.ListResources()
	if err != nil {
		return fmt.Errorf("failed to list resources: %w", err)
	}
	
	fmt.Printf("Available resources: %d\n", len(resources))
	for _, resource := range resources {
		fmt.Printf("  - %s (%s)\n", resource["name"], resource["uri"])
	}
	
	tools := server.ListTools()
	fmt.Printf("Available tools: %d\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  - %s: %s\n", tool["name"], tool["description"])
	}
	
	fmt.Println("MCP server started")
	return nil
}

func main() {
	// Set base directories
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	
	RMD_DIR = filepath.Join(wd, "rmd")
	OUTPUT_DIR = filepath.Join(RMD_DIR, "output")

	// Ensure directories exist
	if err := EnsureDirectoriesExist(RMD_DIR, OUTPUT_DIR); err != nil {
		fmt.Printf("Error creating directories: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Using RMD_DIR: %s\n", RMD_DIR)
	fmt.Printf("Using OUTPUT_DIR: %s\n", OUTPUT_DIR)

	// Start the MCP server
	if err := StartMCPServer(); err != nil {
		fmt.Printf("Error starting MCP server: %v\n", err)
		os.Exit(1)
	}
}
