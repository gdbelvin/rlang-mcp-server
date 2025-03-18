import { jest } from '@jest/globals';
import {
  extractRmdTitle,
  getFileExtension,
  getMimeType,
  ensureRmdExtension,
  createRmdFrontMatter
} from '../utils.js';

describe('Utility Functions', () => {
  describe('extractRmdTitle', () => {
    it('should extract title from R Markdown content', () => {
      const content = `---
title: "Sample Document"
author: "Test Author"
date: "2023-01-01"
---

## Content`;
      
      expect(extractRmdTitle(content)).toBe('Sample Document');
    });
    
    it('should return empty string if title is not found', () => {
      const content = `---
author: "Test Author"
date: "2023-01-01"
---

## Content`;
      
      expect(extractRmdTitle(content)).toBe('');
    });
  });
  
  describe('getFileExtension', () => {
    it('should return the file extension in lowercase', () => {
      expect(getFileExtension('file.txt')).toBe('txt');
      expect(getFileExtension('file.TXT')).toBe('txt');
      expect(getFileExtension('file.Rmd')).toBe('rmd');
      expect(getFileExtension('path/to/file.HTML')).toBe('html');
    });
    
    it('should return empty string if there is no extension', () => {
      expect(getFileExtension('file')).toBe('');
      expect(getFileExtension('')).toBe('');
    });
  });
  
  describe('getMimeType', () => {
    it('should return the correct MIME type for HTML files', () => {
      expect(getMimeType('file.html')).toBe('text/html');
      expect(getMimeType('file.HTML')).toBe('text/html');
    });
    
    it('should return the correct MIME type for PDF files', () => {
      expect(getMimeType('file.pdf')).toBe('application/pdf');
    });
    
    it('should return the correct MIME type for Word files', () => {
      expect(getMimeType('file.docx')).toBe('application/vnd.openxmlformats-officedocument.wordprocessingml.document');
    });
    
    it('should return the correct MIME type for Markdown files', () => {
      expect(getMimeType('file.md')).toBe('text/markdown');
      expect(getMimeType('file.Rmd')).toBe('text/markdown');
      expect(getMimeType('file.rmd')).toBe('text/markdown');
    });
    
    it('should return text/plain for unknown file types', () => {
      expect(getMimeType('file.unknown')).toBe('text/plain');
      expect(getMimeType('file')).toBe('text/plain');
    });
  });
  
  describe('ensureRmdExtension', () => {
    it('should not modify filenames that already have .Rmd extension', () => {
      expect(ensureRmdExtension('file.Rmd')).toBe('file.Rmd');
      expect(ensureRmdExtension('file.rmd')).toBe('file.rmd');
    });
    
    it('should add .Rmd extension to filenames without it', () => {
      expect(ensureRmdExtension('file')).toBe('file.Rmd');
      expect(ensureRmdExtension('file.txt')).toBe('file.txt.Rmd');
    });
  });
  
  describe('createRmdFrontMatter', () => {
    it('should not modify content that already has front matter', () => {
      const content = `---
title: "Existing Title"
---

## Content`;
      
      expect(createRmdFrontMatter('New Title', content)).toBe(content);
    });
    
    it('should add front matter to content without it', () => {
      const content = '## Content';
      const result = createRmdFrontMatter('New Title', content);
      
      expect(result).toContain('title: "New Title"');
      expect(result).toContain('author: "R-Server"');
      expect(result).toContain('date: "`r Sys.Date()`"');
      expect(result).toContain('## Content');
    });
  });
});
