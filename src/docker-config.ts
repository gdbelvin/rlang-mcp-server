import * as path from "path";

/**
 * Docker configuration for R Markdown rendering
 */
export interface DockerConfig {
  image: string;
  command: string[];
  binds: string[];
}

/**
 * Create Docker configuration for rendering R Markdown
 * 
 * @param rmdDir Directory containing R Markdown files
 * @param inputFile Path to the R Markdown file to render
 * @param outputFormat Output format (html, pdf, or word)
 * @returns Docker configuration
 */
export function createDockerConfig(
  rmdDir: string,
  inputFile: string,
  outputFormat: "html" | "pdf" | "word" = "html"
): DockerConfig {
  // Inside the container, the file will be at /rmd/filename
  const containerInputFile = path.join("/rmd", path.basename(inputFile));
  
  // Determine the render command based on output format
  const renderCommand = outputFormat === "html" 
    ? "render" 
    : outputFormat === "pdf" 
      ? "render_pdf" 
      : "render_word";
  
  return {
    image: "r-server-rmd",
    command: [renderCommand, containerInputFile],
    binds: [`${rmdDir}:/rmd`]
  };
}

/**
 * Generate docker-compose environment variables
 * 
 * @param rmdDir Directory containing R Markdown files
 * @param inputFile Path to the R Markdown file to render
 * @param outputFormat Output format (html, pdf, or word)
 * @returns Environment variables for docker-compose
 */
export function createDockerComposeEnv(
  rmdDir: string,
  inputFile: string,
  outputFormat: "html" | "pdf" | "word" = "html"
): Record<string, string> {
  // Inside the container, the file will be at /rmd/filename
  const containerInputFile = path.join("/rmd", path.basename(inputFile));
  
  return {
    RMD_DIR: rmdDir,
    INPUT_FILE: containerInputFile,
    OUTPUT_FORMAT: outputFormat
  };
}
