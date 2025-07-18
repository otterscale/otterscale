#!/bin/bash

get_interface_through_ip() {
    local CIDR=$1
    OTTERSCALE_NETWORK_INTERFACE=$(ip -br addr show to $CIDR | awk '{print $1}')
    if [ -z $OTTERSCALE_NETWORK_INTERFACE ]; then
        error_exit "Failed get network interface from $CIDR"
    fi
}

is_interface_bridge() {
    local INTERFACE=$1
    if ! brctl show $INTERFACE > /dev/null 2>&1; then
        error_exit "Network interface $INTERFACE is not a network bridge"
    fi
    return 0
}

get_maas_cidr() {
    while true; do
        read -p "Please enter the CIDR that you want to use on MAAAS: " $OTTERSCALE_CONFIG_MAAS_CIDR
	if validate_cidr $OTTERSCALE_CONFIG_MAAS_CIDR; then
            break
        fi
        log "WARN" "Invild CIDR. Please try agein" "OS network" 
    done
}

check_bridge() {
    if [ -z $OTTERSCALE_CONFIG_MAAS_CIDR ]; then
        get_maas_cidr
    fi
    get_interface_through_ip $OTTERSCALE_CONFIG_MAAS_CIDR

    if is_interface_bridge $OTTERSCALE_NETWORK_INTERFACE ;then
        OTTERSCALE_BRIDGE_NAME=$OTTERSCALE_NETWORK_INTERFACE
        OTTERSCALE_INTERFACE_IP=$(echo $OTTERSCALE_CONFIG_MAAS_CIDR | cut -d'/' -f1)
        OTTERSCALE_INTERFACE_IP_MASK=$(echo $OTTERSCALE_CONFIG_MAAS_CIDR | cut -d'/' -f2)
	get_current_dns $OTTERSCALE_BRIDGE_NAME
    fi
}

get_current_dns() {
    local INTERFACE=$1
    OTTERSCALE_INTERFACE_DNS=$(resolvectl -i $INTERFACE | grep "Current DNS Server" | awk '{print $4}' | paste -sd, -)
    if [ -z "$OTTERSCALE_INTERFACE_DNS" ]; then
        log "WARN" "No dns found for $INTERFACE, used 8.8.8.8 instead" "OS network"
	OTTERSCALE_INTERFACE_DNS="8.8.8.8"
    fi
}

# Function to convert an IP address to a number
ip_to_number() {
    local ip=$1
    local -a octets=(${ip//./ })
    echo $((octets[0] * 256**3 + octets[1] * 256**2 + octets[2] * 256 + octets[3]))
}

# Function to convert a network and mask to a number
network_to_number() {
    local network=$1
    local mask=$2
    local -a octets=(${network//./ })
    local -a mask_octets=(${mask//./ })
    local network_number=0
    for i in {0..3}; do
        network_number=$((network_number + (octets[i] & mask_octets[i]) * 256**(3-i)))
    done
    echo $network_number
}

# Function to check if an IP is in the network
is_ip_in_network() {
    local ip=$1
    local network=$2
    local mask=$3
    local ip_number=$(ip_to_number $ip)
    local network_number=$(network_to_number $network $mask)
    local mask_number=$(ip_to_number $mask)

    if [ $((ip_number & mask_number)) -eq $network_number ]; then
        return 0
    else
        return 1
    fi
}

check_ip_range() {
    local network=$(echo $MAAS_NETWORK_SUBNET | cut -d'/' -f1)
    local mask=$(echo $MAAS_NETWORK_SUBNET | cut -d'/' -f2)
    local mask_dotted=$(printf "%d.%d.%d.%d" \
        $((0xFF << (32 - mask) >> 24 & 0xFF)) \
        $((0xFF << (32 - mask) >> 16 & 0xFF)) \
        $((0xFF << (32 - mask) >> 8 & 0xFF)) \
        $((0xFF << (32 - mask) & 0xFF)))

    # Check if start_ip and end_ip are in the network
    if is_ip_in_network $DHCP_START_IP $network $mask_dotted; then
        if is_ip_in_network $DHCP_END_IP $network $mask_dotted; then
            log "INFO" "IP range $DHCP_START_IP to $DHCP_END_IP is within the network $MAAS_NETWORK_SUBNET"
            return 0
        else
            log "WARN" "End IP $DHCP_END_IP is not in the network $MAAS_NETWORK_SUBNET"
            return 1
        fi
    else
        log "WARN" "Start IP $DHCP_START_IP is not in the network $MAAS_NETWORK_SUBNET"
        return 1
    fi
}
