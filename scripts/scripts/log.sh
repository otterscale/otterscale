#!/bin/bash

# Cleanup on exit
cleanup() {
    log "INFO" "Remove unused yaml file."
    rm -f $INSTALLER_DIR/lxd-config.yaml
    rm -f $INSTALLER_DIR/cloud.yaml
    rm -f $INSTALLER_DIR/credential.yaml

    rm -f "$TEMP_LOG"
    echo "Cleaning up temporary files..."
}

# Enhanced logging functions
log() {
    local level=$1
    local message=$2
    echo "$(date '+%Y-%m-%d %H:%M:%S') [${level}] ${message}" | tee -a $INSTALLER_DIR/setup.log
}

error_exit() {
    log "ERROR" "$1"
    if [ -s "$TEMP_LOG" ]; then
#        log "DEBUG" "Last 50 lines of error output:"
#        tail -n 50 "$TEMP_LOG" | while read line; do log "DEBUG" "$line"; done
        log "DEBUG" "Full error output:"
        cat "$TEMP_LOG" | while read line; do log "DEBUG" "$line"; done
    fi
    trap cleanup EXIT
    exit 1
}
