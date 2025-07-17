#!/bin/bash

apt_update() {
    log "INFO" "APT updating package lists..."
    if ! apt-get update --fix-missing >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update apt package lists. Check your network connection."
    fi
}

apt_install() {
    local PKG_LIST=$1
    log "INFO" "Installing required apt packages: $PKG_LIST"
    if ! DEBIAN_FRONTEND=noninteractive apt-get install -y $PKG_LIST >>"$TEMP_LOG" 2>&1; then
        error_exit "APT package installation failed."
    fi
    log "INFO" "Apt packages installed successfully"
}

add_ceph_repository() {
    local CEPH_CLI="/usr/local/bin/cephadm"
    local CEPH_GITHUB_LINK="https://github.com/ceph/ceph/raw/pacific/src/cephadm/cephadm"
    if [[ -f $CEPH_CLI ]]; then
        log "INFO" "Binary cephadm already exist, skipping..."
    else
        log "INFO" "Curl get cephadm from github"
        if ! curl --silent --remote-name --location $CEPH_GITHUB_LINK; then
            error_exit "Failed curl $CEPH_GITHUB_LINK."
        else
            chmod +x $INSTALLER_DIR/cephadm
            mv $INSTALLER_DIR/cephadm $CEPH_CLI
	fi
    fi

    if [[ -f /etc/apt/sources.list.d/ceph.list ]]; then
        log "INFO" "Add ceph apt repository"
        if ! $CEPH_CLI add-repo --release $CEPH_VERSION; then
            error_exit "Failed cephadm add-repo."
        fi
    fi
}
