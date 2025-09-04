# OtterScale Installation Manual

This document provides instructions for installing **OtterScale** on an **Ubuntu environment**.  
All installation steps **must be executed with `sudo` privileges**.

## Prerequisites

OtterScale requires a **network bridge** on the host system.  

- If you already have a bridge prepared, the installation will proceed more smoothly.  
- If you are not familiar with bridge configuration, the provided `install.sh` script can automatically create one using the system's default interface and gateway.  
  - The bridge will be created with the name: **`br-otters`**

---

## Installation Methods

There are two installation methods available:

### 1. Interactive Installation

This method prompts for necessary configurations during execution.

```bash
bash -c "$(curl -fsSL https://raw.githubusercontent.com/otterscale/otterscale/refs/heads/441-update-install-script-execute-command/scripts/install.sh)" -- url=YOUR_OTTERSCALE_ENDPOINT
```
Replace YOUR_OTTERSCALE_ENDPOINT with your actual OtterScale endpoint.

### 2. One-Shot Installation
This method uses a predefined configuration file, allowing for automated setup.
Steps:
Prepare a configuration file
Copy and modify the provided install.cfg.example file according to your environment.
Run the installation script

```bash
bash -c "$(curl -fsSL https://raw.githubusercontent.com/otterscale/otterscale/refs/heads/441-update-install-script-execute-command/scripts/install.sh)" -- url=YOUR_OTTERSCALE_ENDPOINT config=~/install.cfg
```
Replace YOUR_OTTERSCALE_ENDPOINT with your actual OtterScale endpoint, and ensure ~/install.cfg points to your customized configuration file.

## Notes

- Ensure you run all commands with **`sudo`**.  
- If a network bridge named `br-otters` already exists, the script will attempt to reuse it.  
- Incorrect configuration may prevent OtterScale from functioning properly. Please review your `install.cfg` carefully.

## Support

For additional help or troubleshooting, please open an issue on the [OtterScale GitHub repository](https://github.com/otterscale/otterscale).
