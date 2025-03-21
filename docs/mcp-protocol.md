# R-Server MCP for ggplot: MCP Protocol Implementation Strategy

## 1. Introduction to MCP

The Model Context Protocol (MCP) is a standardized protocol that enables AI models to interact with external tools and resources. For the R-Server ggplot project, we'll implement an MCP server that provides a tool for rendering statistical visualizations using R's ggplot2 library.

## 2. MCP Protocol Overview

### 2.1 Core Concepts

The MCP protocol defines two primary types of capabilities:

1. **Tools**: Functions that can be called to perform actions
2. **Resources**: Data sources that can be accessed

For our R-Server ggplot project, we'll focus on implementing a single tool (`render_ggplot`) that generates visualizations from R code.

### 2.2 Protocol Structure

MCP uses a JSON-RPC 2.0 based protocol with the following key methods:

| Method | Description |
|--------|-------------|
| `listTools` | Lists available tools provided by the server |
| `callTool` | Calls a specific tool with provided arguments |
| `listResources` | Lists available resources (optional) |
| `readResource` | Reads a specific resource (optional) |

Each method follows the JSON-RPC 2.0 request/response format:

**Request Format:**
```json
{
  "jsonrpc": "2.0",
  "id": "request-id",
  "method": "method-name",
  "params": {
    // Method-specific parameters
  }
}
```

**Response Format:**
```json
{
  "jsonrpc": "2.0",
  "id": "request-id",
  "result": {
    // Method-specific result
  }
}
```

**Error Response Format:**
```json
{
  "jsonrpc": "2.0",
  "id": "request-id",
  "error": {
    "code": -32000,
    "message": "Error message"
  }
}
```

## 3. MCP Implementation Strategy

### 3.1 Library Selection

We'll use the `github.com/metoro-io/mcp-golang` library for implementing the MCP server. This library provides:

- Core MCP protocol implementation
- Transport abstractions (HTTP, stdio)
- Request/response handling
- Tool and resource registration

### 3.2 Server Implementation

The MCP server will be implemented in Go with the following components:

1. **Server Initialization**: Create and configure the MCP server
2. **Tool Registration**: Register the render_ggplot tool
3. **Request Handling**: Process incoming MCP requests
4. **Response Generation**: Generate appropriate responses

```go
// Example server initialization
func NewMCPServer(transport transport.Transport) (*MCPServer, error) {
    server := &MCPServer{
        Server: mcp.NewServer(transport),
    }
    
    // Register the render_ggplot tool
    if err := server.RegisterTool("render_ggplot", "Render a ggplot2 visualization", RenderGGPlot); err != nil {
        return nil, fmt.Errorf("failed to register render_ggplot tool: %w", err)
    }
    
    return server, nil
}
```

### 3.3 Tool Implementation

The `render_ggplot` tool will be implemented as a Go function that:

1. Validates input parameters
2. Executes R code to generate a visualization
3. Processes the resulting image
4. Returns the image data in the requested format

```go
// Example tool implementation
func RenderGGPlot(args GGPlotRenderArgs) (*mcp.ToolResponse, error) {
    // Validate arguments
    if args.Code == "" {
        return nil, fmt.Errorf("code is required")
    }
    
    // Set default values
    if args.OutputType == "" {
        args.OutputType = "png"
    }
    if args.Width == 0 {
        args.Width = 800
    }
    if args.Height == 0 {
        args.Height = 600
    }
    if args.Resolution == 0 {
        args.Resolution = 96
    }
    
    // Execute R code and generate image
    imagePath, err := ExecuteRCode(args)
    if err != nil {
        return nil, fmt.Errorf("failed to execute R code: %w", err)
    }
    
    // Read and process the image
    imageData, err := ProcessImage(imagePath, args.OutputType)
    if err != nil {
        return nil, fmt.Errorf("failed to process image: %w", err)
    }
    
    // Create and return the response
    return mcp.NewToolResponse(mcp.NewImageContent(
        GetMimeType(args.OutputType),
        imageData,
    )), nil
}
```

## 4. MCP Request/Response Formats

