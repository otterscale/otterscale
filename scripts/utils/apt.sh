#!/bin/bash

apt_update() {
    log "INFO" "Executing command apt update..." "APT update"
    if ! apt-get update --fix-missing >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update apt package lists. Please check your network connection"
    fi
}

apt_install() {
    local PKG_LIST=$1
    log "INFO" "Installing required apt packages: $PKG_LIST" "APT Install"
    if ! DEBIAN_FRONTEND=noninteractive apt-get install -y $PKG_LIST >>"$TEMP_LOG" 2>&1; then
        error_exit "APT package installation failed"
    fi
    log "INFO" "Apt packages installed successfully" "APT Install"
}
