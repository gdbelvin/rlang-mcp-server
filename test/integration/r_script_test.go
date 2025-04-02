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

// TestExecuteRScript tests executing an R script and returning the result
func TestExecuteRScript(t *testing.T) {
	// Skip this test if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=1 to run.")
	}

	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "r-script-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Define the R script to execute
	rScript := `
# Simple R script to test the execute_r_script tool

# Create a data frame
data <- data.frame(
  x = 1:10,
  y = (1:10)^2
)

# Calculate summary statistics
summary_stats <- summary(data)

# Print the summary statistics
print(summary_stats)

# Calculate correlation
correlation <- cor(data$x, data$y)
cat("Correlation between x and y:", correlation, "\n")

# Create a linear model
model <- lm(y ~ x, data = data)
cat("\nLinear model summary:\n")
print(summary(model))

# Print a message
cat("\nR script execution completed successfully!\n")
`

	// Create the arguments for the execute_r_script tool
	args := map[string]interface{}{
		"code": rScript,
	}

	// Create a request for the execute_r_script tool
	request := map[string]interface{}{
		"name":      "execute_r_script",
		"arguments": args,
	}

	// Write the request to a file for debugging
	requestFile := filepath.Join(tempDir, "request.json")
	requestData, err := json.MarshalIndent(request, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(requestFile, requestData, 0644)
	require.NoError(t, err)

	// Call the ExecuteRScriptTool function directly
	fmt.Printf("Integration test sending request: %s\n", requestFile)

	// Convert the arguments to the expected type
	scriptArgs := mcp.RScriptArgs{
		Code: args["code"].(string),
	}

	// Call the ExecuteRScriptTool function
	response, err := mcp.ExecuteRScriptTool(scriptArgs)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Content)
	require.Len(t, response.Content, 1)

	// Verify the content type is text
	assert.Equal(t, "text", string(response.Content[0].Type))

	// Verify the text content
	require.NotNil(t, response.Content[0].TextContent)
	textContent := response.Content[0].TextContent.Text
	assert.Contains(t, textContent, "Correlation between x and y:")
	assert.Contains(t, textContent, "Linear model summary:")
	assert.Contains(t, textContent, "R script execution completed successfully!")

	fmt.Println("Successfully executed R script")
}
