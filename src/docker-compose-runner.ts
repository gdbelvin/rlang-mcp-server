import { promisify } from "util";
import { exec } from "child_process";
import * as path from "path";
import { createDockerComposeEnv } from "./docker-config.js";

const execPromise = promisify(exec);

/**
 * Render R Markdown using docker-compose
 * 
 * @param rmdDir Directory containing R Markdown files
 * @param filename Filename of the R Markdown file to render
 * @param outputFormat Output format (html, pdf, or word)
 * @returns Output filename
 */
export async function renderWithDockerCompose(
  rmdDir: string,
  filename: string,
  outputFormat: "html" | "pdf" | "word" = "html"
): Promise<string> {
  const inputFile = path.join(rmdDir, filename);
  
  // Create environment variables for docker-compose
  const env = createDockerComposeEnv(rmdDir, filename, outputFormat);
  
  // Add the render command based on output format
  const renderCommand = outputFormat === "html" 
    ? "render" 
    : outputFormat === "pdf" 
      ? "render_pdf" 
      : "render_word";
  
  env.RENDER_COMMAND = renderCommand;
  
  // Build the environment variables string for the command
  const envString = Object.entries(env)
    .map(([key, value]) => `${key}=${value}`)
    .join(" ");
  
  try {
    // Run docker-compose with the environment variables
    const { stdout, stderr } = await execPromise(`${envString} docker-compose up --build`);
    
    // Extract the output file path from logs
    const outputFileMatch = stdout.match(/Rendered file: ([^\n]+)/);
    if (!outputFileMatch) {
      throw new Error("Failed to extract output file path from logs");
    }
    
    const outputFilePath = outputFileMatch[1];
    const outputFileName = path.basename(outputFilePath);
    
    return outputFileName;
  } catch (error: any) {
    console.error("Error rendering with docker-compose:", error);
    throw new Error(`Failed to render R Markdown file with docker-compose: ${error.message || String(error)}`);
  } finally {
    // Clean up by stopping and removing containers
    try {
      await execPromise("docker-compose down");
    } catch (cleanupError) {
      console.error("Error cleaning up docker-compose:", cleanupError);
    }
  }
}
