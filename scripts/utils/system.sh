#!/bin/bash

disable_ipv6() {
    log "INFO" "Disable ipv6 from sysctl, it will resume after reboot" "OS config"
    sysctl -w net.ipv6.conf.all.disable_ipv6=1 >/dev/null 2>&1
    sysctl -w net.ipv6.conf.default.disable_ipv6=1 >/dev/null 2>&1
}

find_first_non_user() {
    local USER_HOME=""
    for USER in $(ls /home); do
        if [ -d "/home/$USER" ]; then
            USER_HOME="/home/$USER"
            break
        fi
    done

    if [ -z "$USER_HOME" ]; then
        error_exit "No non-root user found for SSH key setup"
    fi

    NON_ROOT_USER=$(basename "$USER_HOME")
    log "INFO" "Non-root user is $NON_ROOT_USER"
}

generate_ssh_key() {
    if [[ ! -f "/home/$NON_ROOT_USER/.ssh/id_rsa" ]]; then
        if ! su "$NON_ROOT_USER" -c 'mkdir -p $HOME/.ssh; ssh-keygen -q -t rsa -N "" -f "$HOME/.ssh/id_rsa" >>"$LOG" 2>&1'; then
            error_exit "SSH key generation failed"
        fi
    fi

    chown -R "$NON_ROOT_USER:$NON_ROOT_USER" "/home/$NON_ROOT_USER/.ssh"
    chmod 600 "/home/$NON_ROOT_USER/.ssh/id_rsa"
    chmod 644 "/home/$NON_ROOT_USER/.ssh/id_rsa.pub"
}

add_key_to_maas() {
    if [[ $(maas admin sshkeys read | jq -r 'length') -eq 0 ]]; then
        if ! maas admin sshkeys create key="$(cat "/home/$NON_ROOT_USER/.ssh/id_rsa.pub")" >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed to add SSH key to MAAS"
        fi
    fi
}

set_sshkey() {
    find_first_non_user
    generate_ssh_key
    add_key_to_maas
}

execute_non_user_cmd() {
    local USERNAME="$1"
    local COMMAND="$2"
    local DESCRIPTION="$3"
    if ! su "$USERNAME" -c "${COMMAND} >>$LOG 2>&1"; then
        log "WARN" "Failed to $DESCRIPTION, check $LOG for details" "Non-root cmd"
        return 1
    fi
    return 0
}

execute_cmd() {
    local CMD=$1
    local MSG=$2
    log "INFO" "Execute command: $CMD" "$MSG"
    if ! $CMD >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed $MSG"
    fi
    return 0
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

config_modules() {
    if ! echo 'rbd' | tee /etc/modules; then
        error_exit "Failed tee rbd into /etc/modules"
    fi
}
