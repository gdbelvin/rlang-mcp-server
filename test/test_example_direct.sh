#!/bin/bash
set -x  # Enable debug mode to print each command

# This script tests the example server and client directly without named pipes

# Clean up any previous test artifacts
rm -f server_output.log client_output.log

# Start the example server in the background
echo "Starting example server..."
./test/examples/server/server > server_output.log 2>&1 &
SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"

# Give the server a moment to start
sleep 2
echo "Server had time to initialize"

# Run the client with --help to see available options
echo "Checking client help..."
./test/examples/client/client --help > client_output.log 2>&1

# Display logs
echo "Server output:"
cat server_output.log

echo "Client output:"
cat client_output.log

# Clean up
echo "Cleaning up..."
kill $SERVER_PID || echo "Server already terminated"
rm -f server_output.log client_output.log

echo "Done!"
