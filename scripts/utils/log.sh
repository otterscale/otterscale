#!/bin/bash

log() {
    local LOG_LEVEL=$1
    local MESSAGE=$2
    local PHASE=$3
    echo "$(date '+%Y-%m-%d %H:%M:%S') [${LOG_LEVEL}] ${MESSAGE}" | tee -a $OTTERSCALE_INSTALL_DIR/setup.log
    send_status_data "$PHASE" "$MESSAGE"
}

error_exit() {
    local MESSAGE=$1
    log "ERROR" "$MESSAGE" "ERROR"
    if [ -s "$TEMP_LOG" ]; then
        log "DEBUG" "Full error output:" "ERROR"
        cat "$TEMP_LOG" | while read line; do log "DEBUG" "$line" "ERROR"; done
    fi
    trap cleanup EXIT
    exit 1
}

# Cleanup on exit
cleanup() {
    echo "Cleaning up temporary files..."
    rm -f "$TEMP_LOG"
}
