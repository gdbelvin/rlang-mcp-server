import Dockerode from "dockerode";
import * as path from "path";
import { createDockerConfig } from "./docker-config.js";

const docker = new Dockerode();

/**
 * Render R Markdown using Dockerode
 * 
 * @param rmdDir Directory containing R Markdown files
 * @param filename Filename of the R Markdown file to render
 * @param outputFormat Output format (html, pdf, or word)
 * @returns Output filename
 */
export async function renderWithDockerode(
  rmdDir: string,
  filename: string,
  outputFormat: "html" | "pdf" | "word" = "html"
): Promise<string> {
  // Get Docker configuration
  const config = createDockerConfig(rmdDir, filename, outputFormat);
  
  try {
    console.error(`Creating Docker container with bind: ${config.binds[0]}`);
    
    // Create a container using the configuration
    const container = await docker.createContainer({
      Image: config.image,
      Cmd: config.command,
      HostConfig: {
        Binds: config.binds
      }
    });

    // Start the container
    await container.start();

    // Wait for the container to finish
    await container.wait();

    // Get the logs
    const logs = await container.logs({
      stdout: true,
      stderr: true
    });

    // Remove the container
    await container.remove();

    // Convert logs to string
    const logsString = logs.toString();
    
    // Extract the output file path
    const outputFileMatch = logsString.match(/Rendered file: ([^\n]+)/);
    if (!outputFileMatch) {
      throw new Error("Failed to extract output file path from logs");
    }

    const outputFilePath = outputFileMatch[1];
    const outputFileName = path.basename(outputFilePath);
    
    return outputFileName;
  } catch (error: any) {
    console.error("Error rendering with Dockerode:", error);
    throw new Error(`Failed to render R Markdown file with Dockerode: ${error.message || String(error)}`);
  }
}
