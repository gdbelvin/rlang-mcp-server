package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// MCPServer represents the MCP server for R Markdown
type MCPServer struct {
	// Note: This is a simplified implementation for demonstration purposes.
	// In a real implementation with a Go MCP SDK, this would include the MCP SDK server instance
	// For example:
	// mcpServer *mcp.Server
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() *MCPServer {
	return &MCPServer{}
}

// ListResources lists available R Markdown files as resources
func (s *MCPServer) ListResources() ([]map[string]string, error) {
	var resources []map[string]string

	// Get R Markdown files
	rmdFiles, err := GetRMarkdownFiles(RMD_DIR)
	if err != nil {
		return nil, fmt.Errorf("failed to get R Markdown files: %w", err)
	}

	// Add R Markdown files as resources
	for _, file := range rmdFiles {
		resources = append(resources, map[string]string{
			"uri":         fmt.Sprintf("rmd:///%s", file.Filename),
			"mimeType":    "text/markdown",
			"name":        file.Title,
			"description": fmt.Sprintf("R Markdown file: %s", file.Title),
		})
	}

	// Add rendered output files as resources
	files, err := os.ReadDir(OUTPUT_DIR)
	if err != nil {
		return nil, fmt.Errorf("failed to read output directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext != ".html" && ext != ".pdf" && ext != ".docx" {
			continue
		}

		resources = append(resources, map[string]string{
			"uri":         fmt.Sprintf("rmd-output:///%s", filename),
			"mimeType":    GetMimeType(filename),
			"name":        fmt.Sprintf("Rendered: %s", filename),
			"description": fmt.Sprintf("Rendered output: %s", filename),
		})
	}

	return resources, nil
}

// ReadResource reads an R Markdown file or rendered output
func (s *MCPServer) ReadResource(uri string) (string, string, error) {
	// Parse URI
	scheme := ""
	path := ""

	// Simple URI parsing (in a real implementation, use url.Parse)
	if len(uri) > 6 && uri[:6] == "rmd://" {
		scheme = "rmd"
		path = uri[6:]
	} else if len(uri) > 13 && uri[:13] == "rmd-output://" {
		scheme = "rmd-output"
		path = uri[13:]
	} else {
		return "", "", fmt.Errorf("unsupported URI scheme: %s", uri)
	}

	// Remove leading slashes
	for len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if scheme == "rmd" {
		// Reading an R Markdown file
		filePath := filepath.Join(RMD_DIR, path)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", "", fmt.Errorf("R Markdown file not found: %s", path)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", "", fmt.Errorf("failed to read R Markdown file: %w", err)
		}

		return string(content), "text/markdown", nil
	} else if scheme == "rmd-output" {
		// Reading a rendered output file
		filePath := filepath.Join(OUTPUT_DIR, path)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", "", fmt.Errorf("rendered output file not found: %s", path)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", "", fmt.Errorf("failed to read rendered output file: %w", err)
		}

		return string(content), GetMimeType(path), nil
	}

	return "", "", fmt.Errorf("unsupported URI scheme: %s", scheme)
}

// ListTools lists available tools for R Markdown
func (s *MCPServer) ListTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "create_rmd",
			"description": "Create a new R Markdown file",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "Filename for the R Markdown file (without extension)",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Title for the R Markdown document",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Content of the R Markdown file",
					},
				},
				"required": []string{"filename", "title", "content"},
			},
		},
		{
			"name":        "render_rmd",
			"description": "Render an R Markdown file",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filename": map[string]interface{}{
						"type":        "string",
						"description": "Filename of the R Markdown file to render",
					},
					"format": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"html", "pdf", "word"},
						"description": "Output format (html, pdf, or word)",
						"default":     "html",
					},
					"use_docker_compose": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to use docker-compose (true) or Dockerode (false)",
						"default":     false,
					},
				},
				"required": []string{"filename"},
			},
		},
	}
}

// CallTool calls an R Markdown tool
func (s *MCPServer) CallTool(name string, args map[string]interface{}) (string, error) {
	switch name {
	case "create_rmd":
		// Extract arguments
		filename, ok := args["filename"].(string)
		if !ok || filename == "" {
			return "", fmt.Errorf("filename is required")
		}

		title, ok := args["title"].(string)
		if !ok || title == "" {
			return "", fmt.Errorf("title is required")
		}

		content, ok := args["content"].(string)
		if !ok || content == "" {
			return "", fmt.Errorf("content is required")
		}

		// Ensure filename has .Rmd extension
		fullFilename := EnsureRmdExtension(filename)
		filePath := filepath.Join(RMD_DIR, fullFilename)

		// Create YAML front matter if not present
		finalContent := CreateRmdFrontMatter(title, content)

		// Write the file
		if err := os.WriteFile(filePath, []byte(finalContent), 0644); err != nil {
			return "", fmt.Errorf("failed to write R Markdown file: %w", err)
		}

		return fmt.Sprintf("Created R Markdown file: %s", fullFilename), nil

	case "render_rmd":
		// Extract arguments
		filename, ok := args["filename"].(string)
		if !ok || filename == "" {
			return "", fmt.Errorf("filename is required")
		}

		format := "html"
		if f, ok := args["format"].(string); ok && (f == "html" || f == "pdf" || f == "word") {
			format = f
		}

		useDockerCompose := false
		if udc, ok := args["use_docker_compose"].(bool); ok {
			useDockerCompose = udc
		}

		// Ensure filename has .Rmd extension
		fullFilename := EnsureRmdExtension(filename)
		filePath := filepath.Join(RMD_DIR, fullFilename)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", fmt.Errorf("R Markdown file not found: %s", fullFilename)
		}

		// Build Docker image
		if err := BuildDockerImage(); err != nil {
			return "", fmt.Errorf("failed to build Docker image: %w", err)
		}

		// Render the R Markdown file
		outputFile, err := RenderRMarkdown(fullFilename, format, useDockerCompose)
		if err != nil {
			return "", fmt.Errorf("failed to render R Markdown file: %w", err)
		}

		return fmt.Sprintf("Successfully rendered %s to %s", fullFilename, outputFile), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
