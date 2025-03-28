#!/bin/bash
cd /Users/gdb/dev/Cline/MCP/r-server 
./r-server
# For debugging
# mkfifo mcp.fifo
# tee mcp.fifo | ./r-server | tee mcp.fifo 
