#!/bin/bash

disable_ipv6() {
    log "INFO" "Disable ipv6 from sysctl, it will resume after reboot."
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
    log "INFO" "System validation check"
    check_root
    check_os
    check_memory
    check_disk
    disable_ipv6
    log "INFO" "System validation passed"
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

    log "INFO" "Validate URL: $URL"
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
    local network=$1
    local mask=$2
    local -a octets=(${network//./ })
    local -a mask_octets=(${mask//./ })
    local network_number=0
    for i in {0..3}; do
        network_number=$((network_number + (octets[i] & mask_octets[i]) * 256**(3-i)))
    done
    echo $network_number
}

# Function to check if an IP is in the network
is_ip_in_network() {
    local ip=$1
    local network=$2
    local mask=$3
    local ip_number=$(ip_to_number $ip)
    local network_number=$(network_to_number $network $mask)
    local mask_number=$(ip_to_number $mask)

    if [ $((ip_number & mask_number)) -eq $network_number ]; then
        return 0
    else
        return 1
    fi
}

check_ip_range() {
    # Extract network and mask from subnet
    local network=$(echo $subnet | cut -d'/' -f1)
    local mask=$(echo $subnet | cut -d'/' -f2)

    # Convert mask to dotted decimal format
    local mask_dotted=$(printf "%d.%d.%d.%d" \
        $((0xFF << (32 - mask) >> 24 & 0xFF)) \
        $((0xFF << (32 - mask) >> 16 & 0xFF)) \
        $((0xFF << (32 - mask) >> 8 & 0xFF)) \
        $((0xFF << (32 - mask) & 0xFF)))

    # Check if start_ip and end_ip are in the network
    if is_ip_in_network $start_ip $network $mask_dotted; then
        if is_ip_in_network $end_ip $network $mask_dotted; then
            log "INFO" "IP range $start_ip to $end_ip is within the network $subnet"
            return 0
        else
            log "WARN" "End IP $end_ip is not in the network $subnet"
            return 1
        fi
    else
        log "WARN" "Start IP $start_ip is not in the network $subnet"
        return 1
    fi
}
