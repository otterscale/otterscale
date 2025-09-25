# Build Notes

## Development Requirements

### System Dependencies

Some components require system libraries to be installed:

```bash
# Ubuntu/Debian
sudo apt-get install libcephfs-dev librbd-dev librados-dev build-essential

# CentOS/RHEL/Fedora  
sudo dnf install libcephfs-devel librbd-devel librados-devel gcc gcc-c++ make
```

### Build Issues

If you encounter Ceph-related build errors like:
```
fatal error: rados/librados.h: No such file or directory
```

You can:

1. Install the system dependencies above, or
2. Use the standalone capabilities tool: `go build -o ./bin/capabilities ./cmd/capabilities/...`
3. Or run specific modules that don't depend on Ceph: `go test ./internal/enum ./internal/wrap`

### Protobuf Generation

To regenerate protobuf files after changes:
```bash
make proto
```

This requires `protoc` to be installed. On Ubuntu/Debian:
```bash
sudo apt-get install protobuf-compiler
```

## Alternative Testing

For testing capabilities without full dependencies:
```bash
# Build standalone capabilities CLI
go build -o ./bin/capabilities ./cmd/capabilities/...

# Test capabilities command
./bin/capabilities --help
./bin/capabilities -l zh  # Chinese
./bin/capabilities -f json # JSON output

# Build test CLI wrapper
go build -o ./bin/test-capabilities ./cmd/test-capabilities/...
./bin/test-capabilities capabilities
```