#!/bin/bash

get_snap_channel() {
    local SNAP_NAME=$1
    snap list | grep "^${SNAP_NAME}[[:space:]]" | awk '{print $4}'
}

retry_snap_install() {
    local SNAP_NAME="$1"
    local MAX_RETRIES="$2"
    local SNAP_OPTION="$3"
    local RETRIES=0

    while [ $RETRIES -lt $MAX_RETRIES ]; do
        log "INFO" "Installing snap $SNAP_NAME... (Attempt $((RETRIES)))"
        if snap install $SNAP_NAME $SNAP_OPTION >>"$TEMP_LOG" 2>&1; then
            break
        else
            log "WARN" "Failed to install snap $SNAP_NAME. Retrying... (Attempt $((RETRIES)))"
            RETRIES=$((RETRIES+1))
        fi
    done

    if [ $RETRIES -eq $MAX_RETRIES ]; then
        error_exit "Failed to install snap $SNAP_NAME after $MAX_RETRIES attempts."
    fi
}

retry_snap_refresh() {
    local SNAP_NAME="$1"
    local SNAP_CHANNEL="$2"
    local MAX_RETRIES="$3"
    local RETRIES=0

    while [ $RETRIES -lt $MAX_RETRIES ]; do
        log "INFO" "Refreshing snap $SNAP_NAME to $MAX_RETRIES... (Attempt $((RETRIES)))"
        if snap refresh $SNAP_NAME --channel=$MAX_RETRIES >>"$TEMP_LOG" 2>&1; then
            break
        else
            log "WARN" "Failed to refresh snap $snSNAP_NAMEap to $MAX_RETRIES. Retrying... (Attempt $((RETRIES)))"
            RETRIES=$((RETRIES+1))
        fi
    done

    if [ $RETRIES -eq $MAX_RETRIES ]; then
        error_exit "Failed to refresh snap $SNAP_NAME to $MAX_RETRIES after $MAX_RETRIES attempts."
    fi
}

install_or_update_snap() {
    local SNAP_NAME=$1
    local SNAP_CHANNEL=$2
    if snap list | grep -q "^${SNAP_NAME}[[:space:]]"; then
        if [[ $(get_snap_channel "$SNAP_NAME") != "$SNAP_CHANNEL" ]]; then
            retry_snap_refresh "$SNAP_NAME" "$SNAP_CHANNEL" "$OTTERSCALE_MAX_RETRIES"
        fi
    else
        if [[ $SNAP_NAME == "microk8s" ]]; then
            retry_snap_install "$SNAP_NAME" "$OTTERSCALE_MAX_RETRIES" "--classic --channel=$SNAP_CHANNEL"
        else
            retry_snap_install "$SNAP_NAME" "$OTTERSCALE_MAX_RETRIES" "--channel=$SNAP_CHANNEL"
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

    for SNAP_NAME in $SNAP_PACKAGES; do
        CHANNEL=${SNAP_CHANNELS[$SNAP_NAME]}
        if [[ -z $CHANNEL ]]; then
            CHANNEL=""
        fi
        install_or_update_snap "$SNAP_NAME" "$CHANNEL"
    done

    log "INFO" "Holding all snaps..."
    snap refresh --hold >>"$TEMP_LOG" 2>&1
}
