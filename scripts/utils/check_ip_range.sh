#!/bin/bash

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
    # Extract network and mask from subnet
    local network=$(echo $subnet | cut -d'/' -f1)
    local mask=$(echo $subnet | cut -d'/' -f2)

    # Convert mask to dotted decimal format
    local mask_dotted=$(printf "%d.%d.%d.%d" \
        $((0xFF << (32 - mask) >> 24 & 0xFF)) \
        $((0xFF << (32 - mask) >> 16 & 0xFF)) \
        $((0xFF << (32 - mask) >> 8 & 0xFF)) \
        $((0xFF << (32 - mask) & 0xFF)))

    # Check if start_ip and end_ip are in the network
    if is_ip_in_network $start_ip $network $mask_dotted; then
        log "INFO" "Start IP $start_ip is in the network $subnet"
        return 0
    else
        log "WARN" "Start IP $start_ip is not in the network $subnet"
        return 1
    fi

    if is_ip_in_network $end_ip $network $mask_dotted; then
        log "INFO" "End IP $end_ip is in the network $subnet"
        return 0
    else
        log "WARN" "End IP $end_ip is not in the network $subnet"
        return 1
    fi

    log "INFO" "IP range $start_ip to $end_ip is within the network $subnet"
}
