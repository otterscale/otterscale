#!/bin/bash

# =============================================================================
# OtterScale Installation Script - Optimized Version
# This script installs and configures OtterScale with MAAS, Juju, LXD, and MicroK8s
# Author: OtterScale Team
# Version: 2.0
# =============================================================================

set -euo pipefail

# =============================================================================
# CONSTANTS AND CONFIGURATION
# =============================================================================

# System requirements
readonly MIN_MEMORY_GB=8
readonly MIN_DISK_GB=100
readonly OTTERSCALE_OS="24.04"
readonly OTTERSCALE_MAAS_VERSION="2.0"
readonly OTTERSCALE_BASE_IMAGE="ubuntu@24.04"

# LXD configuration
readonly LXD_STORAGE_SIZE_GB="60GB"
readonly LXD_CORES=2
readonly LXD_MEMORY_MB=4096
readonly LXD_DISK_GB="50G"

# Package configuration
readonly APT_PACKAGES="jq openssh-server bridge-utils openvswitch-switch"
readonly SNAP_PACKAGES="core24 maas maas-test-db juju lxd microk8s"

# Snap channel versions
readonly CORE24_CHANNEL="latest/stable"
readonly MAAS_CHANNEL="3.6/stable"
readonly MAAS_DB_CHANNEL="3.6/stable"
readonly JUJU_CHANNEL="3.6/stable"
readonly LXD_CHANNEL="5.0/stable"
readonly MICROK8S_CHANNEL="1.33/stable"
readonly CONTROLLER_CHARM_CHANNEL="3.6/stable"

# OtterScale configuration
readonly OTTERSCALE_MAX_RETRIES=5
readonly OTTERSCALE_CHARMHUB_URL="https://api.charmhub.io"
readonly OTTERSCALE_MAAS_ADMIN_USER="admin"
readonly OTTERSCALE_MAAS_ADMIN_PASS="admin"
readonly OTTERSCALE_MAAS_ADMIN_EMAIL="admin@example.com"
readonly OTTERSCALE_INSTALL_DIR=$(dirname "$(readlink -f "$0")")

# Runtime variables
readonly TEMP_LOG=$(mktemp)
readonly LOG="$OTTERSCALE_INSTALL_DIR/setup.log"
OTTERSCALE_ENDPOINT=""
NON_ROOT_USER=""
OTTERSCALE_BRIDGE_NAME="br-otters"
OTTERSCALE_INTERFACE_IP=""
OTTERSCALE_CIDR=""
MAAS_DHCP_START_IP=""
MAAS_DHCP_END_IP=""

# =============================================================================
# UTILITY FUNCTIONS
# =============================================================================

# Initialize logging system
init_logging() {
    touch "$LOG"
    chmod 644 "$LOG"
    exec 3>&1 4>&2
    trap cleanup EXIT INT TERM
}

# Enhanced logging with structured output
log() {
    local level="$1"
    local message="$2"
    local phase="${3:-GENERAL}"
    local timestamp
    timestamp="$(date '+%Y-%m-%d %H:%M:%S')"

    local formatted_message
    formatted_message=$(printf "%s [%s] %-15s %s\n" "$timestamp" "$level" "$phase" "$message")

    # local formatted_message="$timestamp [$level] [$phase] $message"
    echo "$formatted_message" | tee -a "$LOG"

    # Send status update (non-debug messages only)
    if [[ "$level" != "DEBUG" && -n "$OTTERSCALE_ENDPOINT" ]]; then
        send_status_data "$phase" "$message" || true
    fi
}

# Enhanced error handling with cleanup
error_exit() {
    local message="$1"
    local exit_code="${2:-1}"

    log "ERROR" "$message" "ERROR"

    # Show detailed error output if available
    if [[ -s "$TEMP_LOG" ]]; then
        log "DEBUG" "Detailed error output:" "ERROR"
        while IFS= read -r line; do
            log "DEBUG" "$line" "ERROR"
        done < "$TEMP_LOG"
    fi

    cleanup
    exit "$exit_code"
}

# Cleanup function
cleanup() {
    log "INFO" "Performing cleanup..." "CLEANUP"

    # Remove temporary files
    if [[ -f "$TEMP_LOG" ]]; then
        rm -f "$TEMP_LOG"
    fi

    # Restore file descriptors
    exec 2>&4 1>&3

    log "INFO" "Cleanup completed" "CLEANUP"
}

# Execute command with proper error handling
execute_cmd() {
    local cmd="$1"
    local description="$2"
    local phase="${3:-COMMAND}"

    log "INFO" "Executing: $cmd" "$phase"

    if ! eval "$cmd" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to $description"
    fi
}

# Execute command as specific user
execute_as_user() {
    local user="$1"
    local cmd="$2"
    local description="${3:-execute command as $user}"

    log "INFO" "Executing as $user: $cmd" "USER_COMMAND"

    if ! su "$user" -c "$cmd" >>"$TEMP_LOG" 2>&1; then
        log "WARN" "Failed to $description" "USER_COMMAND"
        return 1
    fi

    return 0
}

# =============================================================================
# VALIDATION FUNCTIONS
# =============================================================================

