# Build stage
FROM golang:1.24 AS build

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o r-server ./cmd/r-server

# Final stage
FROM rocker/r-ver:4.4.3

# Install required system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install required R packages
RUN R -e "install.packages(c('ggplot2'), repos='https://cran.rstudio.com/')"

# Create a non-root user
RUN useradd -m -s /bin/bash -u 1000 rserver

# Create directories for the application
RUN mkdir -p /app/output && chown -R rserver:rserver /app

# Set working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build --chown=rserver:rserver /app/r-server /app/

# Switch to non-root user
USER rserver

# Expose the port
EXPOSE 22011

# Set the entrypoint
ENTRYPOINT ["/app/r-server"]
