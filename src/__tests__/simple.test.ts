import { jest } from '@jest/globals';
import {
  extractRmdTitle,
  getFileExtension,
  getMimeType,
  ensureRmdExtension,
  createRmdFrontMatter
} from '../utils.js';

describe('Simple Tests', () => {
  it('should pass a basic test', () => {
    expect(1 + 1).toBe(2);
  });
  
  it('should correctly extract title from R Markdown content', () => {
    const content = `---
title: "Test Title"
---

Content`;
    
    expect(extractRmdTitle(content)).toBe('Test Title');
  });
  
  it('should correctly get file extension', () => {
    expect(getFileExtension('file.txt')).toBe('txt');
    expect(getFileExtension('file.Rmd')).toBe('rmd');
  });
  
  it('should correctly determine MIME type', () => {
    expect(getMimeType('file.html')).toBe('text/html');
    expect(getMimeType('file.pdf')).toBe('application/pdf');
    expect(getMimeType('file.Rmd')).toBe('text/markdown');
  });
  
  it('should ensure .Rmd extension', () => {
    expect(ensureRmdExtension('file')).toBe('file.Rmd');
    expect(ensureRmdExtension('file.Rmd')).toBe('file.Rmd');
  });
  
  it('should create YAML front matter', () => {
    const content = 'Content';
    const result = createRmdFrontMatter('Test Title', content);
    
    expect(result).toContain('title: "Test Title"');
    expect(result).toContain('Content');
  });
});
