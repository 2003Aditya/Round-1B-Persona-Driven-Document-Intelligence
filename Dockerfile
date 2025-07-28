
# =====================
# 🧱 Stage 1: Build Go binary
# =====================
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Only copy necessary files
COPY go/go.mod  ./
RUN go mod download

COPY go/ ./
RUN go build -o main ./cmd/main.go

# ============================
# 🐍 Stage 2: Runtime - Python + Go Binary
# ============================

FROM python:3.11-slim AS runtime

# Install required system dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        libgl1-mesa-glx \
        libglib2.0-0 \
        libpoppler-cpp0v5 \
    && apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy and install only necessary Python dependencies
COPY python/requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir -r /tmp/requirements.txt && \
    python -m spacy download en_core_web_md && \
    rm -rf /root/.cache/pip

# Copy Go binary
COPY --from=builder /app/main /usr/local/bin/doc-intelligence

# Set working directory
WORKDIR /app

# Copy only needed runtime files
COPY input/ ./input/
COPY output/ ./output/
COPY go/temp_output/ ./temp_output/
COPY python/extractor/ ./python/extractor/
COPY python/analyzer/ ./python/analyzer/

# Run the Go binary
ENTRYPOINT ["doc-intelligence"]

