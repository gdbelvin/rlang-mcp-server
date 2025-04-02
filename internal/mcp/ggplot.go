package mcp

import (
	"fmt"
	"os"
	"path/filepath"

	mcp "github.com/metoro-io/mcp-golang"
)

// GGPlotRenderArgs represents the arguments for rendering a ggplot image
type GGPlotRenderArgs struct {
	Code       string `json:"code" jsonschema:"required,description=R code containing ggplot commands"`
	OutputType string `json:"output_type" jsonschema:"description=Output format (png, jpeg, pdf, svg)"`
	Width      int    `json:"width" jsonschema:"description=Width of the output image in pixels"`
	Height     int    `json:"height" jsonschema:"description=Height of the output image in pixels"`
	Resolution int    `json:"resolution" jsonschema:"description=Resolution of the output image in dpi"`
}

// RenderGGPlot renders a ggplot2 visualization and returns the image directly in the response
func RenderGGPlot(args GGPlotRenderArgs) (*mcp.ToolResponse, error) {
	// Validate arguments
	if args.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	// Set default values
	outputType := args.OutputType
	if outputType == "" {
		outputType = "png"
	}

	width := args.Width
	if width == 0 {
		width = 800
	} else if width < 100 || width > 5000 {
		return nil, fmt.Errorf("width must be between 100 and 5000")
	}

	height := args.Height
	if height == 0 {
		height = 600
	} else if height < 100 || height > 5000 {
		return nil, fmt.Errorf("height must be between 100 and 5000")
	}

	resolution := args.Resolution
	if resolution == 0 {
		resolution = 96
	} else if resolution < 72 || resolution > 600 {
		return nil, fmt.Errorf("resolution must be between 72 and 600")
	}

	// Create a temporary directory for the R script and output
	tempDir, err := os.MkdirTemp("", "ggplot-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the R script
	scriptPath := filepath.Join(tempDir, "script.R")
	outputPath := filepath.Join(tempDir, fmt.Sprintf("output.%s", outputType))

	// Generate the R script content
	scriptContent := fmt.Sprintf(`
# Load required libraries
library(ggplot2)
library(cowplot)

# Set output parameters
width <- %d
height <- %d
dpi <- %d
output_file <- "%s"
pdf(NULL)

# Execute the provided code
%s

# Save the last plot
ggsave(output_file, width = width/dpi, height = height/dpi, dpi = dpi)
`, width, height, resolution, outputPath, args.Code)

	// Write the R script to a file
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write R script: %w", err)
	}

	// Execute the R script
	config := RExecutionConfig{
		ScriptPath:   scriptPath,
		OutputPath:   outputPath,
		OutputFormat: outputType,
		Width:        width,
		Height:       height,
		Resolution:   resolution,
	}

	imageData, err := ExecuteRScript(config)
	if err != nil {
		return nil, fmt.Errorf("failed to execute R script: %w", err)
	}

	// Create the image content
	imageContent := mcp.NewImageContent(
		EncodeImageToBase64(imageData),
		GetMimeType(outputType))

	return mcp.NewToolResponse(imageContent), nil
}
