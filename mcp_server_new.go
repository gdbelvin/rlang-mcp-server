package main

import (
	"fmt"
	"os"
	"path/filepath"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
	"github.com/metoro-io/mcp-golang/transport/http"
)

// CreateMCPServerWithTransport creates a new MCP server with the given transport and registers all tools and resources
func CreateMCPServerWithTransport(transport transport.Transport) (*mcp.Server, error) {
	// Create a new MCP server
	server := mcp.NewServer(transport)

	// Register the create_rmd tool
	err := server.RegisterTool("create_rmd", "Create a new R Markdown file", func(args RMarkdownCreateArgs) (*mcp.ToolResponse, error) {
		// Ensure filename has .Rmd extension
		fullFilename := EnsureRmdExtension(args.Filename)
		filePath := filepath.Join(RMD_DIR, fullFilename)

		// Create YAML front matter if not present
		finalContent := CreateRmdFrontMatter(args.Title, args.Content)

		// Write the file
		if err := os.WriteFile(filePath, []byte(finalContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write R Markdown file: %w", err)
		}

		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Created R Markdown file: %s", fullFilename))), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register create_rmd tool: %w", err)
	}

	// Register the render_rmd tool
	err = server.RegisterTool("render_rmd", "Render an R Markdown file", func(args RMarkdownRenderArgs) (*mcp.ToolResponse, error) {
		// Ensure filename has .Rmd extension
		fullFilename := EnsureRmdExtension(args.Filename)
		filePath := filepath.Join(RMD_DIR, fullFilename)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("R Markdown file not found: %s", fullFilename)
		}

		// Build Docker image
		if err := BuildDockerImage(); err != nil {
			return nil, fmt.Errorf("failed to build Docker image: %w", err)
		}

		// Set default format if not provided
		format := args.Format
		if format == "" {
			format = "html"
		}

		// Render the R Markdown file
		outputFile, err := RenderRMarkdown(fullFilename, format, args.UseDockerCompose)
		if err != nil {
			return nil, fmt.Errorf("failed to render R Markdown file: %w", err)
		}

		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Successfully rendered %s to %s", fullFilename, outputFile))), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register render_rmd tool: %w", err)
	}

	// Register resources for R Markdown files
	RegisterRMarkdownResources(server)
	
	// Register resources for rendered output files
	RegisterOutputResources(server)

	return server, nil
}

// RegisterRMarkdownResources registers all R Markdown files as resources
func RegisterRMarkdownResources(server *mcp.Server) {
	files, err := os.ReadDir(RMD_DIR)
	if err != nil {
		fmt.Printf("Warning: failed to read RMD_DIR: %v\n", err)
		return
	}
	
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext != ".Rmd" && ext != ".rmd" {
			continue
		}

		uri := fmt.Sprintf("rmd:///%s", filename)
		resourceName := fmt.Sprintf("rmd_%s", filename)
		description := fmt.Sprintf("R Markdown file: %s", filename)

		// Create a closure to capture the filename
		resourceFunc := func(capturedFilename, capturedUri string) func() (*mcp.ResourceResponse, error) {
			return func() (*mcp.ResourceResponse, error) {
				filePath := filepath.Join(RMD_DIR, capturedFilename)
				content, err := os.ReadFile(filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to read R Markdown file: %w", err)
				}
				return mcp.NewResourceResponse(mcp.NewTextEmbeddedResource(capturedUri, string(content), "text/markdown")), nil
			}
		}(filename, uri)

		err = server.RegisterResource(uri, resourceName, description, "text/markdown", resourceFunc)
		if err != nil {
			fmt.Printf("Warning: failed to register resource %s: %v\n", uri, err)
		}
	}
}

// RegisterOutputResources registers all rendered output files as resources
func RegisterOutputResources(server *mcp.Server) {
	outputFiles, err := os.ReadDir(OUTPUT_DIR)
	if err != nil {
		fmt.Printf("Warning: failed to read OUTPUT_DIR: %v\n", err)
		return
	}
	
	for _, file := range outputFiles {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext != ".html" && ext != ".pdf" && ext != ".docx" {
			continue
		}

		uri := fmt.Sprintf("rmd-output:///%s", filename)
		resourceName := fmt.Sprintf("output_%s", filename)
		description := fmt.Sprintf("Rendered output: %s", filename)
		mimeType := GetMimeType(filename)

		// Create a closure to capture the filename
		resourceFunc := func(capturedFilename, capturedUri, capturedMimeType string) func() (*mcp.ResourceResponse, error) {
			return func() (*mcp.ResourceResponse, error) {
				filePath := filepath.Join(OUTPUT_DIR, capturedFilename)
				content, err := os.ReadFile(filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to read rendered output file: %w", err)
				}
				return mcp.NewResourceResponse(mcp.NewTextEmbeddedResource(capturedUri, string(content), capturedMimeType)), nil
			}
		}(filename, uri, mimeType)

		err = server.RegisterResource(uri, resourceName, description, mimeType, resourceFunc)
		if err != nil {
			fmt.Printf("Warning: failed to register resource %s: %v\n", uri, err)
		}
	}
}

// StartMCPServerWithHttp starts the MCP server with HTTP transport using the metoro-io/mcp-golang library
func StartMCPServerWithHttp(httpAddr string) error {
	// Create an HTTP transport
	transport := http.NewHTTPTransport("/mcp")
	transport.WithAddr(httpAddr)

	// Create and configure the server
	server, err := CreateMCPServerWithTransport(transport)
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	// Start the server
	fmt.Printf("Starting MCP server with HTTP transport on %s\n", httpAddr)
	return server.Serve()
}
