#!/bin/bash

retry_snap_install() {
    local snap="$1"
    local max_retries="$2"
    local retries=0

    while [ $retries -lt $max_retries ]; do
        log "INFO" "Installing snap $snap... (Attempt $((retries+1)))"
        if snap install $snap >"$TEMP_LOG" 2>&1; then
            break
        else
            log "WARN" "Failed to install snap $snap. Retrying... (Attempt $((retries+1)))"
            retries=$((retries+1))
        fi
    done

    if [ $retries -eq $max_retries ]; then
        error_exit "Failed to install snap $snap after $max_retries attempts."
    fi
}

retry_snap_refresh() {
    local snap="$1"
    local channel="$2"
    local max_retries="$3"
    local retries=0

    while [ $retries -lt $max_retries ]; do
        log "INFO" "Refreshing snap $snap to $channel... (Attempt $((retries+1)))"
        if snap refresh $snap --channel=$channel >"$TEMP_LOG" 2>&1; then
            break
        else
            log "WARN" "Failed to refresh snap $snap to $channel. Retrying... (Attempt $((retries+1)))"
            retries=$((retries+1))
        fi
    done

    if [ $retries -eq $max_retries ]; then
        error_exit "Failed to refresh snap $snap to $channel after $max_retries attempts."
    fi
}

install_snaps() {
    for snap in $SNAP_PACKAGES; do
        if snap list | grep -q "^${snap}[[:space:]]"; then
            log "INFO" "Snap $snap is already installed. Skipping..."
            continue
        fi
        retry_snap_install "$snap" "$MAX_RETRIES"
    done

    retry_snap_refresh "lxd" "5.0/stable" "$MAX_RETRIES"

    log "INFO" "Holding all snaps..."
    snap refresh --hold >"$TEMP_LOG" 2>&1
}
