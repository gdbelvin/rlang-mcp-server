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

# Install required R packages
#RUN R -q -e "install.packages(c('tidyverse'), repos='https://cran.rstudio.com/', Ncpus=4)"
RUN R -q -e "options(warn=2); install.packages('ggplot2', repos = 'https://packagemanager.rstudio.com/cran/2024-03-01')"
RUN R -q -e "options(warn=2); install.packages('cowplot', repos = 'https://packagemanager.rstudio.com/cran/2024-03-01')"
# Install common R packages
#RUN R -q -e "install.packages(c('data.table', 'lubridate', 'scales', 'reshape2', 'magrittr'), repos='https://cran.rstudio.com/')"
# Install Rmarkdown packages
#RUN R -q -e "install.packages(c('quarto', 'knitr', 'rmarkdown'), repos='https://cran.rstudio.com/')"
# Install AI packages
#RUN R -q -e "install.packages(c('caret', 'randomForest'), repos='https://cran.rstudio.com/')"

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
