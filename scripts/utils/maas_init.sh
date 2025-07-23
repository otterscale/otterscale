#!/bin/bash

# MAAS initialization with validation
init_maas() {
    # Check if MAAS admin user already exists
    if maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "INFO" "MAAS is already initialized ($OTTERSCALE_MAAS_ADMIN_USER exist). Skipping initialization..."
        return 0
    fi

    log "INFO" "Initializing MAAS..." "MAAS init"
    execute_cmd "maas init region+rack --database-uri maas-test-db:/// --maas-url http://$OTTERSCALE_INTERFACE_IP:5240/MAAS" "maas initialization"
}

create_maas_admin() {
    log "INFO" "Creating MAAS admin user..." "MAAS init"
    if maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "WARN" "Admin user '$OTTERSCALE_MAAS_ADMIN_USER' already exists. Using existing credentials" "MAAS init"
    else
        execute_cmd "maas createadmin --username $OTTERSCALE_MAAS_ADMIN_USER --password $OTTERSCALE_MAAS_ADMIN_PASS --email $OTTERSCALE_MAAS_ADMIN_EMAIL" "create MAAS admin user"
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

    if [[ $RETRIES -eq $OTTERSCALE_MAX_RETRIES ]]; then
        error_exit "Failed to get login MAAS after $OTTERSCALE_MAX_RETRIES attempts"
    fi
}

get_maas_dns() {
    local maas_current_dns=$(maas admin maas get-config name=upstream_dns | jq -r)
    if [[ -z $maas_current_dns ]]; then
        dns_value="$OTTERSCALE_INTERFACE_DNS"
    elif [[ "$maas_current_dns" =~ "$OTTERSCALE_INTERFACE_DNS" ]]; then
        log "INFO" "Current dns already existed, skipping..."
    elif [[ $maas_current_dns != "null" ]]; then
        dns_value="$maas_current_dns $OTTERSCALE_INTERFACE_DNS"
    fi
}

set_config() {
    local NAME=$1
    local VALUE=$2
    log "INFO" "Update Config: $NAME, Value: $VALUE" "MAAS config update"
    execute_cmd "maas admin maas set-config name=$NAME value=$VALUE" "set maas $NAME config"
}

update_maas_config() {
    get_maas_dns
    set_config "upstream_dns" "$dns_value"
    set_config "boot_images_auto_import" "false"
    set_config "enable_http_proxy" "false"
    set_config "enable_analytics" "false"
    set_config "network_discovery" "disabled"
    set_config "release_notifications" "false"
}

enter_dhcp_subnet() {
    while true; do
        read -p "Enter DHCP subnet in CIDR notation (e.g. $OTTERSCALE_CONFIG_MAAS_CIDR): " MAAS_NETWORK_SUBNET
        if validate_cidr "$MAAS_NETWORK_SUBNET"; then
            break
        fi
        echo "Invalid CIDR format. Please try again."
    done
}

enter_dhcp_start_ip() {
    while true; do
        read -p "Enter DHCP start IP: " DHCP_START_IP
        if validate_ip "$DHCP_START_IP"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

enter_dhcp_end_ip() {
    while true; do
        read -p "Enter DHCP end IP: " DHCP_END_IP
        if validate_ip "$DHCP_END_IP"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

update_fabric_dns() {
    local FABRIC_DNS=$(maas admin subnet read $MAAS_NETWORK_SUBNET | jq -r '.dns_servers')
    log "INFO" "Update dns $OTTERSCALE_INTERFACE_DNS to fabric $MAAS_NETWORK_SUBNET" "MAAS config update"

    if [[ "$FABRIC_DNS" =~ "$OTTERSCALE_INTERFACE_DNS" ]]; then
        log "INFO" "Current dns already existed, skipping..." "MAAS config update"
    elif [[ ! -z $maas_current_dns ]]; then
        dns_value="$FABRIC_DNS $OTTERSCALE_INTERFACE_DNS"
    fi

    execute_cmd "maas admin subnet update $MAAS_NETWORK_SUBNET dns_servers=$dns_value" "update maas dns to fabric"
}

get_fabric() {
    log "INFO" "Getting fabric and VLAN information..." "MAAS config update"
    FABRIC_ID=$(maas admin subnet read "$MAAS_NETWORK_SUBNET" | jq -r ".vlan.fabric_id")
    VLAN_TAG=$(maas admin subnet read "$MAAS_NETWORK_SUBNET" | jq -r ".vlan.vid")
    PRIMARY_RACK=$(maas admin rack-controllers read | jq -r ".[] | .system_id")
    if [ -z "$FABRIC_ID" ] || [ -z "$VLAN_TAG" ] || [ -z "$PRIMARY_RACK" ]; then
        error_exit "Failed to get network configuration details"
    fi
}

create_dhcp_iprange() {
    log "INFO" "Creating DHCP IP range..." "MAAS config update"
    if ! maas admin ipranges create type=dynamic start_ip=$DHCP_START_IP end_ip=$DHCP_END_IP >>"$TEMP_LOG" 2>&1; then
        log "WARN" "Please confirm if address is within subnet $MAAS_NETWORK_SUBNET, or maybe it conflicts with an existing IP address or range" "MAAS config update"
        error_exit "Failed to create DHCP range"
    fi
}

update_dhcp_config() {
    log "INFO" "Enabling DHCP on VLAN..." "MAAS config update"
    if ! maas admin vlan update $FABRIC_ID $VLAN_TAG dhcp_on=True primary_rack=$PRIMARY_RACK >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to enable DHCP"
    fi
}

get_dhcp_subnet_and_ip() {
    if [ -z $OTTERSCALE_CNOFIG_MAAS_DHCP_CIDR ]; then
        enter_dhcp_subnet
    else
        MAAS_NETWORK_SUBNET=$OTTERSCALE_CNOFIG_MAAS_DHCP_CIDR
    fi

    if [ -z $OTTERSCALE_CONFIG_MAAS_DHCP_START_IP ]; then
        enter_dhcp_start_ip
    else
        DHCP_START_IP=$OTTERSCALE_CONFIG_MAAS_DHCP_START_IP
    fi

    if [ -z $OTTERSCALE_CONFIG_MAAS_DHCP_END_IP ]; then
        enter_dhcp_end_ip
    else
        DHCP_END_IP=$OTTERSCALE_CONFIG_MAAS_DHCP_END_IP
    fi
}

enable_maas_dhcp() {
    if [ $(maas admin ipranges read | jq '. | length') -ne 0 ]; then
        log "INFO" "MAAS already has dynamic IP ranges, skipping..." "MAAS config update"
        return 0
    fi

    log "INFO" "Configuring MAAS DHCP..." "MAAS config update"
    while true; do
        get_dhcp_subnet_and_ip
        if check_ip_range ; then
            update_fabric_dns
            get_fabric
            create_dhcp_iprange
            update_dhcp_config
	    break
        else
            if ! -z $OTTERSCALE_CNOFIG_MAAS_DHCP_CIDR && ! -z $OTTERSCALE_CONFIG_MAAS_DHCP_START_IP && ! -z $OTTERSCALE_CONFIG_MAAS_DHCP_END_IP ]]; then
                break
            fi
        fi
    done
    log "INFO" "DHCP configuration completed" "MAAS config update"
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
