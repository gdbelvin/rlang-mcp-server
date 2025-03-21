# R-Server MCP for ggplot: Testing Strategy

## 1. Introduction

This document outlines the comprehensive testing strategy for the R-Server MCP for ggplot project. The goal is to ensure that the server correctly implements the MCP protocol, properly renders ggplot2 visualizations, and handles errors gracefully.

## 2. Testing Objectives

The primary objectives of our testing strategy are to:

1. **Verify Protocol Conformance**: Ensure the server correctly implements the MCP protocol
2. **Validate Tool Functionality**: Confirm that the `render_ggplot` tool works as expected
3. **Test Error Handling**: Verify proper handling of various error scenarios
4. **Measure Performance**: Assess the performance characteristics of the server
5. **Ensure Security**: Validate that security measures are effective

## 3. Testing Levels

### 3.1 Unit Testing

Unit tests will focus on testing individual components in isolation:

| Component | Test Focus | Tools |
|-----------|------------|-------|
| MCP Server | Server initialization, tool registration | Go testing package |
| Tool Implementation | Parameter validation, response generation | Go testing package |
| R Code Execution | Script generation, execution | Go testing package, mocks |
| Image Processing | Image reading, format conversion | Go testing package |

### 3.2 Integration Testing

Integration tests will verify that components work together correctly:

| Integration Point | Test Focus | Tools |
|-------------------|------------|-------|
| Server + Tool | Tool registration and execution | Go testing package |
| Tool + R Execution | End-to-end tool execution | Go testing package, Docker |
| R Execution + Image Processing | Image generation and processing | Go testing package, Docker |

### 3.3 System Testing

System tests will validate the entire system as a whole:

| System Aspect | Test Focus | Tools |
|---------------|------------|-------|
| MCP Protocol | Protocol conformance | Custom MCP client |
| End-to-End Workflow | Complete request-response cycle | Custom MCP client |
| Error Scenarios | Error handling and recovery | Custom MCP client |

## 4. Test Environment

### 4.1 Local Development Environment

For local development and testing, we'll use:

- Go testing framework
- Docker for R execution
- Task for test automation
- VSCode for test execution and debugging

### 4.2 Continuous Integration Environment

For automated testing in CI, we'll use:

- Task for test execution
- Docker for containerized testing
- Test result reporting and aggregation

## 5. Test Categories

### 5.1 Functional Tests

#### 5.1.1 MCP Protocol Tests

Tests to verify correct implementation of the MCP protocol:

- **listTools Method**: Verify correct tool listing and schema
- **callTool Method**: Verify correct tool execution and response
- **Error Handling**: Verify correct error responses for various scenarios

```go
func TestListTools(t *testing.T) {
    // Setup test server
    server, transport := setupTestServer(t)
    
    // Send listTools request
    req := mcp.NewListToolsRequest("test-id")
    resp, err := transport.SendRequest(req)
    
    // Verify response
    require.NoError(t, err)
    require.Equal(t, "test-id", resp.ID)
    require.NotNil(t, resp.Result)
    
    // Verify tools list
    tools := resp.Result.(map[string]interface{})["tools"].([]interface{})
    require.Len(t, tools, 1)
    
    // Verify render_ggplot tool
    tool := tools[0].(map[string]interface{})
    require.Equal(t, "render_ggplot", tool["name"])
    require.NotNil(t, tool["inputSchema"])
}
```

#### 5.1.2 Tool Functionality Tests

Tests to verify the functionality of the `render_ggplot` tool:

- **Parameter Validation**: Test validation of all parameters
- **Default Values**: Test application of default values
- **R Code Execution**: Test execution of various R code snippets
- **Output Formats**: Test all supported output formats

```go
func TestRenderGGPlot(t *testing.T) {
    // Setup test server
    server, transport := setupTestServer(t)
    
    // Create callTool request
    req := mcp.NewCallToolRequest("test-id", "render_ggplot", map[string]interface{}{
        "code": "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
        "output_type": "png",
        "width": 800,
        "height": 600,
    })
    
    // Send request
    resp, err := transport.SendRequest(req)
    
    // Verify response
    require.NoError(t, err)
    require.Equal(t, "test-id", resp.ID)
    require.NotNil(t, resp.Result)
    
    // Verify content
    content := resp.Result.(map[string]interface{})["content"].([]interface{})
    require.Len(t, content, 1)
    
    // Verify image
    image := content[0].(map[string]interface{})
    require.Equal(t, "image", image["type"])
    require.Equal(t, "image/png", image["mimeType"])
    require.NotEmpty(t, image["data"])
}
```

#### 5.1.3 Error Handling Tests

Tests to verify proper handling of various error scenarios:

- **Invalid Parameters**: Test handling of invalid parameters
- **R Code Errors**: Test handling of R code errors
- **System Errors**: Test handling of system-level errors

