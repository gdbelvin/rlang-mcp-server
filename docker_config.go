package main

import (
	"path/filepath"
)

// DockerConfig represents Docker configuration for R Markdown rendering
type DockerConfig struct {
	Image   string
	Command []string
	Binds   []string
}

// CreateDockerConfig creates Docker configuration for rendering R Markdown
func CreateDockerConfig(rmdDir, inputFile string, outputFormat string) DockerConfig {
	// Inside the container, the file will be at /rmd/filename
	containerInputFile := filepath.Join("/rmd", filepath.Base(inputFile))

	// Determine the render command based on output format
	renderCommand := "render"
	if outputFormat == "pdf" {
		renderCommand = "render_pdf"
	} else if outputFormat == "word" {
		renderCommand = "render_word"
	}

	return DockerConfig{
		Image:   "r-server-rmd",
		Command: []string{renderCommand, containerInputFile},
		Binds:   []string{rmdDir + ":/rmd"},
	}
}

// CreateDockerComposeEnv generates docker-compose environment variables
func CreateDockerComposeEnv(rmdDir, inputFile string, outputFormat string) map[string]string {
	// Inside the container, the file will be at /rmd/filename
	containerInputFile := filepath.Join("/rmd", filepath.Base(inputFile))

	return map[string]string{
		"RMD_DIR":       rmdDir,
		"INPUT_FILE":    containerInputFile,
		"OUTPUT_FORMAT": outputFormat,
	}
}
