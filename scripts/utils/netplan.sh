#!/bin/bash

check_bridge() {
    if [ -z $OTTERSCALE_CONFIG_BRIDGE_NAME ]; then
        select_bridge
    else
        OTTERSCALE_BRIDGE_NAME=$OTTERSCALE_CONFIG_BRIDGE_NAME
        if ! check_interface_exist $OTTERSCALE_BRIDGE_NAME; then
            error_exit "Network interface $OTTERSCALE_BRIDGE_NAME not exist."
        fi

        if ! check_ip_in_interface $OTTERSCALE_BRIDGE_NAME $OTTERSCALE_CONFIG_BRIDGE_CIDR; then
            error_exit "CIDR $OTTERSCALE_CONFIG_BRIDGE_CIDR not found in network interface $OTTERSCALE_BRIDGE_NAME"
        fi
        OTTERSCALE_INTERFACE_IP=$(echo $OTTERSCALE_CONFIG_BRIDGE_CIDR | cut -d'/' -f1)
        OTTERSCALE_INTERFACE_IP_MASK=$(echo $OTTERSCALE_CONFIG_BRIDGE_CIDR | cut -d'/' -f2)
        get_current_dns $OTTERSCALE_BRIDGE_NAME
    fi
}

select_bridge() {
    while true; do
        log "INFO" "Detecting available bridges..." "OS network"
        local AVAILABLE_BRIDGES=($(brctl show 2>/dev/null | awk 'NR>1 {print $1}' | grep -v '^$'))

        echo "Available network bridges:"
        echo "0) Create new bridge"
        for i in "${!AVAILABLE_BRIDGES[@]}"; do
            echo "$((i+1))) ${AVAILABLE_BRIDGES[$i]}"
        done

        read -p "Select bridge (0-${#AVAILABLE_BRIDGES[@]}): " choice
        case $choice in
            0)
                create_new_bridge
                return
                ;;
            [1-9]*)
                if [ $choice -le ${#AVAILABLE_BRIDGES[@]} ]; then
                    OTTERSCALE_BRIDGE_NAME=${AVAILABLE_BRIDGES[$((choice-1))]}
                    get_current_ip $OTTERSCALE_BRIDGE_NAME
                    get_current_dns $OTTERSCALE_BRIDGE_NAME
                    return
                fi
                ;;
        esac
        log "WARN" "Invalid selection. Please try again." "OS network"
    done
}

backup_netplan() {
    NETPLAN_FILE=$(ls /etc/netplan/*.yaml | head -n1)
    if [ -z "$NETPLAN_FILE" ]; then
        touch "$NETPLAN_FILE"
    else
        log "INFO" "Backed up network config to ${NETPLAN_FILE}.backup" "OS network"
        cp "$NETPLAN_FILE" "${NETPLAN_FILE}.backup"
    fi
}

select_interfaces() {
    interfaces=($(ip -o link show | awk -F': ' '{print $2}' | grep -v 'lo'))
    echo "Please select a network interface for bridge $OTTERSCALE_BRIDGE_NAME"
    echo "Available network interfaces:"
    for i in "${!interfaces[@]}"; do
        echo "$((i+1))) ${interfaces[$i]}"
    done

    while true; do
        read -p "Select interface to bridge (1-${#interfaces[@]}): " iface_choice
        if [[ $iface_choice =~ ^[0-9]+$ ]] && [ $iface_choice -ge 1 ] && [ $iface_choice -le ${#interfaces[@]} ]; then
            OTTERSCALE_NETWORK_INTERFACE=${interfaces[$((iface_choice-1))]}
            break
        fi
        log "WARN" "Invalid selection. Please try again." "OS network"
    done
}

enter_bridge_name() {
    while true; do
        read -p "Please enter bridge name, (e.g. br-pxe): " OTTERSCALE_BRIDGE_NAME
        read -p "You entered: $OTTERSCALE_BRIDGE_NAME. Is this correct? [y/n]: " confirm
        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            break
        elif [[ "$confirm" =~ ^[Nn]$ ]]; then
            echo "Please re-enter the bridge name."
        else
            echo "Invalid input. Please enter y or n."
        fi
    done
}

create_netplan() {
    cat > "$NETPLAN_FILE" <<EOF
network:
  version: 2
  renderer: networkd
  ethernets:
    $OTTERSCALE_NETWORK_INTERFACE:
      dhcp4: no
      dhcp6: no
  bridges:
    $OTTERSCALE_BRIDGE_NAME:
      link-local: []
      interfaces: [$OTTERSCALE_NETWORK_INTERFACE]
      addresses: [$OTTERSCALE_INTERFACE_IP/$OTTERSCALE_INTERFACE_IP_MASK]
      routes:
      - to: default
        via: $OTTERSCALE_INTERFACE_GATEWAY
      nameservers:
        addresses: [$OTTERSCALE_INTERFACE_DNS]
EOF
    chmod 600 /etc/netplan/*.yaml
}

get_current_dns() {
    local INTERFACE=$1
    OTTERSCALE_INTERFACE_DNS=$(resolvectl -i $INTERFACE | grep "Current DNS Server" | awk '{print $4}' | paste -sd, -)
    if [ -z "$OTTERSCALE_INTERFACE_DNS" ]; then
        log "WARN" "No dns found for $INTERFACE, used 8.8.8.8 instead" "OS network"
	OTTERSCALE_INTERFACE_DNS="8.8.8.8"
    fi
}

get_current_ip() {
    local INTERFACE=$1
    local INTERFACE_IPS=($(ip -o -4 addr show dev "$INTERFACE" | awk '{print $4}'))
    local INTERFACE_MASK=($(ip -o -4 addr show dev "$INTERFACE" | awk '{print $4}' | cut -d'/' -f2))

    case ${#INTERFACE_IPS[@]} in
        0)
            error_exit "Network interface $INTERFACE has no IP address assigned"
            ;;
        1)
            OTTERSCALE_INTERFACE_IP=$(echo ${INTERFACE_IPS[0]} | cut -d'/' -f1)
            OTTERSCALE_INTERFACE_IP_MASK=$(echo ${INTERFACE_IPS[0]} | cut -d'/' -f2)
            ;;
        2)
            log "INFO" "Detect multiple IPs on network interface $INTERFACE" "OS network"
            for i in "${!INTERFACE_IPS[@]}"; do
                echo "$((i+1))) ${INTERFACE_IPS[$i]}"
            done

            while true; do
                read -p "Please select the IP you want to use on MAAS: " USER_IP_SELECT
                if validate_ip ${INTERFACE_IPS[$((USER_IP_SELECT-1))]}; then
                    OTTERSCALE_INTERFACE_IP=$(echo ${INTERFACE_IPS[$((USER_IP_SELECT-1))]} | cut -d'/' -f1)
                    OTTERSCALE_INTERFACE_IP_MASK=$(echo ${INTERFACE_IPS[$((USER_IP_SELECT-1))]} | cut -d'/' -f2)
                else
                    log "WARN" "Invalid selection. Please try again." "OS network"
                fi
            done
            ;;
    esac

    log "INFO" "Using bridge $INTERFACE with IP $OTTERSCALE_INTERFACE_IP/$OTTERSCALE_INTERFACE_IP_MASK"
}

get_current_gateway() {
    local INTERFACE=$1
    OTTERSCALE_INTERFACE_GATEWAY=$(ip route show dev $INTERFACE | awk '/default/ {print $3}' | head -1)
    if [ -z "$OTTERSCALE_INTERFACE_GATEWAY" ]; then
        OTTERSCALE_INTERFACE_GATEWAY=$(ip route show | awk '/default/ {print $3}' | head -1)
        log "WARN" "No gateway found for $INTERFACE, using system default: $OTTERSCALE_INTERFACE_GATEWAY" "OS network"
    fi
}

create_new_bridge() {
    log "INFO" "Preparing to create new bridge..." "OS network"
    enter_bridge_name
    select_interfaces
    backup_netplan

    get_current_dns $OTTERSCALE_NETWORK_INTERFACE
    get_current_ip $OTTERSCALE_NETWORK_INTERFACE
    get_current_gateway $OTTERSCALE_NETWORK_INTERFACE
    log "INFO" "Creating bridge $OTTERSCALE_BRIDGE_NAME with network interface $OTTERSCALE_NETWORK_INTERFACE..." "OS network"
    log "INFO" "Using existing IP: $OTTERSCALE_INTERFACE_IP, Gateway: $OTTERSCALE_INTERFACE_GATEWAY, DNS: $OTTERSCALE_INTERFACE_DNS" "OS network"

    create_netplan
    stop_service "NetworkManager"
    disable_service "NetworkManager"
    start_service "systemd-networkd"
    enable_service "systemd-networkd"
    netplan apply || error_exit "Failed to apply netplan configuration"

    log "INFO" "Successfully created bridge $OTTERSCALE_BRIDGE_NAME with IP $OTTERSCALE_INTERFACE_IP" "OS network"
}
