# 🔧 Troubleshooting Guide

This guide helps you resolve common issues when installing, configuring, or running OtterScale.

## 🚀 Installation Issues

### Cannot Access Web Interface (localhost:3000)

**Problem**: Unable to access `http://localhost:3000` after starting OtterScale.

**Solutions**:

1. **Check if port is in use**:

   ```bash
   # Check what's using port 3000
   lsof -i :3000
   # Or use netstat
   netstat -an | grep 3000
   ```

2. **Change the port**:

   ```bash
   # Edit your .env file
   echo "PORT=3001" >> .env
   # Or set in docker-compose
   docker compose up -d
   ```

3. **Check container status**:

   ```bash
   docker compose ps
   docker compose logs
   ```

4. **Firewall issues**:
   ```bash
   # Ubuntu/Debian
   sudo ufw allow 3000
   # CentOS/RHEL
   sudo firewall-cmd --add-port=3000/tcp --permanent
   sudo firewall-cmd --reload
   ```

### Docker Compose Fails to Start

**Problem**: `docker compose up -d` fails with errors.

**Solutions**:

1. **Check Docker daemon**:

   ```bash
   docker --version
   docker compose --version
   sudo systemctl status docker
   ```

2. **Insufficient resources**:

   - Ensure you have at least 4GB RAM available
   - Free up disk space (minimum 10GB required)

   ```bash
   df -h
   free -h
   ```

3. **Permission issues**:
   ```bash
   # Add user to docker group
   sudo usermod -aG docker $USER
   # Restart terminal or run
   newgrp docker
   ```

### Configuration File Issues

**Problem**: `otterscale.yaml` initialization fails.

**Solutions**:

1. **Docker image pull issues**:

   ```bash
   # Pull image manually
   docker pull ghcr.io/otterscale/otterscale/service:latest

   # Alternative: use specific version
   docker pull ghcr.io/otterscale/otterscale/service:v0.5.0
   ```

2. **Network connectivity**:

   ```bash
   # Test connectivity to GitHub Container Registry
   ping ghcr.io
   curl -I https://ghcr.io
   ```

3. **Authentication issues**:
   ```bash
   # Login to GitHub Container Registry if needed
   echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
   ```

## 🔧 Development Issues

### Build Failures

**Problem**: `make build` fails with compilation errors.

**Solutions**:

1. **Check Go version**:

   ```bash
   go version
   # Should be Go 1.25.3 or later
   ```

2. **Install system dependencies**:

   ```bash
   # Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install libcephfs-dev librbd-dev librados-dev build-essential

   # CentOS/RHEL/Fedora
   sudo dnf install libcephfs-devel librbd-devel librados-devel gcc gcc-c++ make
   ```

3. **Clean and rebuild**:
   ```bash
   make clean
   go mod tidy
   make build
   ```

### Protobuf Compilation Issues

**Problem**: `make proto` fails with protobuf errors.

**Solutions**:

1. **Install protobuf compiler**:

   ```bash
   # Ubuntu/Debian
   sudo apt-get install protobuf-compiler

   # macOS
   brew install protobuf

   # Or install from releases
   wget https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip
   ```

