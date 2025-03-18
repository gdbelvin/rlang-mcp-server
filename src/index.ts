#!/usr/bin/env node

/**
 * R-Server MCP server that executes R Markdown files.
 * It demonstrates core MCP concepts like resources and tools by allowing:
 * - Listing R Markdown files as resources
 * - Reading R Markdown files
 * - Executing R Markdown files via Docker
 * - Retrieving rendered output
 */

import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListResourcesRequestSchema,
  ListToolsRequestSchema,
  ReadResourceRequestSchema,
  ErrorCode,
  McpError,
} from "@modelcontextprotocol/sdk/types.js";
import * as fs from "fs";
import * as path from "path";
import { promisify } from "util";
import { exec } from "child_process";
import { renderWithDockerode } from "./dockerode-runner.js";
import { renderWithDockerCompose } from "./docker-compose-runner.js";
import {
  RMarkdownFile,
  extractRmdTitle,
  getMimeType,
  ensureRmdExtension,
  createRmdFrontMatter
} from "./utils.js";

const execPromise = promisify(exec);

// Base directory for R Markdown files
const RMD_DIR = path.resolve(process.cwd());
const OUTPUT_DIR = path.resolve(RMD_DIR, "output");

// Ensure directories exist
if (!fs.existsSync(RMD_DIR)) {
  fs.mkdirSync(RMD_DIR, { recursive: true });
}
if (!fs.existsSync(OUTPUT_DIR)) {
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

console.error(`Using RMD_DIR: ${RMD_DIR}`);
console.error(`Using OUTPUT_DIR: ${OUTPUT_DIR}`);

/**
 * Get all R Markdown files in the RMD_DIR
 */
function getRMarkdownFiles(): RMarkdownFile[] {
  const files = fs.readdirSync(RMD_DIR);
  return files
    .filter(file => file.endsWith(".Rmd") || file.endsWith(".rmd"))
    .map(filename => {
      const filePath = path.join(RMD_DIR, filename);
      const content = fs.readFileSync(filePath, "utf-8");
      
      // Extract title from YAML front matter
      const title = extractRmdTitle(content) || filename;
      
      return {
        filename,
        title,
        path: filePath
      };
    });
}

/**
 * Function to build the Docker image for R Markdown rendering
 */
async function buildDockerImage(): Promise<void> {
  try {
    console.error("Building Docker image for R Markdown rendering...");
    const { stdout, stderr } = await execPromise("docker build -t r-server-rmd .");
    console.error("Docker image built successfully");
  } catch (error: any) {
    console.error("Error building Docker image:", error);
    throw new McpError(
      ErrorCode.InternalError,
      `Failed to build Docker image for R Markdown rendering: ${error.message || String(error)}`
    );
  }
}

/**
 * Function to render an R Markdown file using Docker
 * 
 * @param filename Filename of the R Markdown file to render
 * @param outputFormat Output format (html, pdf, or word)
 * @param useDockerCompose Whether to use docker-compose (true) or Dockerode (false)
 * @returns Output filename
 */
async function renderRMarkdown(
  filename: string,
  outputFormat: "html" | "pdf" | "word" = "html",
  useDockerCompose: boolean = false
): Promise<string> {
  try {
    // Choose the appropriate Docker runner based on the useDockerCompose flag
    if (useDockerCompose) {
      console.error(`Rendering ${filename} using docker-compose`);
      return await renderWithDockerCompose(RMD_DIR, filename, outputFormat);
    } else {
      console.error(`Rendering ${filename} using Dockerode`);
      return await renderWithDockerode(RMD_DIR, filename, outputFormat);
    }
  } catch (error: any) {
    console.error("Error rendering R Markdown:", error);
    throw new McpError(
      ErrorCode.InternalError,
      `Failed to render R Markdown file: ${error.message || String(error)}`
    );
  }
}

/**
 * Create an MCP server with capabilities for resources and tools
 */
const server = new Server(
  {
    name: "R-Server",
    version: "0.1.0",
  },
  {
    capabilities: {
      resources: {},
      tools: {},
    },
  }
);

/**
 * Handler for listing available R Markdown files as resources
 */
server.setRequestHandler(ListResourcesRequestSchema, async () => {
  const rmdFiles = getRMarkdownFiles();
  
  return {
    resources: [
      ...rmdFiles.map(file => ({
        uri: `rmd:///${file.filename}`,
        mimeType: "text/markdown",
        name: file.title,
        description: `R Markdown file: ${file.title}`
      })),
      ...fs.readdirSync(OUTPUT_DIR)
        .filter(file => file.endsWith(".html") || file.endsWith(".pdf") || file.endsWith(".docx"))
        .map(file => ({
          uri: `rmd-output:///${file}`,
          mimeType: getMimeType(file),
          name: `Rendered: ${file}`,
          description: `Rendered output: ${file}`
        }))
    ]
  };
});

/**
 * Handler for reading R Markdown files and rendered outputs
 */
server.setRequestHandler(ReadResourceRequestSchema, async (request) => {
  const url = new URL(request.params.uri);
  const scheme = url.protocol.replace(/:$/, "");
  const filename = url.pathname.replace(/^\//, "");
  
  if (scheme === "rmd") {
    // Reading an R Markdown file
    const filePath = path.join(RMD_DIR, filename);
    
    if (!fs.existsSync(filePath)) {
      throw new McpError(
        ErrorCode.InvalidRequest,
        `R Markdown file not found: ${filename}`
      );
    }
    
    const content = fs.readFileSync(filePath, "utf-8");
    
    return {
      contents: [{
        uri: request.params.uri,
        mimeType: "text/markdown",
        text: content
      }]
    };
  } else if (scheme === "rmd-output") {
    // Reading a rendered output file
    const filePath = path.join(OUTPUT_DIR, filename);
    
    if (!fs.existsSync(filePath)) {
      throw new McpError(
        ErrorCode.InvalidRequest,
        `Rendered output file not found: ${filename}`
      );
    }
    
    const mimeType = getMimeType(filename);
    const content = fs.readFileSync(filePath, "utf-8");
    
    return {
      contents: [{
        uri: request.params.uri,
        mimeType,
        text: content
      }]
    };
  } else {
    throw new McpError(
      ErrorCode.InvalidRequest,
      `Unsupported URI scheme: ${scheme}`
    );
  }
});

/**
 * Handler that lists available tools for R Markdown
 */
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "create_rmd",
        description: "Create a new R Markdown file",
        inputSchema: {
          type: "object",
          properties: {
            filename: {
              type: "string",
              description: "Filename for the R Markdown file (without extension)"
            },
            title: {
              type: "string",
              description: "Title for the R Markdown document"
            },
            content: {
              type: "string",
              description: "Content of the R Markdown file"
            }
          },
          required: ["filename", "title", "content"]
        }
      },
      {
        name: "render_rmd",
        description: "Render an R Markdown file",
        inputSchema: {
          type: "object",
          properties: {
            filename: {
              type: "string",
              description: "Filename of the R Markdown file to render"
            },
            format: {
              type: "string",
              enum: ["html", "pdf", "word"],
              description: "Output format (html, pdf, or word)",
              default: "html"
            },
            use_docker_compose: {
              type: "boolean",
              description: "Whether to use docker-compose (true) or Dockerode (false)",
              default: false
            }
          },
          required: ["filename"]
        }
      }
    ]
  };
});

