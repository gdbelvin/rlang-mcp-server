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

The MCP client will automatically communicate with the server using stdio transport, which is the recommended approach for stability and reliability. The dockerized version maintains this communication pattern while providing isolation and dependency management.


## License

Creative Commons Attribution-NonCommercial 4.0 International (CC-BY-NC 4.0)

This work is licensed under a [Creative Commons Attribution-NonCommercial 4.0 International License](http://creativecommons.org/licenses/by-nc/4.0/).

See the [LICENSE](LICENSE) file for details.
