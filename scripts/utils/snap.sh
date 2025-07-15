#!/bin/bash

get_snap_channel() {
    local snap=$1
    snap list | grep "^${snap}[[:space:]]" | awk '{print $4}'
}

retry_snap_install() {
    local snap="$1"
    local max_retries="$2"
    local option="$3"
    local retries=0

    while [ $retries -lt $max_retries ]; do
        log "INFO" "Installing snap $snap... (Attempt $((retries+1)))"
        if snap install $snap $option >>"$TEMP_LOG" 2>&1; then
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
        if snap refresh $snap --channel=$channel >>"$TEMP_LOG" 2>&1; then
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

install_or_update_snap() {
    local snap=$1
    local channel=$2
    if snap list | grep -q "^${snap}[[:space:]]"; then
        if [[ $(get_snap_channel "$snap") != "$channel" ]]; then
            retry_snap_refresh "$snap" "$channel" "$OTTERSCALE_MAX_RETRIES"
        fi
    else
        if [[ $snap == "microk8s" ]]; then
            retry_snap_install "$snap" "$OTTERSCALE_MAX_RETRIES" "--classic --channel=$channel"
        else
            retry_snap_install "$snap" "$OTTERSCALE_MAX_RETRIES" "--channel=$channel"
        fi
    fi
}

snap_install() {
    declare -A SNAP_CHANNELS
    SNAP_CHANNELS[core24]=$CORE24_CHANNEL
    SNAP_CHANNELS[maas]=$MAAS_CHANNEL
    SNAP_CHANNELS[maas-test-db]=$MAAS_DB_CHANNEL
    SNAP_CHANNELS[juju]=$JUJU_CHANNEL
    SNAP_CHANNELS[lxd]=$LXD_CHANNEL
    SNAP_CHANNELS[microk8s]=$MICROK8S_CHANNEL

    for snap in $SNAP_PACKAGES; do
        CHANNEL=${SNAP_CHANNELS[$snap]}
        if [[ -z $CHANNEL ]]; then
            CHANNEL=""
        fi
        install_or_update_snap "$snap" "$CHANNEL"
    done

    log "INFO" "Holding all snaps..."
    snap refresh --hold >>"$TEMP_LOG" 2>&1
}
