# OtterScale

[![Go Reference](https://pkg.go.dev/badge/github.com/openhdc/otterscale.svg)](https://pkg.go.dev/github.com/openhdc/otterscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/openhdc/otterscale?style=flat-square)](https://goreportcard.com/report/github.com/openhdc/otterscale)
[![GitHub Build Status](https://github.com/openhdc/otterscale/actions/workflows/go.yml/badge.svg?style=flat-square)](https://github.com/openhdc/otterscale/actions/workflows/go.yml)
[![GitHub Release](https://img.shields.io/github/v/release/openhdc/otterscale?style=flat-square)](https://github.com/openhdc/otterscale/releases)
[![GitHub License](https://img.shields.io/github/license/openhdc/otterscale?style=flat-square)](https://opensource.org/license/agpl-v3)

**OtterScale** is a comprehensive **hyper-converged infrastructure platform** that unifies compute, storage, and networking resources into a single, scalable solution. Built for modern data centers, it seamlessly integrates **virtual machine management**, **software-defined networking**, **distributed storage**, **container orchestration**, **GPU resource management**, and **application marketplace** capabilities.

Designed to simplify complex infrastructure operations, OtterScale provides a **unified control plane** for managing heterogeneous workloads across **bare metal**, **virtualized**, and **containerized environments**. Whether you're running traditional VMs, cloud-native applications, or GPU-accelerated workloads, OtterScale delivers **enterprise-grade performance** with **operational simplicity**.

## ✨ Features

### Virtualization and Compute

- **Bare Metal Automation**: Integration with MAAS for streamlined physical server lifecycle management.
- **Virtual Machine Management**: Comprehensive KVM/QEMU virtualization supporting live migration and dynamic resource allocation.
- **GPU Resource Management**: Dynamic allocation and sharing of GPU resources for AI/ML and HPC workloads.

### Container and Orchestration

- **Container Orchestration**: Native Kubernetes integration for managing containerized workloads.
- **Service Orchestration**: Juju charm deployment for simplified application modeling and service orchestration.

### Storage and Data Management

- **Storage Management**: Built-in Ceph storage cluster provisioning and management for scalable storage solutions.
- **Backup and Disaster Recovery**: Automated snapshot management and cross-site replication for data protection.

### Monitoring and Diagnostics

- **Monitoring and Observability**: Integrated Prometheus, Grafana, and distributed tracing for real-time system insights.
- **Built-in Self Test (BIST)**: Comprehensive system health monitoring and diagnostics for proactive maintenance.

### Security and Access Control

- **Identity and Access Management**: Role-Based Access Control (RBAC) with LDAP/Active Directory integration and Single Sign-On (SSO) support.

### Application and API Ecosystem

- **Application Marketplace**: Curated catalog of pre-configured applications and services for easy deployment.
- **API-First Architecture**: gRPC APIs with Protocol Buffers specification and SDK support for seamless integration.

### Scalability and Reliability

- **High Availability**: Multi-node deployment with automatic failover for uninterrupted operations.
- **Extensible Architecture**: Plugin system enabling custom integrations and workflows.

## 🍺 Quick Start

- Download the [latest release](https://github.com/openhdc/otterscale/releases/latest) or compile from source

  ```sh
  make build && cd bin
  ```

- Initialize configuration and launch server

  ```sh
  ./otterscale init > otterscale.yaml
  ./otterscale serve --address :8299 --config otterscale.yaml
  ```

## 🔨 Development

Ensure you have the following environment setup:

- Go 1.24.3 or later
- Protobuf compiler (`protoc`)
- Make

### Prerequisites

Before building from source, you must install the following system dependencies:

```sh
# Ubuntu/Debian
sudo apt-get install libcephfs-dev librbd-dev librados-dev build-essential

# CentOS/RHEL/Fedora
sudo yum install libcephfs-devel librbd-devel librados-devel gcc gcc-c++ make
# or
sudo dnf install libcephfs-devel librbd-devel librados-devel gcc gcc-c++ make
```

## 🔍 Documentation

For detailed documentation, please visit [docs](/docs) directory.

## 🦮 Help

If you need help, feel free to open an issue on GitHub or use the discussions feature to contact the maintainers.

We'll do our best to assist you promptly.

## 📢 Roadmap

- [x] [v1.0.0](https://github.com/openhdc/otterscale/milestone/1)
- [ ] [v1.1.0](https://github.com/openhdc/otterscale/milestone/2)
- [ ] [v1.2.0](https://github.com/openhdc/otterscale/milestone/3)

## ⛔ Rules

Please review and adhere to the contribution guidelines outlined in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## ⚖️ License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
