FROM golang:1.24 AS go-builder

# Set working directory for Go build
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY *.go ./
COPY *.mod ./

# Build the application
RUN go build -o r-server

# Final image with R and the Go binary
FROM rocker/r-ver:4.4.3

# Install pandoc (required for R Markdown)
RUN apt-get update && apt-get install -y --no-install-recommends \
    pandoc \
    && rm -rf /var/lib/apt/lists/*

# Install R packages needed for R Markdown
RUN R -e "install.packages(c('rmarkdown', 'knitr', 'tinytex'), repos='https://cran.rstudio.com/')"
RUN R -e "install.packages(c('ggplot2'), repos='https://cran.rstudio.com/')"

# Create a working directory for R Markdown files
WORKDIR /rmd

# Create output directory
RUN mkdir -p /rmd/output

# Create app directory and copy Go binary from builder stage
RUN mkdir -p /app
COPY --from=go-builder /app/r-server /app/

# Expose the port that the server listens on
EXPOSE 22011

# Run the Go application
ENTRYPOINT ["/app/r-server"]
