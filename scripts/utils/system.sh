#!/bin/bash

start_service() {
    local serviceName=$1
    systemctl daemon-reload >/dev/null 2>&1
    if systemctl start $serviceName >/dev/null 2>&1; then
        log "INFO" "Start service $serviceName"
    else
        error_exit "Failed start service $serviceName"
    fi
}

stop_service() {
    local serviceName=$1
    systemctl daemon-reload >/dev/null 2>&1
    if systemctl stop $serviceName >/dev/null 2>&1; then
        log "INFO" "Stop service $serviceName"
    else
        error_exit "Failed stop service $serviceName"
    fi
}

enable_service() {
    local serviceName=$1
    systemctl daemon-reload >/dev/null 2>&1
    if ! systemctl enable $serviceName >/dev/null 2>&1; then
        error_exit "Failed enable service $serviceName"
    fi
}

disable_service() {
    local serviceName=$1
    systemctl daemon-reload >/dev/null 2>&1
    if ! systemctl disable $serviceName >/dev/null 2>&1; then
        error_exit "Failed disable service $serviceName"
    fi
}

# Cleanup on exit
cleanup() {
    echo "Cleaning up temporary files..."
    rm -f "$TEMP_LOG"
}

execute_non_user_cmd() {
    local username="$1"
    local command="$2"
    local description="$3"
    if ! su "$username" -c "${command} >>$LOG 2>&1"; then
        log "WARN" "Failed to $description, check $LOG for details."
        return 1
    fi
    return 0
}
