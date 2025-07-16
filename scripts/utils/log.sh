#!/bin/bash

# Enhanced logging functions
log() {
    local LOG_LEVEL=$1
    local MESSAGE=$2
    echo "$(date '+%Y-%m-%d %H:%M:%S') [${LOG_LEVEL}] ${MESSAGE}" | tee -a $OTTERSCALE_INSTALL_DIR/setup.log
    send_status_data "Otterscale" "$MESSAGE"
}

error_exit() {
    log "ERROR" "$1"
    if [ -s "$TEMP_LOG" ]; then
        log "DEBUG" "Full error output:"
        cat "$TEMP_LOG" | while read line; do log "DEBUG" "$line"; done
    fi
    trap cleanup EXIT
    exit 1
}