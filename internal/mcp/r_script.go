package mcp

import (
	"fmt"
	"os"
	"path/filepath"

	mcp "github.com/metoro-io/mcp-golang"
)

// RScriptArgs represents the arguments for executing an R script
type RScriptArgs struct {
	Code string `json:"code" jsonschema:"required,description=R code to execute"`
}

// ExecuteRScriptTool executes an R script and returns the result as text
func ExecuteRScriptTool(args RScriptArgs) (*mcp.ToolResponse, error) {
	// Validate arguments
	if args.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	// Create a temporary directory for the R script and output
	tempDir, err := os.MkdirTemp("", "r-script-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the R script
	scriptPath := filepath.Join(tempDir, "script.R")
	outputPath := filepath.Join(tempDir, "output.txt")

	// Generate the R script content with output capture
	scriptContent := fmt.Sprintf(`
# Redirect output to a file
sink("%s")

# Execute the provided code
%s

# Make sure the output file is created
cat("Output file path: %s\n")

# Close the output file
sink()

# Write directly to the output file as a fallback
cat("R script execution completed successfully!\n", file="%s", append=TRUE)
`, outputPath, args.Code, outputPath, outputPath)

	// Write the R script to a file
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write R script: %w", err)
	}

	// Execute the R script
	config := RExecutionConfig{
		ScriptPath: scriptPath,
		OutputPath: outputPath,
	}

	// Execute the R script using the existing ExecuteRScript function
	// This will return the output data directly
	outputData, err := ExecuteRScript(config)
	if err != nil {
		return nil, fmt.Errorf("failed to execute R script: %w", err)
	}

	// Create the text content
	textContent := mcp.NewTextContent(string(outputData))

	return mcp.NewToolResponse(textContent), nil
}
