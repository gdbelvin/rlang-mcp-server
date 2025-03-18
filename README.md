# R-Server MCP (Go Implementation)

This is a Go implementation of the R-Server MCP server that executes R Markdown files. It demonstrates core MCP concepts like resources and tools by allowing:

- Listing R Markdown files as resources
- Reading R Markdown files
- Executing R Markdown files via Docker
- Retrieving rendered output

## Features

- **MCP Resources**: Access R Markdown files and rendered outputs
- **MCP Tools**: Create and render R Markdown files
- **Docker Integration**: Render R Markdown files using Docker
- **Multiple Rendering Methods**: Support for both Docker API and docker-compose

## Project Structure

- `main.go`: Main application entry point
- `utils.go`: Utility functions for R Markdown handling
- `docker_config.go`: Docker configuration
- `dockerode_runner.go`: Docker API integration for rendering
- `docker_compose_runner.go`: Docker Compose integration for rendering
- `mcp_server.go`: MCP server implementation

## Requirements

- Go 1.22 or later
- Docker
- Docker Compose (optional, for docker-compose rendering)

## Building

```bash
go build -o r-server
```

## Running

```bash
./r-server
```

## Docker Setup

The server uses Docker to render R Markdown files. The Docker image is built automatically when rendering an R Markdown file, but you can also build it manually:

```bash
docker build -t r-server-rmd .
```

## MCP Integration

To use this server with an MCP client, you can configure it in your MCP settings file:

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

## Usage

### Creating an R Markdown File

Use the `create_rmd` tool to create a new R Markdown file:

```json
{
  "filename": "example",
  "title": "Example R Markdown",
  "content": "This is an example R Markdown file.\n\n```{r}\nplot(cars)\n```"
}
```

### Rendering an R Markdown File

Use the `render_rmd` tool to render an R Markdown file:

```json
{
  "filename": "example",
  "format": "html",
  "use_docker_compose": false
}
```

## License

MIT
