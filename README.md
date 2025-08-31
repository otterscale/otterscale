# 🦦 OtterScale

[![Go Reference](https://pkg.go.dev/badge/github.com/otterscale/otterscale.svg)](https://pkg.go.dev/github.com/otterscale/otterscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/otterscale/otterscale?style=flat-square)](https://goreportcard.com/report/github.com/otterscale/otterscale)
[![GitHub Build Status](https://github.com/otterscale/otterscale/actions/workflows/go.yml/badge.svg?style=flat-square)](https://github.com/otterscale/otterscale/actions/workflows/go.yml)
[![GitHub Release](https://img.shields.io/github/v/release/otterscale/otterscale?style=flat-square)](https://github.com/otterscale/otterscale/releases)
[![GitHub License](https://img.shields.io/github/license/otterscale/otterscale?style=flat-square)](https://opensource.org/license/agpl-v3)
[![Contributors](https://img.shields.io/github/contributors/otterscale/otterscale?style=flat-square)](https://github.com/otterscale/otterscale/graphs/contributors)

**OtterScale** is a powerful **hyper-converged infrastructure platform (HCI)** that combines **compute**, **storage**, and **networking** into one scalable solution. Tailored for **modern data centers**, it seamlessly manages **VMs**, **software-defined networking**, **distributed storage**, **containers**, **GPUs**, and an **app marketplace**.

![Login](/assets/screenshot-login.png)

Simplify complex operations with OtterScale's **unified control plane**, effortlessly handling diverse workloads across **bare metal**, **virtual**, and **containerized environments**. Run **traditional VMs**, **cloud-native apps**, or **GPU-accelerated tasks** with **enterprise-grade performance** and unmatched ease.

## ✨ Features

### 🖥️ Virtualization and Compute

> - **Bare Metal Automation**: Seamless server lifecycle management through MAAS integration
> - **Virtual Machine Management**: Full-featured KVM/QEMU with live migration and dynamic scaling
> - **GPU Resource Management**: Intelligent allocation and sharing for AI/ML and HPC workloads

### 🐳 Container and Orchestration

> - **Container Orchestration**: Native Kubernetes integration for cloud-native applications
> - **Service Orchestration**: Simplified application modeling with Juju charm deployment

### 💾 Storage and Data Management

> - **Storage Management**: Built-in Ceph clusters for scalable, distributed storage
> - **Backup and Recovery**: Automated snapshots with cross-site replication

### 📊 Monitoring and Diagnostics

> - **Observability Stack**: Integrated Prometheus and Grafana for real-time insights
> - **Health Monitoring**: Comprehensive BIST for proactive system maintenance

### 🔐 Security and Access Control

> - **Identity Management**: RBAC with LDAP/AD integration and SSO support

### 🛒 Application and API Ecosystem

> - **Application Marketplace**: Curated catalog of ready-to-deploy applications
> - **API-First Design**: gRPC APIs with Protocol Buffers and SDK support

### ⚡ Scalability and Reliability

> - **High Availability**: Multi-node deployment with automatic failover
> - **Extensible Architecture**: Plugin system for custom integrations

## 🚀 Quick Start

### System Requirements

- Docker and Docker Compose
- At least 4GB of available RAM
- 10GB of free disk space

### Installation

1. **Initialize configuration:**

   ```sh
   docker run ghcr.io/otterscale/otterscale/service:latest init > otterscale.yaml
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
   git clone https://github.com/otterscale/otterscale.git
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

- **[Getting Started Guide](https://otterscale.github.io/getting-started)** - Complete setup and configuration guide
- **[API Reference](https://otterscale.github.io/api)** - gRPC and REST API documentation
- **[Architecture Overview](https://otterscale.github.io/architecture)** - System architecture and components
- **[Deployment Guide](https://otterscale.github.io/deployment)** - Production deployment best practices
- **[Configuration Reference](https://otterscale.github.io/configuration)** - Configuration options and examples
- **[Troubleshooting](https://otterscale.github.io/troubleshooting)** - Common issues and solutions

## 🆘 Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/otterscale/otterscale/issues)
- **GitHub Discussions**: [Community discussions and Q&A](https://github.com/otterscale/otterscale/discussions)
- **Documentation**: [Comprehensive guides and API docs](/docs)
- **Email Support**: For enterprise support inquiries, contact [support@otterscale.io](mailto:support@otterscale.io)

## 🗺️ Roadmap

- [ ] [v1.0.0](https://github.com/otterscale/otterscale/milestone/1): MAAS, Juju, Ceph, Kubernetes, Helm
- [ ] [v1.1.0](https://github.com/otterscale/otterscale/milestone/2): Virtual Machine, GPU Operator, Open Policy Agent
- [ ] [v1.2.0](https://github.com/otterscale/otterscale/milestone/3): Helm Upload, API Interface, AI Agent

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
2. **Email us** at [security@otterscale.io](mailto:security@otterscale.io) with details
3. **Include** steps to reproduce and potential impact
4. **Wait** for our response before disclosing publicly

We aim to respond to security reports within 24 hours and provide regular updates on resolution progress.

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0** (AGPL-3.0).

- **Open Source**: You can use, modify, and distribute this software freely
- **Copyleft**: Modifications must also be open source under AGPL-3.0
- **Network Use**: If you run OtterScale as a service, you must provide source code to users

For the complete license terms, see the [LICENSE](LICENSE) file.

## 🐾 FOSSA Status

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fotterscale%2Fotterscale.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fotterscale%2Fotterscale?ref=badge_large&issueType=license)

---

Built with ❤️ by the otterscale community

[Website](https://otterscale.io) • [GitHub](https://github.com/otterscale/otterscale) • [Documentation](/docs) • [Community](https://github.com/otterscale/otterscale/discussions)
