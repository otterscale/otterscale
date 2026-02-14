# Build stage
FROM golang:1.25.4-trixie@sha256:a02d35efc036053fdf0da8c15919276bf777a80cbfda6a35c5e9f087e652adfc AS builder

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
FROM debian:trixie-slim@sha256:91e29de1e4e20f771e97d452c8fa6370716ca4044febbec4838366d459963801

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