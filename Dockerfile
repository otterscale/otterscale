# Build stage
FROM golang:1.26.5@sha256:079e59808d2d252516e27e3f3a9c003740dee7f75e55aa71528766d52bcfc16a AS builder

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
FROM gcr.io/distroless/static:nonroot@sha256:963fa6c544fe5ce420f1f54fb88b6fb01479f054c8056d0f74cc2c6000df5240

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