import * as fs from "fs";
import * as path from "path";

/**
 * Interface for R Markdown file metadata
 */
export interface RMarkdownFile {
  filename: string;
  title: string;
  path: string;
}

/**
 * Extract title from R Markdown content
 */
export function extractRmdTitle(content: string): string {
  const titleMatch = content.match(/title:\s*"([^"]+)"/);
  return titleMatch ? titleMatch[1] : "";
}

/**
 * Get file extension
 */
export function getFileExtension(filename: string): string {
  const parts = filename.split(".");
  return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : "";
}

/**
 * Determine MIME type based on file extension
 */
export function getMimeType(filename: string): string {
  const ext = getFileExtension(filename);
  
  switch (ext) {
    case "html":
      return "text/html";
    case "pdf":
      return "application/pdf";
    case "docx":
      return "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
    case "rmd":
    case "md":
      return "text/markdown";
    default:
      return "text/plain";
  }
}

/**
 * Ensure filename has .Rmd extension
 */
export function ensureRmdExtension(filename: string): string {
  if (filename.endsWith(".Rmd") || filename.endsWith(".rmd")) {
    return filename;
  }
  return `${filename}.Rmd`;
}

/**
 * Create YAML front matter for R Markdown
 */
export function createRmdFrontMatter(title: string, content: string): string {
  if (content.startsWith("---")) {
    return content;
  }
  
  return `---
title: "${title}"
author: "R-Server"
date: "\`r Sys.Date()\`"
output: html_document
---

${content}`;
}
