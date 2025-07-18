#!/bin/bash

# MAAS initialization with validation
init_maas() {
    # Check if MAAS admin user already exists
    if maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "INFO" "MAAS is already initialized (admin user exists). Skipping initialization." "MAAS init"
        return 0
    fi

    log "INFO" "Initializing MAAS..." "MAAS init"
    if ! maas init region+rack \
        --database-uri maas-test-db:/// \
        --maas-url "http://$OTTERSCALE_INTERFACE_IP:5240/MAAS" \
        >>"$TEMP_LOG" 2>&1; then
        error_exit "MAAS initialization failed"
    fi
    log "INFO" "MAAS initialized successfully" "MAAS init"
}

create_maas_admin() {
    log "INFO" "Creating MAAS admin user..." "MAAS init"
    if maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "WARN" "Admin user '$OTTERSCALE_MAAS_ADMIN_USER' already exists. Using existing credentials" "MAAS init"
    else
        if ! maas createadmin \
            --username "$OTTERSCALE_MAAS_ADMIN_USER" \
            --password "$OTTERSCALE_MAAS_ADMIN_PASS" \
            --email "$OTTERSCALE_MAAS_ADMIN_EMAIL" \
            >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed to create MAAS admin user"
        fi
    fi
    log "INFO" "MAAS web url: http://$OTTERSCALE_INTERFACE_IP:5240/MAAS" "MAAS init"
    log "INFO" "MAAS Username: $OTTERSCALE_MAAS_ADMIN_USER" "MAAS init"
    log "INFO" "MAAS Password: $OTTERSCALE_MAAS_ADMIN_PASS" "MAAS init"
}

login_maas() {
    log "INFO" "Attempting to login maas..." "MAAS init"
    local RETRIES=0

    APIKEY=$(maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER")
    while [ $RETRIES -lt $OTTERSCALE_MAX_RETRIES ]; do
        if maas login admin "http://localhost:5240/MAAS/" "$APIKEY" >>"$TEMP_LOG" 2>&1; then
            log "INFO" "MAAS login successfully" "MAAS init"
            break
        else
            log "WARN" "Failed to login to MAAS, retry in 10 secs. (Attempt $RETRIES)" "MAAS init"
            RETRIES=$((RETRIES+1))
            sleep 10
        fi
    done

    if [ $RETRIES -eq $OTTERSCALE_MAX_RETRIES ]; then
        error_exit "Failed to get login MAAS after $OTTERSCALE_MAX_RETRIES attempts"
    fi
}

check_maas() {
    OTTERSCALE_MAAS_ADMIN_USER=${OTTERSCALE_CONFIG_MAAS_ADMIN_USER:-admin}
    OTTERSCALE_MAAS_ADMIN_PASS=${OTTERSCALE_CONFIG_MAAS_ADMIN_PASS:-admin}
    OTTERSCALE_MAAS_ADMIN_EMAIL=${OTTERSCALE_CONFIG_MAAS_ADMIN_EMAIL:-admin@example.com}

    ## Init, create, and login
    init_maas
    create_maas_admin
    login_maas

    ## Generate ssh
    set_sshkey

    ## Configure
    update_maas_dns
    update_maas_config
    download_maas_img
    enable_maas_dhcp

    ## Lxd
    init_lxd
    create_maas_lxd_project
    create_lxd_vm
    create_vm_from_maas
    set_vm_static_ip
}
