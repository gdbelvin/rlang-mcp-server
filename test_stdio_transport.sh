#!/bin/bash

# This script demonstrates how to use the MCP server and client with stdio transport
# using named pipes for communication

# Create named pipes for communication
echo "Creating named pipes..."
mkfifo server_in server_out

# Start the server in the background, reading from server_in and writing to server_out
echo "Starting server..."
./cmd/r-server/r-server < server_in > server_out &
SERVER_PID=$!

# Give the server a moment to start
sleep 1

# Start the client, reading from server_out and writing to server_in
echo "Starting client..."
./test/examples/client/client < server_out > server_in

# Clean up
echo "Cleaning up..."
kill $SERVER_PID
rm server_in server_out

echo "Done!"