```go
func TestInvalidParameters(t *testing.T) {
    // Setup test server
    server, transport := setupTestServer(t)
    
    // Test cases for invalid parameters
    testCases := []struct{
        name string
        args map[string]interface{}
        expectedError string
    }{
        {
            name: "Missing code",
            args: map[string]interface{}{
                "width": 800,
                "height": 600,
            },
            expectedError: "code is required",
        },
        {
            name: "Invalid width",
            args: map[string]interface{}{
                "code": "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
                "width": 10000,
            },
            expectedError: "width must be between 100 and 5000",
        },
        // Additional test cases...
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create callTool request
            req := mcp.NewCallToolRequest("test-id", "render_ggplot", tc.args)
            
            // Send request
            resp, err := transport.SendRequest(req)
            
            // Verify error response
            require.NoError(t, err)
            require.NotNil(t, resp.Error)
            require.Contains(t, resp.Error.Message, tc.expectedError)
        })
    }
}
```

### 5.2 Non-Functional Tests

#### 5.2.1 Performance Tests

Tests to measure the performance characteristics of the server:

- **Response Time**: Measure time to process requests
- **Throughput**: Measure requests per second
- **Resource Usage**: Measure CPU, memory, and disk usage

```go
func BenchmarkRenderGGPlot(b *testing.B) {
    // Setup test server
    server, transport := setupTestServer(b)
    
    // Create request
    req := mcp.NewCallToolRequest("bench-id", "render_ggplot", map[string]interface{}{
        "code": "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point()",
        "output_type": "png",
        "width": 800,
        "height": 600,
    })
    
    // Run benchmark
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        resp, err := transport.SendRequest(req)
        require.NoError(b, err)
        require.Nil(b, resp.Error)
    }
}
```

#### 5.2.2 Security Tests

Tests to validate the security measures:

- **Input Validation**: Test handling of malicious input
- **Resource Limits**: Test enforcement of resource limits
- **Error Information**: Test that error messages don't expose sensitive information

```go
func TestSecurityValidation(t *testing.T) {
    // Setup test server
    server, transport := setupTestServer(t)
    
    // Test cases for security validation
    testCases := []struct{
        name string
        code string
        expectedError string
    }{
        {
            name: "File system access",
            code: "write.csv(mtcars, '/etc/passwd')",
            expectedError: "permission denied",
        },
        {
            name: "Network access",
            code: "download.file('http://example.com', 'file.txt')",
            expectedError: "network access not allowed",
        },
        // Additional test cases...
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create callTool request
            req := mcp.NewCallToolRequest("test-id", "render_ggplot", map[string]interface{}{
                "code": tc.code,
            })
            
            // Send request
            resp, err := transport.SendRequest(req)
            
            // Verify error response
            require.NoError(t, err)
            require.NotNil(t, resp.Error)
            require.Contains(t, resp.Error.Message, tc.expectedError)
        })
    }
}
```

## 6. Test Data

### 6.1 Sample R Code

We'll create a library of sample R code for testing:

- **Basic Plots**: Simple ggplot2 visualizations
- **Complex Plots**: Multi-layer, faceted visualizations
- **Error Cases**: R code with various errors
- **Edge Cases**: R code that tests limits of the system

Example:

```r
# Basic plot
ggplot(mtcars, aes(x = mpg, y = hp)) + 
  geom_point() + 
  theme_minimal() + 
  labs(title = "MPG vs Horsepower")

# Complex plot
ggplot(mtcars, aes(x = mpg, y = hp, color = factor(cyl))) + 
  geom_point() + 
  geom_smooth(method = "lm") + 
  facet_wrap(~gear) + 
  theme_bw() + 
  labs(title = "MPG vs Horsepower by Cylinder Count", 
       subtitle = "Faceted by Gear Count",
       x = "Miles Per Gallon", 
       y = "Horsepower")

# Error case
ggplot(non_existent_data, aes(x = mpg, y = hp)) + 
  geom_point()
```

### 6.2 Expected Outputs

We'll create a library of expected outputs for comparison:

- **Reference Images**: Pre-generated images for comparison
- **Error Messages**: Expected error messages for error cases
- **Response Formats**: Expected response formats for various scenarios

## 7. Test Automation

### 7.1 Task Runner Configuration

We'll configure Task to automate test execution:

```yaml
version: '3'

tasks:
  test:
    desc: Run all tests
    cmds:
      - go test -v ./...

  test:unit:
    desc: Run unit tests
    cmds:
      - go test -v -short ./...

  test:integration:
    desc: Run integration tests
    cmds:
      - go test -v -tags=integration ./...

  test:protocol:
    desc: Run protocol conformance tests
    cmds:
      - go test -v -tags=protocol ./...

  test:coverage:
    desc: Generate test coverage report
    cmds:
      - go test -coverprofile=coverage.out ./...
      - go tool cover -html=coverage.out -o coverage.html

  benchmark:
    desc: Run benchmarks
    cmds:
      - go test -bench=. -benchmem ./...
```

### 7.2 VSCode Integration

