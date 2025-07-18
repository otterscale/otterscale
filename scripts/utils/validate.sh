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
    fi
    return 0
}
