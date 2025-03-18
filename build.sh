#!/bin/bash
set -e

echo "Building R-Server Go implementation..."

# Build the application
echo "Building application..."
go build -o r-server

echo "Build complete. You can run the server with:"
echo "./r-server"
