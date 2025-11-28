# Build stage
FROM golang:1.25.4-trixie AS builder

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
    git \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Create a home directory for nobody
RUN mkdir -p /home/nobody && chown nobody:nogroup /home/nobody
ENV HOME=/home/nobody

# Copy binary with correct ownership
COPY --from=builder --chown=nobody:nogroup /src/bin/otterscale ./otterscale
RUN chmod 550 ./otterscale

# Switch to non-root user
USER nobody

# Set environment variable
ENV OTTERSCALE_CONTAINER=true

# Expose port
EXPOSE 8299

# Labels
LABEL maintainer="Chung-Hsuan Tsai <paul_tsai@phison.com>"

ENTRYPOINT ["/app/otterscale"]