We'll configure VSCode for test execution and debugging:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "type": "task",
      "taskName": "test",
      "label": "Run All Tests",
      "group": {
        "kind": "test",
        "isDefault": true
      },
      "problemMatcher": ["$go-test"]
    },
    {
      "type": "task",
      "taskName": "test:unit",
      "label": "Run Unit Tests",
      "problemMatcher": ["$go-test"]
    },
    {
      "type": "task",
      "taskName": "test:integration",
      "label": "Run Integration Tests",
      "problemMatcher": ["$go-test"]
    },
    {
      "type": "task",
      "taskName": "test:protocol",
      "label": "Run Protocol Tests",
      "problemMatcher": ["$go-test"]
    },
    {
      "type": "task",
      "taskName": "test:coverage",
      "label": "Generate Coverage Report",
      "problemMatcher": []
    },
    {
      "type": "task",
      "taskName": "benchmark",
      "label": "Run Benchmarks",
      "problemMatcher": []
    }
  ]
}
```

## 8. Test Implementation

### 8.1 Test Structure

We'll organize tests into the following structure:

```
r-server/
├── internal/
│   ├── mcp/
│   │   ├── server_test.go       # Unit tests for MCP server
│   │   ├── tools_test.go        # Unit tests for tools
│   ├── r/
│   │   ├── executor_test.go     # Unit tests for R execution
│   │   ├── ggplot_test.go       # Unit tests for ggplot functionality
│   ├── image/
│   │   ├── processor_test.go    # Unit tests for image processing
├── test/
│   ├── integration/
│   │   ├── mcp_test.go          # Integration tests for MCP
│   │   ├── ggplot_test.go       # Integration tests for ggplot
│   ├── protocol/
│   │   ├── conformance_test.go  # Protocol conformance tests
│   ├── performance/
│   │   ├── benchmark_test.go    # Performance benchmarks
│   ├── testdata/
│   │   ├── r_scripts/           # Sample R scripts
│   │   ├── expected_images/     # Expected output images
```

### 8.2 Test Helpers

We'll create helper functions to simplify test implementation:

```go
// setupTestServer creates a test server with an in-memory transport
func setupTestServer(t testing.TB) (*MCPServer, *mcp.InMemoryTransport) {
    transport := mcp.NewInMemoryTransport()
    server, err := NewMCPServer(transport)
    require.NoError(t, err)
    return server, transport
}

// compareImages compares two images for similarity
func compareImages(t *testing.T, expected, actual []byte) {
    // Image comparison logic
    // ...
}

// loadTestScript loads a test script from testdata
func loadTestScript(t *testing.T, name string) string {
    data, err := os.ReadFile(filepath.Join("testdata", "r_scripts", name))
    require.NoError(t, err)
    return string(data)
}
```

## 9. Test Execution

### 9.1 Local Test Execution

For local development, tests will be executed using:

```bash
# Run all tests
task test

# Run unit tests
task test:unit

# Run integration tests
task test:integration

# Run protocol tests
task test:protocol

# Generate coverage report
task test:coverage

# Run benchmarks
task benchmark
```

### 9.2 Continuous Integration

For CI, tests will be executed automatically on each commit:

```yaml
# Example CI workflow
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - name: Install Task
        run: go install github.com/go-task/task/v3/cmd/task@latest
      - name: Run tests
        run: task test
      - name: Generate coverage report
        run: task test:coverage
      - name: Upload coverage report
        uses: actions/upload-artifact@v3
        with:
          name: coverage-report
          path: coverage.html
```

## 10. Test Reporting

### 10.1 Test Results

Test results will be reported in the following formats:

- **Console Output**: For local development
- **JUnit XML**: For CI integration
- **HTML Reports**: For detailed test results

### 10.2 Coverage Reports

Coverage reports will be generated in the following formats:

- **Console Summary**: For quick feedback
- **HTML Report**: For detailed coverage analysis
- **Coverage Badges**: For README display

## 11. Test-Driven Development Process

We'll follow a test-driven development process:

1. **Write Test**: Write a failing test for the feature
2. **Implement Feature**: Implement the feature to make the test pass
3. **Refactor**: Refactor the code while keeping tests passing
4. **Repeat**: Continue with the next feature

This process ensures that all code is covered by tests and that the implementation meets the requirements.

## 12. Testing Timeline

| Phase | Focus | Timeline |
|-------|-------|----------|
| 1 | Unit test framework setup | Week 1 |
| 2 | Core MCP protocol tests | Week 1-2 |
| 3 | Tool functionality tests | Week 2 |
| 4 | Integration tests | Week 2-3 |
| 5 | Performance and security tests | Week 3 |
| 6 | Test automation and CI setup | Week 3-4 |
| 7 | Final test review and refinement | Week 4 |

## 13. Conclusion

This testing strategy provides a comprehensive approach to testing the R-Server MCP for ggplot project. By following this strategy, we'll ensure that the server correctly implements the MCP protocol, properly renders ggplot2 visualizations, and handles errors gracefully.

The key aspects of this strategy are:

1. **Comprehensive Test Coverage**: Testing all aspects of the system
2. **Multiple Testing Levels**: Unit, integration, and system testing
3. **Automated Testing**: Using Task and CI for test automation
4. **Test-Driven Development**: Writing tests before implementing features
5. **Performance and Security Testing**: Ensuring non-functional requirements are met

By implementing this testing strategy, we'll create a robust, reliable, and secure MCP server for ggplot2 visualizations.