/**
 * Handler for R Markdown tools
 */
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  switch (request.params.name) {
    case "create_rmd": {
      const args = request.params.arguments as any;
      const filename = String(args?.filename);
      const title = String(args?.title);
      const content = String(args?.content);
      
      if (!filename || !title || !content) {
        throw new McpError(
          ErrorCode.InvalidParams,
          "Filename, title, and content are required"
        );
      }
      
      // Ensure filename has .Rmd extension
      const fullFilename = ensureRmdExtension(filename);
      const filePath = path.join(RMD_DIR, fullFilename);
      
      // Create YAML front matter if not present
      const finalContent = createRmdFrontMatter(title, content);
      
      fs.writeFileSync(filePath, finalContent, "utf-8");
      
      return {
        content: [{
          type: "text",
          text: `Created R Markdown file: ${fullFilename}`
        }]
      };
    }
    
    case "render_rmd": {
      const args = request.params.arguments as any;
      const filename = String(args?.filename);
      const format = (args?.format || "html") as "html" | "pdf" | "word";
      const useDockerCompose = Boolean(args?.use_docker_compose);
      
      if (!filename) {
        throw new McpError(
          ErrorCode.InvalidParams,
          "Filename is required"
        );
      }
      
      // Ensure filename has .Rmd extension
      const fullFilename = ensureRmdExtension(filename);
      const filePath = path.join(RMD_DIR, fullFilename);
      
      if (!fs.existsSync(filePath)) {
        throw new McpError(
          ErrorCode.InvalidRequest,
          `R Markdown file not found: ${fullFilename}`
        );
      }
      
      // Build Docker image if needed
      try {
        await buildDockerImage();
      } catch (error: any) {
        return {
          content: [{
            type: "text",
            text: `Error building Docker image: ${error.message}`
          }],
          isError: true
        };
      }
      
      // Render the R Markdown file
      try {
        const outputFile = await renderRMarkdown(fullFilename, format, useDockerCompose);
        
        return {
          content: [{
            type: "text",
            text: `Successfully rendered ${fullFilename} to ${outputFile}`
          }]
        };
      } catch (error: any) {
        return {
          content: [{
            type: "text",
            text: `Error rendering R Markdown: ${error.message || String(error)}`
          }],
          isError: true
        };
      }
    }
    
    default:
      throw new McpError(
        ErrorCode.MethodNotFound,
        `Unknown tool: ${request.params.name}`
      );
  }
});

/**
 * Start the server using stdio transport.
 * This allows the server to communicate via standard input/output streams.
 */
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((error) => {
  console.error("Server error:", error);
  process.exit(1);
});
