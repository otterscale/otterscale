# OtterScale

[![Go Reference](https://pkg.go.dev/badge/github.com/openhdc/otterscale.svg)](https://pkg.go.dev/github.com/openhdc/otterscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/openhdc/otterscale?style=flat-square)](https://goreportcard.com/report/github.com/openhdc/otterscale)
[![GitHub Build Status](https://github.com/openhdc/otterscale/actions/workflows/go.yml/badge.svg?style=flat-square)](https://github.com/openhdc/otterscale/actions/workflows/go.yml)
[![GitHub Release](https://img.shields.io/github/v/release/openhdc/otterscale?style=flat-square)](https://github.com/openhdc/otterscale/releases)
[![GitHub License](https://img.shields.io/github/license/openhdc/otterscale?style=flat-square)](https://opensource.org/license/agpl-v3)
[![Docker Pulls](https://img.shields.io/docker/pulls/otterscale/otterscale?style=flat-square)](https://hub.docker.com/r/otterscale/otterscale)
[![Contributors](https://img.shields.io/github/contributors/openhdc/otterscale?style=flat-square)](https://github.com/openhdc/otterscale/graphs/contributors)

**OtterScale** is a comprehensive **hyper-converged infrastructure platform** that unifies compute, storage, and networking resources into a single, scalable solution. Built for modern data centers, it seamlessly integrates virtual machine management, software-defined networking, distributed storage, container orchestration, GPU resource management, and application marketplace capabilities.

Designed to simplify complex infrastructure operations, OtterScale provides a unified control plane for managing heterogeneous workloads across bare metal, virtualized, and containerized environments. Whether you're running traditional VMs, cloud-native applications, or GPU-accelerated workloads, OtterScale delivers enterprise-grade performance with operational simplicity.

## 🏗️ Architecture

OtterScale follows a microservices architecture with the following key components:

- **Control Plane**: Centralized management and orchestration layer
- **Compute Engine**: VM lifecycle management with KVM/QEMU integration
- **Storage Layer**: Distributed Ceph storage with automatic provisioning
- **Network Fabric**: Software-defined networking with overlay support
- **Container Runtime**: Kubernetes integration for cloud-native workloads
- **API Gateway**: gRPC and REST APIs with authentication and authorization
- **Web Interface**: Modern React-based dashboard for system management

## ✨ Features

### 🖥️ Virtualization and Compute

- **Bare Metal Automation**: Integration with MAAS for streamlined physical server lifecycle management
- **Virtual Machine Management**: Comprehensive KVM/QEMU virtualization with live migration and dynamic resource allocation
- **GPU Resource Management**: Dynamic allocation and sharing of GPU resources for AI/ML and HPC workloads

### 🐳 Container and Orchestration

- **Container Orchestration**: Native Kubernetes integration for managing containerized workloads
- **Service Orchestration**: Juju charm deployment for simplified application modeling and service orchestration

### 💾 Storage and Data Management

- **Storage Management**: Built-in Ceph storage cluster provisioning and management for scalable storage solutions
- **Backup and Disaster Recovery**: Automated snapshot management and cross-site replication for data protection

### 📊 Monitoring and Diagnostics

- **Monitoring and Observability**: Integrated Prometheus, Grafana, and distributed tracing for real-time system insights
- **Built-in Self Test (BIST)**: Comprehensive system health monitoring and diagnostics for proactive maintenance

### 🔐 Security and Access Control

- **Identity and Access Management**: Role-Based Access Control (RBAC) with LDAP/Active Directory integration and Single Sign-On (SSO) support

### 🛒 Application and API Ecosystem

- **Application Marketplace**: Curated catalog of pre-configured applications and services for easy deployment
- **API-First Architecture**: gRPC APIs with Protocol Buffers specification and SDK support for seamless integration

### ⚡ Scalability and Reliability

- **High Availability**: Multi-node deployment with automatic failover for uninterrupted operations
- **Extensible Architecture**: Plugin system enabling custom integrations and workflows

## 🚀 Quick Start

### System Requirements

- Docker and Docker Compose
- At least 4GB of available RAM
- 10GB of free disk space

### Installation

1. **Initialize configuration:**

   ```sh
   docker run ghcr.io/otterscale/otterscale:latest init > otterscale.yaml
   ```

2. **Set up environment variables:**

   ```sh
   cp .env.example .env
   ```

   Edit the `.env` file to configure your deployment settings.

3. **Start OtterScale:**

   ```sh
   docker compose up -d
   ```

4. **Access the web interface:**

   Open your browser and navigate to `http://localhost:3000` (or your configured port).

For production deployments, please refer to our [deployment guide](/docs/deployment.md).

## 🔧 Development

### Requirements

- Go 1.24.6 or later
- Docker and Docker Compose
- Protobuf compiler (`protoc`)
- Make
- Git

### System Dependencies

Install system dependencies before building:

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install libcephfs-dev librbd-dev librados-dev build-essential

# CentOS/RHEL/Fedora
sudo yum install libcephfs-devel librbd-devel librados-devel gcc gcc-c++ make
# or for newer versions
sudo dnf install libcephfs-devel librbd-devel librados-devel gcc gcc-c++ make
```

### Building from Source

1. **Clone the repository:**

   ```bash
   git clone https://github.com/openhdc/otterscale.git
   cd otterscale
   ```

2. **Build the project:**

   ```bash
   make build
   ```

3. **Run tests:**

   ```bash
   make test
   ```

### Development Commands

- `make` - Show available targets
- `make build` - Build binary
- `make test` - Run tests
- `make lint` - Run linters
- `make proto` - Generate protobuf code
- `make openapi` - Generate API spec
- `make image` - Build Docker images

## 📚 Documentation

- **[Getting Started Guide](https://openhdc.github.io/getting-started)** - Complete setup and configuration guide
- **[API Reference](https://openhdc.github.io/api)** - gRPC and REST API documentation
- **[Architecture Overview](https://openhdc.github.io/architecture)** - System architecture and components
- **[Deployment Guide](https://openhdc.github.io/deployment)** - Production deployment best practices
- **[Configuration Reference](https://openhdc.github.io/configuration)** - Configuration options and examples
- **[Troubleshooting](https://openhdc.github.io/troubleshooting)** - Common issues and solutions

## 🆘 Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/openhdc/otterscale/issues)
- **GitHub Discussions**: [Community discussions and Q&A](https://github.com/openhdc/otterscale/discussions)
- **Documentation**: [Comprehensive guides and API docs](/docs)
- **Email Support**: For enterprise support inquiries, contact [support@openhdc.io](mailto:support@openhdc.io)

## 🗺️ Roadmap

- [x] [v1.0.0](https://github.com/openhdc/otterscale/milestone/1)
- [ ] [v1.1.0](https://github.com/openhdc/otterscale/milestone/2)
- [ ] [v1.2.0](https://github.com/openhdc/otterscale/milestone/3)

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

1. **Fork the repository** and create your feature branch
2. **Make your changes** and add tests where appropriate
3. **Ensure your code** passes all tests and linting
4. **Submit a pull request** with a clear description of your changes

Please review our [Contributing Guidelines](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

### Development Workflow

- Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification
- All pull requests require code review and passing CI checks
- Include tests for new features and bug fixes
- Update documentation for API changes

## 🔒 Security

Security is a top priority for OtterScale. If you discover a security vulnerability, please:

1. **Do not** create a public GitHub issue
2. **Email us** at [security@openhdc.io](mailto:security@openhdc.io) with details
3. **Include** steps to reproduce and potential impact
4. **Wait** for our response before disclosing publicly

We aim to respond to security reports within 24 hours and provide regular updates on resolution progress.

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0** (AGPL-3.0).

- **Open Source**: You can use, modify, and distribute this software freely
- **Copyleft**: Modifications must also be open source under AGPL-3.0
- **Network Use**: If you run OtterScale as a service, you must provide source code to users

For the complete license terms, see the [LICENSE](LICENSE) file.

---

Built with ❤️ by the OpenHDC community

[Website](https://openhdc.io) • [GitHub](https://github.com/openhdc/otterscale) • [Documentation](/docs) • [Community](https://github.com/openhdc/otterscale/discussions)
