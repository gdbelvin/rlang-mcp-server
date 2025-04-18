# R-Server MCP API Documentation

This document provides a comprehensive reference for the Model Context Protocol (MCP) implementation in R-Server, including available resources, resource templates, and tools.

## Table of Contents

- [R-Server MCP API Documentation](#r-server-mcp-api-documentation)
  - [Table of Contents](#table-of-contents)
  - [Introduction to MCP](#introduction-to-mcp)
    - [MCP Protocol Overview](#mcp-protocol-overview)
      - [Core MCP Request Types](#core-mcp-request-types)
      - [MCP Server Implementation](#mcp-server-implementation)
      - [MCP Request and Response Formats](#mcp-request-and-response-formats)
    - [R-Server MCP Implementation](#r-server-mcp-implementation)
  - [MCP Resources](#mcp-resources)
    - [Resource URIs](#resource-uris)
    - [Available Resources](#available-resources)
      - [R Markdown Files](#r-markdown-files)
      - [Rendered Outputs](#rendered-outputs)
  - [MCP Tools](#mcp-tools)
    - [render\_ggplot](#render_ggplot)
      - [Input Schema](#input-schema)
      - [Example Input](#example-input)
      - [Response](#response)
      - [Implementation Details](#implementation-details)
    - [execute\_r\_script](#execute_r_script)
      - [Input Schema](#input-schema-1)
      - [Example Input](#example-input-1)
      - [Response](#response-1)
      - [Implementation Details](#implementation-details-1)
    - [create\_rmd](#create_rmd)
      - [Input Schema](#input-schema-2)
      - [Example Input](#example-input-2)
      - [Response](#response-2)
      - [Implementation Details](#implementation-details-2)
    - [render\_rmd](#render_rmd)
      - [Input Schema](#input-schema-3)
      - [Example Input](#example-input-3)
      - [Response](#response-3)
      - [Implementation Details](#implementation-details-3)
  - [Implementation Details](#implementation-details-4)
    - [Server Architecture](#server-architecture)
      - [MCP Protocol Implementation Details](#mcp-protocol-implementation-details)
    - [Docker Integration](#docker-integration)
      - [Docker API (Dockerode)](#docker-api-dockerode)
      - [Docker Compose](#docker-compose)
  - [Usage Examples](#usage-examples)
    - [Executing an R Script](#executing-an-r-script)
    - [Creating and Rendering an R Markdown File](#creating-and-rendering-an-r-markdown-file)
    - [Using Docker Compose for Rendering](#using-docker-compose-for-rendering)
  - [Error Handling](#error-handling)
    - [Common Error Scenarios](#common-error-scenarios)
    - [Error Response Example](#error-response-example)
    - [Error Handling Implementation](#error-handling-implementation)

## Introduction to MCP

The Model Context Protocol (MCP) is a standardized protocol for communication between AI models and external services. It enables AI models to access external data sources and execute operations through a consistent interface.

### MCP Protocol Overview

MCP defines a structured way for AI models to interact with external systems through:

1. **Resources**: Data sources that can be accessed by the model
2. **Tools**: Operations that can be executed by the model

The protocol uses a request-response pattern with JSON-based messages. Each request has a specific schema and corresponding response format.

#### Core MCP Request Types

- **ListResources**: Lists available resources from the server
- **ReadResource**: Reads the content of a specific resource
- **ListResourceTemplates**: Lists available resource templates (for dynamic resources)
- **ListTools**: Lists available tools provided by the server
- **CallTool**: Calls a specific tool with provided arguments

#### MCP Server Implementation

An MCP server typically:
1. Defines its capabilities (resources and tools)
2. Sets up request handlers for each MCP request type
3. Connects to a transport (e.g., stdio, WebSocket)
4. Processes requests and returns responses

#### MCP Request and Response Formats

MCP uses JSON-RPC 2.0 as its underlying protocol. Each request and response follows this format:

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

Examples of specific MCP requests and responses:

**ListResources Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "listResources",
  "params": {}
}
```

**ListResources Response:**
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "resources": [
      {
        "uri": "rmd:///example.Rmd",
        "mimeType": "text/markdown",
        "name": "Example R Markdown",
        "description": "R Markdown file: Example R Markdown"
      }
    ]
  }
}
```

**CallTool Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "method": "callTool",
  "params": {
    "name": "render_rmd",
    "arguments": {
      "filename": "example",
      "format": "html"
    }
  }
}
```

**CallTool Response:**
```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Successfully rendered example.Rmd to example.html"
      }
    ]
  }
}
```

### R-Server MCP Implementation

R-Server implements the MCP protocol to provide access to R Markdown files and rendering capabilities. The server exposes:

- **Resources**: R Markdown files and rendered outputs
- **Tools**: Operations to create and render R Markdown files

## MCP Resources

### Resource URIs

R-Server uses the following URI schemes for resources:

- `rmd:///filename.Rmd` - Access to R Markdown source files
- `rmd-output:///filename.html` - Access to rendered HTML output
- `rmd-output:///filename.pdf` - Access to rendered PDF output
- `rmd-output:///filename.docx` - Access to rendered Word document output

### Available Resources

The server dynamically lists available resources based on the R Markdown files in the `rmd` directory and rendered outputs in the `rmd/output` directory.

#### R Markdown Files

Each R Markdown file is exposed as a resource with:

- **URI**: `rmd:///filename.Rmd`
- **MIME Type**: `text/markdown`
- **Name**: Title extracted from the R Markdown front matter, or the filename if no title is found
- **Description**: A description of the R Markdown file

#### Rendered Outputs

Each rendered output file is exposed as a resource with:

- **URI**: `rmd-output:///filename.ext`
- **MIME Type**: Appropriate MIME type based on the file extension:
  - HTML: `text/html`
  - PDF: `application/pdf`
  - DOCX: `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- **Name**: `Rendered: filename.ext`
- **Description**: `Rendered output: filename.ext`

## MCP Tools

R-Server provides the following tools through the MCP interface:

### render_ggplot

Renders a ggplot2 visualization.

#### Input Schema

```json
{
  "type": "object",
  "properties": {
    "code": {
      "type": "string",
      "description": "R code containing ggplot commands"
    },
    "output_type": {
      "type": "string",
      "description": "Output format (png, jpeg, pdf, svg)",
      "default": "png"
    },
    "width": {
      "type": "integer",
      "description": "Width of the output image in pixels",
      "default": 800
    },
    "height": {
      "type": "integer",
      "description": "Height of the output image in pixels",
      "default": 600
    },
    "resolution": {
      "type": "integer",
      "description": "Resolution of the output image in dpi",
      "default": 96
    }
  },
  "required": ["code"]
}
```

#### Example Input

```json
{
  "code": "ggplot(mtcars, aes(x = wt, y = mpg)) + geom_point() + theme_minimal()",
  "output_type": "png",
  "width": 800,
  "height": 600,
  "resolution": 96
}
```

#### Response

An image of the rendered ggplot visualization.

#### Implementation Details

- The R code must include ggplot2 commands
- The output format can be:
  - PNG (default)
  - JPEG
  - PDF
  - SVG
- The width, height, and resolution parameters control the size and quality of the output image

### execute_r_script

Executes an R script and returns the result as text.

#### Input Schema

```json
{
  "type": "object",
  "properties": {
    "code": {
      "type": "string",
      "description": "R code to execute"
    }
  },
  "required": ["code"]
}
```

#### Example Input

```json
{
  "code": "data <- data.frame(x = 1:10, y = (1:10)^2)\nsummary(data)\ncor(data$x, data$y)"
}
```

#### Response

The text output of the executed R script.

#### Implementation Details

- The R code can include any valid R commands
- The output is captured and returned as text
- The R environment includes common packages like ggplot2, dplyr, etc.
- The execution is performed in a temporary directory that is cleaned up after execution

### create_rmd

Creates a new R Markdown file.

#### Input Schema

```json
{
  "type": "object",
  "properties": {
    "filename": {
      "type": "string",
      "description": "Filename for the R Markdown file (without extension)"
    },
    "title": {
      "type": "string",
      "description": "Title for the R Markdown document"
    },
    "content": {
      "type": "string",
      "description": "Content of the R Markdown file"
    }
  },
  "required": ["filename", "title", "content"]
}
```

#### Example Input

```json
{
  "filename": "example",
  "title": "Example R Markdown",
  "content": "This is an example R Markdown file.\n\n```{r}\nplot(cars)\n```"
}
```

#### Response

A success message indicating the file was created:

```
Created R Markdown file: example.Rmd
```

#### Implementation Details

- The filename will automatically have the `.Rmd` extension added if not provided
- If the content doesn't include YAML front matter, it will be automatically added with the provided title
- The file is saved to the `rmd` directory

### render_rmd

Renders an R Markdown file to HTML, PDF, or Word format.

#### Input Schema

```json
{
  "type": "object",
  "properties": {
    "filename": {
      "type": "string",
      "description": "Filename of the R Markdown file to render"
    },
    "format": {
      "type": "string",
      "enum": ["html", "pdf", "word"],
      "description": "Output format (html, pdf, or word)",
      "default": "html"
    },
    "use_docker_compose": {
      "type": "boolean",
      "description": "Whether to use docker-compose (true) or Dockerode (false)",
      "default": false
    }
  },
  "required": ["filename"]
}
```

#### Example Input

```json
{
  "filename": "example",
  "format": "html",
  "use_docker_compose": false
}
```

#### Response

A success message indicating the file was rendered:

```
Successfully rendered example.Rmd to example.html
```

#### Implementation Details

- The filename will automatically have the `.Rmd` extension added if not provided
- The file must exist in the `rmd` directory
- The rendering is performed using Docker, with two possible methods:
  - Docker API (Dockerode) - default
  - Docker Compose - optional
- The rendered output is saved to the `rmd/output` directory
- The output format can be:
  - HTML (default)
  - PDF
  - Word (DOCX)

## Implementation Details

### Server Architecture

R-Server is implemented in Go and follows a modular architecture:

- **MCPServer**: The main server class that implements the MCP protocol
- **Resource Handling**: Lists and reads R Markdown files and rendered outputs
- **Tool Handling**: Implements the `render_ggplot` and `execute_r_script` tools
- **Docker Integration**: Uses Docker to render R Markdown files

#### MCP Protocol Implementation Details

The R-Server implements the MCP protocol through the following components:

| MCP Request Type | R-Server Implementation | Description |
|-----------------|-------------------------|-------------|
| `listResources` | `MCPServer.ListResources()` | Lists available R Markdown files and rendered outputs |
| `readResource` | `MCPServer.ReadResource()` | Reads the content of an R Markdown file or rendered output |
| `listTools` | `MCPServer.ListTools()` | Lists available tools for creating and rendering R Markdown files |
| `callTool` | `MCPServer.CallTool()` | Executes a tool with the provided arguments |

The server handles these requests through the following process:

1. **Request Parsing**: Parses the incoming JSON-RPC request
2. **Request Routing**: Routes the request to the appropriate handler based on the method
3. **Request Handling**: Executes the requested operation
4. **Response Generation**: Generates a JSON-RPC response with the result or error

For example, when a `listResources` request is received:

1. The server calls `MCPServer.ListResources()`
2. This method scans the `rmd` directory for R Markdown files and the `output` directory for rendered outputs
3. It constructs resource objects for each file with appropriate URIs, MIME types, names, and descriptions
4. It returns these resources in the response

Similarly, when a `callTool` request for `execute_r_script` is received:

1. The server calls `MCPServer.CallTool("execute_r_script", args)`
2. This method validates the arguments and ensures the R code is provided
3. It creates a temporary directory for the R script and output
4. It executes the R script and captures the output
5. It returns the output as text content

### Docker Integration

R-Server uses Docker to render R Markdown files, providing a consistent environment for R and its dependencies. Two methods are supported:

#### Docker API (Dockerode)

- Uses the Docker API to create and run containers
- Provides direct control over container creation and execution
- Implemented in `dockerode_runner.go`

#### Docker Compose

- Uses docker-compose to manage containers
- Provides a more declarative approach to container configuration
- Implemented in `docker_compose_runner.go`

Both methods use the same Docker image, which includes:

- R 4.4.3 (from rocker/r-ver)
- Pandoc (for document conversion)
- R packages: rmarkdown, knitr, tinytex, ggplot2
- A custom entrypoint script that handles rendering

## Usage Examples

### Executing an R Script

To execute an R script and get the result:

```json
{
  "name": "execute_r_script",
  "arguments": {
    "code": "# Create a data frame\ndata <- data.frame(x = 1:10, y = (1:10)^2)\n\n# Calculate summary statistics\nsummary_stats <- summary(data)\n\n# Print the summary statistics\nprint(summary_stats)\n\n# Calculate correlation\ncorrelation <- cor(data$x, data$y)\ncat(\"Correlation between x and y:\", correlation, \"\\n\")\n\n# Create a linear model\nmodel <- lm(y ~ x, data = data)\ncat(\"\\nLinear model summary:\\n\")\nprint(summary(model))"
  }
}
```

This will execute the R script and return the output as text.

### Creating and Rendering an R Markdown File

1. Create a new R Markdown file:

```json
{
  "name": "create_rmd",
  "arguments": {
    "filename": "example",
    "title": "Example R Markdown",
    "content": "---\ntitle: \"Example R Markdown\"\nauthor: \"R-Server\"\ndate: \"`r Sys.Date()`\"\noutput: html_document\n---\n\n## R Markdown Example\n\nThis is an example R Markdown file.\n\n```{r}\nplot(cars)\n```\n\n## Summary Statistics\n\n```{r}\nsummary(cars)\n```"
  }
}
```

2. Render the R Markdown file to HTML:

```json
{
  "name": "render_rmd",
  "arguments": {
    "filename": "example",
    "format": "html"
  }
}
```

3. Access the rendered HTML:

```
Resource URI: rmd-output:///example.html
```

### Using Docker Compose for Rendering

To use Docker Compose instead of the Docker API:

```json
{
  "name": "render_rmd",
  "arguments": {
    "filename": "example",
    "format": "pdf",
    "use_docker_compose": true
  }
}
```

This will render the R Markdown file to PDF using Docker Compose.

## Error Handling

The R-Server MCP implementation includes comprehensive error handling to provide clear feedback when operations fail.

### Common Error Scenarios

| Error Scenario | Error Code | Description |
|----------------|------------|-------------|
| Invalid URI format | InvalidRequest | The URI format does not match the expected pattern (`rmd:///` or `rmd-output:///`) |
| Resource not found | NotFound | The requested R Markdown file or rendered output does not exist |
| Invalid tool name | MethodNotFound | The requested tool is not supported by the server |
| Missing required arguments | InvalidParams | Required arguments for a tool are missing or invalid |
| R script execution failure | InternalError | The R script execution process failed |
| Rendering failure | InternalError | The R Markdown rendering process failed |
| Docker error | InternalError | An error occurred while interacting with Docker |

### Error Response Example

When an error occurs, the server returns a JSON-RPC error response:

```json
{
  "jsonrpc": "2.0",
  "id": "request-id",
  "error": {
    "code": -32602,
    "message": "Invalid params: filename is required"
  }
}
```

### Error Handling Implementation

The R-Server implements error handling through:

1. **Input Validation**: Validates all input parameters before executing operations
2. **Resource Checking**: Verifies that requested resources exist before attempting to access them
3. **Docker Error Handling**: Captures and reports errors from Docker operations
4. **Detailed Error Messages**: Provides specific error messages to help diagnose issues

For example, when executing an R script:

1. The server checks that the R code is provided
2. It creates a temporary directory for the script and output
3. It executes the R script and captures any errors
4. It returns a detailed error message if the execution fails

Similarly, when rendering a ggplot visualization:

1. The server checks that the R code is provided
2. It validates the output format, width, height, and resolution parameters
3. It executes the R script to generate the visualization
4. It captures any errors during the rendering process
5. It returns a detailed error message if the rendering fails
