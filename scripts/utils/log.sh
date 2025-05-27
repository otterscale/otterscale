#!/bin/bash

# Enhanced logging functions
log() {
    local level=$1
    local message=$2
    echo "$(date '+%Y-%m-%d %H:%M:%S') [${level}] ${message}" | tee -a $INSTALLER_DIR/setup.log
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
