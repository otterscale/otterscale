# OtterScale Installation Manual

## Interactive Installation

To install OtterScale interactively, follow these steps:

1. **Run the Installation Script**
   Execute the installation script with superuser privileges:
   ```bash
   $ sudo bash install.sh
   ```

2. **Input OtterScale Endpoint (CIDR)**
   When prompted, enter the OtterScale endpoint, below is an examaple:
   ```
   http://127.0.0.1:8299
   ```
3. **Input MAAS IP (CIDR Format)**
Provide the IP address for MAAS in CIDR format:
   ```
   172.20.10.5/28
   ```
4. **Configure MAAS DHCP Dynamic Range**
- **Input the DHCP Dynamic Range Name**
  Enter the name of the DHCP dynamic range used for Network PXE boot:
  ```
  Enter DHCP subnet in CIDR notation: 172.20.10.5/28
  ```
- **Input the DHCP Dynamic Start IP**
  Specify the start IP of the DHCP dynamic range:
  ```
  MAAS DHCP dynamic start IP: 172.20.10.8
  ```
- **Input the DHCP Dynamic End IP**
  Specify the end IP of the DHCP dynamic range:
  ```
  MAAS DHCP dynamic end IP: 172.20.10.12
  ```

5. **Input Juju VM Fixed IP**
Provide the static IP address that Juju VM will use:
  ```
  Enter the IP that juju-vm will used: 172.20.10.6
  ```

## One-Time Installation

For a one-time installation, follow these steps:

1. **Prepare an Existing Network Bridge**
Ensure that you have an existing network bridge configured on your system. This bridge will be used for network connectivity during the installation process.

2. **Prepare the Configuration File**
Create a configuration file (you can name it anything, e.g., `install.cfg`) with the required parameters. Here is an example of what the configuration file should look like:

```cfg
## Otterscale endpoint
OTTERSCALE_ENDPOINT="http://127.0.0.1:8299"

## Type network bridge name that will be used
OTTERSCALE_CONFIG_BRIDGE_CIDR="172.20.10.5/28"

## Type the ip range for network pxe boot used
OTTERSCALE_CONFIG_MAAS_DHCP_CIDR="172.20.10.5/28"
OTTERSCALE_CONFIG_MAAS_DHCP_START_IP="172.20.10.8"
OTTERSCALE_CONFIG_MAAS_DHCP_END_IP="172.20.10.12"

## Type JuJu IP
OTTERSCALE_CONFIG_JUJU_IP="172.20.10.6"

## MAAS configure
OTTERSCALE_CONFIG_MAAS_ADMIN_USER="admin"
OTTERSCALE_CONFIG_MAAS_ADMIN_PASS="admin"
OTTERSCALE_CONFIG_MAAS_ADMIN_EMAIL="admin@example.com"
  ```

3. Run the Installation Script with Configuration File Execute the installation script with the configuration file using the --config flag:
  ```
$ sudo bash install.sh --config=./install.cfg
  ```