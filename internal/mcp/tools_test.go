package mcp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRExecutor is a mock implementation of RExecutor for testing
type MockRExecutor struct {
	MockExecuteRScript func(config RExecutionConfig) ([]byte, error)
}

// ExecuteRScript is the mock implementation
func (m *MockRExecutor) ExecuteRScript(config RExecutionConfig) ([]byte, error) {
	if m.MockExecuteRScript != nil {
		return m.MockExecuteRScript(config)
	}
	return []byte{}, nil
}

// SetupMockExecutor sets up a mock executor for testing and returns a cleanup function
func SetupMockExecutor(mock RExecutor) func() {
	original := DefaultExecutor
	DefaultExecutor = mock
	return func() {
		DefaultExecutor = original
	}
}

// TestRenderGGPlotExecution tests the R code execution functionality
func TestRenderGGPlotExecution(t *testing.T) {
	// Create a mock executor
	mockExecutor := &MockRExecutor{
		MockExecuteRScript: func(config RExecutionConfig) ([]byte, error) {
			// Verify the script content and parameters
			assert.Contains(t, config.ScriptPath, "script.R")
			assert.Contains(t, config.OutputPath, "output.png")
			assert.Equal(t, "png", config.OutputFormat)
			assert.Equal(t, 800, config.Width)
			assert.Equal(t, 600, config.Height)
			assert.Equal(t, 96, config.Resolution)

			// Return mock image data
			return []byte("mock-image-data"), nil
		},
	}

	// Set up the mock executor
	cleanup := SetupMockExecutor(mockExecutor)
	defer cleanup()

	// Call RenderGGPlot with valid arguments
	args := GGPlotRenderArgs{
		Code: "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
	}

	response, err := RenderGGPlot(args)

	// Verify the response
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Content)
	require.Len(t, response.Content, 1)
	
	// Verify the content type is image
	assert.Equal(t, "image", string(response.Content[0].Type))
	
	// Verify the image content
	require.NotNil(t, response.Content[0].ImageContent)
	assert.Equal(t, "image/png", response.Content[0].ImageContent.MimeType)
	assert.Equal(t, "bW9jay1pbWFnZS1kYXRh", response.Content[0].ImageContent.Data) // base64 encoded "mock-image-data"
}

// TestRenderGGPlotExecutionError tests error handling in R code execution
func TestRenderGGPlotExecutionError(t *testing.T) {
	// Create a mock executor that returns an error
	mockExecutor := &MockRExecutor{
		MockExecuteRScript: func(config RExecutionConfig) ([]byte, error) {
			return nil, errors.New("mock execution error")
		},
	}

	// Set up the mock executor
	cleanup := SetupMockExecutor(mockExecutor)
	defer cleanup()

	// Call RenderGGPlot with valid arguments
	args := GGPlotRenderArgs{
		Code: "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
	}

	response, err := RenderGGPlot(args)

	// Verify the error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute R script")
	assert.Contains(t, err.Error(), "mock execution error")
	assert.Nil(t, response)
}

// TestRenderGGPlotValidation tests the validation logic in RenderGGPlot
func TestRenderGGPlotValidation(t *testing.T) {
	tests := []struct {
		name        string
		args        GGPlotRenderArgs
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty code",
			args:        GGPlotRenderArgs{},
			expectError: true,
			errorMsg:    "code is required",
		},
		{
			name: "Width too small",
			args: GGPlotRenderArgs{
				Code:  "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
				Width: 50,
			},
			expectError: true,
			errorMsg:    "width must be between 100 and 5000",
		},
		{
			name: "Width too large",
			args: GGPlotRenderArgs{
				Code:  "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
				Width: 6000,
			},
			expectError: true,
			errorMsg:    "width must be between 100 and 5000",
		},
		{
			name: "Height too small",
			args: GGPlotRenderArgs{
				Code:   "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
				Height: 50,
			},
			expectError: true,
			errorMsg:    "height must be between 100 and 5000",
		},
		{
			name: "Height too large",
			args: GGPlotRenderArgs{
				Code:   "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
				Height: 6000,
			},
			expectError: true,
			errorMsg:    "height must be between 100 and 5000",
		},
		{
			name: "Resolution too small",
			args: GGPlotRenderArgs{
				Code:       "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
				Resolution: 50,
			},
			expectError: true,
			errorMsg:    "resolution must be between 72 and 600",
		},
		{
			name: "Resolution too large",
			args: GGPlotRenderArgs{
				Code:       "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
				Resolution: 700,
			},
			expectError: true,
			errorMsg:    "resolution must be between 72 and 600",
		},
		{
			name: "Valid arguments with defaults",
			args: GGPlotRenderArgs{
				Code: "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
			},
			expectError: false,
		},
		{
			name: "Valid arguments with custom values",
			args: GGPlotRenderArgs{
				Code:       "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
				OutputType: "pdf",
				Width:      1200,
				Height:     800,
				Resolution: 300,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := RenderGGPlot(tt.args)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, response)
				// In a real implementation, we would also check the response content
			}
		})
	}
}
