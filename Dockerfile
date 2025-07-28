# 🔧 Stage 1: Build Go binary
FROM golang:1.21 AS builder

WORKDIR /app
COPY go/go.mod ./go/
RUN cd go && go mod download

COPY go/ ./go/
WORKDIR /app/go
RUN go build -o /main ./cmd/main.go

# 🐍 Stage 2: Final runtime image
FROM python:3.10-slim

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential poppler-utils && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy Python files
COPY python/ ./python/
RUN pip install --no-cache-dir -r python/requirements.txt && \
    python -m spacy download en_core_web_md

# Copy Go binary
COPY --from=builder /main /main

# Create mount points for volumes
RUN mkdir input output temp_output

# Set entrypoint
ENTRYPOINT ["/main"]

