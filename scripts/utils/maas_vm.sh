#!/bin/bash

create_maas_lxd_project() {
    if ! lxc project list --format json | jq --exit-status '.[] | select(.name == "maas")' >>"$TEMP_LOG" 2>&1; then
        lxc project create maas >>"$TEMP_LOG" 2>&1
        log "INFO" "Create lxd project maas" "MAAS lxd create"
    fi

    if ! lxc profile show default | lxc profile edit default --project maas >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update LXD profile"
    fi
}

search_available_vmhost() {
    while IFS= read -r host; do
        VM_HOST_ID=$(echo "$host" | jq -r '.id')
        local AVAILABLE_CORES=$(echo "$host" | jq -r '.available.cores')
        local AVAILABLE_MEMORY_GB=$(echo "$host" | jq -r '.available.memory / 1024'| bc -l | xargs printf "%.2f\n") # Convert to GB
        local AVAILABLE_DISK_GB=$(echo "$host" | jq -r '.available.local_storage / (1024*1024*1024)' | bc -l | xargs printf "%.2f\n") # Convert to GB

        if [ $(echo "$AVAILABLE_CORES >= 1" | bc -l) -eq 1 ] && \
           [ $(echo "$AVAILABLE_MEMORY_GB >= 4" | bc -l) -eq 1 ] && \
           [ $(echo "$AVAILABLE_DISK_GB >= 8" | bc -l) -eq 1 ]; then
            log "INFO" "Using existing VM host $VM_HOST_ID with sufficient resources" "MAAS vmhost create"
            log "DEBUG" "Available resources - Cores: $AVAILABLE_CORES, Memory: ${AVAILABLE_MEMORY_GB}GB, Disk: ${AVAILABLE_DISK_GB}GB" "MAAS vmhost create"
            return 0
        fi
    done < <(echo "$MAAS_VM_HOSTS" | jq -c '.[]')
}

# LXD VM host creation with validation
create_lxd_vm() {
    log "INFO" "Checking for existing LXD VM hosts..." "MAAS lxd create"
    local MAAS_VM_HOSTS=$(maas admin vm-hosts read)
    local MAAS_VM_HOST_COUNT=$(echo "$MAAS_VM_HOSTS" | jq '. | length')

    if [ "$MAAS_VM_HOST_COUNT" -gt 0 ]; then
        log "INFO" "Found existing VM hosts, checking resources..." "MAAS lxd create"
        search_available_vmhost
    else
        log "INFO" "Creating new LXD VM host..." "MAAS lxd create"
        execute_cmd "maas admin vm-hosts create password=password type=lxd power_address=https://$OTTERSCALE_INTERFACE_IP:8443 project=maas" "create maas LXD VM host"
        VM_HOST_ID=$(maas admin vm-hosts read | jq -r '.[0].id')
    fi
    log "INFO" "LXD VM host (ID: $VM_HOST_ID) is ready" "MAAS lxd create"
}

rename_machine() {
    local MACHINE_ID=$1
    local MACHINE_NAME=$2
    execute_cmd "maas admin machine update $MACHINE_ID hostname=$MACHINE_NAME" "rename maas machine $MACHINE_ID"
}

wait_commissioning() {
    log "INFO" "Waiting for the machine to transition from commissioning to ready state" "MAAS prepare machine"
    while true; do
        local MACHINE_STATUS=$(maas admin machine read $JUJU_MACHINE_ID | jq -r '.status_name')
        case $MACHINE_STATUS in
            "Ready")
                log "INFO" "Machine $JUJU_MACHINE_ID created successfully" "MAAS prepare machine"
                log "INFO" "Machine juju-vm id is $JUJU_MACHINE_ID" "MAAS prepare machine"
                rename_machine $JUJU_MACHINE_ID "juju-vm"
                break
                ;;
            "Failed commissioning")
                error_exit "Failed commissioning machine $JUJU_MACHINE_ID"
                ;;
            "Failed testing")
                error_exit "Failed testing machine $JUJU_MACHINE_ID"
                ;;
	esac
        sleep 10
    done
}

