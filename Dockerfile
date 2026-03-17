# Build stage
FROM golang:1.26.1@sha256:318ba178e04ea7655d4e4b1f3f0e81da0da5ff28a2c48681ff0418fb75f5e189 AS builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the Go source (relies on .dockerignore to filter)
COPY . .

# Build the application with FIPS 140-3 support.
# GOFIPS140=latest selects the Go Cryptographic Module and enables
# FIPS mode by default. The module is pure Go (no cgo required).
ARG VERSION=devel
RUN CGO_ENABLED=0 VERSION=${VERSION} make build

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot@sha256:e3f945647ffb95b5839c07038d64f9811adf17308b9121d8a2b87b6a22a80a39

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /workspace/bin/otterscale .

# Switch to non-root user
USER 65532:65532

# Set environment variables
ENV OTTERSCALE_SERVER_ADDRESS=0.0.0.0:8299
ENV OTTERSCALE_SERVER_TUNNEL_ADDRESS=0.0.0.0:8300
ENV GODEBUG=fips140=on

# Expose ports (8299: HTTP/gRPC API, 8300: Tunnel)
EXPOSE 8299 8300

# Labels
LABEL maintainer="Chung-Hsuan Tsai <paul_tsai@phison.com>"

ENTRYPOINT ["/otterscale"]
CMD ["server"]