# R-Server MCP for ggplot

A specialized Model Context Protocol (MCP) server that enables AI models to generate data visualizations using R's ggplot2 library.

## Overview

This MCP server provides a streamlined interface for creating statistical visualizations without requiring direct access to an R environment. It exposes a single MCP tool (`render_ggplot`) that generates visualizations from R code containing ggplot2 commands.

## Features

- **ggplot2 Rendering**: Execute R code containing ggplot2 commands and return the resulting visualization
- **Format Options**: Support for PNG, JPEG, PDF, and SVG output formats
- **Customization**: Control image dimensions and resolution
- **Error Handling**: Clear error messages for invalid R code or rendering failures
- **MCP Protocol Compliance**: Full implementation of the Model Context Protocol
- **Docker Integration**: Secure execution of R code in isolated containers

## Requirements

- Go 1.22 or later
- R 4.0 or later with ggplot2 package
- Docker (for containerized execution)

## Building

### Local Build

```bash
# Build the server
task build

# Run tests
task test
```

### Docker Build

```bash
# Build the Docker image
task docker:build

# Run the server in Docker
task docker:run
```

## Usage

### Starting the Server

```bash
# Start the server on the default port (22011)
./r-server

# Start the server on a custom port
./r-server --port 8080
```

### MCP Integration

To use this server with an MCP client, configure it in your MCP settings file:

```json
{
  "mcpServers": {
    "r-server": {
      "command": "/path/to/r-server",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

### Using the render_ggplot Tool

The `render_ggplot` tool accepts the following parameters:

- `code` (required): R code containing ggplot2 commands
- `output_type`: Output format (png, jpeg, pdf, svg), default: png
- `width`: Width of the output image in pixels, default: 800
- `height`: Height of the output image in pixels, default: 600
- `resolution`: Resolution of the output image in dpi, default: 96

Example request:

```json
{
  "name": "render_ggplot",
  "arguments": {
    "code": "ggplot(mtcars, aes(x = mpg, y = hp)) + geom_point() + theme_minimal() + labs(title = 'MPG vs Horsepower')",
    "output_type": "png",
    "width": 800,
    "height": 600,
    "resolution": 96
  }
}
```

## Development

### Project Structure

```
r-server/
├── cmd/
│   └── r-server/         # Main application entry point
├── docs/                 # Documentation
├── internal/
│   ├── mcp/              # MCP server implementation
│   ├── r/                # R code execution
│   └── image/            # Image processing
├── test/
│   ├── integration/      # Integration tests
│   ├── protocol/         # Protocol conformance tests
│   └── testdata/         # Test data
├── tools/
│   └── mcp-validator/    # MCP protocol validator
├── Dockerfile            # Docker configuration
├── Taskfile.yml          # Task runner configuration
└── README.md             # This file
```

### Task Runner

This project uses [Task](https://taskfile.dev/) for automation. Available tasks:

```bash
# List all available tasks
task

# Build the server
task build

# Run tests
task test

# Run integration tests
task test:integration

# Run protocol conformance tests
task test:protocol

# Generate test coverage report
task coverage

# Lint the code
task lint

# Run the full CI pipeline
task ci
```

### Testing

The project includes several types of tests:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test components working together
- **Protocol Conformance Tests**: Test MCP protocol implementation
- **Performance Tests**: Measure performance characteristics

To run the tests:

```bash
# Run all tests
task test

# Run unit tests
task test:unit

# Run integration tests
task test:integration

# Run protocol conformance tests
task test:protocol

# Generate coverage report
task coverage
```

## License

MIT
