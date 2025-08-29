# OtterScale Installation Manual

## Prerequisites

Before starting the installation process, ensure the following prerequisites are met:

- A network bridge must be set up in your environment and bind a group of IP addresses to the bridge.
- The IP address used by MAAS must be part of the network configured on the bridge.

## On-line Installation

1. **Prepare the install.cfg**

2. **Run the Installation Script with Configuration File**

   ```bash
   curl -L https://raw.githubusercontent.com/otterscale/otterscale/refs/heads/main/scripts/install.sh | sudo bash -s -- --config=~/install.cfg
   ```

## Interactive Installation

To perform an interactive installation, follow these steps:

1. **Run the Installation Script**

   ```bash
   sudo bash install.sh
   ```

2. **Input OtterScale Endpoint**
   - You will be prompted to enter the OtterScale endpoint URL.
   - Example: `http://127.0.0.1:8299`

3. **Input MAAS IP (CIDR Format)**
   - You will be prompted to enter the IP address (in CIDR format) that MAAS will use.
   - Example: `172.20.10.5/28`

4. **Input MAAS DHCP Dynamic Range**
   a. **Name of the IP Range**
      - You will be prompted to enter the name of the IP range for MAAS DHCP.
      - Example: `172.20.10.5/28`
   b. **Start IP for MAAS DHCP**
      - You will be prompted to enter the start IP address for the MAAS DHCP dynamic range.
      - Example: `172.20.10.8`
   c. **End IP for MAAS DHCP**
      - You will be prompted to enter the end IP address for the MAAS DHCP dynamic range.
      - Example: `172.20.10.12`

5. **Input Juju VM Static IP**
   - You will be prompted to enter the static IP address that the Juju VM will use.
   - Example: `172.20.10.6`

## One-Time Installation

For a one-time installation, follow these steps:

1. **Prepare an Existing Network Bridge**
   - Ensure that a network bridge is already set up and configured in your environment.

2. **Prepare the Configuration File**
   - Create a configuration file named `install.cfg` (or any name you prefer) with the required parameters.

3. **Run the Installation Script with Configuration File**

   ```bash
   sudo bash install.sh --config=./install.cfg
   ```

### Example `install.cfg` File

```bash
## OtterScale endpoint
OTTERSCALE_ENDPOINT="http://127.0.0.1:8299"

## Type network bridge name that will be used
OTTERSCALE_CONFIG_BRIDGE_CIDR="172.20.10.5/28"

## Type the IP range for network PXE boot used
OTTERSCALE_CNOFIG_MAAS_DHCP_CIDR="172.20.10.5/28"
OTTERSCALE_CONFIG_MAAS_DHCP_START_IP="172.20.10.8"
OTTERSCALE_CONFIG_MAAS_DHCP_END_IP="172.20.10.12"

## Type JuJu IP
OTTERSCALE_CONFIG_JUJU_IP="172.20.10.6"

## MAAS configure
OTTERSCALE_CONFIG_MAAS_ADMIN_USER="admin"
OTTERSCALE_CONFIG_MAAS_ADMIN_PASS="admin"
OTTERSCALE_CONFIG_MAAS_ADMIN_EMAIL="admin@example.com"
```

### Notes

- Ensure that the network bridge is correctly configured and the IP addresses specified are within the same subnet.
- The `OTTERSCALE_CONFIG_MAAS_DHCP_CIDR` should match the `OTTERSCALE_CONFIG_BRIDGE_CIDR`.
- The `OTTERSCALE_CONFIG_MAAS_DHCP_START_IP` and `OTTERSCALE_CONFIG_MAAS_DHCP_END_IP` should be within the range specified in `OTTERSCALE_CONFIG_MAAS_DHCP_CIDR`.
- The `OTTERSCALE_CONFIG_JUJU_IP` should be a static IP address outside the DHCP range.

## Troubleshooting

- **Network Issues**: Ensure that the network bridge is properly configured and that the IP addresses are not in use.
- **Script Errors**: If you encounter any errors during the installation, check the script's logs for more detailed information.
- **MAAS Configuration**: Verify that the MAAS admin credentials and IP configurations are correct.

For further assistance, refer to the OtterScale documentation or contact the support team.
