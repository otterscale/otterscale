# OtterScale

[![GitHub Release](https://img.shields.io/github/v/release/otterscale/otterscale?logo=github)](https://github.com/otterscale/otterscale/releases)
[![GitHub License](https://img.shields.io/github/license/otterscale/otterscale?logo=github)](https://opensource.org/license/agpl-v3)
[![Go CI](https://github.com/otterscale/otterscale/actions/workflows/ci-go.yml/badge.svg)](https://github.com/otterscale/otterscale/actions/workflows/ci-go.yml)
[![SvelteKit CI](https://github.com/otterscale/otterscale/actions/workflows/ci-sveltekit.yml/badge.svg)](https://github.com/otterscale/otterscale/actions/workflows/ci-sveltekit.yml)
[![Codecov](https://codecov.io/gh/otterscale/otterscale/graph/badge.svg?token=I7R0YEMXER)](https://codecov.io/gh/otterscale/otterscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/otterscale/otterscale)](https://goreportcard.com/report/github.com/otterscale/otterscale)
[![Go Reference](https://pkg.go.dev/badge/github.com/otterscale/otterscale.svg)](https://pkg.go.dev/github.com/otterscale/otterscale)

**OtterScale** is a hyper-converged infrastructure (HCI) platform that unifies compute, storage, and networking into a scalable solution. Seamlessly manage VMs, containers, GPUs, and applications through a single control plane with enterprise-grade performance.

![Login](/assets/screenshot-login.png)

> [!NOTE]
> Click below to view screenshots and explore the interface.

<details>
  <summary><b>📸 Screenshots</b></summary>

|                                                                                           Home                                                                                           |                                                                                      Scope Selector                                                                                       |
| :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
|                   ![Home](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-home.jpeg)                   |         ![Scope Selector](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-scope-selector.jpeg)          |
|                                                                                     **Create Scope**                                                                                     |                                                                                 **Create Scope Settings**                                                                                 |
|         ![Create Scope 1](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-scope-create-1.jpeg)         |         ![Create Scope 2](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-scope-create-2.jpeg)          |
|                                                                                    **Scope Settings**                                                                                    |                                                                                   **Application Store**                                                                                   |
|         ![Scope Settings](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-scope-settings.jpeg)         | ![Application Store](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-application-management-store.jpeg) |
|                                                                                       **Machines**                                                                                       |                                                                                     **Machines Dark**                                                                                     |
|               ![Machines](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-machines.jpeg)               |          ![Machines Dark](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-machines-dark.jpeg)           |
|                                                                                       **Storage**                                                                                        |                                                                                      **Networking**                                                                                       |
|                ![Storage](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-storage.jpeg)                |             ![Networking](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-networking.jpeg)              |
|                                                                                **Application Management**                                                                                |                                                                                 **Application Workloads**                                                                                 |
| ![Application Management](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-application-management.jpeg) |   ![Workloads](https://raw.githubusercontent.com/otterscale/otterscale.github.io/refs/heads/feat/add_screenshot/static/img/screenshot/screenshot-application-management-workloads.jpeg)   |

</details>

## ✨ Key Features

- **🖥️ Virtualization**: KVM/QEMU VMs with live migration and GPU management
- **🐳 Container Orchestration**: Native Kubernetes and Juju charm deployment
- **💾 Distributed Storage**: Built-in Ceph clusters with automated backup
- **📊 Monitoring**: Integrated Prometheus and Grafana stack
- **🔐 Security**: RBAC with LDAP/AD integration and SSO
- **🛒 Application Marketplace**: Curated catalog of ready-to-deploy apps
- **⚡ High Availability**: Multi-node deployment with automatic failover

## 🚀 Quick Start

> [!IMPORTANT]
> Requirements: `Git`, `Docker`, `8GB RAM`, `100GB disk space`

1. **Clone the repository:**

   ```sh
   git clone --depth 1 https://github.com/otterscale/otterscale.git
   cd otterscale
   ```

2. **Set up environment:**

   ```sh
   cp .env.example .env
   # Edit .env file with your settings
   ```

3. **Start OtterScale:**

   ```sh
   docker compose up -d
   ```

4. **Access web interface:**

Open your browser and navigate to `http://localhost:3000` (or your configured port).

> [!TIP]
> If you cannot access `http://localhost:3000`, check if the port is in use or refer to the [troubleshooting guide](/docs/troubleshooting.md).

For production deployments, see our [deployment guide](/docs/deployment.md).

## 🔧 Development

> [!IMPORTANT]
> Requirements: `Go 1.25+`, `Docker`, `Protobuf compiler`, `Make`, `Git`

### System Dependencies

```bash
# Ubuntu/Debian
sudo apt-get install libcephfs-dev librbd-dev librados-dev build-essential

# CentOS/RHEL/Fedora
sudo dnf install libcephfs-devel librbd-devel librados-devel gcc gcc-c++ make
```

### Build from Source

```bash
git clone https://github.com/otterscale/otterscale.git
cd otterscale
make build
make test
```

**Development Commands**: `make` (show targets), `make build`, `make test`, `make lint`, `make proto`, `make image`

## 📚 Documentation & Support

- **[Getting Started](https://otterscale.github.io/getting-started)** - Setup and configuration guide
- **[API Reference](https://otterscale.github.io/api)** - gRPC and REST API docs
- **[Architecture](https://otterscale.github.io/architecture)** - System components overview
- **[Deployment](https://otterscale.github.io/deployment)** - Production best practices

**Need Help?**

- [GitHub Issues](https://github.com/otterscale/otterscale/issues) - Bug reports and feature requests
- [GitHub Discussions](https://github.com/otterscale/otterscale/discussions) - Community Q&A
- [Enterprise Support](mailto:support@otterscale.com) - Commercial inquiries

## 🗺️ Roadmap

|                              Version                               | Topic                        | Status         |
| :----------------------------------------------------------------: | ---------------------------- | -------------- |
| **[v0.5.0](https://github.com/otterscale/otterscale/milestone/1)** | Infrastructure Core          | ✅ Complete    |
| **[v0.6.0](https://github.com/otterscale/otterscale/milestone/2)** | Compute Resources & Policies | ⏳ In Progress |
| **[v0.7.0](https://github.com/otterscale/otterscale/milestone/3)** | Developer Experience         | 📅 Planned     |

## 🤝 Contributing

We welcome contributions! Please:

1. Fork the repository and create a feature branch
2. Make changes and add tests
3. Ensure code passes tests and linting
4. Submit a pull request

Follow [Conventional Commits](https://www.conventionalcommits.org/) and review our [Contributing Guidelines](CONTRIBUTING.md).

## 🔒 Security

For security vulnerabilities, email [security@otterscale.com](mailto:security@otterscale.com) instead of creating public issues.

## 🙏 Acknowledgements

We extend our heartfelt gratitude to those who make OtterScale possible:

- **Open Source Community**: Thanks to contributors, early adopters, and active members in [GitHub Issues](https://github.com/otterscale/otterscale/issues) and [Discussions](https://github.com/otterscale/otterscale/discussions) for code, feedback, and ideas.
- **Core Technologies**:
  - **[Kubernetes](https://kubernetes.io/)**: Container orchestration
  - **[Ceph](https://ceph.io/)**: Distributed storage
  - **[Juju](https://juju.is/)**: Application deployment
  - **[MAAS](https://maas.io/)**: Metal provisioning

Your support drives OtterScale’s mission to build better hyper-converged solutions! 🌟

## 📄 License

Licensed under [GNU Affero General Public License v3.0](LICENSE) (AGPL-3.0). Open source with copyleft requirements.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fotterscale%2Fotterscale.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fotterscale%2Fotterscale?ref=badge_large&issueType=license)
