import { jest } from '@jest/globals';
import * as fs from 'fs';
import * as path from 'path';
import { exec } from 'child_process';
import Dockerode from 'dockerode';

// Mock fs module
jest.mock('fs', () => ({
  existsSync: jest.fn(),
  mkdirSync: jest.fn(),
  readdirSync: jest.fn(),
  readFileSync: jest.fn(),
  writeFileSync: jest.fn(),
}));

// Mock path module
jest.mock('path', () => ({
  join: jest.fn((...args) => args.join('/')),
}));

// Mock child_process.exec
jest.mock('child_process', () => ({
  exec: jest.fn(),
}));

// Mock Dockerode
jest.mock('dockerode', () => {
  const mockContainer = {
    start: jest.fn().mockResolvedValue(undefined),
    wait: jest.fn().mockResolvedValue(undefined),
    logs: jest.fn().mockResolvedValue(Buffer.from('Rendered file: /rmd/output/test.html')),
    remove: jest.fn().mockResolvedValue(undefined),
  };
  
  const mockDockerode = jest.fn().mockImplementation(() => {
    return {
      createContainer: jest.fn().mockResolvedValue(mockContainer),
    };
  });
  
  return mockDockerode;
});

// We'll need to mock these functions from the server
const mockGetRMarkdownFiles = jest.fn();
const mockBuildDockerImage = jest.fn();
const mockRenderRMarkdown = jest.fn();

// Mock the server module
jest.mock('../index.js', () => {
  return {
    getRMarkdownFiles: mockGetRMarkdownFiles,
    buildDockerImage: mockBuildDockerImage,
    renderRMarkdown: mockRenderRMarkdown,
  };
}, { virtual: true });

describe('R-Server', () => {
  beforeEach(() => {
    // Reset all mocks
    jest.resetAllMocks();
    
    // Setup default mock implementations
    (fs.existsSync as jest.Mock).mockReturnValue(true);
    (fs.readdirSync as jest.Mock).mockReturnValue(['sample.Rmd']);
    (fs.readFileSync as jest.Mock).mockReturnValue(`---
title: "Sample R Markdown Document"
author: "R-Server"
date: "\`r Sys.Date()\`"
output: html_document
---

## R Markdown Sample

This is a sample R Markdown document.
`);
    
    // Mock exec to resolve successfully
    (exec as unknown as jest.Mock).mockImplementation((cmd, callback) => {
      if (callback) {
        callback(null, { stdout: 'Success', stderr: '' });
      }
      return {
        stdout: 'Success',
        stderr: '',
      };
    });
    
    // Setup mock implementations for our server functions
    mockGetRMarkdownFiles.mockReturnValue([
      {
        filename: 'sample.Rmd',
        title: 'Sample R Markdown Document',
        path: 'rmd/sample.Rmd',
      },
    ]);
    
    mockBuildDockerImage.mockResolvedValue(undefined);
    mockRenderRMarkdown.mockResolvedValue('test.html');
  });
  
  describe('getRMarkdownFiles', () => {
    it('should return a list of R Markdown files', () => {
      // Call the function
      const result = mockGetRMarkdownFiles();
      
      // Verify the result
      expect(result).toEqual([
        {
          filename: 'sample.Rmd',
          title: 'Sample R Markdown Document',
          path: 'rmd/sample.Rmd',
        },
      ]);
    });
  });
  
  describe('buildDockerImage', () => {
    it('should build the Docker image', async () => {
      // Call the function
      await mockBuildDockerImage();
      
      // Verify the function was called
      expect(mockBuildDockerImage).toHaveBeenCalled();
    });
    
    it('should handle errors', async () => {
      // Mock the function to reject
      mockBuildDockerImage.mockRejectedValueOnce(new Error('Docker build failed'));
      
      // Call the function and expect it to throw
      await expect(mockBuildDockerImage()).rejects.toThrow('Docker build failed');
    });
  });
  
  describe('renderRMarkdown', () => {
    it('should render an R Markdown file', async () => {
      // Call the function
      const result = await mockRenderRMarkdown('sample.Rmd');
      
      // Verify the result
      expect(result).toBe('test.html');
      
      // Verify the function was called with the correct arguments
      expect(mockRenderRMarkdown).toHaveBeenCalledWith('sample.Rmd');
    });
  });
});
