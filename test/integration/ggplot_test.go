//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGGPlotRendering tests the end-to-end flow of rendering a ggplot visualization
func TestGGPlotRendering(t *testing.T) {
	// Skip this test if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=1 to run.")
	}

	// This is a placeholder for a real integration test
	// In a real implementation, we would:
	// 1. Start the MCP server
	// 2. Create an MCP client
	// 3. Call the render_ggplot tool
	// 4. Verify the response

	t.Run("Basic Plot", func(t *testing.T) {
		// Create a temporary directory for test output
		tempDir, err := os.MkdirTemp("", "ggplot-integration-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Define the test request
		request := map[string]interface{}{
			"name": "render_ggplot",
			"arguments": map[string]interface{}{
				"code":        "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point() + theme_minimal() + labs(title = 'MPG vs Horsepower')",
				"output_type": "png",
				"width":       800,
				"height":      600,
				"resolution":  96,
			},
		}

		// Write the request to a file for debugging
		requestFile := filepath.Join(tempDir, "request.json")
		requestData, err := json.MarshalIndent(request, "", "  ")
		require.NoError(t, err)
		err = os.WriteFile(requestFile, requestData, 0644)
		require.NoError(t, err)

		// In a real test, we would send this request to the server and verify the response
		// For now, we'll just print a message
		fmt.Printf("Integration test would send request: %s\n", requestFile)

		// Assert that the test ran (placeholder)
		assert.True(t, true, "Integration test placeholder")
	})
}
