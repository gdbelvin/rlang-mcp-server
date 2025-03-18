package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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

// Base directory for R Markdown files
var RMD_DIR string
var OUTPUT_DIR string

// RMarkdownCreateArgs represents the arguments for creating an R Markdown file
type RMarkdownCreateArgs struct {
	Filename string `json:"filename" jsonschema:"required,description=Filename for the R Markdown file (without extension)"`
	Title    string `json:"title" jsonschema:"required,description=Title for the R Markdown document"`
	Content  string `json:"content" jsonschema:"required,description=Content of the R Markdown file"`
}

// RMarkdownRenderArgs represents the arguments for rendering an R Markdown file
type RMarkdownRenderArgs struct {
	Filename         string `json:"filename" jsonschema:"required,description=Filename of the R Markdown file to render"`
	Format           string `json:"format" jsonschema:"description=Output format (html, pdf, or word)"`
	UseDockerCompose bool   `json:"use_docker_compose" jsonschema:"description=Whether to use docker-compose (true) or Dockerode (false)"`
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

func main() {
	// Define command-line flags
	port := flag.Int("port", 22011, "Port number for the MCP server to listen on")
	flag.Parse()

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

	// Start the MCP server with HTTP transport using the metoro-io/mcp-golang library
	httpAddr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting MCP server on port %d\n", *port)

	// Start the server
	if err := StartMCPServerWithHttp(httpAddr); err != nil {
		fmt.Printf("Error starting MCP server: %v\n", err)
		os.Exit(1)
	}
}
