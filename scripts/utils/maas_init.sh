#!/bin/bash

# MAAS initialization with validation
init_maas() {
    # Check if MAAS admin user already exists
    if maas apikey --username "$MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "INFO" "MAAS is already initialized (admin user exists). Skipping initialization."
        return 0
    fi

    log "INFO" "Initializing MAAS..."
    if ! maas init region+rack \
        --database-uri maas-test-db:/// \
        --maas-url "http://$BRIDGE_IP:5240/MAAS" \
        >"$TEMP_LOG" 2>&1; then
        error_exit "MAAS initialization failed."
    fi
    log "INFO" "MAAS initialized successfully"
    log "INFO" "Access MAAS at: http://$BRIDGE_IP:5240/MAAS"
    log "INFO" "MAAS Username: $MAAS_ADMIN_USER"
    log "INFO" "MAAS Password: $MAAS_ADMIN_PASS"
}

create_maas_admin() {
    log "INFO" "Creating MAAS admin user..."
    if maas apikey --username "$MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "WARN" "Admin user '$MAAS_ADMIN_USER' already exists. Using existing credentials."
    else
        if ! maas createadmin \
            --username "$MAAS_ADMIN_USER" \
            --password "$MAAS_ADMIN_PASS" \
            --email "$MAAS_ADMIN_EMAIL" \
            >"$TEMP_LOG" 2>&1; then
            error_exit "Failed to create MAAS admin user."
        fi
    fi
}

login_maas() {
    log "INFO" "Attempting to login maas..."
    local retries=0

    ## proxy will cause maas login failed
    del_proxy

    APIKEY=$(maas apikey --username "$MAAS_ADMIN_USER")
    while [ $retries -lt $MAX_RETRIES ]; do
        if maas login admin "http://localhost:5240/MAAS/" "$APIKEY" >"$TEMP_LOG" 2>&1; then
            log "INFO" "MAAS login successfully"
            add_proxy
            break
        else
            log "WARN" "Failed to login to MAAS, retry in 10 secs. (Attempt $((retries+1)))"
            retries=$((retries+1))
            sleep 10
        fi
    done

    if [ $retries -eq $MAX_RETRIES ]; then
        error_exit "Failed to get login MAAS after $MAX_RETRIES attempts"
    fi
}
