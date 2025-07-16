#!/bin/bash

disable_ipv6() {
    log "INFO" "Disable ipv6 from sysctl, it will resume after reboot" "OS config"
    sysctl -w net.ipv6.conf.all.disable_ipv6=1 >/dev/null 2>&1
    sysctl -w net.ipv6.conf.default.disable_ipv6=1 >/dev/null 2>&1
}

check_root() {
    [ "$(id -u)" -eq 0 ] || error_exit "This script must be run as root"
}

check_os() {
    local OS_ID=$(lsb_release -si)
    local OS_VERSION=$(lsb_release -sr)
    if [ "$OS_ID" != "Ubuntu" ] || [ "$OS_VERSION" != "$OTTERSCALE_OS" ]; then
        error_exit "This script requires Ubuntu $OTTERSCALE_OS. Detected: $OS_ID $OS_VERSION"
    fi
}

check_memory() {
    local OS_MEMORY_GB=$(free -g | awk '/Mem:/ {print $2}')
    if [ "$OS_MEMORY_GB" -lt "$MIN_MEMORY_GB" ]; then
        error_exit "Insufficient memory. Minimum required: ${MIN_MEMORY_GB}GB, Available: ${OS_MEMORY_GB}GB"
    fi
}

check_disk() {
    local OS_DISK_AVAILABLE_GB=$(df -BG / | awk 'NR==2 {print $4}' | tr -d 'G')
    if [ "$OS_DISK_AVAILABLE_GB" -lt "$MIN_DISK_GB" ]; then
        error_exit "Insufficient disk space. Minimum required: ${MIN_DISK_GB}GB, Available: ${OS_DISK_AVAILABLE_GB}GB"
    fi
}

# System validation checks
validate_system() {
    log "INFO" "System validation check" "OS check"
    check_root
    check_os
    check_memory
    check_disk
    disable_ipv6
    log "INFO" "System validation passed" "OS check finished"
}

validate_url() {
    local URL=$1
    local IP=$(echo "$URL" | awk -F '[/:]' '{print $4}')
    local PORT=$(echo "$URL" | awk -F '[/:]' '{print $5}')

    if ! validate_ip $IP; then
        error_exit "Invalid IP format: $IP"
    fi

    if ! validate_port $PORT; then
        error_exit "Invalid Port format: $PORT"
    fi

    log "INFO" "Validate URL: $URL" "URL check"
}

validate_ip() {
    local IP=$1
    if [[ ! $IP =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        return 1
    else
        return 0
    fi
}

validate_port() {
    local PORT=$1
    if [[ ! $PORT =~ ^[0-9]+$ ]]; then
        return 1
    fi

    if [[ "$PORT" -lt 1 || "$PORT" -gt 65535 ]]; then
        return 1
    fi
    return 0
}

validate_cidr() {
    local CIDR=$1
    if [[ ! $CIDR =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+$ ]]; then
        return 1
    else
        return 0
    fi
}

# Function to convert an IP address to a number
ip_to_number() {
    local IP=$1
    local -a octets=(${IP//./ })
    echo $((octets[0] * 256**3 + octets[1] * 256**2 + octets[2] * 256 + octets[3]))
}

# Function to convert a network and mask to a number
network_to_number() {
    local NETWORK=$1
    local MASK=$2
    local -a octets=(${NETWORK//./ })
    local -a mask_octets=(${MASK//./ })
    local NETWORK_NUMBER=0
    for i in {0..3}; do
        NETWORK_NUMBER=$((NETWORK_NUMBER + (octets[i] & mask_octets[i]) * 256**(3-i)))
    done
    echo $NETWORK_NUMBER
}

# Function to check if an IP is in the network
is_ip_in_network() {
    local IP=$1
    local NETWORK=$2
    local MASK=$3
    local IP_NUMBER=$(ip_to_number $IP)
    local NETWORK_NUMBER=$(network_to_number $NETWORK $MASK)
    local MASK_NUMBER=$(ip_to_number $mask)

    if [ $((NETWORK_NUMBER & MASK_NUMBER)) -eq $NETWORK_NUMBER ]; then
        return 0
    else
        return 1
    fi
}

check_ip_range() {
    # Extract network and mask from subnet
    local NETWORK=$(echo $MAAS_NETWORK_SUBNET | cut -d'/' -f1)
    local MASK=$(echo $MAAS_NETWORK_SUBNET | cut -d'/' -f2)

    # Convert mask to dotted decimal format
    local MASK_DOTTED=$(printf "%d.%d.%d.%d" \
        $((0xFF << (32 - MASK) >> 24 & 0xFF)) \
        $((0xFF << (32 - MASK) >> 16 & 0xFF)) \
        $((0xFF << (32 - MASK) >> 8 & 0xFF)) \
        $((0xFF << (32 - MASK) & 0xFF)))

    # Check if DHCP_START_IP and DHCP_END_IP are in the network
    if is_ip_in_network $DHCP_START_IP $NETWORK $MASK_DOTTED; then
        if is_ip_in_network $DHCP_END_IP $NETWORK $MASK_DOTTED; then
            log "INFO" "IP range $DHCP_START_IP to $DHCP_END_IP is within the network $MAAS_NETWORK_SUBNET" "IP range check"
            return 0
        else
            log "WARN" "End IP $DHCP_END_IP is not in the network $MAAS_NETWORK_SUBNET" "IP range check"
            return 1
        fi
    else
        log "WARN" "Start IP $DHCP_START_IP is not in the network $MAAS_NETWORK_SUBNET" "IP range check"
        return 1
    fi
}
