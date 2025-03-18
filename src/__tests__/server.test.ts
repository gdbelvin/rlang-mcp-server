import { jest } from '@jest/globals';
import * as fs from 'fs';
import * as path from 'path';
import { promisify } from 'util';
import { exec } from 'child_process';

// Mock the fs module
jest.mock('fs', () => ({
  existsSync: jest.fn(),
  mkdirSync: jest.fn(),
  readdirSync: jest.fn(),
  readFileSync: jest.fn(),
  writeFileSync: jest.fn(),
}));

// Mock the path module
jest.mock('path', () => ({
  join: jest.fn((...args) => args.join('/')),
  basename: jest.fn((p) => p.split('/').pop()),
}));

// Mock the child_process.exec function
jest.mock('child_process', () => ({
  exec: jest.fn(),
}));

// Mock the dockerode module
jest.mock('dockerode', () => {
  return function() {
    return {
      createContainer: jest.fn().mockResolvedValue({
        start: jest.fn().mockResolvedValue(undefined),
        wait: jest.fn().mockResolvedValue(undefined),
        logs: jest.fn().mockResolvedValue(Buffer.from('Rendered file: /rmd/output/test.html')),
        remove: jest.fn().mockResolvedValue(undefined),
      }),
    };
  };
});

// Mock the MCP SDK
jest.mock('@modelcontextprotocol/sdk/server/index.js', () => ({
  Server: jest.fn().mockImplementation(() => ({
    setRequestHandler: jest.fn(),
    connect: jest.fn().mockResolvedValue(undefined),
    onerror: jest.fn(),
    close: jest.fn().mockResolvedValue(undefined),
  })),
}));

jest.mock('@modelcontextprotocol/sdk/server/stdio.js', () => ({
  StdioServerTransport: jest.fn(),
}));

jest.mock('@modelcontextprotocol/sdk/types.js', () => ({
  ListResourcesRequestSchema: 'ListResourcesRequestSchema',
  ReadResourceRequestSchema: 'ReadResourceRequestSchema',
  ListToolsRequestSchema: 'ListToolsRequestSchema',
  CallToolRequestSchema: 'CallToolRequestSchema',
  ErrorCode: {
    InvalidRequest: 'InvalidRequest',
    InvalidParams: 'InvalidParams',
    MethodNotFound: 'MethodNotFound',
    InternalError: 'InternalError',
  },
  McpError: class McpError extends Error {
    constructor(code, message) {
      super(message);
      this.code = code;
    }
  },
}));

