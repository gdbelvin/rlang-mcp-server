#!/bin/bash
set -x  # Enable debug mode to print each command

# This script tests the communication between the example server and client using stdio transport
# using named pipes for communication

# Create named pipes for communication
echo "Creating named pipes..."
mkfifo server_in server_out

# Start the example server in the background, reading from server_in and writing to server_out
echo "Starting example server..."
./test/examples/server/server < server_in > server_out 2>server_stderr.log &
SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"

# Give the server a moment to start
sleep 2
echo "Server had time to initialize"

# Start the example client, reading from server_out and writing to server_in
echo "Starting example client..."
./test/examples/client/client < server_out > server_in 2>client_stderr.log
CLIENT_EXIT=$?
echo "Client exited with status: $CLIENT_EXIT"

# Display logs
echo "Server stderr log:"
cat server_stderr.log

echo "Client stderr log:"
cat client_stderr.log

# Clean up
echo "Cleaning up..."
kill $SERVER_PID || echo "Server already terminated"
rm -f server_in server_out server_stderr.log client_stderr.log

echo "Done!"
