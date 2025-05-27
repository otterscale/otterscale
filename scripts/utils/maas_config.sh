#!/bin/bash

update_maas_dns() {
    get_current_dns $bridge
    local dns_value="$current_dns"
    local maas_current_dns=$(maas admin maas get-config name=upstream_dns | jq -r)

    log "INFO" "Update $dns_value to maas dns."
    if [[ "$maas_current_dns" =~ "$dns_value" ]]; then
        log "INFO" "Current dns already existed, skipping..."
    elif [[ $maas_current_dns != "null" ]]; then
        dns_value="$maas_current_dns $dns_value"
    fi

    set_config "upstream_dns" "$dns_value"
}

set_config() {
    local name=$1
    local value=$2
    if ! maas admin maas set-config name=$name value=$value >"$TEMP_LOG" 2>&1; then
        error_exit "Failed to set config $name to $value."
    fi
}

update_maas_config() {
    set_config "boot_images_auto_import" "false"
    set_config "enable_http_proxy" "false"
    set_config "enable_analytics" "false"
    set_config "network_discovery" "disabled"
    set_config "release_notifications" "false"
}

enter_dhcp_subnet() {
    while true; do
        read -p "Enter DHCP subnet in CIDR notation (e.g., $current_ip): " subnet
        if validate_cidr "$subnet"; then
            break
        fi
        echo "Invalid CIDR format. Please try again."
    done
}

enter_dhcp_start_ip() {
    while true; do
        read -p "Enter DHCP start IP: " start_ip
        if validate_ip "$start_ip"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

enter_dhcp_end_ip() {
    while true; do
        read -p "Enter DHCP end IP: " end_ip
        if validate_ip "$end_ip"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

update_fabric_dns() {
    local dns_value=$current_dns
    local fabric_dns=$(maas admin subnet read $subnet | jq -r '.dns_servers')
    log "INFO" "Update dns $dns_value to fabric $subnet."

    if [[ "$fabric_dns" =~ "$dns_value" ]]; then
        log "INFO" "Current dns already existed, skipping..."
    elif [[ ! -z $maas_current_dns ]]; then
        dns_value="$fabric_dns $dns_value"
    fi

    if ! maas admin subnet update "$subnet" dns_servers="$dns_value" >"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update dns to fabric."
    fi
}

get_fabric() {
    log "INFO" "Getting fabric and VLAN information..."
    FABRIC_ID=$(maas admin subnet read "$subnet" | jq -r ".vlan.fabric_id")
    VLAN_TAG=$(maas admin subnet read "$subnet" | jq -r ".vlan.vid")
    PRIMARY_RACK=$(maas admin rack-controllers read | jq -r ".[] | .system_id")
    if [ -z "$FABRIC_ID" ] || [ -z "$VLAN_TAG" ] || [ -z "$PRIMARY_RACK" ]; then
        error_exit "Failed to get network configuration details"
    fi
}

create_dhcp_iprange() {
    log "INFO" "Creating DHCP IP range..."
    if ! maas admin ipranges create type=dynamic start_ip=$start_ip end_ip=$end_ip >"$TEMP_LOG" 2>&1; then
        log "WARN" "Please confirm if address is within subnet $subnet, or maybe it already exist."
        error_exit "Failed to create DHCP range."
    fi
}

update_dhcp_config() {
    log "INFO" "Enabling DHCP on VLAN..."
    if ! maas admin vlan update $FABRIC_ID $VLAN_TAG dhcp_on=True primary_rack=$PRIMARY_RACK >"$TEMP_LOG" 2>&1; then
        error_exit "Failed to enable DHCP."
    fi
}

enable_maas_dhcp() {
    dynamic_ipranges_count=$(maas admin ipranges read | jq '. | length')
    if [ "$dynamic_ipranges_count" -ne 0 ]; then
        log "INFO" "MAAS already has dynamic IP ranges, skipping..."
        return 0
    fi

    log "INFO" "Configuring MAAS DHCP..."
    get_current_ip $bridge
    while true; do
        enter_dhcp_subnet
        enter_dhcp_start_ip
        enter_dhcp_end_ip
        if check_ip_range ; then
            update_fabric_dns
            get_fabric
            create_dhcp_iprange
            update_dhcp_config
	    break
        fi
    done
    log "INFO" "DHCP configuration completed"
}
