#!/bin/bash
cd /Users/gdb/dev/Cline/MCP/r-server 
tee mcp.fifo | ./r-server | tee mcp.fifo 
