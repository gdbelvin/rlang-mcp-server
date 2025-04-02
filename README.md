# R-Server MCP

A specialized Model Context Protocol (MCP) server that enables AI models to generate data visualizations using R's ggplot2 library and execute R scripts.

## Overview

This MCP server provides a streamlined interface for creating statistical visualizations and executing R scripts without requiring direct access to an R environment. It exposes two MCP tools:
- `render_ggplot`: Generates visualizations from R code containing ggplot2 commands
- `execute_r_script`: Executes any R script and returns the text output

## Features

- **ggplot2 Rendering**: Execute R code containing ggplot2 commands and return the resulting visualization
- **R Script Execution**: Execute any R script and return the text output
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

### Using Docker with stdin/stdout

The server can be run in Docker while preserving stdin/stdout communication, which is essential for MCP:

```bash
# Build and run using docker-compose
./start_server.sh --docker
```

Or set an environment variable:

```bash
USE_DOCKER=true ./start_server.sh
```

This approach ensures that stdin and stdout are properly connected between the host and the container, allowing seamless MCP communication.

## Usage

### Starting the Server

The server can be started in two modes: stdio transport (recommended) or HTTP transport.

#### Stdio Transport

```bash
# Start the server with stdio transport
./r-server
```

Using stdio transport is recommended as it avoids potential deadlock issues that can occur with HTTP transport. The server reads from stdin and writes to stdout, making it suitable for use with named pipes or direct process communication.

#### Using Named Pipes

You can use named pipes to connect the server with a client:

```bash
# Create named pipes
mkfifo server_in server_out

# Start the server in the background
./r-server < server_in > server_out &

# Start the client (in another terminal or script)
./client < server_out > server_in
```

A test script is provided to demonstrate this setup:

```bash
# Run the test script
./test_stdio_transport.sh
```

### MCP Integration

To use this server with an MCP client, configure it in your MCP settings file:

#### Local Execution

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

#### Docker Execution

```json
{
  "mcpServers": {
    "r-server": {
      "command": "/path/to/start_server.sh",
      "args": ["--docker"],
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

Or using environment variables:

```json
{
  "mcpServers": {
    "r-server": {
      "command": "/path/to/start_server.sh",
      "env": {
        "USE_DOCKER": "true"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

The MCP client will automatically communicate with the server using stdio transport, which is the recommended approach for stability and reliability. The dockerized version maintains this communication pattern while providing isolation and dependency management.

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

Creative Commons Attribution-NonCommercial 4.0 International (CC-BY-NC 4.0)

This work is licensed under a [Creative Commons Attribution-NonCommercial 4.0 International License](http://creativecommons.org/licenses/by-nc/4.0/).

See the [LICENSE](LICENSE) file for details.
