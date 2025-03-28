package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Define test cases for MCP protocol validation
var testCases = []struct {
	name     string
	method   string
	params   interface{}
	validate func(map[string]interface{}) error
}{
	{
		name:   "List Tools",
		method: "listTools",
		params: map[string]interface{}{},
		validate: func(response map[string]interface{}) error {
			// Check if the response has a result field
			result, ok := response["result"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("response does not have a result field")
			}

			// Check if the result has a tools field
			tools, ok := result["tools"].([]interface{})
			if !ok {
				return fmt.Errorf("result does not have a tools field")
			}

			// Check if the tools field is an array
			if len(tools) == 0 {
				return fmt.Errorf("tools array is empty")
			}

			// Check if the render_ggplot tool is present
			found := false
			for _, tool := range tools {
				toolMap, ok := tool.(map[string]interface{})
				if !ok {
					continue
				}

				name, ok := toolMap["name"].(string)
				if !ok {
					continue
				}

				if name == "render_ggplot" {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("render_ggplot tool not found")
			}

			return nil
		},
	},
	{
		name:   "Call Tool Schema Validation",
		method: "callTool",
		params: map[string]interface{}{
			"name": "render_ggplot",
			"arguments": map[string]interface{}{
				// Missing "code" field which should be required
				"output_type": "png",
				"width":       800,
				"height":      600,
			},
		},
		validate: func(response map[string]interface{}) error {
			// Check if the response has an error field
			errorObj, ok := response["error"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("response does not have an error field")
			}

			// Check if the error has a code field
			code, ok := errorObj["code"].(float64)
			if !ok {
				return fmt.Errorf("error does not have a code field")
			}

			// Check if the error code is InvalidParams (code -32602)
			if code != -32602 {
				return fmt.Errorf("expected error code -32602, got %v", code)
			}

			return nil
		},
	},
}

func main() {
	fmt.Println("MCP Protocol Validator")
	fmt.Println("=====================")

	// Get the server URL from the command line or use the default
	serverURL := "http://localhost:22011/mcp"
	if len(os.Args) > 1 {
		serverURL = os.Args[1]
	}

	fmt.Printf("Testing MCP server at %s\n\n", serverURL)

	// Run all test cases
	failedTests := 0
	for _, tc := range testCases {
		fmt.Printf("Running test: %s\n", tc.name)

		// Create the JSON-RPC request
		request := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      "test-" + tc.name,
			"method":  tc.method,
			"params":  tc.params,
		}

		// Convert the request to JSON
		requestJSON, err := json.Marshal(request)
		if err != nil {
			fmt.Printf("❌ Failed to marshal request: %v\n", err)
			failedTests++
			continue
		}

		// Send the request to the server
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(requestJSON))
		if err != nil {
			fmt.Printf("❌ Failed to send request: %v\n", err)
			failedTests++
			continue
		}
		defer resp.Body.Close()

		// Read the response
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("❌ Failed to read response: %v\n", err)
			failedTests++
			continue
		}

		// Parse the response
		var response map[string]interface{}
		if err := json.Unmarshal(responseBody, &response); err != nil {
			fmt.Printf("❌ Failed to parse response: %v\n", err)
			failedTests++
			continue
		}

		// Validate the response
		if err := tc.validate(response); err != nil {
			fmt.Printf("❌ Validation failed: %v\n", err)
			failedTests++
			continue
		}

		fmt.Printf("✅ Test passed\n")
	}

	// Print summary
	fmt.Println("\nTest Summary")
	fmt.Println("===========")
	fmt.Printf("Total tests: %d\n", len(testCases))
	fmt.Printf("Failed tests: %d\n", failedTests)

	if failedTests > 0 {
		os.Exit(1)
	}
}
