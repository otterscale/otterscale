#!/bin/bash

update_boot_source(){
    # Existing source found - modify it
    MAAS_BOOT_SOURCE_ID=$(echo "$MAAS_BOOT_SOURCES" | jq -r '.[0].id')
    MAAS_BOOT_SELECTION_ID=$(maas admin boot-source-selections read "$MAAS_BOOT_SOURCE_ID" | jq -r '.[0].id')
    log "INFO" "Modifying existing boot source (ID: $MAAS_BOOT_SOURCE_ID)" "MAAS boot image"

    # Remove any additional sources if present
    if [ "$MAAS_BOOT_SOURCE_COUNT" -gt 1 ]; then
        log "INFO" "Removing duplicate boot sources..." "MAAS boot image"
        for id in $(echo "$sources" | jq -r '.[].id' | tail -n +2); do
            execute_cmd "maas admin boot-source delete $id" "remove maas boot source $id"
        done
    fi

    # Update to use only Jammy amd64
    log "INFO" "Updating boot source to Ubuntu Jammy (22.04) amd64..." "MAAS boot image"
    if ! maas admin boot-source-selection update "$MAAS_BOOT_SOURCE_ID" "$MAAS_BOOT_SELECTION_ID" \
        release=jammy \
        arches=amd64 \
        subarches="*" \
        labels="*" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update maas boot source"
    fi
}

create_boot_source() {
    # No existing source - create new one
    log "INFO" "Creating new boot source for Ubuntu Jammy (22.04) amd64..." "MAAS boot image"
    if ! maas admin boot-sources create \
        url="http://images.maas.io/ephemeral-v3/stable/" \
        keyring_filename="/usr/share/keyrings/ubuntu-cloudimage-keyring.gpg" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to create maas boot source"
    fi

    # Get the new source ID
    MAAS_BOOT_SOURCE_ID=$(maas admin boot-sources read | jq -r '.[0].id')

    # Create selection for Jammy amd64
    if ! maas admin boot-source-selections create "$MAAS_BOOT_SOURCE_ID" \
        release=jammy \
        arches=amd64 \
        subarches="*" \
        labels="*" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed create maas boot source selection"
    fi
}

start_import() {
    log "INFO" "Starting download MAAS boot image..." "MAAS boot image"
    maas admin boot-resources stop-import >>"$TEMP_LOG" 2>&1
    sleep 10

    execute_cmd "maas admin boot-resources import" "start maas image download"
    sleep 10

    log "INFO" "Waiting for image download to complete..." "MAAS boot image"
    while true; do
	if [ $(maas admin boot-resources is-importing | jq -r) != "true" ]; then
            break
        fi
        sleep 10
    done
}

download_maas_img() {
    log "INFO" "Configuring MAAS boot sources..." "MAAS boot image"

    MAAS_BOOT_SOURCES=$(maas admin boot-sources read)
    MAAS_BOOT_SOURCE_COUNT=$(echo "$MAAS_BOOT_SOURCES" | jq '. | length')
    if [ "$MAAS_BOOT_SOURCE_COUNT" -gt 0 ]; then
        update_boot_source
    else
        create_boot_source
    fi

    start_import
    set_config "commissioning_distro_series" "jammy"
    set_config "default_distro_series" "jammy"
    set_config "default_osystem" "ubuntu"
    log "INFO" "MAAS images downloaded successfully" "MAAS boot image"
}