### 4.1 listTools Request/Response

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "listTools",
  "params": {}
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "tools": [
      {
        "name": "render_ggplot",
        "description": "Render a ggplot2 visualization",
        "inputSchema": {
          "type": "object",
          "properties": {
            "code": {
              "type": "string",
              "description": "R code containing ggplot2 commands"
            },
            "output_type": {
              "type": "string",
              "enum": ["png", "jpeg", "pdf", "svg"],
              "description": "Output format for the image",
              "default": "png"
            },
            "width": {
              "type": "integer",
              "description": "Width of the output image in pixels",
              "default": 800,
              "minimum": 100,
              "maximum": 5000
            },
            "height": {
              "type": "integer",
              "description": "Height of the output image in pixels",
              "default": 600,
              "minimum": 100,
              "maximum": 5000
            },
            "resolution": {
              "type": "integer",
              "description": "Resolution of the output image in dpi",
              "default": 96,
              "minimum": 72,
              "maximum": 600
            }
          },
          "required": ["code"]
        }
      }
    ]
  }
}
```

### 4.2 callTool Request/Response

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "method": "callTool",
  "params": {
    "name": "render_ggplot",
    "arguments": {
      "code": "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point() + theme_minimal() + labs(title = 'MPG vs Horsepower')",
      "output_type": "png",
      "width": 800,
      "height": 600,
      "resolution": 96
    }
  }
}
```

**Success Response:**
```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "result": {
    "content": [
      {
        "type": "image",
        "mimeType": "image/png",
        "data": "base64-encoded-image-data"
      }
    ]
  }
}
```

**Error Response:**
```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "error": {
    "code": -32603,
    "message": "Error executing R code: object 'non_existent_data' not found"
  }
}
```

## 5. Error Handling Strategy

### 5.1 Error Types

We'll handle the following types of errors:

1. **Protocol Errors**: Invalid JSON, missing required fields, etc.
2. **Validation Errors**: Invalid tool arguments, unsupported formats, etc.
3. **Execution Errors**: R code errors, file system errors, etc.
4. **System Errors**: Out of memory, timeout, etc.

### 5.2 Error Mapping

Errors will be mapped to appropriate JSON-RPC error codes:

| Error Type | JSON-RPC Error Code | Description |
|------------|---------------------|-------------|
| Parse Error | -32700 | Invalid JSON |
| Invalid Request | -32600 | Invalid request object |
| Method Not Found | -32601 | Method not found |
| Invalid Params | -32602 | Invalid method parameters |
| Internal Error | -32603 | Internal server error |
| Server Error | -32000 to -32099 | Custom server errors |

### 5.3 Error Response Format

All error responses will include:

1. The appropriate error code
2. A descriptive error message
3. Additional data when available (e.g., R error details)

```json
{
  "jsonrpc": "2.0",
  "id": "request-id",
  "error": {
    "code": -32602,
    "message": "Invalid parameters: width must be between 100 and 5000",
    "data": {
      "parameter": "width",
      "value": 10000,
      "constraint": "100 <= width <= 5000"
    }
  }
}
```

## 6. Security Considerations

### 6.1 Input Validation

All input parameters will be strictly validated:

1. **R Code**: Validate for potential security issues
2. **Image Dimensions**: Enforce reasonable limits
3. **Output Format**: Restrict to supported formats

### 6.2 Resource Limits

Resource usage will be limited to prevent abuse:

1. **Execution Time**: Limit R code execution time
2. **Memory Usage**: Limit memory consumption
3. **Disk Usage**: Limit temporary file storage

### 6.3 Error Information

Error messages will be informative but will not expose sensitive information:

1. **R Errors**: Include R error messages but filter out system paths
2. **System Errors**: Provide generic messages for system-level errors
3. **Stack Traces**: Never include stack traces in responses

## 7. Testing Strategy

### 7.1 Protocol Conformance Tests

We'll implement tests to verify protocol conformance:

1. **Request Validation**: Test handling of valid and invalid requests
2. **Response Format**: Verify response format correctness
3. **Error Handling**: Test error response format and codes

### 7.2 Tool Functionality Tests

We'll test the `render_ggplot` tool functionality:

1. **Parameter Validation**: Test handling of valid and invalid parameters
2. **R Code Execution**: Test execution of various R code snippets
3. **Image Generation**: Verify correct image generation
4. **Format Support**: Test all supported output formats

### 7.3 Integration Tests

We'll implement integration tests that simulate real-world usage:

1. **End-to-End Tests**: Test the complete request-response cycle
2. **Client Simulation**: Simulate MCP client behavior
3. **Error Scenarios**: Test handling of various error scenarios

## 8. Implementation Roadmap

### 8.1 Phase 1: Core Protocol Implementation

1. Set up the basic MCP server structure
2. Implement the `listTools` method
3. Create the tool schema for `render_ggplot`
4. Implement basic request/response handling

### 8.2 Phase 2: Tool Implementation

1. Implement R code execution
2. Implement image processing
3. Implement the `callTool` method for `render_ggplot`
4. Add parameter validation and error handling

