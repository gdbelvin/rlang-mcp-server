package main

import (
	"fmt"
	"os"
	"path/filepath"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

// MCPServer represents the MCP server for R Markdown
type MCPServer struct {
	*mcp.Server
}

type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}
type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai, google, claude, etc)"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

// CreateMCPServerWithTransport creates a new MCP server with the given transport and registers all tools and resources
func NewMCPServerWithTransport(transport transport.Transport) (*MCPServer, error) {
	// Create a new MCP server
	server := &MCPServer{
		Server: mcp.NewServer(transport),
	}

	err := server.RegisterTool("hello", "Say hello to a person", func(arguments MyFunctionsArguments) (*mcp.ToolResponse, error) {
		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Hello, %server!", arguments.Submitter))), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterPrompt("promt_test", "This is a test prompt", func(arguments Content) (*mcp.PromptResponse, error) {
		return mcp.NewPromptResponse("description", mcp.NewPromptMessage(mcp.NewTextContent(fmt.Sprintf("Hello, %server!", arguments.Title)), mcp.RoleUser)), nil
	})
	if err != nil {
		panic(err)
	}

	// Register the create_rmd tool
	if err := server.RegisterTool("create_rmd", "Create a new R Markdown file", CreateMarkdownFile); err != nil {
		return nil, fmt.Errorf("failed to register create_rmd tool: %w", err)
	}

	// Register the render_rmd tool
	if err := server.RegisterTool("render_rmd", "Render an R Markdown file", RenderMarkdownFile); err != nil {
		return nil, fmt.Errorf("failed to register render_rmd tool: %w", err)
	}

	return server, nil
}

// StartMCPServerWithStdio starts the MCP server with stdio transport using the metoro-io/mcp-golang library
func StartMCPServerWithStdio() error {
	// Create a stdio transport
	transport := stdio.NewStdioServerTransport()

	// Create and configure the server
	server, err := NewMCPServerWithTransport(transport)
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	// Start the server
	fmt.Printf("Starting MCP server with stdio transport\n")
	return server.Serve()
}

func CreateMarkdownFile(args RMarkdownCreateArgs) (*mcp.ToolResponse, error) {
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
}

func RenderMarkdownFile(args RMarkdownRenderArgs) (*mcp.ToolResponse, error) {
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
}
