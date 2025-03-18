#!/bin/bash
set -e

echo "Setting up containerized MCP server..."

# Build and start the containers
echo "Building and starting containers..."
docker-compose up -d r-server-mcp

# Get the container ID
CONTAINER_ID=$(docker-compose ps -q r-server-mcp)
echo "MCP server container ID: $CONTAINER_ID"

# Create the updated MCP settings
echo "Creating updated MCP settings..."
cat << EOF
{
  "mcpServers": {
    "r-server": {
      "command": "docker",
      "args": ["exec", "-i", "$CONTAINER_ID", "node", "build/index.js"],
      "disabled": false,
      "autoApprove": []
    }
  }
}
EOF

echo ""
echo "To use the containerized MCP server, update your MCP settings file at:"
echo "~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json"
echo "with the configuration above."
echo ""
echo "For Claude desktop app, update:"
echo "~/Library/Application Support/Claude/claude_desktop_config.json"
echo ""
echo "The containerized MCP server is now running. You can stop it with:"
echo "docker-compose down"
