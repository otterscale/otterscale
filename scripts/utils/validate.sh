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
    local os_id=$(lsb_release -si)
    local os_version=$(lsb_release -sr)
    if [ "$os_id" != "Ubuntu" ] || [ "$os_version" != "$REQUIRED_UBUNTU_VERSION" ]; then
        error_exit "This script requires Ubuntu $REQUIRED_UBUNTU_VERSION. Detected: $os_id $os_version"
    fi
}

check_memory() {
    local total_mem=$(free -g | awk '/Mem:/ {print $2}')
    if [ "$total_mem" -lt "$MIN_MEMORY_GB" ]; then
        error_exit "Insufficient memory. Minimum required: ${MIN_MEMORY_GB}GB, Available: ${total_mem}GB"
    fi
}

check_disk() {
    local disk_space=$(df -BG / | awk 'NR==2 {print $4}' | tr -d 'G')
    if [ "$disk_space" -lt "$MIN_DISK_GB" ]; then
        error_exit "Insufficient disk space. Minimum required: ${MIN_DISK_GB}GB, Available: ${disk_space}GB"
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
    local url=$1
    local ip=$(echo "$url" | awk -F '[/:]' '{print $4}')
    local port=$(echo "$url" | awk -F '[/:]' '{print $5}')

    if ! validate_ip $ip; then
        error_exit "Invalid IP format: $ip"
    fi

    if ! validate_port $port; then
        error_exit "Invalid Port format: $port"
    fi

    log "INFO" "Validate URL: $url"
}

validate_ip() {
    local ip=$1
    if [[ ! $ip =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        return 1
    fi
    return 0
}

validate_port() {
    local port=$1
    if [[ ! $port =~ ^[0-9]+$ ]]; then
        return 1
    fi

    if [[ "$port" -lt 1 || "$port" -gt 65535 ]]; then
        return 1
    fi
    return 0
}

validate_cidr() {
    local cidr=$1
    if [[ ! $cidr =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+$ ]]; then
        return 1
    fi
    return 0
}

validate_selected_bridge() {
    BRIDGE_IP=$(ip -o -4 addr show dev "$bridge" | awk '{print $4}' | cut -d'/' -f1)
    if [ -z "$BRIDGE_IP" ]; then
        error_exit "Selected bridge $bridge has no IP address assigned"
    fi

    log "INFO" "Using bridge $bridge with IP $BRIDGE_IP"
}

# Function to convert an IP address to a number
ip_to_number() {
    local ip=$1
    local -a octets=(${ip//./ })
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