### 8.3 Phase 3: Testing and Refinement

1. Implement protocol conformance tests
2. Implement tool functionality tests
3. Implement integration tests
4. Refine error handling and security measures

### 8.4 Phase 4: Documentation and Finalization

1. Document the API
2. Create usage examples
3. Optimize performance
4. Prepare for deployment

## 9. Example Implementation

### 9.1 Server Initialization

```go
package main

import (
    "fmt"
    "os"
    
    mcp "github.com/metoro-io/mcp-golang"
    "github.com/metoro-io/mcp-golang/transport"
    "github.com/metoro-io/mcp-golang/transport/http"
)

func main() {
    // Create HTTP transport
    httpTransport := http.NewHTTPTransport("/mcp").WithAddr(":22011")
    
    // Create and configure the server
    server, err := NewMCPServer(httpTransport)
    if err != nil {
        fmt.Printf("Error creating MCP server: %v\n", err)
        os.Exit(1)
    }
    
    // Start the server
    fmt.Println("Starting MCP server on port 22011")
    if err := server.Serve(); err != nil {
        fmt.Printf("Error starting MCP server: %v\n", err)
        os.Exit(1)
    }
}
```

### 9.2 Tool Registration

```go
func NewMCPServer(transport transport.Transport) (*MCPServer, error) {
    // Create a new MCP server
    server := &MCPServer{
        Server: mcp.NewServer(transport),
    }
    
    // Register the render_ggplot tool
    if err := server.RegisterTool("render_ggplot", "Render a ggplot2 visualization", RenderGGPlot); err != nil {
        return nil, fmt.Errorf("failed to register render_ggplot tool: %w", err)
    }
    
    return server, nil
}
```

### 9.3 Tool Implementation

```go
// GGPlotRenderArgs represents the arguments for rendering a ggplot image
type GGPlotRenderArgs struct {
    Code       string `json:"code" jsonschema:"required,description=R code containing ggplot commands"`
    OutputType string `json:"output_type" jsonschema:"description=Output format (png, jpeg, pdf, svg)"`
    Width      int    `json:"width" jsonschema:"description=Width of the output image in pixels"`
    Height     int    `json:"height" jsonschema:"description=Height of the output image in pixels"`
    Resolution int    `json:"resolution" jsonschema:"description=Resolution of the output image in dpi"`
}

// RenderGGPlot renders a ggplot2 visualization
func RenderGGPlot(args GGPlotRenderArgs) (*mcp.ToolResponse, error) {
    // Validate and set default values
    if args.Code == "" {
        return nil, fmt.Errorf("code is required")
    }
    
    outputType := args.OutputType
    if outputType == "" {
        outputType = "png"
    }
    
    width := args.Width
    if width == 0 {
        width = 800
    } else if width < 100 || width > 5000 {
        return nil, fmt.Errorf("width must be between 100 and 5000")
    }
    
    height := args.Height
    if height == 0 {
        height = 600
    } else if height < 100 || height > 5000 {
        return nil, fmt.Errorf("height must be between 100 and 5000")
    }
    
    resolution := args.Resolution
    if resolution == 0 {
        resolution = 96
    } else if resolution < 72 || resolution > 600 {
        return nil, fmt.Errorf("resolution must be between 72 and 600")
    }
    
    // Execute R code and generate image
    imagePath, err := ExecuteRCode(args.Code, outputType, width, height, resolution)
    if err != nil {
        return nil, fmt.Errorf("failed to execute R code: %w", err)
    }
    
    // Read and process the image
    imageData, err := ProcessImage(imagePath)
    if err != nil {
        return nil, fmt.Errorf("failed to process image: %w", err)
    }
    
    // Create and return the response
    return mcp.NewToolResponse(mcp.NewImageContent(
        GetMimeType(outputType),
        imageData,
    )), nil
}
```

## 10. Conclusion

This MCP protocol implementation strategy provides a comprehensive plan for implementing an MCP server that renders ggplot2 visualizations. By following this strategy, we'll create a robust, secure, and efficient MCP server that conforms to the protocol specifications and meets the project requirements.

The implementation will focus on:

1. **Protocol Conformance**: Ensuring full compliance with the MCP protocol
2. **Tool Functionality**: Implementing a powerful and flexible ggplot2 rendering tool
3. **Security**: Implementing strict validation and resource limits
4. **Error Handling**: Providing clear and informative error messages
5. **Testing**: Verifying protocol conformance and tool functionality

This strategy will guide the implementation process and ensure that the resulting MCP server meets the project goals and requirements.
