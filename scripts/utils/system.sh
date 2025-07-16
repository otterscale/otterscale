#!/bin/bash

start_service() {
    local SERVICE_NAME=$1
    systemctl daemon-reload >/dev/null 2>&1
    if systemctl start $SERVICE_NAME >/dev/null 2>&1; then
        log "INFO" "Start service $SERVICE_NAME"
    else
        error_exit "Failed start service $SERVICE_NAME"
    fi
}

stop_service() {
    local SERVICE_NAME=$1
    systemctl daemon-reload >/dev/null 2>&1
    if systemctl stop $SERVICE_NAME >/dev/null 2>&1; then
        log "INFO" "Stop service $SERVICE_NAME"
    else
        error_exit "Failed stop service $SERVICE_NAME"
    fi
}

enable_service() {
    local SERVICE_NAME=$1
    systemctl daemon-reload >/dev/null 2>&1
    if ! systemctl enable $SERVICE_NAME >/dev/null 2>&1; then
        error_exit "Failed enable service $SERVICE_NAME"
    fi
}

disable_service() {
    local SERVICE_NAME=$1
    systemctl daemon-reload >/dev/null 2>&1
    if ! systemctl disable $SERVICE_NAME >/dev/null 2>&1; then
        error_exit "Failed disable service $SERVICE_NAME"
    fi
}

# Cleanup on exit
cleanup() {
    echo "Cleaning up temporary files..."
    rm -f "$TEMP_LOG"
}

execute_non_user_cmd() {
    local USERNAME="$1"
    local COMMAND="$2"
    local DESCRIPTION="$3"
    if ! su "$USERNAME" -c "${COMMAND} >>$LOG 2>&1"; then
        log "WARN" "Failed to $DESCRIPTION, check $LOG for details."
        return 1
    fi
    return 0
}
