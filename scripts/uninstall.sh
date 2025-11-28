#!/bin/bash

APT_PACKAGES=""
SNAP_PACKAGES="core24 maas maas-test-db juju lxd microk8s"

log() {
    local LOG_LEVEL=$1
    local MESSAGE=$2
    echo "$(date '+%Y-%m-%d %H:%M:%S') [${LOG_LEVEL}] ${MESSAGE}"
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
        log "Info" "No non-root user found"
        exit 1
    fi

    NON_ROOT_USER=$(basename "$USER_HOME")
}

get_units_from_models() {
    JUJU_UNITS=$(su "$NON_ROOT_USER" -c \
        "juju status -m $JUJU_MODEL --format json 2>/dev/null | \
        jq -r 'select(.\"applications\" != {}) | \
        .applications[] | \
        select(.charm == \"ceph-osd\" and .units != null) | \
        .units | \
        keys[]'")
}

dd_ceph_osd_device() {
    local JUJU_APPLICATIONS=$(echo $JUJU_UNIT | cut -d'/' -f1)
    OSD_DEVICES=$(su "$NON_ROOT_USER" -c "juju config -m $JUJU_MODEL $JUJU_APPLICATIONS osd-devices")
    log "INFO" "OSD devices: $OSD_DEVICES"
    for osd_device in "${OSD_DEVICES[@]}"; do
        log "INFO" "dd $JUJU_UNIT disk $osd_device..."
        local GET_SZ_COMMAND="sudo blockdev --getsz $osd_device"
        local SECTOR=$(su "$NON_ROOT_USER" -c "juju ssh -m $JUJU_MODEL $JUJU_UNIT $GET_SZ_COMMAND 2>/dev/null" | tr -dc '0-9')
        local TARGET_SECTOR=$((SECTOR - 20480))
        local DD_COMMAND_FIRST="sudo dd if=/dev/zero of=$osd_device bs=1M count=10 conv=fsync status=progress 2>/dev/null"
        local DD_COMMAND_END="sudo dd if=/dev/zero of=$osd_device bs=512 count=20480 seek=$TARGET_SECTOR conv=fsync status=progress 2>/dev/null"
        su "$NON_ROOT_USER" -c "juju ssh -m $JUJU_MODEL $JUJU_UNIT $DD_COMMAND_FIRST"
        su "$NON_ROOT_USER" -c "juju ssh -m $JUJU_MODEL $JUJU_UNIT $DD_COMMAND_END"
    done
}

remove_juju_model() {
    JUJU_MODELS=$(su "$NON_ROOT_USER" -c "juju models --format json | jq -r '.models[] | select(.\"is-controller\" == false and .\"cloud\" != \"cos-k8s\") | .name' 2>/dev/null")
    for juju_model in $JUJU_MODELS; do
        export JUJU_MODEL=$juju_model
        log "INFO" "Target juju model: $JUJU_MODEL"
        get_units_from_models
        for juju_unit in "$JUJU_UNITS"; do
            if [ -z $juju_unit ]; then
                continue
            fi
            export JUJU_UNIT=$juju_unit
            log "INFO" "Target juju unit: $JUJU_UNIT"
            dd_ceph_osd_device
        done

        log "INFO" "Removing juju model $JUJU_MODEL..."
        su "$NON_ROOT_USER" -c "juju destroy-model --no-prompt $JUJU_MODEL --force --no-wait 2>/dev/null"
    done

    unset JUJU_MODEL
    unset JUJU_UNIT
}

remove_pkg() {
    for SNAP_NAME in $SNAP_PACKAGES; do
        log "INFO" "Removing snap $SNAP_NAME..."
        snap remove --purge "$SNAP_NAME" >/dev/null 2>&1
    done

    log "INFO" "Removing apt $APT_PACKAGES..."
    DEBIAN_FRONTEND=noninteractive apt-get purge -y $APT_PACKAGES >/dev/null 2>&1
    DEBIAN_FRONTEND=noninteractive apt-get autoremove -y >/dev/null 2>&1
}

remove_juju_file() {
    if [ -d /home/$NON_ROOT_USER/.local/share/juju ]; then
        log "INFO" "Removing folder /home/$NON_ROOT_USER/.local/share/juju"
        rm -rf /home/$NON_ROOT_USER/.local/share/juju
    fi
}

remove_iptables() {
    if iptables -C FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT 2>/dev/null; then
        iptables -D FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT
    fi

    if iptables -C FORWARD -i br-otters -j ACCEPT 2>/dev/null; then
        iptables -D FORWARD -i br-otters -j ACCEPT
    fi
}

main() {
    find_first_non_user
    if command -v juju >/dev/null 2>&1; then
        remove_juju_model
    fi
    remove_pkg
    remove_juju_file
    remove_iptables
}

main