# Validate IP address format
validate_ip() {
    local ip="$1"

    if [[ ! $ip =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
        return 1
    fi

    IFS='.' read -ra octets <<< "$ip"
    for octet in "${octets[@]}"; do
        if ((octet > 255)); then
            return 1
        fi
    done

    return 0
}

# Validate port number
validate_port() {
    local port="$1"

    if [[ ! $port =~ ^[0-9]+$ ]] || ((port < 1 || port > 65535)); then
        return 1
    fi

    return 0
}

# Validate URL format
validate_url() {
    local url="$1"
    local ip port

    if [[ $url =~ ^https?://([^:]+):([0-9]+) ]]; then
        ip="${BASH_REMATCH[1]}"
        port="${BASH_REMATCH[2]}"
    else
        echo "URL format not recognized: $url"
        return 1
    fi

    if ! validate_ip "$ip"; then
        echo "Invalid IP in URL: $ip"
        return 1
    fi

    if ! validate_port "$port"; then
        echo "Invalid port in URL: $port"
        return 1
    fi

    return 0
}

# Check if IP is in network range
is_ip_in_network() {
    local ip=$1
    local network=$2
    local mask=$3
    local ip_number=$(ip_to_number $ip)
    local network_number=$(network_to_number $network $mask)
    local mask_number=$(ip_to_number $mask)

    if [[ $((ip_number & mask_number)) -eq "$network_number" ]]; then
        return 0
    fi

    return 1
}

# Convert IP to number for calculations
ip_to_number() {
    local ip="$1"
    IFS='.' read -ra octets <<< "$ip"
    echo $((octets[0] * 256**3 + octets[1] * 256**2 + octets[2] * 256 + octets[3]))
}

# Convert network and mask to number
network_to_number() {
    local network="$1"
    local mask="$2"
    IFS='.' read -ra net_octets <<< "$network"
    IFS='.' read -ra mask_octets <<< "$mask"

    local network_number=0
    for i in {0..3}; do
        network_number=$((network_number + (net_octets[i] & mask_octets[i]) * 256**(3-i)))
    done

    echo "$network_number"
}

# Check IP range validity
check_ip_range() {
    local network subnet_cidr mask_dotted
    network=$(echo "$MAAS_NETWORK_SUBNET" | cut -d'/' -f1)
    subnet_cidr=$(echo "$MAAS_NETWORK_SUBNET" | cut -d'/' -f2)

    # Convert CIDR to dotted decimal mask
    mask_dotted=$(printf "%d.%d.%d.%d" \
        $((0xFF << (32 - subnet_cidr) >> 24 & 0xFF)) \
        $((0xFF << (32 - subnet_cidr) >> 16 & 0xFF)) \
        $((0xFF << (32 - subnet_cidr) >> 8 & 0xFF)) \
        $((0xFF << (32 - subnet_cidr) & 0xFF)))

    # Check if both IPs are in the network
    if is_ip_in_network "$MAAS_DHCP_START_IP" "$network" "$mask_dotted" && \
       is_ip_in_network "$MAAS_DHCP_END_IP" "$network" "$mask_dotted"; then
        log "INFO" "IP range $MAAS_DHCP_START_IP-$MAAS_DHCP_END_IP is valid for network $MAAS_NETWORK_SUBNET" "VALIDATION"
        return 0
    else
        error_exit "IP range $MAAS_DHCP_START_IP-$MAAS_DHCP_END_IP is not within the network $MAAS_NETWORK_SUBNET"
    fi
}

# =============================================================================
# SYSTEM VALIDATION FUNCTIONS
# =============================================================================

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        error_exit "This script must be run as root"
    fi
    log "INFO" "Root access validated" "VALIDATION"
}

# Validate Ubuntu version
check_os() {
    local os_id os_release
    os_id=$(lsb_release -si)
    os_release=$(lsb_release -sr)

    if [[ "$os_id" != "Ubuntu" ]]; then
        error_exit "This script requires Ubuntu $OTTERSCALE_OS, found: $os_id"
    fi

    if [[ "$os_release" != "$OTTERSCALE_OS" ]]; then
        error_exit "This script requires Ubuntu version $OTTERSCALE_OS, found: $os_release"
    fi

    log "INFO" "OS validation passed: $os_id $os_release" "VALIDATION"
}

# Check memory requirements
check_memory() {
    local memory_gb
    memory_gb=$(free -g | awk '/^Mem:/ {print $2}')

    if ((memory_gb < MIN_MEMORY_GB)); then
        error_exit "Insufficient memory. Required: ${MIN_MEMORY_GB}GB, Available: ${memory_gb}GB"
    fi

    log "INFO" "Memory validation passed: ${memory_gb}GB available" "VALIDATION"
}

# Check disk space requirements
check_disk() {
    local disk_gb
    disk_gb=$(df -BG / | awk 'NR==2 {print $4}' | tr -d 'G')

    if ((disk_gb < MIN_DISK_GB)); then
        error_exit "Insufficient disk space. Required: ${MIN_DISK_GB}GB, Available: ${disk_gb}GB"
    fi

    log "INFO" "Disk space validation passed: ${disk_gb}GB available" "VALIDATION"
}

# Disable IPv6 temporarily
disable_ipv6() {
    log "INFO" "Temporarily disabling IPv6 (restored after reboot)" "SYSTEM_CONFIG"
    sysctl -w net.ipv6.conf.all.disable_ipv6=1 >/dev/null 2>&1
    sysctl -w net.ipv6.conf.default.disable_ipv6=1 >/dev/null 2>&1
}

# =============================================================================
# COMMUNICATION FUNCTIONS
# =============================================================================

# Send HTTP request to OtterScale endpoint
send_request() {
    local url_path="$1"
    local data="$2"
    local max_retries="${3:-3}"
    local retry_count=1

    while ((retry_count <= max_retries)); do
        if curl -s --max-time 30 \
                --header "Content-Type: application/json" \
                --data "$data" \
                "$OTTERSCALE_ENDPOINT$url_path" >/dev/null 2>&1; then
            return 0
        fi

        echo "HTTP request failed (attempt $retry_count/$max_retries)"
        retry_count=$((retry_count+1))
    done

    echo "Failed to send HTTP request after $max_retries attempts"
    exit 1
}

# Send status update to OtterScale
send_status_data() {
    local phase="$1"
    local message="$2"
    local new_url="${3:-}"

    # Skip if endpoint is not configured
    if [[ -z "$OTTERSCALE_ENDPOINT" ]]; then
        return 0
    fi

    local data
    data=$(cat <<EOF
{
    "phase": "$phase",
    "message": "$message",
    "new_url": "$new_url"

}
EOF
)

    send_request "/otterscale.environment.v1.EnvironmentService/UpdateStatus" "$data"
}

# =============================================================================
# PACKAGE MANAGEMENT FUNCTIONS
# =============================================================================

check_curl() {
    if ! command -v curl &> /dev/null; then
        echo "Please apt install curl first"
        exit 1
    fi
}

# Update APT package lists
apt_update() {
    log "INFO" "Updating APT package lists..." "APT_UPDATE"

    if ! apt-get update --fix-missing >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update APT package lists - check network connectivity"
    fi

    log "INFO" "APT package lists updated successfully" "APT_UPDATE"
}

# Install APT packages
apt_install() {
    local package_list="$1"

    log "INFO" "Installing APT packages: $package_list" "APT_INSTALL"

    if ! DEBIAN_FRONTEND=noninteractive apt-get install -y $package_list >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to install APT packages: $package_list"
    fi

    log "INFO" "APT packages installed successfully" "APT_INSTALL"
}

# Get current snap channel
get_snap_channel() {
    local snap_name="$1"
    snap list | grep "^${snap_name}[[:space:]]" | awk '{print $4}'
}

# Install snap package with retry mechanism
retry_snap_install() {
    local snap_name="$1"
    local max_retries="$2"
    local snap_options="$3"
    local retry_count=0

    while ((retry_count < max_retries)); do
        log "INFO" "Installing snap $snap_name (attempt $((retry_count + 1))/$max_retries)" "SNAP_INSTALL"

        if snap install "$snap_name" $snap_options >>"$TEMP_LOG" 2>&1; then
            log "INFO" "Successfully installed snap: $snap_name" "SNAP_INSTALL"
            return 0
        fi

        log "WARN" "Failed to install snap $snap_name (attempt $retry_count/$max_retries)" "SNAP_INSTALL"
        retry_count=$((retry_count+1))
    done

    error_exit "Failed to install snap $snap_name after $max_retries attempts"
}

# Refresh snap package with retry mechanism
retry_snap_refresh() {
    local snap_name="$1"
    local snap_channel="$2"
    local max_retries="$3"
    local retry_count=0

    while ((retry_count < max_retries)); do
        log "INFO" "Refreshing snap $snap_name to $snap_channel (attempt $((retry_count + 1))/$max_retries)" "SNAP_REFRESH"

        if snap refresh "$snap_name" --channel="$snap_channel" >>"$TEMP_LOG" 2>&1; then
            log "INFO" "Successfully refreshed snap: $snap_name" "SNAP_REFRESH"
            return 0
        fi

        log "WARN" "Failed to refresh snap $snap_name (attempt $retry_count/$max_retries)" "SNAP_REFRESH"
        retry_count=$((retry_count+1))
    done

    error_exit "Failed to refresh snap $snap_name after $max_retries attempts"
}

# Install or update snap package
install_or_update_snap() {
    local snap_name="$1"
    local snap_channel="$2"

    if snap list | grep -q "^${snap_name}[[:space:]]"; then
        local current_channel
        current_channel=$(get_snap_channel "$snap_name")

        if [[ "$current_channel" != "$snap_channel" ]]; then
            retry_snap_refresh "$snap_name" "$snap_channel" "$OTTERSCALE_MAX_RETRIES"
        else
            log "INFO" "Snap $snap_name already at correct channel: $snap_channel" "SNAP_CHECK"
        fi
    else
        local snap_options=""
        if [[ "$snap_name" == "microk8s" ]]; then
            snap_options="--classic --channel=$snap_channel"
        else
            snap_options="--channel=$snap_channel"
        fi

        retry_snap_install "$snap_name" "$OTTERSCALE_MAX_RETRIES" "$snap_options"
    fi
}

# Install all required snap packages
snap_install() {
    log "INFO" "Installing snap packages..." "SNAP_INSTALL"

    # Define snap channels
    declare -A snap_channels=(
        ["core24"]="$CORE24_CHANNEL"
        ["maas"]="$MAAS_CHANNEL"
        ["maas-test-db"]="$MAAS_DB_CHANNEL"
        ["juju"]="$JUJU_CHANNEL"
        ["lxd"]="$LXD_CHANNEL"
        ["microk8s"]="$MICROK8S_CHANNEL"
    )

    # Install each snap package
    for snap_name in $SNAP_PACKAGES; do
        install_or_update_snap "$snap_name" "${snap_channels[$snap_name]}"
    done

    # Hold all snaps to prevent auto-updates
    log "INFO" "Setting snap packages to hold..." "SNAP_HOLD"
    if ! snap refresh --hold >>"$TEMP_LOG" 2>&1; then
        log "WARN" "Failed to hold snap packages" "SNAP_HOLD"
    else
        log "INFO" "Snap packages held successfully" "SNAP_HOLD"
    fi
}

# =============================================================================
# NETWORK MANAGEMENT FUNCTIONS
# =============================================================================

# Get default network interface
get_default_interface() {
    CURRENT_INTERFACE=$(ip route show default | awk '{print $5}' | head -n 1)
    if [[ -z $CURRENT_INTERFACE ]]; then
        error_exit "Default route network interface not found"
    fi
    log "INFO" "Detected default interface: $CURRENT_INTERFACE" "NETWORK"
}

# Get default gateway
get_default_gateway() {
    CURRENT_GATEWAY=$(ip route show default | awk '{print $3}' | head -n 1)
    if [[ -z $CURRENT_GATEWAY ]]; then
        error_exit "Default gateway not found"
    fi
    log "INFO" "Detected default gateway: $CURRENT_GATEWAY" "NETWORK"
}

# Get interface CIDR
get_default_cidr() {
    CURRENT_CIDR=$(ip -o -4 addr show dev "$CURRENT_INTERFACE" | awk '{print $4}')
    if [[ -z $CURRENT_CIDR ]]; then
        error_exit "Interface $CURRENT_INTERFACE CIDR not found"
    fi
    log "INFO" "Detected CIDR for $CURRENT_INTERFACE: $CURRENT_CIDR" "NETWORK"
}

# Get DNS servers
get_default_dns() {
    local interface="$1"
    CURRENT_DNS=$(resolvectl status "$interface" | grep "DNS Servers" | awk '{print $3}' | head -1)
    if [[ -z $CURRENT_DNS ]]; then
        log "WARN" "No DNS found for $interface, using fallback 8.8.8.8" "NETWORK"
        CURRENT_DNS="8.8.8.8"
    else
        log "INFO" "Detected DNS server for $interface: $CURRENT_DNS" "NETWORK"
    fi
}

# Select network bridge
select_bridge() {
    local bridges
    readarray -t bridges < <(brctl show 2>/dev/null | awk 'NR>1 && $1!="" {print $1}')

    if [[ ${#bridges[@]} -eq 0 ]]; then
        return 1
    fi

    echo "Available network bridges:"
    for i in "${!bridges[@]}"; do
        echo "$((i+1))) ${bridges[$i]}"
    done

    while true; do
        read -p "Select bridge (1-${#bridges[@]}): " choice
        if [[ $choice =~ ^[0-9]+$ ]] && ((choice >= 1 && choice <= ${#bridges[@]})); then
            OTTERSCALE_BRIDGE_NAME="${bridges[$((choice-1))]}"
            log "INFO" "User selected bridge: $OTTERSCALE_BRIDGE_NAME" "NETWORK"
            return 0
        fi
        echo "Invalid selection. Please try again."
    done
}

# Prompt for bridge creation
prompt_bridge_creation() {
    if ! systemctl is-active --quiet NetworkManager; then
        error_exit "NetworkManager is not active, stop network bridge creation."
    fi

    get_default_interface
    get_default_gateway
    get_default_cidr
    get_default_dns "$CURRENT_INTERFACE"

    log "INFO" "Create network bridge $OTTERSCALE_BRIDGE_NAME" "NETWORK"
    local CURRENT_CONNECTION=$(nmcli -t -f NAME,DEVICE connection show --active | grep "$CURRENT_INTERFACE" | cut -d: -f1)
    if ! nmcli connection add type bridge ifname "$OTTERSCALE_BRIDGE_NAME" con-name "$OTTERSCALE_BRIDGE_NAME" \
        ipv4.method manual \
        ipv4.addresses "$CURRENT_CIDR" \
        ipv4.gateway "$CURRENT_GATEWAY" \
        ipv4.dns "$CURRENT_DNS" > /dev/null; then
        error_exit "Failed network bridge creation"
    fi

    log "INFO" "Connect network bridge $OTTERSCALE_BRIDGE_NAME to interface $CURRENT_INTERFACE"
    nmcli connection add type bridge-slave con-name br-otters-slave ifname "$CURRENT_INTERFACE" master "$OTTERSCALE_BRIDGE_NAME" > /dev/null

    log "INFO" "Start up network bridge $OTTERSCALE_BRIDGE_NAME" "NETWORK"
    nmcli connection up "$OTTERSCALE_BRIDGE_NAME" > /dev/null && nmcli connection down "$CURRENT_CONNECTION" > /dev/null

    sleep 10
}


# Check/setup network bridge
check_bridge() {
    log "INFO" "Checking network bridge configuration..." "NETWORK"

    if ip link show "$OTTERSCALE_BRIDGE_NAME" &>/dev/null; then
        log "INFO" "Bridge $OTTERSCALE_BRIDGE_NAME exists" "NETWORK"
    else
        log "INFO" "Bridge $OTTERSCALE_BRIDGE_NAME not found" "NETWORK"

        if ! select_bridge; then
            prompt_bridge_creation
        fi
    fi

    # Get bridge network information
    OTTERSCALE_CIDR=$(ip -o -4 addr show dev "$OTTERSCALE_BRIDGE_NAME" | awk '{print $4}')
    OTTERSCALE_INTERFACE_IP=$(echo "$OTTERSCALE_CIDR" | cut -d'/' -f1)
    get_default_dns "$OTTERSCALE_BRIDGE_NAME"

    log "INFO" "Bridge configuration - IP: $OTTERSCALE_INTERFACE_IP, CIDR: $OTTERSCALE_CIDR" "NETWORK"
}

# =============================================================================
# USER MANAGEMENT FUNCTIONS
# =============================================================================

# Find first non-root user
find_first_non_root_user() {
    for user_home in /home/*; do
        if [[ -d "$user_home" ]]; then
            NON_ROOT_USER=$(basename "$user_home")
            break
        fi
    done

    if [[ -z "$NON_ROOT_USER" ]]; then
        error_exit "No non-root user found"
    fi

    log "INFO" "Using non-root user: $NON_ROOT_USER" "USER_SETUP"
}

# Generate SSH key for user
generate_ssh_key() {
    local ssh_dir="/home/$NON_ROOT_USER/.ssh"
    local private_key="$ssh_dir/id_rsa"

    if [[ -f "$private_key" ]]; then
        log "INFO" "SSH key already exists for user $NON_ROOT_USER" "SSH_SETUP"
        return 0
    fi

    log "INFO" "Generating SSH key for user $NON_ROOT_USER" "SSH_SETUP"

    if ! su "$NON_ROOT_USER" -c "mkdir -p '$ssh_dir' && ssh-keygen -q -t rsa -N '' -f '$private_key'"; then
        error_exit "SSH key generation failed"
    fi

    # Set proper permissions
    chown -R "$NON_ROOT_USER:$NON_ROOT_USER" "$ssh_dir"
    chmod 700 "$ssh_dir"
    chmod 600 "$private_key"
    chmod 644 "${private_key}.pub"

    log "INFO" "SSH key generated successfully" "SSH_SETUP"
}

# Add SSH key to MAAS
add_key_to_maas() {
    local public_key="/home/$NON_ROOT_USER/.ssh/id_rsa.pub"

    # Check if keys already exist
    local key_count
    key_count=$(maas admin sshkeys read | jq -r 'length')

    if ((key_count > 0)); then
        log "INFO" "SSH keys already exist in MAAS" "SSH_SETUP"
        return 0
    fi

    log "INFO" "Adding SSH key to MAAS" "SSH_SETUP"

    if ! maas admin sshkeys create key="$(cat "$public_key")" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to add SSH key to MAAS"
    fi

    log "INFO" "SSH key added to MAAS successfully" "SSH_SETUP"
}

# =============================================================================
# MAAS MANAGEMENT FUNCTIONS
# =============================================================================

# Initialize MAAS
init_maas() {
    log "INFO" "Initializing MAAS..." "MAAS_INIT"

    # Check if MAAS admin user already exists
    if maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "INFO" "MAAS is already initialized (user $OTTERSCALE_MAAS_ADMIN_USER exists). Skipping initialization" "MAAS_INIT"
        return 0
    fi

    execute_cmd "maas init region+rack --database-uri maas-test-db:/// --maas-url http://$OTTERSCALE_INTERFACE_IP:5240/MAAS" "initialize MAAS"
    log "INFO" "MAAS initialized successfully" "MAAS_INIT"
}

# Create MAAS admin user
create_maas_admin() {
    log "INFO" "Creating MAAS admin user..." "MAAS_ADMIN"

    if maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "INFO" "Admin user '$OTTERSCALE_MAAS_ADMIN_USER' already exists. Using existing credentials" "MAAS_ADMIN"
    else
        execute_cmd "maas createadmin --username $OTTERSCALE_MAAS_ADMIN_USER --password $OTTERSCALE_MAAS_ADMIN_PASS --email $OTTERSCALE_MAAS_ADMIN_EMAIL" "create MAAS admin user"
    fi

    log "INFO" "MAAS web URL: http://$OTTERSCALE_INTERFACE_IP:5240/MAAS" "MAAS_ADMIN"
    log "INFO" "MAAS Username: $OTTERSCALE_MAAS_ADMIN_USER" "MAAS_ADMIN"
    log "INFO" "MAAS Password: $OTTERSCALE_MAAS_ADMIN_PASS" "MAAS_ADMIN"
}

# Login to MAAS with retry mechanism
login_maas() {
    log "INFO" "Attempting to login to MAAS..." "MAAS_LOGIN"
    local retry_count=0

    APIKEY=$(maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER")
    sleep 10

    while ((retry_count < OTTERSCALE_MAX_RETRIES)); do
        if maas login admin "http://localhost:5240/MAAS/" "$APIKEY" >>"$TEMP_LOG" 2>&1; then
            log "INFO" "MAAS login successful" "MAAS_LOGIN"
            return 0
        else
            log "WARN" "Failed to login to MAAS, retrying in 10 seconds (attempt $((retry_count + 1)))" "MAAS_LOGIN"
            retry_count=$((retry_count+1))
            sleep 10
        fi
    done

    error_exit "Failed to login to MAAS after $OTTERSCALE_MAX_RETRIES attempts"
}

# Get MAAS DNS configuration
get_maas_dns() {
    maas_current_dns=$(maas admin maas get-config name=upstream_dns | jq -r)

    if [[ "$maas_current_dns" == null ]]; then
        dns_value="$CURRENT_DNS"
        log "INFO" "MAAS upstream DNS not set, will use system DNS: $dns_value" "MAAS_CONFIG"
    else
        dns_value=$maas_current_dns
    fi
}

# Set MAAS configuration
set_config() {
    local name="$1"
    local value="$2"

    log "INFO" "Setting MAAS config: $name = $value" "MAAS_CONFIG"
    execute_cmd "maas admin maas set-config name=$name value=$value" "set MAAS $name config"
}

# Update MAAS configuration values
update_maas_config() {
    log "INFO" "Updating MAAS configuration..." "MAAS_CONFIG"

    get_maas_dns
    set_config "upstream_dns" "$dns_value"
    set_config "boot_images_auto_import" "false"
    set_config "enable_http_proxy" "false"
    set_config "enable_analytics" "false"
    set_config "network_discovery" "disabled"
    set_config "release_notifications" "false"

    log "INFO" "MAAS configuration updated successfully" "MAAS_CONFIG"
}

# Download MAAS images
download_maas_img() {
    log "INFO" "Configuring MAAS boot sources..." "MAAS_IMAGES"

    local boot_sources
    boot_sources=$(maas admin boot-sources read)
    local boot_source_count
    boot_source_count=$(echo "$boot_sources" | jq '. | length')

    if ((boot_source_count > 0)); then
        update_boot_source
    else
        create_boot_source
    fi

    start_import
    set_config "commissioning_distro_series" "noble"
    set_config "default_distro_series" "noble"
    set_config "default_osystem" "ubuntu"

    log "INFO" "MAAS images configured successfully" "MAAS_IMAGES"
}

# Update existing boot source
update_boot_source() {
    local boot_source_id boot_selection_id
    boot_source_id=$(echo "$boot_sources" | jq -r '.[0].id')
    boot_selection_id=$(maas admin boot-source-selections read "$boot_source_id" | jq -r '.[0].id')

    log "INFO" "Modifying existing boot source (ID: $boot_source_id)" "MAAS_IMAGES"

    if ! maas admin boot-source-selection update "$boot_source_id" "$boot_selection_id" \
        release=noble arches=amd64 subarches="*" labels="*" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update MAAS boot source"
    fi
}

# Create new boot source
create_boot_source() {
    log "INFO" "Creating new boot source for Ubuntu Noble (24.04) amd64..." "MAAS_IMAGES"

    if ! maas admin boot-sources create \
        url="http://images.maas.io/ephemeral-v3/stable/" \
        keyring_filename="/usr/share/keyrings/ubuntu-cloudimage-keyring.gpg" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to create MAAS boot source"
    fi

    # Get the new source ID
    local boot_source_id
    boot_source_id=$(maas admin boot-sources read | jq -r '.[0].id')

    # Create selection for Noble amd64
    if ! maas admin boot-source-selections create "$boot_source_id" \
        release=noble arches=amd64 subarches="*" labels="*" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to create MAAS boot source selection"
    fi
}

# Start image import
start_import() {
    log "INFO" "Starting MAAS boot image download..." "MAAS_IMAGES"

    maas admin boot-resources stop-import >>"$TEMP_LOG" 2>&1 || true
    sleep 10

    execute_cmd "maas admin boot-resources import" "start MAAS image download"
    sleep 10

    log "INFO" "Waiting for image download to complete..." "MAAS_IMAGES"
    while true; do
        local is_importing
        is_importing=$(maas admin boot-resources is-importing | jq -r)

        if [[ "$is_importing" != "true" ]]; then
            break
        fi

        sleep 10
    done

    log "INFO" "MAAS image download completed" "MAAS_IMAGES"
}

enter_dhcp_start_ip() {
    while true; do
        read -p "Enter MAAS dhcp start IP: " MAAS_DHCP_START_IP
        if validate_ip "$MAAS_DHCP_START_IP"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

enter_dhcp_end_ip() {
    while true; do
        read -p "Enter MAAS dhcp end IP: " MAAS_DHCP_END_IP
        if validate_ip "$MAAS_DHCP_END_IP"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

get_dhcp_subnet_and_ip() {
    if [[ -z "$MAAS_DHCP_START_IP" ]]; then
        enter_dhcp_start_ip
    else
        log "INFO" "MAAS dhcp start ip: $MAAS_DHCP_START_IP" "MAAS_DHCP"
    fi

    if [[ -z "$MAAS_DHCP_END_IP" ]]; then
        enter_dhcp_end_ip
    else
        log "INFO" "MAAS dhcp start ip: $MAAS_DHCP_END_IP" "MAAS_DHCP"
    fi
}

update_fabric_dns() {
    local FABRIC_DNS=$(maas admin subnet read "$MAAS_NETWORK_SUBNET" | jq -r '.dns_servers[]')
    log "INFO" "Update MAAS dns $CURRENT_DNS to fabric $MAAS_NETWORK_SUBNET" "MAAS_DNS"

    if [[ "$FABRIC_DNS" =~ "$CURRENT_DNS" ]]; then
        log "INFO" "Current dns already existed, skipping..." "MAAS_DNS"
    elif [ ! -n "$maas_current_dns" ]; then
        if [ -z "$FABRIC_DNS" ]; then
            dns_value="$CURRENT_DNS"
        else
            dns_value="$FABRIC_DNS $CURRENT_DNS"
        fi
    fi

    execute_cmd "maas admin subnet update $MAAS_NETWORK_SUBNET dns_servers=$dns_value" "update maas dns to fabric"
}

get_fabric() {
    log "INFO" "Getting fabric and VLAN information..." "MAAS_CONFIG"
    FABRIC_ID=$(maas admin subnet read "$MAAS_NETWORK_SUBNET" | jq -r ".vlan.fabric_id")
    VLAN_TAG=$(maas admin subnet read "$MAAS_NETWORK_SUBNET" | jq -r ".vlan.vid")
    PRIMARY_RACK=$(maas admin rack-controllers read | jq -r ".[] | .system_id")
    if [ -z "$FABRIC_ID" ] || [ -z "$VLAN_TAG" ] || [ -z "$PRIMARY_RACK" ]; then
        error_exit "Failed to get network configuration details"
    fi
}

create_dhcp_iprange() {
    log "INFO" "Enable MAAS dhcp" "MAAS_DHCP"
    if ! maas admin ipranges create type=dynamic start_ip="$MAAS_DHCP_START_IP" end_ip="$MAAS_DHCP_END_IP" >>"$TEMP_LOG" 2>&1; then
        log "WARN" "Please confirm if address is within subnet $MAAS_NETWORK_SUBNET, or maybe it conflicts with an existing IP address or range" "MAAS_DHCP"
        error_exit "Failed to create DHCP range"
    fi
}

update_dhcp_config() {
    local ENABLED=$1
    log "INFO" "Set MAAS VLAN DHCP to $ENABLED" "MAAS_DHCP"
    if ! maas admin vlan update $FABRIC_ID $VLAN_TAG dhcp_on=$ENABLED primary_rack=$PRIMARY_RACK >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to set MAAS DHCP to $ENABLED"
    fi
}

# Enable MAAS DHCP
enable_maas_dhcp() {
    log "INFO" "Configuring MAAS DHCP..." "MAAS_DHCP"

    local ip_ranges_count
    ip_ranges_count=$(maas admin ipranges read | jq '. | length')

    if ((ip_ranges_count > 0)); then
        log "INFO" "MAAS already has dynamic IP ranges - skipping DHCP configuration" "MAAS_DHCP"
        return 0
    fi

    MAAS_NETWORK_SUBNET=$OTTERSCALE_CIDR
    while true; do
        get_dhcp_subnet_and_ip
        if check_ip_range ; then
            update_fabric_dns
            get_fabric
            create_dhcp_iprange
            update_dhcp_config "True"
            log "INFO" "DHCP configuration completed" "MAAS_DHCP"
            break
        else
            if [ -n "$MAAS_NETWORK_SUBNET" ] && [ -n "$MAAS_DHCP_START_IP" ] && [ -n "$MAAS_DHCP_END_IP" ]; then
                break
            fi
        fi
    done
}

# =============================================================================
# LXD MANAGEMENT FUNCTIONS
# =============================================================================

# Generate LXD pre-seed configuration
generate_lxd_config() {
    local lxd_file="$1"

    log "INFO" "Generating LXD pre-seed configuration: $lxd_file" "LXD_INIT"

    cat > "$lxd_file" <<EOF
config:
  core.https_address: '[::]:8443'
  core.trust_password: password
storage_pools:
- config:
    size: $LXD_STORAGE_SIZE_GB
  description: ""
  name: default
  driver: zfs
profiles:
- name: default
  config:
    boot.autostart: "true"
  description: ""
  devices:
    eth0:
      name: eth0
      nictype: bridged
      parent: $OTTERSCALE_BRIDGE_NAME
      type: nic
    root:
      path: /
      pool: default
      type: disk
projects: []
cluster: null
EOF

    log "INFO" "LXD pre-seed configuration generated" "LXD_INIT"
}

# Initialize LXD
init_lxd() {
    local lxd_file="$OTTERSCALE_INSTALL_DIR/lxd-config.yaml"

    log "INFO" "Initializing LXD with bridge $OTTERSCALE_BRIDGE_NAME..." "LXD_INIT"

    generate_lxd_config "$lxd_file"

    if ! cat "$lxd_file" | lxd init --preseed >>"$TEMP_LOG" 2>&1; then
        error_exit "LXD initialization failed"
    else
        log "INFO" "LXD initialized successfully" "LXD_INIT"
        rm -f "$lxd_file"
    fi
}

# Create LXD project for MAAS
create_maas_lxd_project() {
    log "INFO" "Creating LXD project for MAAS..." "LXD_PROJECT"

    if ! lxc project list --format json | jq --exit-status '.[] | select(.name == "maas")' >>"$TEMP_LOG" 2>&1; then
        if ! lxc project create maas >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed to create LXD project 'maas'"
        fi
        log "INFO" "Created LXD project 'maas'" "LXD_PROJECT"
    else
        log "INFO" "LXD project 'maas' already exists" "LXD_PROJECT"
    fi

    if ! lxc profile show default | lxc profile edit default --project maas >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update LXD profile for MAAS project"
    fi

    log "INFO" "LXD project configuration completed" "LXD_PROJECT"
}

# Create LXD VM host in MAAS
create_lxd_vm() {
    log "INFO" "Creating LXD VM host in MAAS..." "LXD_VM"

    local vm_hosts
    vm_hosts=$(maas admin vm-hosts read)
    local vm_host_count
    vm_host_count=$(echo "$vm_hosts" | jq '. | length')

    if ((vm_host_count > 0)); then
        log "INFO" "Found existing VM hosts, checking resources..." "LXD_VM"
        search_available_vmhost "$vm_hosts"
    else
        log "INFO" "Creating new LXD VM host..." "LXD_VM"
        execute_cmd "maas admin vm-hosts create password=password type=lxd power_address=https://$OTTERSCALE_INTERFACE_IP:8443 project=maas" "create MAAS LXD VM host"
        VM_HOST_ID=$(maas admin vm-hosts read | jq -r '.[0].id')
    fi

    log "INFO" "LXD VM host (ID: $VM_HOST_ID) is ready" "LXD_VM"
}

# Search for available VM host
search_available_vmhost() {
    local vm_hosts="$1"

    while IFS= read -r host; do
        VM_HOST_ID=$(echo "$host" | jq -r '.id')
        local available_cores available_memory_gb available_disk_gb
        available_cores=$(echo "$host" | jq -r '.available.cores')
        available_memory_gb=$(echo "$host" | jq -r '.available.memory / 1024' | bc -l | xargs printf "%.2f\n")
        available_disk_gb=$(echo "$host" | jq -r '.available.local_storage / (1024*1024*1024)' | bc -l | xargs printf "%.2f\n")

        if [[ $(echo "$available_cores >= 1" | bc -l) -eq 1 ]] && \
           [[ $(echo "$available_memory_gb >= 4" | bc -l) -eq 1 ]] && \
           [[ $(echo "$available_disk_gb >= 8" | bc -l) -eq 1 ]]; then
            log "INFO" "Using existing VM host $VM_HOST_ID with sufficient resources" "LXD_VM"
            log "INFO" "Available resources - Cores: $available_cores, Memory: ${available_memory_gb}GB, Disk: ${available_disk_gb}GB" "LXD_VM"
            return 0
        fi
    done < <(echo "$vm_hosts" | jq -c '.[]')

    error_exit "No VM host with sufficient resources found"
}

# Create VM from MAAS
create_vm_from_maas() {
    log "INFO" "Creating VM from MAAS..." "VM_CREATE"
    local juju_machine_id

    # Check if juju-vm already exists
    juju_machine_id=$(maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm") | .system_id')
    if [[ ! -z $juju_machine_id ]]; then
        log "INFO" "Machine juju-vm (id: $juju_machine_id) already exists - skipping creation" "VM_CREATE"
        return 0
    fi

    log "INFO" "Creating VM from LXD host ID $VM_HOST_ID..." "VM_CREATE"
    juju_machine_id=$(maas admin vm-host compose "$VM_HOST_ID" cores="$LXD_CORES" memory="$LXD_MEMORY_MB" disk=1:size="$LXD_DISK_GB" | jq -r '.system_id')
    if [[ -z "$juju_machine_id" ]]; then
        error_exit "Failed to create VM from LXD host ID $VM_HOST_ID"
    fi

    wait_commissioning "$juju_machine_id"
}

# Wait for machine commissioning to complete
wait_commissioning() {
    local machine_id="$1"

    log "INFO" "Waiting for machine $machine_id to transition from commissioning to ready state..." "VM_CREATE"

    while true; do
        local machine_status
        machine_status=$(maas admin machine read "$machine_id" | jq -r '.status_name')

        case "$machine_status" in
            "Ready")
                log "INFO" "Machine $machine_id created successfully" "VM_CREATE"
                rename_machine "$machine_id" "juju-vm"
                break
                ;;
            "Failed commissioning"|"Failed testing")
                error_exit "Machine $machine_id failed: $machine_status"
                ;;
            *)
                log "INFO" "Machine $machine_id status: $machine_status" "VM_CREATE"
                sleep 10
                ;;
        esac
    done
}

# Rename machine in MAAS
rename_machine() {
    local machine_id="$1"
    local machine_name="$2"

    execute_cmd "maas admin machine update $machine_id hostname=$machine_name" "rename machine $machine_id to $machine_name"
    log "INFO" "Machine $machine_id renamed to $machine_name" "VM_CREATE"
}

# Continue with additional functions...

# =============================================================================
# JUJU MANAGEMENT FUNCTIONS
# =============================================================================

# Write the Juju cloud definition file (executed as the non‑root user)
generate_clouds_yaml() {
    log "INFO" "Generating Juju cloud definition file ($JUJU_CLOUD)" "JUJU_BOOTSTRAP"
    export OTTERSCALE_INTERFACE_IP=$OTTERSCALE_INTERFACE_IP
    su "$NON_ROOT_USER" -c 'cat > $JUJU_CLOUD <<EOF
clouds:
  maas-cloud:
    type: maas
    description: Metal As A Service
    auth-types: [oauth1]
    endpoint: http://$OTTERSCALE_INTERFACE_IP:5240/MAAS/api/2.0/
    regions:
      default: {}
EOF'
}

# Write the Juju credential definition file (executed as the non‑root user)
generate_credentials_yaml() {
    log "INFO" "Generating Juju credential file ($JUJU_CREDENTIAL)" "JUJU_BOOTSTRAP"
    export APIKEY=$APIKEY
    su "$NON_ROOT_USER" -c 'cat > $JUJU_CREDENTIAL <<EOF
credentials:
  maas-cloud:
    maas-cloud-credential:
      auth-type: oauth1
      maas-oauth: $APIKEY
EOF'
}

add_clouds_yaml() {
    log "INFO" "Check juju clouds" "JUJU_CONFIG"
    if su "$NON_ROOT_USER" -c 'juju clouds 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Juju cloud maas-cloud already exists – skipping creation" "JUJU_CLOUD"
    else
        log "INFO" "Adding MAAS cloud to Juju..." "JUJU_BOOTSTRAP"
        if ! su "$NON_ROOT_USER" -c 'juju add-cloud --client maas-cloud $JUJU_CLOUD >>/dev/null 2>&1'; then
            error_exit "Failed to add MAAS cloud to Juju"
        fi
    fi
}

add_generates_yaml() {
    log "INFO" "Check juju credentials" "JUJU_CONFIG"
    if su "$NON_ROOT_USER" -c 'juju credentials 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Juju credential for maas-cloud already exists – skipping creation" "JUJU_CREDENTIAL"
    else
        log "INFO" "Adding MAAS credential to Juju..." "JUJU_BOOTSTRAP"
        if ! su "$NON_ROOT_USER" -c 'juju add-credential --client maas-cloud -f $JUJU_CREDENTIAL >>/dev/null 2>&1'; then
            error_exit "Failed to add MAAS credentials to Juju"
        fi
    fi
}

# Bootstrap Juju controller
bootstrap_juju() {
    log "INFO" "Prepare Juju controller..." "JUJU_BOOTSTRAP"
    local juju_machine_id

    su "$NON_ROOT_USER" -c 'mkdir -p ~/otterscale-tmp'
    export JUJU_CLOUD=/home/$NON_ROOT_USER/otterscale-tmp/cloud.yaml
    export JUJU_CREDENTIAL=/home/$NON_ROOT_USER/otterscale-tmp/credential.yaml
    generate_clouds_yaml
    generate_credentials_yaml
    add_clouds_yaml
    add_generates_yaml
    rm -rf /home/"$NON_ROOT_USER"/otterscale-tmp

    juju_machine_id=$(maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm") | .system_id')
    if [[ -z "$juju_machine_id" ]]; then
        error_exit "juju-vm machine not found"
    fi

    if [[ $(maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | jq -r '.status_name') == Deployed ]]; then
        log "INFO" "Machine juju-vm is bootstrapped" "JUJU_BOOTSTRAP"
    else
        log "INFO" "Bootstrapping Juju controller..." "JUJU_BOOTSTRAP"
        local bootstrap_cmd="juju bootstrap maas-cloud maas-cloud-controller --bootstrap-base=$OTTERSCALE_BASE_IMAGE"
        local bootstrap_config="--config default-base=$OTTERSCALE_BASE_IMAGE --config bootstrap-timeout=7200 --controller-charm-channel=$CONTROLLER_CHARM_CHANNEL"
        if ! su "$NON_ROOT_USER" -c "$bootstrap_cmd $bootstrap_config --to juju-vm --debug"; then
            rm -rf /home/"$NON_ROOT_USER"/.local/share/juju
            error_exit "Failed to bootstrap Juju controller"
        fi
        log "INFO" "Juju controller bootstrapped successfully" "JUJU_BOOTSTRAP"
    fi
}

# Configure MicroK8s settings
prepare_microk8s_config() {
    log "INFO" "Configuring MicroK8s for user $NON_ROOT_USER..." "MICROK8S_CONFIG"

    usermod -aG microk8s "$NON_ROOT_USER"
    log "INFO" "Added $NON_ROOT_USER to microk8s group" "MICROK8S_CONFIG"

    local kube_folder="/home/$NON_ROOT_USER/.kube"
    if [[ ! -d "$kube_folder" ]]; then
        mkdir -p "$kube_folder"
    fi
    chown "$NON_ROOT_USER":"$NON_ROOT_USER" "$kube_folder"

    log "INFO" "Updating MicroK8s Calico DaemonSet for bridge $OTTERSCALE_BRIDGE_NAME" "MICROK8S_CONFIG"
    if [[ -z "$OTTERSCALE_BRIDGE_NAME" ]]; then
        error_exit "Bridge name is empty"
    fi

    if ! microk8s kubectl set env -n kube-system daemonset.apps/calico-node -c calico-node IP_AUTODETECTION_METHOD="interface=$OTTERSCALE_BRIDGE_NAME" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update MicroK8s Calico environment"
    fi
}

# Enable MicroK8s add-ons
enable_microk8s_option() {
    log "INFO" "Enabling MicroK8s add-ons..." "MICROK8S_ADDONS"

    if ! microk8s status --wait-ready >>"$TEMP_LOG" 2>&1; then
        error_exit "MicroK8s is not ready"
    fi

    log "INFO" "MicroK8s is ready" "MICROK8S_ADDONS"

    local kube_folder="/home/$NON_ROOT_USER/.kube"
    microk8s config > "$kube_folder/config"
    chown "$NON_ROOT_USER":"$NON_ROOT_USER" "$kube_folder/config"

    execute_cmd "microk8s enable dns" "enable MicroK8s DNS"
    execute_cmd "microk8s enable hostpath-storage" "enable MicroK8s hostpath-storage"
    execute_cmd "microk8s enable metallb:$OTTERSCALE_INTERFACE_IP-$OTTERSCALE_INTERFACE_IP" "enable MicroK8s MetalLB"
    execute_cmd "microk8s enable helm3" "enable MicroK8s helm3"

    log "INFO" "MicroK8s add-ons enabled successfully" "MICROK8S_ADDONS"
}

# Extend MicroK8s certificates
extend_microk8s_cert() {
    log "INFO" "Extending MicroK8s certificates to 3650 days..." "MICROK8S_CERT"

    local snap_ssl="/snap/microk8s/current/usr/bin/openssl"
    local snap_data="/var/snap/microk8s/current"

    log "INFO" "Generating server CSR for MicroK8s certificate" "MICROK8S_CERT"
    if ! "$snap_ssl" req -new -sha256 -key "$snap_data/certs/server.key" -out "$snap_data/certs/server.csr" -config "$snap_data/certs/csr.conf" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to generate server CSR for MicroK8s certificate"
    fi

    log "INFO" "Signing server certificate (3650 days)" "MICROK8S_CERT"
    if ! "$snap_ssl" x509 -req -sha256 -in "$snap_data/certs/server.csr" -CA "$snap_data/certs/ca.crt" -CAkey "$snap_data/certs/ca.key" -CAcreateserial -out "$snap_data/certs/server.crt" -days 3650 -extensions v3_ext -extfile "$snap_data/certs/csr.conf" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to sign server certificate for MicroK8s"
    fi

    log "INFO" "Generating front-proxy CSR" "MICROK8S_CERT"
    if ! "$snap_ssl" req -new -sha256 -key "$snap_data/certs/front-proxy-client.key" -out "$snap_data/certs/front-proxy-client.csr" -config <(sed '/^prompt = no/d' "$snap_data/certs/csr.conf") -subj "/CN=front-proxy-client" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to generate front-proxy CSR for MicroK8s certificate"
    fi

    log "INFO" "Signing front-proxy certificate (3650 days)" "MICROK8S_CERT"
    if ! "$snap_ssl" x509 -req -sha256 -in "$snap_data/certs/front-proxy-client.csr" -CA "$snap_data/certs/front-proxy-ca.crt" -CAkey "$snap_data/certs/front-proxy-ca.key" -CAcreateserial -out "$snap_data/certs/front-proxy-client.crt" -days 3650 -extensions v3_ext -extfile "$snap_data/certs/csr.conf" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to sign front-proxy certificate for MicroK8s"
    fi

    log "INFO" "MicroK8s certificates extended successfully" "MICROK8S_CERT"
}

# Add Kubernetes cluster to Juju
juju_add_k8s() {
    log "INFO" "Adding Kubernetes cluster to Juju..." "JUJU_K8S"

    local kubeconfig="/home/$NON_ROOT_USER/.kube/config"
    if [[ ! -f "$kubeconfig" ]]; then
        error_exit "Kubernetes config file not found at $kubeconfig"
    fi

    if ! su "$NON_ROOT_USER" -c "juju add-k8s cos-k8s --controller maas-cloud-controller --client" >>"$TEMP_LOG" 2>&1; then
        log "WARN" "Controller cos-k8s already exist" "Cos-k8s exist"
    fi

    if ! su "$NON_ROOT_USER" -c "juju show-model cos >/dev/null 2>&1"; then
        su "$NON_ROOT_USER" -c "juju add-model cos cos-k8s >/dev/null 2>&1"
        su "$NON_ROOT_USER" -c "juju deploy -m cos cos-lite --trust >/dev/null 2>&1"
        su "$NON_ROOT_USER" -c "juju deploy -m cos prometheus-scrape-target-k8s --channel=2/edge >/dev/null 2>&1"
    fi

    ## Config prometheus
    su "$NON_ROOT_USER" -c "juju config -m cos prometheus metrics_retention_time=180d >/dev/null 2>&1"
    su "$NON_ROOT_USER" -c "juju config -m cos prometheus maximum_retention_size=70% >/dev/null 2>&1"

    ## Offer
    su "$NON_ROOT_USER" -c "juju offer cos.grafana:grafana-dashboard global-grafana >/dev/null 2>&1"
    su "$NON_ROOT_USER" -c "juju offer cos.prometheus:receive-remote-write global-prometheus >/dev/null 2>&1"

    ## Relate (integrate)
    if ! su "$NON_ROOT_USER" -c 'juju status -m cos --relations 2>/dev/null | grep -Eq "(^|[[:space:]])(prometheus(:[^[:space:]]+)?[[:space:]]+prometheus-scrape-target-k8s(:[^[:space:]]+)?|prometheus-scrape-target-k8s(:[^[:space:]]+)?[[:space:]]+prometheus(:[^[:space:]]+)?)"'; then
      su "$NON_ROOT_USER" -c 'juju relate -m cos prometheus prometheus-scrape-target-k8s >/dev/null 2>&1'
    fi

    ## Config prometheus scrape target
    su "$NON_ROOT_USER" -c "juju config -m cos prometheus-scrape-target-k8s job_name=federate >/dev/null 2>&1"
    su "$NON_ROOT_USER" -c "juju config -m cos prometheus-scrape-target-k8s scheme=http >/dev/null 2>&1"
    su "$NON_ROOT_USER" -c "juju config -m cos prometheus-scrape-target-k8s metrics_path='/federate' >/dev/null 2>&1"
    su "$NON_ROOT_USER" -c "juju config -m cos prometheus-scrape-target-k8s params='match[]:
  - \"{__name__!=''}\"'"

    ## Cos-lite resource, default is not limit
    # Grafana
    #su "$NON_ROOT_USER" -c "juju config -m cos grafana cpu=500m >/dev/null 2>&1"
    #su "$NON_ROOT_USER" -c "juju config -m cos grafana memory=512Mi >/dev/null 2>&1"

    # Prometheus
    #su "$NON_ROOT_USER" -c "juju config -m cos prometheus cpu=2 >/dev/null 2>&1"
    #su "$NON_ROOT_USER" -c "juju config -m cos prometheus memory=4Gi >/dev/null 2>&1"

    # Loki
    #su "$NON_ROOT_USER" -c "juju config -m cos loki cpu=250m >/dev/null 2>&1"
    #su "$NON_ROOT_USER" -c "juju config -m cos loki memory=256Mi >/dev/null 2>&1"

    log "INFO" "Kubernetes cluster added to Juju successfully" "JUJU_K8S"
}

# Configure kernel modules
config_modules() {
    local module="rbd"
    local modules_file="/etc/modules"

    log "INFO" "Configuring kernel modules..." "KERNEL_CONFIG"

    if ! grep -q "^$module$" "$modules_file"; then
        echo "$module" >> "$modules_file"
        log "INFO" "Added $module to $modules_file" "KERNEL_CONFIG"
    else
        log "INFO" "Module $module already configured" "KERNEL_CONFIG"
    fi
}

# =============================================================================
# UTILITY FUNCTIONS
# =============================================================================

# Check if machine is deployed
is_machine_deployed() {
    local machine_status
    machine_status=$(maas admin machine read "$juju_machine_id" | jq -r '.status_name')
    [[ "$machine_status" == "Deployed" ]]
}

# Wait for machine deployment
wait_for_deployment() {
    local machine_id="$1"

    log "INFO" "Waiting for machine $machine_id deployment..." "DEPLOY_WAIT"

    while true; do
        local machine_status
        machine_status=$(maas admin machine read "$machine_id" | jq -r '.status_name')

        case "$machine_status" in
            "Deployed")
                log "INFO" "Machine $machine_id deployed successfully" "DEPLOY_WAIT"
                break
                ;;
            "Failed deployment")
                error_exit "Machine $machine_id failed deployment"
                ;;
            *)
                log "INFO" "Machine $machine_id deployment status: $machine_status" "DEPLOY_WAIT"
                sleep 30
                ;;
        esac
    done
}

otterscale_helm_deploy() {
    log "INFO" "Process microk8s helm3" "HELM_CHECK"
    local repository_url="https://otterscale.github.io/otterscale-charts/docs"
    local repository_name="otterscale-charts"
    local deploy_name="otterscale"
    local namespace="otterscale"

    log "INFO" "Add and update helm repository $repository_name" "HELM_REPO"
    execute_cmd "microk8s helm3 repo add $repository_name $repository_url" "helm repository add"
    execute_cmd "microk8s helm3 repo update" "helm repository update"

    if ! microk8s helm3 list -n "$namespace" | grep -qw "$deploy_name"; then
        log "INFO" "Collecting configuration data for helm deployment" "HELM_CONFIG"

        # Collect MAAS configuration
        local maas_endpoint="http://$OTTERSCALE_INTERFACE_IP:5240/MAAS"
        local kube_folder="/home/$NON_ROOT_USER/.kube"
        local maas_key
        maas_key=$(su "$NON_ROOT_USER" -c "juju show-credentials maas-cloud maas-cloud-credential --show-secrets --client | grep maas-oauth | awk '{print \$2}'")

        # Get Juju controller information
        local controller_name controller_details
        controller_name=$(su "$NON_ROOT_USER" -c "juju controllers --format json | jq -r '.\"current-controller\"'")
        controller_details=$(su "$NON_ROOT_USER" -c "juju show-controller '$controller_name' --show-password --format=json")

        # Extract controller details
        local juju_endpoints juju_username juju_password juju_cacert
        juju_endpoints=$(echo "$controller_details" | jq -r ".\"$controller_name\".details.\"api-endpoints\"[0]")
        juju_username=$(echo "$controller_details" | jq -r ".\"$controller_name\".account.user")
        juju_password=$(echo "$controller_details" | jq -r ".\"$controller_name\".account.password")
        juju_cacert=$(echo "$controller_details" | jq -r ".\"$controller_name\".details.\"ca-cert\"")

        log "INFO" "Deploy OtterScale helm chart with configuration" "HELM_DEPLOY"

        # Create temporary file for CA certificate to handle multiline content properly
        local ca_cert_file="/tmp/juju-ca-cert.pem"
        echo "$juju_cacert" > "$ca_cert_file"

        # Create temporary values file
        values_file="/tmp/otterscale-values.yaml"
        cat > "$values_file" << EOF
configContent: |
  maas:
    url: $maas_endpoint
    key: $maas_key
    version: "$OTTERSCALE_MAAS_VERSION"
  juju:
    controller: $controller_name
    controller_addresses:
    - $juju_endpoints
    username: $juju_username
    password: $juju_password
    ca_cert: |
$(echo "$juju_cacert" | sed 's/^/      /')
    cloud_name: maas-cloud
    cloud_region: default
    charmhub_api_url: $OTTERSCALE_CHARMHUB_URL
  kube:
    helm_repository_urls:
    - https://charts.bitnami.com/bitnami
  ceph:
    rados_timeout: 0s
EOF

        execute_cmd "microk8s helm3 install $deploy_name $repository_name/otterscale -n $namespace --create-namespace -f $values_file" "Deply OtterScale chart"

        # Clean up temporary file
        rm -f "$ca_cert_file"

        # Get Kubernetes configuration
        local k8s_endpoints otterscale_svc_port otterscale_web
        k8s_endpoints=$(microk8s kubectl get EndpointSlice kubernetes -o json | jq -r '.endpoints[0].addresses[0]')
        otterscale_svc_port=$(microk8s kubectl get svc -n $namespace otterscale-web -o json | jq -r '.spec.ports[] | select(.name == "http") | .nodePort')
        otterscale_web="http://$k8s_endpoints:$otterscale_svc_port"
        send_status_data "FINISHED" "OtterScale endpoint is $otterscale_web" "$otterscale_web"
        log "INFO" "OtterScale Install Finished" "FINISHED"
    else
        log "INFO" "Helm releases already has otterscale" "HELM_CHECK"
    fi
}

# =============================================================================
# This is getting quite long. Let me continue with the main() function at the end.

# =============================================================================
# MAIN INSTALLATION FUNCTION
# =============================================================================

main() {
    # Initialize logging
    init_logging

    log "INFO" "Starting OtterScale installation..." "INSTALLATION"
    log "INFO" "Target endpoint: $OTTERSCALE_ENDPOINT" "INSTALLATION"

    # System validation
    log "INFO" "Performing system validation..." "VALIDATION"
    check_root
    check_os
    check_memory
    check_disk
    disable_ipv6
    log "INFO" "All pre-checks passed" "VALIDATION"

    # Package installation
    log "INFO" "Installing packages..." "PACKAGES"
    apt_update
    apt_install "$APT_PACKAGES"
    snap_install

    # Network setup
    log "INFO" "Setting up network..." "NETWORK"
    check_bridge

    # User setup
    log "INFO" "Setting up users and SSH..." "USER_SETUP"
    find_first_non_root_user
    generate_ssh_key

    # MAAS Setup
    log "INFO" "Setting up MAAS..." "MAAS_SETUP"
    init_maas
    create_maas_admin
    login_maas
    add_key_to_maas
    update_maas_config
    download_maas_img
    enable_maas_dhcp

    # LXD Setup
    log "INFO" "Setting up LXD..." "LXD_SETUP"
    init_lxd
    create_maas_lxd_project
    create_lxd_vm
    create_vm_from_maas

    # Juju Setup
    log "INFO" "Setting up Juju..." "JUJU_SETUP"
    prepare_microk8s_config
    enable_microk8s_option
    extend_microk8s_cert
    bootstrap_juju
    juju_add_k8s

    # Final configuration
    log "INFO" "Finalizing configuration..." "FINAL_CONFIG"
    config_modules
    otterscale_helm_deploy

    log "INFO" "OtterScale installation completed successfully!" "INSTALLATION"
}

# =============================================================================
# ARGUMENT PARSING AND SCRIPT ENTRY POINT
# =============================================================================

# Parse command line arguments
parse_arguments() {
    # Default endpoint if no arguments provided
    if [[ $# -eq 0 ]]; then
        while true; do
            read -p "Enter OtterScale endpoint (default: http://127.0.0.1:8299): " user_endpoint
            OTTERSCALE_ENDPOINT="${user_endpoint:-http://127.0.0.1:8299}"

            if validate_url "$OTTERSCALE_ENDPOINT"; then
                break
            else
                echo "Invalid URL format. Please try again."
            fi
        done
        return
    fi

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --url=*|url=*)
                OTTERSCALE_ENDPOINT="${1#*=}"
                if ! validate_url "$OTTERSCALE_ENDPOINT"; then
                    error_exit "Invalid URL: $OTTERSCALE_ENDPOINT"
                fi
                ;;
            --config=*|config=*)
                local config_path="${1#*=}"
                if [[ ! -f "$config_path" ]]; then
                    error_exit "Config file not found: $config_path"
                fi
                # shellcheck source=/dev/null
                source "$config_path"
                ;;
            -h|--help|help)
                echo "Usage: sudo bash $0 [options]"
                echo ""
                echo "Options:"
                echo "  -h, --help, help     Show this help message"
                echo "  --url=URL            Specify OtterScale endpoint"
                echo "  --config=FILE        Specify configuration file"
                echo ""
                echo "Examples:"
                echo "  sudo bash $0"
                echo "  sudo bash $0 --url=http://192.168.1.100:8299"
                echo "  sudo bash $0 --config=/path/to/config"
                exit 0
                ;;
            *)
                error_exit "Invalid option: $1. Use --help for usage information."
                ;;
        esac
        shift
    done
}

# Script entry point
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    parse_arguments "$@"
    check_curl
    main "$@"
fi
