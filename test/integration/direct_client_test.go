//go:build integration

package integration

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"r-server/internal/mcp"
)

// TestDirectClientGGPlotRendering tests rendering a ggplot visualization by directly calling the RenderGGPlot function
// This is a simpler approach that avoids the complexity of setting up an HTTP server and client
func TestDirectClientGGPlotRendering(t *testing.T) {
	// Skip this test if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=1 to run.")
	}

	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "ggplot-direct-client-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Read the test R script
	scriptPath := filepath.Join("..", "testdata", "r_scripts", "basic_plot.R")
	scriptData, err := os.ReadFile(scriptPath)
	require.NoError(t, err)

	// Create the arguments for the render_ggplot tool
	args := mcp.GGPlotRenderArgs{
		Code:       string(scriptData),
		OutputType: "png",
		Width:      800,
		Height:     600,
		Resolution: 96,
	}

	// Call the RenderGGPlot function directly
	fmt.Println("Calling RenderGGPlot function directly...")
	response, err := mcp.RenderGGPlot(args)
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

	// Save the image to a file for inspection
	imageFile := filepath.Join(tempDir, "output.png")
	imageData, err := base64.StdEncoding.DecodeString(response.Content[0].ImageContent.Data)
	require.NoError(t, err)
	err = os.WriteFile(imageFile, imageData, 0644)
	require.NoError(t, err)

	fmt.Printf("Successfully rendered ggplot image to %s\n", imageFile)
}
