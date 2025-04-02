# Build stage
FROM golang:1.24 AS build

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o r-server ./cmd/r-server

# Final stage
FROM rocker/r-ver:4.4.3

# Install required system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install tidyverse
RUN R -q -e "install.packages('tidyverse', repos='https://packagemanager.rstudio.com/cran/2024-03-01', Ncpus=3)"
# Install required R packages
RUN R -q -e "install.packages('cowplot',   repos='https://packagemanager.rstudio.com/cran/2024-03-01', Ncpus=3)"

# Inst-q all RMarkdown packages 
RUN R -q -e "install.packages('quarto', repos = 'https://packagemanager.rstudio.com/cran/2024-03-01', Ncpus=3)"
RUN R -q -e "install.packages('knitr', repos = 'https://packagemanager.rstudio.com/cran/2024-03-01', Ncpus=3)"
RUN R -q -e "install.packages('rmarkdown', repos = 'https://packagemanager.rstudio.com/cran/2024-03-01', Ncpus=3)"

# Create a non-root user
RUN useradd -m -s /bin/bash -u 1000 rserver

# Create directories for the application and ensure proper permissions
RUN mkdir -p /app/output && chown -R rserver:rserver /app

# Set working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build --chown=rserver:rserver /app/r-server /app/

# Switch to non-root user
USER rserver

# Set the entrypoint
# Using exec form to ensure signals are properly passed and stdin/stdout are correctly handled
ENTRYPOINT ["/app/r-server"]

# No CMD is needed as we want to use the ENTRYPOINT directly
# The server communicates via stdin/stdout, so no ports are exposed
