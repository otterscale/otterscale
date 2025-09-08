# Build stage
FROM golang:1.24.6-trixie AS builder

WORKDIR /src

# Install build dependencies
RUN apt-get update && apt-get install -y \
    libcephfs-dev \
    librbd-dev \
    librados-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Runtime stage
FROM debian:trixie-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libcephfs2 \
    librbd1 \
    librados2 \
    pciutils \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Copy binary with correct ownership
COPY --from=builder --chown=nobody:nobody /src/bin/otterscale ./otterscale
RUN chmod 550 ./otterscale

# Switch to non-root user
USER nobody

# Set environment variable
ENV OTTERSCALE_CONTAINER=true

# Expose port
EXPOSE 8299

# Labels
LABEL maintainer="Chung-Hsuan Tsai <zx.c@nycu.edu.tw>"

ENTRYPOINT ["/app/otterscale"]