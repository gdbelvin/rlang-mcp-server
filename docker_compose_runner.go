package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// RenderWithDockerCompose renders R Markdown using docker-compose
func RenderWithDockerCompose(rmdDir, filename, outputFormat string) (string, error) {
	// We don't need to use inputFile directly as it's handled by the Docker environment
	_ = filepath.Join(rmdDir, filename)

	// Create environment variables for docker-compose
	env := CreateDockerComposeEnv(rmdDir, filename, outputFormat)

	// Add the render command based on output format
	renderCommand := "render"
	if outputFormat == "pdf" {
		renderCommand = "render_pdf"
	} else if outputFormat == "word" {
		renderCommand = "render_word"
	}

	env["RENDER_COMMAND"] = renderCommand

	// Build the environment variables for the command
	cmd := exec.Command("docker-compose", "up", "--build")
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Capture stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run docker-compose: %w\nOutput: %s", err, string(output))
	}

	// Extract the output file path from logs
	outputStr := string(output)
	re := regexp.MustCompile(`Rendered file: ([^\n]+)`)
	matches := re.FindStringSubmatch(outputStr)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to extract output file path from logs")
	}

	outputFilePath := matches[1]
	outputFileName := filepath.Base(outputFilePath)

	// Clean up by stopping and removing containers
	cleanupCmd := exec.Command("docker-compose", "down")
	if err := cleanupCmd.Run(); err != nil {
		fmt.Printf("Warning: failed to clean up docker-compose: %v\n", err)
	}

	return outputFileName, nil
}
