package main

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// RenderWithDockerode renders R Markdown using the Docker API
func RenderWithDockerode(rmdDir, filename, outputFormat string) (string, error) {
	// Get Docker configuration
	config := CreateDockerConfig(rmdDir, filename, outputFormat)

	// Create Docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	fmt.Printf("Creating Docker container with bind: %s\n", config.Binds[0])

	// Create container
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: config.Image,
			Cmd:   config.Command,
		},
		&container.HostConfig{
			Binds: config.Binds,
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for container to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for container: %w", err)
		}
	case <-statusCh:
	}

	// Get logs
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer out.Close()

	// Read logs
	logs, err := io.ReadAll(out)
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}

	// Remove container
	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		fmt.Printf("Warning: failed to remove container: %v\n", err)
	}

	// Extract output file path from logs
	logsStr := string(logs)
	re := regexp.MustCompile(`Rendered file: ([^\n]+)`)
	matches := re.FindStringSubmatch(logsStr)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to extract output file path from logs")
	}

	outputFilePath := matches[1]
	outputFileName := filepath.Base(outputFilePath)

	return outputFileName, nil
}
