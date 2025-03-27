//go:build integration

package integration

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestSimpleStdioTransport tests a simple stdio transport between a server and client
// This is to verify that the stdio transport is working correctly
func TestSimpleStdioTransport(t *testing.T) {
	// Skip this test if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=1 to run.")
	}

	// Create a simple echo server script
	echoScript := `#!/bin/bash
while read line; do
  echo "ECHO: $line"
done
`
	scriptPath := "/tmp/echo_server.sh"
	err := os.WriteFile(scriptPath, []byte(echoScript), 0755)
	require.NoError(t, err)
	defer os.Remove(scriptPath)

	// Start the echo server process
	cmd := exec.Command(scriptPath)
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)
	defer cmd.Process.Kill()

	// Send a message to the server
	testMessage := "Hello, stdio transport!"
	_, err = io.WriteString(stdin, testMessage+"\n")
	require.NoError(t, err)

	// Read the response
	buf := make([]byte, 1024)
	n, err := stdout.Read(buf)
	require.NoError(t, err)
	response := string(buf[:n])

	// Check the response
	expectedResponse := fmt.Sprintf("ECHO: %s\n", testMessage)
	require.Equal(t, expectedResponse, response)

	// Clean up
	stdin.Close()
}
