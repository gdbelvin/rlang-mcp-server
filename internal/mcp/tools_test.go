package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
