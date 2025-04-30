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

install_openhdc() {
    local deb_file=$(ls $INSTALLER_DIR/packages/ | grep deb | head -n 1)
    if ! apt-get install -y $INSTALLER_DIR/packages/$deb_file >"$TEMP_LOG" 2>&1; then
        error_exit "Failed apt installed openhdc."
    else
        log "INFO" "OpenHDC installed successfully, check systemctl status openhdc"
        log "INFO" "OpenHDC used 5059 tcp port, you can visit it from browser."
    fi
}
