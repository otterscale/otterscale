# OpenHDC

[![Go Reference](https://pkg.go.dev/badge/github.com/openhdc/openhdc.svg)](https://pkg.go.dev/github.com/openhdc/openhdc)
[![Go Report Card](https://goreportcard.com/badge/github.com/openhdc/openhdc?style=flat-square)](https://goreportcard.com/report/github.com/openhdc/openhdc)
[![GitHub Build Status](https://github.com/openhdc/openhdc/actions/workflows/go.yml/badge.svg?style=flat-square)](https://github.com/openhdc/openhdc/actions/workflows/go.yml)
[![GitHub Release](https://img.shields.io/github/v/release/openhdc/openhdc?style=flat-square)](https://github.com/openhdc/openhdc/releases)
[![GitHub License](https://img.shields.io/github/license/openhdc/openhdc)](https://opensource.org/license/mpl-2-0)

OpenHDC (**Open** **H**ybrid **D**ata **C**enter) is an open-source project designed to provide a robust and local server solution for hybrid data sources.

We preserves data privacy by localizing data and combining siloed information to enhance traceability.

## ✨ Features

- ***Robust Local Server***: Safeguards data with local storage.
- ***Hybrid Data Center***: Consolidates data from multiple sources.
- ***Improved Traceability***: Facilitates data tracking across systems.

## 🍺 Build & Run

1. Install by `go install`

    ```sh
    go install github.com/openhdc/openhdc@latest
    ```

2. Build from source
    1. Clone the repository:

        ```sh
        git clone https://github.com/openhdc/openhdc.git
        ```

    2. Change directory and build:

        ```sh
        cd openhdc && make build
        ```

    3. Run:

        ```sh
        ./bin/openhdc
        ```

## 🔨 Environment

Ensure you have the following environment setup:

- Go 1.23.4 or later
- Protobuf compiler (`protoc`)
- Make

## 🔍 Documentation

For detailed documentation, please visit [docs](/docs) directory.

## 🦮 Help

If you need help, feel free to open an issue on GitHub or use the discussions feature to contact the maintainers.

We'll do our best to assist you promptly.

## 📢 Roadmap

- [ ] v0.0.1
  - [ ] Better error messages
  - [ ] Improved naming conventions
  - [ ] Enhanced application closing procedures

## ⛔ Rules

Please review and adhere to the contribution guidelines outlined in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## ⚖️ License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