create_vm_from_maas() {
    ## if juju-vm already exist, do not create
    if maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q . >/dev/null 2>&1; then
        log "INFO" "juju-vm already existed, skipping create..." "MAAS prepare machine"
    else
        log "INFO" "Creating VM on host $VM_HOST_ID..." "MAAS prepare machine"
        JUJU_MACHINE_ID=$(maas admin vm-host compose $VM_HOST_ID cores=$LXD_CORES memory=$LXD_MEMORY_MB disk=1:size=$LXD_DISK_GB | jq -r '.system_id')
        if [[ -z $JUJU_MACHINE_ID ]]; then
            error_exit "Failed create vm host from kvm lxd $VM_HOST_ID"
        else
            wait_commissioning
        fi
    fi
}

enter_vm_ip() {
    log "INFO" "Please provide an IP address that falls within the range of $MAAS_NETWORK_SUBNET"
    while true; do
        read -p "Enter the IP that juju-vm will used: " JUJU_VM_IP
        if validate_ip "$JUJU_VM_IP"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

set_static_vm_ip() {
    if [ ! -z $OTTERSCALE_CONFIG_JUJU_IP ]; then
        JUJU_VM_IP=$OTTERSCALE_CONFIG_JUJU_IP
    else
        enter_vm_ip
    fi
    update_vm_ip
}


ask_user_type_vm_ip() {
    while true; do
        read -p "The juju-vm is already configured with a static IP $CURRENT_JUJU_IP, do you want to continue using this IP? [y/n]: " confirm
        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            break
        elif [[ "$confirm" =~ ^[Nn]$ ]]; then
            enter_vm_ip
            update_vm_ip
            break
        else
            echo "Invalid input. Please enter y or n."
        fi
    done
}

check_modify_vm_ip() {
    CURRENT_JUJU_IP=$(maas admin interfaces read $JUJU_MACHINE_ID | jq -r '.[].links' | jq '.[] | select(.subnet.name=="'"$MAAS_NETWORK_SUBNET"'") | .ip_address')
    if [ ! -z $OTTERSCALE_CONFIG_JUJU_IP ]; then
        if [ $OTTERSCALE_CONFIG_JUJU_IP != $CURRENT_JUJU_IP ]; then
            JUJU_VM_IP=$OTTERSCALE_CONFIG_JUJU_IP
            update_vm_ip
        fi
    else
        ask_user_type_vm_ip	    
    fi
}

check_vm_ip() {
    local JUJU_MACHINE_SUBNET_MODE=$(maas admin interfaces read $JUJU_MACHINE_ID | jq -r '.[].links' | jq -r '.[] | select(.subnet.name=="'"$MAAS_NETWORK_SUBNET"'") | .mode')
    log "INFO" "Machine $JUJU_MACHINE_ID interfaces mode is $JUJU_MACHINE_SUBNET_MODE"
    if [[ $JUJU_MACHINE_SUBNET_MODE != "static" ]]; then
        set_static_vm_ip
    else
        check_modify_vm_ip
    fi
}

update_vm_ip() {
    log "INFO" "Update $JUJU_VM_IP to juju-vm."
    local JUJU_MACHINE_INTERFACE_NAME=$(maas admin interfaces read $JUJU_MACHINE_ID | jq -r '.[].name')
    local MAAS_SUBNET_CIDR=$(maas admin subnet read $MAAS_NETWORK_SUBNET | jq -r '.cidr')

    # unlink_subnet
    for id in $(maas admin interfaces read $JUJU_MACHINE_ID | jq -r '.[].links | .[].id'); do
        maas admin interface unlink-subnet $JUJU_MACHINE_ID $JUJU_MACHINE_INTERFACE_NAME id=$id >>"$TEMP_LOG" 2>&1
    done

    # link_subnet and give static ip
    execute_cmd "maas admin interface link-subnet $JUJU_MACHINE_ID $JUJU_MACHINE_INTERFACE_NAME mode=static subnet=$MAAS_SUBNET_CIDR ip_address=$JUJU_VM_IP" "update ip $JUJU_MACHINE_INTERFACE_NAME to machine $JUJU_MACHINE_ID"
}

set_vm_static_ip() {
    if maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q .; then
        JUJU_MACHINE_ID=$(maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm") | .system_id')
        MAAS_NETWORK_SUBNET=$(maas admin subnet read $(ip -o -4 addr show dev $OTTERSCALE_BRIDGE_NAME | awk '{print $4}') | jq -r '.name')
	if ! is_machine_deployed; then
            check_vm_ip
        fi
    else
        error_exit "Machine juju-vm not found"
    fi
}
