#!/bin/bash
set -e

# Create the updated MCP settings
echo "Creating updated MCP settings..."
cat << EOF
{
  "mcpServers": {
    "r-server": {
      "command": "docker-compose",
      "args": ["-f", "/path/to/r-server/docker-compose.yml", "run", "--rm", "r-server-mcp"],
      "disabled": false,
      "autoApprove": ["render_ggplot"]
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
