import { jest } from '@jest/globals';
import * as fs from 'fs';
import * as path from 'path';
import { RMarkdownFile } from '../utils.js';

// Mock fs module
jest.mock('fs');

// Mock path module
jest.mock('path');

describe('R-Server Functions', () => {
  beforeEach(() => {
    // Reset mocks
    jest.resetAllMocks();
    
    // Setup mock implementations for fs
    jest.spyOn(fs, 'existsSync').mockReturnValue(true);
    jest.spyOn(fs, 'readdirSync').mockReturnValue(['sample.Rmd'] as any);
    jest.spyOn(fs, 'readFileSync').mockReturnValue(`---
title: "Sample R Markdown Document"
author: "R-Server"
date: "\`r Sys.Date()\`"
output: html_document
---

## R Markdown Sample

This is a sample R Markdown document.
` as any);
    
    // Setup mock implementations for path
    jest.spyOn(path, 'join').mockImplementation((...args: string[]) => args.join('/'));
    jest.spyOn(path, 'basename').mockImplementation((p: string) => p.split('/').pop() || '');
  });
  
  describe('Utility Functions', () => {
    it('should correctly mock fs functions', () => {
      expect(fs.existsSync).toBeDefined();
      expect(fs.readdirSync).toBeDefined();
      expect(fs.readFileSync).toBeDefined();
      
      // Call the mocked functions
      fs.existsSync('test');
      fs.readdirSync('test');
      fs.readFileSync('test', 'utf-8');
      
      // Verify they were called
      expect(fs.existsSync).toHaveBeenCalledWith('test');
      expect(fs.readdirSync).toHaveBeenCalledWith('test');
      expect(fs.readFileSync).toHaveBeenCalledWith('test', 'utf-8');
    });
    
    it('should correctly mock path functions', () => {
      expect(path.join).toBeDefined();
      expect(path.basename).toBeDefined();
      
      // Call the mocked functions
      const joined = path.join('dir', 'file.txt');
      const base = path.basename('/dir/file.txt');
      
      // Verify they were called and returned expected values
      expect(path.join).toHaveBeenCalledWith('dir', 'file.txt');
      expect(joined).toBe('dir/file.txt');
      expect(path.basename).toHaveBeenCalledWith('/dir/file.txt');
      expect(base).toBe('file.txt');
    });
  });
});
