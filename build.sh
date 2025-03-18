#!/bin/bash
set -e

echo "Building R-Server Go implementation..."

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Update go.sum file
echo "Updating go.sum file..."
go get -v

# Build the application
echo "Building application..."
go build -o r-server

echo "Build complete. You can run the server with:"
echo "./r-server"
