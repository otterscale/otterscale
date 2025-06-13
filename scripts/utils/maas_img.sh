#!/bin/bash

update_boot_source(){
    # Existing source found - modify it
    boot_source_id=$(echo "$sources" | jq -r '.[0].id')
    selection_id=$(maas admin boot-source-selections read $boot_source_id | jq -r '.[0].id')
    log "INFO" "Modifying existing boot source (ID: $boot_source_id)"

    # Remove any additional sources if present
    if [ "$source_count" -gt 1 ]; then
        log "INFO" "Removing duplicate boot sources..."
        for id in $(echo "$sources" | jq -r '.[].id' | tail -n +2); do
            if ! maas admin boot-source delete "$id" >>"$TEMP_LOG" 2>&1; then
                error_exit "Failed to remove boot source $id."
            fi
        done
    fi

    # Update to use only Jammy amd64
    log "INFO" "Updating boot source to Ubuntu Jammy (22.04) amd64..."
    if ! maas admin boot-source-selection update "$boot_source_id" "$selection_id" \
        release=jammy \
        arches=amd64 \
        subarches="*" \
        labels="*" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update boot source."
    fi
}

create_boot_source() {
    # No existing source - create new one
    log "INFO" "Creating new boot source for Ubuntu Jammy (22.04) amd64..."
    if ! maas admin boot-sources create \
        url="http://images.maas.io/ephemeral-v3/stable/" \
        keyring_filename="/usr/share/keyrings/ubuntu-cloudimage-keyring.gpg" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to create boot source."
    fi

    # Get the new source ID
    boot_source_id=$(maas admin boot-sources read | jq -r '.[0].id')

    # Create selection for Jammy amd64
    if ! maas admin boot-source-selections create "$boot_source_id" \
        release=jammy \
        arches=amd64 \
        subarches="*" \
        labels="*" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to create boot source selection."
    fi
}

start_import() {
    log "INFO" "Starting image download..."
    maas admin boot-resources stop-import >>"$TEMP_LOG" 2>&1
    sleep 10

    if ! maas admin boot-resources import >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to start image download."
    fi
    sleep 10

    log "INFO" "Waiting for image download to complete..."
    while true; do
        status=$(maas admin boot-resources is-importing | jq)
        if [ "$status" != "true" ]; then
            break
        fi
        sleep 10
    done
}

download_maas_img() {
    log "INFO" "Configuring MAAS boot sources..."

    sources=$(maas admin boot-sources read)
    source_count=$(echo "$sources" | jq '. | length')
    if [ "$source_count" -gt 0 ]; then
        update_boot_source
    else
        create_boot_source
    fi

    start_import
    set_config "commissioning_distro_series" "jammy"
    set_config "default_distro_series" "jammy"
    set_config "default_osystem" "ubuntu"
    log "INFO" "MAAS images downloaded successfully"
}