describe('R-Server MCP', () => {
  beforeEach(() => {
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
  });
  
  describe('Server initialization', () => {
    it('should create directories if they do not exist', () => {
      // Mock existsSync to return false for directories
      (fs.existsSync as jest.Mock)
        .mockReturnValueOnce(false)  // RMD_DIR
        .mockReturnValueOnce(false); // OUTPUT_DIR
      
      // Import the server module
      jest.isolateModules(() => {
        require('../index.js');
      });
      
      // Verify that mkdirSync was called for both directories
      expect(fs.mkdirSync).toHaveBeenCalledTimes(2);
    });
  });
  
  describe('ListResourcesRequestSchema handler', () => {
    it('should list R Markdown files and rendered outputs', async () => {
      // Mock readdirSync for OUTPUT_DIR
      (fs.readdirSync as jest.Mock)
        .mockReturnValueOnce(['sample.Rmd'])  // RMD_DIR
        .mockReturnValueOnce(['sample.html']); // OUTPUT_DIR
      
      // Import the server module
      let serverModule;
      jest.isolateModules(() => {
        serverModule = require('../index.js');
      });
      
      // Get the Server mock
      const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
      
      // Get the setRequestHandler mock
      const setRequestHandler = Server.mock.results[0].value.setRequestHandler;
      
      // Find the handler for ListResourcesRequestSchema
      const handler = setRequestHandler.mock.calls.find(
        call => call[0] === 'ListResourcesRequestSchema'
      )[1];
      
      // Call the handler
      const result = await handler();
      
      // Verify the result
      expect(result).toHaveProperty('resources');
      expect(result.resources).toHaveLength(2); // 1 Rmd file + 1 output file
      
      // Check the R Markdown resource
      const rmdResource = result.resources.find(r => r.uri.startsWith('rmd:///'));
      expect(rmdResource).toBeDefined();
      expect(rmdResource.name).toBe('Sample R Markdown Document');
      
      // Check the output resource
      const outputResource = result.resources.find(r => r.uri.startsWith('rmd-output:///'));
      expect(outputResource).toBeDefined();
      expect(outputResource.name).toBe('Rendered: sample.html');
    });
  });
  
  describe('ReadResourceRequestSchema handler', () => {
    it('should read an R Markdown file', async () => {
      // Import the server module
      let serverModule;
      jest.isolateModules(() => {
        serverModule = require('../index.js');
      });
      
      // Get the Server mock
      const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
      
      // Get the setRequestHandler mock
      const setRequestHandler = Server.mock.results[0].value.setRequestHandler;
      
      // Find the handler for ReadResourceRequestSchema
      const handler = setRequestHandler.mock.calls.find(
        call => call[0] === 'ReadResourceRequestSchema'
      )[1];
      
      // Call the handler with an R Markdown URI
      const result = await handler({
        params: {
          uri: 'rmd:///sample.Rmd'
        }
      });
      
      // Verify the result
      expect(result).toHaveProperty('contents');
      expect(result.contents).toHaveLength(1);
      expect(result.contents[0].mimeType).toBe('text/markdown');
      expect(result.contents[0].uri).toBe('rmd:///sample.Rmd');
    });
    
    it('should read a rendered output file', async () => {
      // Mock readFileSync for the output file
      (fs.readFileSync as jest.Mock).mockReturnValue('<html>Rendered content</html>');
      
      // Import the server module
      let serverModule;
      jest.isolateModules(() => {
        serverModule = require('../index.js');
      });
      
      // Get the Server mock
      const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
      
      // Get the setRequestHandler mock
      const setRequestHandler = Server.mock.results[0].value.setRequestHandler;
      
      // Find the handler for ReadResourceRequestSchema
      const handler = setRequestHandler.mock.calls.find(
        call => call[0] === 'ReadResourceRequestSchema'
      )[1];
      
      // Call the handler with an output URI
      const result = await handler({
        params: {
          uri: 'rmd-output:///sample.html'
        }
      });
      
      // Verify the result
      expect(result).toHaveProperty('contents');
      expect(result.contents).toHaveLength(1);
      expect(result.contents[0].mimeType).toBe('text/html');
      expect(result.contents[0].uri).toBe('rmd-output:///sample.html');
      expect(result.contents[0].text).toBe('<html>Rendered content</html>');
    });
  });
  
  describe('CallToolRequestSchema handler', () => {
    it('should create an R Markdown file', async () => {
      // Import the server module
      let serverModule;
      jest.isolateModules(() => {
        serverModule = require('../index.js');
      });
      
      // Get the Server mock
      const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
      
      // Get the setRequestHandler mock
      const setRequestHandler = Server.mock.results[0].value.setRequestHandler;
      
      // Find the handler for CallToolRequestSchema
      const handler = setRequestHandler.mock.calls.find(
        call => call[0] === 'CallToolRequestSchema'
      )[1];
      
      // Call the handler with create_rmd tool
      const result = await handler({
        params: {
          name: 'create_rmd',
          arguments: {
            filename: 'test',
            title: 'Test Document',
            content: 'This is a test document.'
          }
        }
      });
      
      // Verify the result
      expect(result).toHaveProperty('content');
      expect(result.content).toHaveLength(1);
      expect(result.content[0].text).toContain('Created R Markdown file: test.Rmd');
      
      // Verify that writeFileSync was called
      expect(fs.writeFileSync).toHaveBeenCalled();
      expect((fs.writeFileSync as jest.Mock).mock.calls[0][0]).toContain('test.Rmd');
      expect((fs.writeFileSync as jest.Mock).mock.calls[0][1]).toContain('title: "Test Document"');
    });
    
    it('should render an R Markdown file', async () => {
      // Import the server module
      let serverModule;
      jest.isolateModules(() => {
        serverModule = require('../index.js');
      });
      
      // Get the Server mock
      const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
      
      // Get the setRequestHandler mock
      const setRequestHandler = Server.mock.results[0].value.setRequestHandler;
      
      // Find the handler for CallToolRequestSchema
      const handler = setRequestHandler.mock.calls.find(
        call => call[0] === 'CallToolRequestSchema'
      )[1];
      
      // Call the handler with render_rmd tool
      const result = await handler({
        params: {
          name: 'render_rmd',
          arguments: {
            filename: 'sample.Rmd',
            format: 'html'
          }
        }
      });
      
      // Verify the result
      expect(result).toHaveProperty('content');
      expect(result.content).toHaveLength(1);
      expect(result.content[0].text).toContain('Successfully rendered');
    });
  });
});
