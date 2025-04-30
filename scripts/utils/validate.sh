#!/bin/bash

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
    check_root
    check_os
    check_memory
    check_disk
    log "INFO" "System validation passed"
}

validate_ip() {
    local ip=$1
    if [[ ! $ip =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
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
