#!/bin/bash

apt_update() {
    log "INFO" "APT updating package lists..."
    if ! apt-get update --fix-missing >"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update apt package lists. Check your network connection."
    fi
}

apt_install() {
    log "INFO" "Installing required apt packages: $APT_PACKAGES"
    if ! DEBIAN_FRONTEND=noninteractive apt-get install -y $APT_PACKAGES >"$TEMP_LOG" 2>&1; then
        error_exit "APT package installation failed."
    fi
    log "INFO" "Apt packages installed successfully"
}

install_packages() {
    apt_update
    apt_install
}
