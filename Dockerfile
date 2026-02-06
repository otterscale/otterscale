# Build stage
FROM golang:1.26rc3-trixie@sha256:bc4a32a68778c1f78911fa769a296a8beb01d1c786d13dab93fefcd1985f4004 AS builder

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
FROM debian:trixie-slim@sha256:e711a7b30ec1261130d0a121050b4ed81d7fb28aeabcf4ea0c7876d4e9f5aca2

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
LABEL maintainer="Chung-Hsuan Tsai <paul_tsai@phison.com>"

ENTRYPOINT ["/app/otterscale"]