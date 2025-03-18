package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// RMarkdownFile represents metadata for an R Markdown file
type RMarkdownFile struct {
	Filename string
	Title    string
	Path     string
}

// ExtractRmdTitle extracts the title from R Markdown content
func ExtractRmdTitle(content string) string {
	re := regexp.MustCompile(`title:\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// GetFileExtension returns the file extension
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}
	return strings.ToLower(ext[1:])
}

// GetMimeType determines MIME type based on file extension
func GetMimeType(filename string) string {
	ext := GetFileExtension(filename)

	switch ext {
	case "html":
		return "text/html"
	case "pdf":
		return "application/pdf"
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "rmd", "md":
		return "text/markdown"
	default:
		return "text/plain"
	}
}

// EnsureRmdExtension ensures filename has .Rmd extension
func EnsureRmdExtension(filename string) string {
	if strings.HasSuffix(filename, ".Rmd") || strings.HasSuffix(filename, ".rmd") {
		return filename
	}
	return filename + ".Rmd"
}

// CreateRmdFrontMatter creates YAML front matter for R Markdown
func CreateRmdFrontMatter(title, content string) string {
	if strings.HasPrefix(content, "---") {
		return content
	}

	frontMatter := `---
title: "` + title + `"
author: "R-Server"
date: "` + "`r Sys.Date()`" + `"
output: html_document
---

` + content

	return frontMatter
}

// GetRMarkdownFiles gets all R Markdown files in the specified directory
func GetRMarkdownFiles(rmdDir string) ([]RMarkdownFile, error) {
	files, err := os.ReadDir(rmdDir)
	if err != nil {
		return nil, err
	}

	var rmdFiles []RMarkdownFile
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if !strings.HasSuffix(strings.ToLower(filename), ".rmd") {
			continue
		}

		filePath := filepath.Join(rmdDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		title := ExtractRmdTitle(string(content))
		if title == "" {
			title = filename
		}

		rmdFiles = append(rmdFiles, RMarkdownFile{
			Filename: filename,
			Title:    title,
			Path:     filePath,
		})
	}

	return rmdFiles, nil
}

// EnsureDirectoriesExist ensures that the required directories exist
func EnsureDirectoriesExist(rmdDir, outputDir string) error {
	if err := os.MkdirAll(rmdDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	return nil
}
