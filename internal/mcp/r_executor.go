package mcp

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RExecutionConfig represents configuration for R script execution
type RExecutionConfig struct {
	ScriptPath   string
	OutputPath   string
	OutputFormat string
	Width        int
	Height       int
	Resolution   int
}

// RExecutor defines the interface for executing R scripts
type RExecutor interface {
	ExecuteRScript(config RExecutionConfig) ([]byte, error)
}

// DefaultRExecutor is the default implementation of RExecutor
type DefaultRExecutor struct{}

// ExecuteRScript executes an R script and returns the output image data
func (e *DefaultRExecutor) ExecuteRScript(config RExecutionConfig) ([]byte, error) {
	// Create the output directory if it doesn't exist
	outputDir := filepath.Dir(config.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Make sure to clean up temporary files when done
	defer func() {
		// Remove the script file
		if err := os.Remove(config.ScriptPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove script file: %v\n", err)
		}

		// Remove the output file
		if err := os.Remove(config.OutputPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove output file: %v\n", err)
		}

		// Remove the output directory if it's empty
		if files, err := os.ReadDir(outputDir); err == nil && len(files) == 0 {
			if err := os.Remove(outputDir); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to remove output directory: %v\n", err)
			}
		}
	}()

	// Execute the R script
	cmd := exec.Command("Rscript", config.ScriptPath)

	// Set environment variables for the R script
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("OUTPUT_PATH=%s", config.OutputPath),
		fmt.Sprintf("WIDTH=%d", config.Width),
		fmt.Sprintf("HEIGHT=%d", config.Height),
		fmt.Sprintf("RESOLUTION=%d", config.Resolution),
	)

	// Capture stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute R script: %w\nOutput: %s", err, string(output))
	}

	// Read the output file
	outputData, err := os.ReadFile(config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	return outputData, nil
}

// Default executor instance
var DefaultExecutor RExecutor = &DefaultRExecutor{}

// ExecuteRScript is a convenience function that uses the default executor
func ExecuteRScript(config RExecutionConfig) ([]byte, error) {
	return DefaultExecutor.ExecuteRScript(config)
}

// GetMimeType returns the MIME type for the given output format
func GetMimeType(outputFormat string) string {

	/* Claude supported image types:
	image/jpeg
	image/png
	image/gif
	image/webp
	*/

	switch outputFormat {
	case "png":
		return "image/png"
	case "jpeg", "jpg":
		return "image/jpeg"
	case "pdf":
		return "application/pdf"
	case "svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// EncodeImageToBase64 encodes image data to base64
func EncodeImageToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
