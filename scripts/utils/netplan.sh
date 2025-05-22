#!/bin/bash

select_bridge() {
    while true; do
        log "INFO" "Detecting available bridges..."
        bridges=($(brctl show 2>/dev/null | awk 'NR>1 {print $1}' | grep -v '^$'))

        echo "Available network bridges:"
        echo "0) Create new bridge"
        for i in "${!bridges[@]}"; do
            echo "$((i+1))) ${bridges[$i]}"
        done

        read -p "Select bridge (0-${#bridges[@]}): " choice
        case $choice in
            0)
                create_new_bridge
                return
                ;;
            [1-9]*)
                if [ $choice -le ${#bridges[@]} ]; then
                    bridge=${bridges[$((choice-1))]}
                    validate_selected_bridge
                    return
                fi
                ;;
        esac
        log "WARN" "Invalid selection. Please try again."
    done
}

backup_netplan() {
    NETPLAN_FILE=$(ls /etc/netplan/*.yaml | head -n1)
    if [ -z "$NETPLAN_FILE" ]; then
        touch "$NETPLAN_FILE"
    else
        log "INFO" "Backed up network config to ${NETPLAN_FILE}.backup"
        cp "$NETPLAN_FILE" "${NETPLAN_FILE}.backup"
    fi
}

select_interfaces() {
    interfaces=($(ip -o link show | awk -F': ' '{print $2}' | grep -v 'lo'))
    echo "Available network interfaces:"
    for i in "${!interfaces[@]}"; do
        echo "$((i+1))) ${interfaces[$i]}"
    done

    while true; do
        read -p "Select interface to bridge (1-${#interfaces[@]}): " iface_choice
        if [[ $iface_choice =~ ^[0-9]+$ ]] && [ $iface_choice -ge 1 ] && [ $iface_choice -le ${#interfaces[@]} ]; then
            selected_iface=${interfaces[$((iface_choice-1))]}
            break
        fi
        log "WARN" "Invalid selection. Please try again."
    done
}

enter_bridge_name() {
    while true; do
        read -p "Enter bridge name [$DEFAULT_BRIDGE_NAME]: " bridge_name
        bridge_name=${bridge_name:-$DEFAULT_BRIDGE_NAME}

        read -p "You entered: $bridge_name. Is this correct? [y/n]: " confirm
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
    $selected_iface:
      dhcp4: no
      dhcp6: no
  bridges:
    $bridge_name:
      interfaces: [$selected_iface]
      addresses: [$current_ip]
      routes:
      - to: default
        via: $current_gateway
      nameservers:
        addresses: [$current_dns]
EOF
    chmod 600 /etc/netplan/*.yaml
}

get_current_dns() {
    local interface=$1
    current_dns=$(resolvectl -i $interface | grep "Current DNS Server" | awk '{print $4}' | paste -sd, -)
    if [ -z "$current_dns" ]; then
        log "WARN" "No dns found for $interface."
    fi
}

get_current_ip() {
    local interface=$1
    current_ip=$(ip -o -4 addr show dev $interface | awk '{print $4}')
    if [ -z "$current_ip" ]; then
        error_exit "Selected interface $interface has no IP address"
    fi
}

get_current_gw() {
    local interface=$1
    current_gateway=$(ip route show dev $interface | awk '/default/ {print $3}' | head -1)
    if [ -z "$current_gateway" ]; then
        current_gateway=$(ip route show | awk '/default/ {print $3}' | head -1)
        log "WARN" "No gateway found for $interface, using system default: $current_gateway"
    fi
}

create_new_bridge() {
    log "INFO" "Preparing to create new bridge..."
    backup_netplan
    select_interfaces
    enter_bridge_name

    get_current_dns $selected_iface
    get_current_ip $selected_iface
    get_current_gw $selected_iface
    log "INFO" "Creating bridge $bridge_name with interface $selected_iface..."
    log "INFO" "Using existing IP: $current_ip, Gateway: $current_gateway, DNS: $current_dns"

    create_netplan

    stop_service "NetworkManager"
    disable_service "NetworkManager"
    start_service "systemd-networkd"
    enable_service "systemd-networkd"

    netplan apply || error_exit "Failed to apply netplan configuration"

    bridge="$bridge_name"
    BRIDGE_IP=$(echo "$current_ip" | cut -d'/' -f1)
    log "INFO" "Successfully created bridge $bridge with IP $BRIDGE_IP"
}
