#!/bin/bash
set -e

# Change to the directory containing this script
cd "$(dirname "$0")"

# Check if we should use Docker or run locally
if [ "$USE_DOCKER" = "true" ] || [ "$1" = "--docker" ]; then
  echo "Starting R MCP server using Docker..." >&2
  
  # Build the Docker image if needed
  docker-compose build
  
  # Run the container with stdin/stdout connected
  exec docker-compose run --rm r-server-mcp
else
  echo "Starting R MCP server locally..." >&2
  
  # Run the server locally
  exec ./r-server
fi

# Debugging options (uncomment if needed):
# mkfifo mcp.fifo
# if [ "$USE_DOCKER" = "true" ] || [ "$1" = "--docker" ]; then
#   tee mcp.fifo | docker-compose run --rm r-server-mcp | tee mcp.fifo
# else
#   tee mcp.fifo | ./r-server | tee mcp.fifo
# fi
