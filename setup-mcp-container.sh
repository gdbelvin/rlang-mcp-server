#!/bin/bash
set -e

echo "Setting up containerized MCP server..."

# Build and start the containers
echo "Building and starting containers..."
docker-compose up -d 

# Create the updated MCP settings
echo "Creating updated MCP settings..."
cat << EOF
{
  "mcpServers": {
    "r-server": {
      "command": "docker-compose",
      "args": ["-f", "/Users/gdb/dev/r-server/docker-compose.yml", "up", "-d"]
      "disabled": false,
      "autoApprove": [],
      "transport": {
        "type": "http",
        "url": "http://localhost:22011/mcp"
      }
    }
  }
}
EOF

echo ""
echo "To use the containerized MCP server, update your MCP settings with the configuration above."
echo ""
echo "For Cline Visual Studio Plugin, update:"
echo "~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json"
echo ""
echo "For Claude desktop app, update:"
echo "~/Library/Application Support/Claude/claude_desktop_config.json"
echo ""
echo "The containerized MCP server is now running:"
docker ps --filter name=r-server
echo ""
echo "You can stop it with:"
echo "docker-compose down"