2. **Install Go plugins**:

   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. **Check PATH**:
   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"
   which protoc-gen-go
   ```

### Test Failures

**Problem**: `make test` fails with test errors.

**Solutions**:

1. **Run specific tests**:

   ```bash
   go test ./internal/... -v
   go test ./cmd/... -v
   ```

2. **Check test dependencies**:

   ```bash
   # Ensure Docker is running for integration tests
   docker info
   ```

3. **Skip integration tests**:
   ```bash
   go test -short ./...
   ```

## 🐳 Container & Orchestration Issues

### Kubernetes Integration Problems

**Problem**: Cannot connect to or manage Kubernetes clusters.

**Solutions**:

1. **Check kubectl configuration**:

   ```bash
   kubectl cluster-info
   kubectl get nodes
   ```

2. **Verify RBAC permissions**:

   ```bash
   kubectl auth can-i '*' '*' --as=system:serviceaccount:otterscale:default
   ```

3. **Network connectivity**:
   ```bash
   # Test cluster API server connectivity
   kubectl get --raw /api/v1/namespaces
   ```

### Ceph Storage Issues

**Problem**: Distributed storage cluster failures.

**Solutions**:

1. **Check Ceph cluster health**:

   ```bash
   ceph -s
   ceph health detail
   ```

2. **Monitor disk space**:

   ```bash
   ceph df
   ceph osd df
   ```

3. **Restart OSD services**:
   ```bash
   systemctl restart ceph-osd@*
   ```

### MAAS Integration Issues

**Problem**: Metal provisioning failures.

**Solutions**:

1. **Check MAAS connectivity**:

   ```bash
   maas login admin http://your-maas-server:5240/MAAS/
   maas admin machines read
   ```

2. **Verify API credentials**:

   - Check API key in configuration
   - Ensure user has proper permissions

3. **Network configuration**:
   - Verify DHCP/PXE boot settings
   - Check network fabric configuration

## 🔐 Authentication & Security Issues

### LDAP/AD Integration Problems

**Problem**: Cannot authenticate users via LDAP/Active Directory.

**Solutions**:

1. **Test LDAP connectivity**:

   ```bash
   ldapsearch -x -H ldap://your-ldap-server -D "cn=admin,dc=example,dc=com" -W
   ```

2. **Check SSL certificates**:

   ```bash
   openssl s_client -connect your-ldap-server:636 -showcerts
   ```

3. **Verify bind credentials**:
   - Ensure service account has read permissions
   - Check DN format and search base

### SSO Configuration Issues

**Problem**: Single Sign-On authentication failures.

**Solutions**:

1. **Check OIDC/SAML configuration**:

   - Verify client ID and secret
   - Ensure redirect URIs are correct
   - Check token endpoint accessibility

2. **Certificate validation**:
   ```bash
   curl -v https://your-idp.com/.well-known/openid_configuration
   ```

## 📊 Performance & Monitoring Issues

### High Resource Usage

**Problem**: OtterScale consuming excessive CPU/memory.

**Solutions**:

1. **Monitor resource usage**:

   ```bash
   docker stats
   top -p $(pgrep otterscale)
   ```

2. **Check for memory leaks**:

   ```bash
   # Enable Go profiling
   go tool pprof http://localhost:6060/debug/pprof/heap
   ```

3. **Optimize configuration**:
   - Reduce polling intervals
   - Limit concurrent operations
   - Tune garbage collection settings

### Slow API Responses

**Problem**: Web interface or API calls are slow.

**Solutions**:

1. **Check database performance**:

   ```sql
   -- For PostgreSQL
   SELECT * FROM pg_stat_activity WHERE state = 'active';
   ```

2. **Network latency**:

   ```bash
   ping your-backend-server
   traceroute your-backend-server
   ```

3. **Enable request logging**:
   - Set log level to DEBUG
   - Monitor slow queries
   - Check for N+1 query problems

## 📝 Log Analysis

### Enable Debug Logging

```yaml
# In otterscale.yaml
logging:
  level: debug
  format: json
```

### Common Log Locations

```bash
# Container logs
docker compose logs -f

# System logs
journalctl -u otterscale -f

# Application logs
tail -f /var/log/otterscale/otterscale.log
```

### Log Pattern Analysis

```bash
# Find errors
grep -i error /var/log/otterscale/*.log

# Monitor authentication failures
grep -i "auth.*fail" /var/log/otterscale/*.log

# Check API rate limits
grep -i "rate.*limit" /var/log/otterscale/*.log
```

## 🆘 Getting Help

If you cannot resolve the issue using this guide:

1. **Check existing issues**: [GitHub Issues](https://github.com/otterscale/otterscale/issues)
2. **Search discussions**: [GitHub Discussions](https://github.com/otterscale/otterscale/discussions)
3. **Create detailed bug report** with:

   - OtterScale version
   - Operating system and version
   - Complete error messages
   - Steps to reproduce
   - Configuration files (redact sensitive data)

4. **Enterprise support**: Email [support@otterscale.com](mailto:support@otterscale.com)

## 📋 System Information Collection

Run this script to collect diagnostic information:

```bash
#!/bin/bash
echo "=== OtterScale Diagnostic Information ==="
echo "Date: $(date)"
echo "OS: $(uname -a)"
echo "Docker: $(docker --version)"
echo "Go: $(go version)"
echo ""
echo "=== Container Status ==="
docker compose ps
echo ""
echo "=== Container Logs (last 50 lines) ==="
docker compose logs --tail=50
echo ""
echo "=== System Resources ==="
free -h
df -h
echo ""
echo "=== Network Connectivity ==="
curl -I http://localhost:3000 2>/dev/null || echo "Web interface unreachable"
```

Save this as `diagnostic.sh`, make it executable with `chmod +x diagnostic.sh`, and run it when reporting issues.
