#!/bin/bash
set -x  # Enable debug mode to print each command

# This script tests the example server directly with a simple request

# Create a temporary directory for test files
TEMP_DIR=$(mktemp -d)
echo "Created temporary directory: $TEMP_DIR"

# Create a simple request file for the render_ggplot tool
REQUEST_FILE="$TEMP_DIR/request.json"
cat > "$REQUEST_FILE" << EOF
{
  "name": "render_ggplot",
  "arguments": {
    "code": "library(ggplot2); ggplot(mtcars, aes(x=wt, y=mpg)) + geom_point()",
    "output_type": "png",
    "width": 800,
    "height": 600
  }
}
EOF
echo "Created request file: $REQUEST_FILE"

# Run the r-server with the -test-tool flag to test a specific tool
echo "Testing the tool..."
./cmd/r-server/r-server -test-tool "$REQUEST_FILE"

# Clean up
echo "Cleaning up..."
rm -rf "$TEMP_DIR"

echo "Done!"
