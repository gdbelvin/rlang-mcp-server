#!/bin/bash
set -e

echo "Setting up local Go module proxy for offline builds..."

# Create a directory for the Go module proxy
PROXY_DIR="./go-proxy"
mkdir -p "$PROXY_DIR"

# Download all dependencies to the local cache
echo "Downloading Go dependencies..."
go mod download

# Copy the Go module cache to our proxy directory
echo "Copying modules to proxy directory..."
GO_MOD_CACHE=$(go env GOMODCACHE)
cp -R "$GO_MOD_CACHE"/* "$PROXY_DIR"/ || echo "No modules found in cache"

# Create a simple index file for the proxy
echo "Creating proxy index..."
find "$PROXY_DIR" -type d -name "v*" | while read -r dir; do
  module_path=$(dirname "$dir")
  module_path=${module_path#"$PROXY_DIR/"}
  module_version=$(basename "$dir")
  echo "{\"Version\":\"$module_version\",\"Time\":\"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"}" > "$dir/index.json"
done

echo "Local Go module proxy setup complete at $PROXY_DIR"
echo "You can now build the Docker image with:"
echo "docker build --build-arg GOPROXY=file://$PWD/$PROXY_DIR -t r-server:latest ."
