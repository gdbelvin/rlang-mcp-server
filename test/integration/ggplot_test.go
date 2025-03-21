//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"r-server/internal/mcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBasicGGPlotRendering tests rendering a basic ggplot visualization
func TestBasicGGPlotRendering(t *testing.T) {
	// Skip this test if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=1 to run.")
	}

	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "ggplot-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Read the test R script
	scriptPath := filepath.Join("..", "testdata", "r_scripts", "basic_plot.R")
	scriptData, err := os.ReadFile(scriptPath)
	require.NoError(t, err)

	// Create the arguments for the render_ggplot tool
	args := map[string]interface{}{
		"code":        string(scriptData),
		"output_type": "png",
		"width":       800,
		"height":      600,
		"resolution":  96,
	}

	// Create a request for the render_ggplot tool
	request := map[string]interface{}{
		"name":      "render_ggplot",
		"arguments": args,
	}

	// Write the request to a file for debugging
	requestFile := filepath.Join(tempDir, "request.json")
	requestData, err := json.MarshalIndent(request, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(requestFile, requestData, 0644)
	require.NoError(t, err)

	// Call the RenderGGPlot function directly
	// This is a simplified approach for integration testing
	// In a real-world scenario, we would use an MCP client to call the tool
	fmt.Printf("Integration test sending request: %s\n", requestFile)

	// Convert the arguments to the expected type
	renderArgs := mcp.GGPlotRenderArgs{
		Code:       args["code"].(string),
		OutputType: args["output_type"].(string),
		Width:      args["width"].(int),
		Height:     args["height"].(int),
		Resolution: args["resolution"].(int),
	}

	// Call the RenderGGPlot function
	response, err := mcp.RenderGGPlot(renderArgs)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Content)
	require.Len(t, response.Content, 1)

	// Verify the content type is image
	assert.Equal(t, "image", string(response.Content[0].Type))

	// Verify the image content
	require.NotNil(t, response.Content[0].ImageContent)
	assert.Equal(t, "image/png", response.Content[0].ImageContent.MimeType)
	assert.NotEmpty(t, response.Content[0].ImageContent.Data)

	fmt.Println("Successfully rendered basic ggplot image")
}

// TestComplexGGPlotRendering tests rendering a complex ggplot visualization
func TestComplexGGPlotRendering(t *testing.T) {
	// Skip this test if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=1 to run.")
	}

	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "ggplot-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Read the test R script
	scriptPath := filepath.Join("..", "testdata", "r_scripts", "complex_plot.R")
	scriptData, err := os.ReadFile(scriptPath)
	require.NoError(t, err)

	// Create the arguments for the render_ggplot tool
	args := map[string]interface{}{
		"code":        string(scriptData),
		"output_type": "png",
		"width":       1200,
		"height":      800,
		"resolution":  150,
	}

	// Create a request for the render_ggplot tool
	request := map[string]interface{}{
		"name":      "render_ggplot",
		"arguments": args,
	}

	// Write the request to a file for debugging
	requestFile := filepath.Join(tempDir, "request.json")
	requestData, err := json.MarshalIndent(request, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(requestFile, requestData, 0644)
	require.NoError(t, err)

	// Call the RenderGGPlot function directly
	fmt.Printf("Integration test sending request: %s\n", requestFile)

	// Convert the arguments to the expected type
	renderArgs := mcp.GGPlotRenderArgs{
		Code:       args["code"].(string),
		OutputType: args["output_type"].(string),
		Width:      args["width"].(int),
		Height:     args["height"].(int),
		Resolution: args["resolution"].(int),
	}

	// Call the RenderGGPlot function
	response, err := mcp.RenderGGPlot(renderArgs)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Content)
	require.Len(t, response.Content, 1)

	// Verify the content type is image
	assert.Equal(t, "image", string(response.Content[0].Type))

	// Verify the image content
	require.NotNil(t, response.Content[0].ImageContent)
	assert.Equal(t, "image/png", response.Content[0].ImageContent.MimeType)
	assert.NotEmpty(t, response.Content[0].ImageContent.Data)

	fmt.Println("Successfully rendered complex ggplot image")
}
