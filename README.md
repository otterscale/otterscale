# OtterScale

[![Go Reference](https://pkg.go.dev/badge/github.com/openhdc/otterscale.svg)](https://pkg.go.dev/github.com/openhdc/otterscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/openhdc/otterscale?style=flat-square)](https://goreportcard.com/report/github.com/openhdc/otterscale)
[![GitHub Build Status](https://github.com/openhdc/otterscale/actions/workflows/go.yml/badge.svg?style=flat-square)](https://github.com/openhdc/otterscale/actions/workflows/go.yml)
[![GitHub Release](https://img.shields.io/github/v/release/openhdc/otterscale?style=flat-square)](https://github.com/openhdc/otterscale/releases)
[![GitHub License](https://img.shields.io/github/license/openhdc/otterscale?style=flat-square)](https://opensource.org/license/agpl-v3)

- ***WIP***

## ✨ Features

- ***WIP***

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
  - [x] MAAS
  - [x] Juju
  - [x] Kubernetes
  - [x] Ceph
  - [x] BIST

## ⛔ Rules

Please review and adhere to the contribution guidelines outlined in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## ⚖️ License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
