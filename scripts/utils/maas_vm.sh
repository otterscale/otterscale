#!/bin/bash

create_maas_lxd_project() {
    if ! lxc project list --format json | jq --exit-status '.[] | select(.name == "maas")' >>"$TEMP_LOG" 2>&1; then
        lxc project create maas >>"$TEMP_LOG" 2>&1
        log "INFO" "Create lxd project maas."
    fi

    if ! lxc profile show default | lxc profile edit default --project maas >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update LXD profile."
    fi
}

search_available_vmhost() {
    while IFS= read -r host; do
        VM_HOST_ID=$(echo "$host" | jq -r '.id')
        local available_cores=$(echo "$host" | jq -r '.available.cores')
        local available_memory=$(echo "$host" | jq -r '.available.memory / 1024'| bc -l | xargs printf "%.2f\n") # Convert to GB
        local available_disk=$(echo "$host" | jq -r '.available.local_storage / (1024*1024*1024)' | bc -l | xargs printf "%.2f\n") # Convert to GB

        if [ $(echo "$available_cores >= 1" | bc -l) -eq 1 ] && \
           [ $(echo "$available_memory >= 4" | bc -l) -eq 1 ] && \
           [ $(echo "$available_disk >= 8" | bc -l) -eq 1 ]; then
            log "INFO" "Using existing VM host $VM_HOST_ID with sufficient resources"
            log "DEBUG" "Available resources - Cores: $available_cores, Memory: ${available_memory}GB, Disk: ${available_disk}GB"
            return 0
        fi
    done < <(echo "$existing_hosts" | jq -c '.[]')
}

# LXD VM host creation with validation
create_lxd_vm() {
    log "INFO" "Checking for existing LXD VM hosts..."
    local existing_hosts=$(maas admin vm-hosts read)
    local host_count=$(echo "$existing_hosts" | jq '. | length')

    if [ "$host_count" -gt 0 ]; then
        log "INFO" "Found existing VM hosts, checking resources..."
        search_available_vmhost
    else
        log "INFO" "Creating new LXD VM host..."
        if ! maas admin vm-hosts create \
            password=password \
            type=lxd \
            power_address=https://$BRIDGE_IP:8443 \
            project=maas >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed to create LXD VM host."
        fi
        VM_HOST_ID=$(maas admin vm-hosts read | jq -r '.[0].id')
    fi
    log "INFO" "LXD VM host created successfully (ID: $VM_HOST_ID)"
}

rename_machine() {
    local machine_id=$1
    local new_name=$2
    if ! maas admin machine update $machine_id hostname=$new_name >>"$TEMP_LOG" 2>&1 ; then
        error_exit "Failed to rename machine $machine_id."
    fi
}

wait_commissioning() {
    log "INFO" "Waiting for the machine to transition from commissioning to ready state."
    while true; do
        local status=$(maas admin machine read $machineID | jq -r '.status_name')
        if [ "$status" == "Ready" ]; then
            log "INFO" "Machine $machineID created successfully"
            log "INFO" "Machine juju-vm id is $machineID."
            rename_machine $machineID "juju-vm"
            break
        elif [ "$status" == "Failed commissioning" ]; then
            error_exit "Failed commissioning machine $machineID."
        elif [ "$status" == "Failed testing" ]; then
            error_exit "Failed testing machine $machineID."
        fi
        sleep 10
    done
}

create_vm_from_maas() {
    ## if juju-vm already exist, do not create
    if maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q . >/dev/null 2>&1; then
        log "INFO" "juju-vm already existed, skipping create..."
    else
        log "INFO" "Creating VM on host $VM_HOST_ID..."
	machineID=$(maas admin vm-host compose "$VM_HOST_ID" cores="$LXD_CORES" memory="$LXD_MEMORY" disk=1:size="$LXD_DISK" | jq -r '.system_id')
	if [[ -z $machineID ]]; then
            error_exit "Failed create vm host from kvm lxd $VM_HOST_ID."
	else
            wait_commissioning
	fi
    fi
}

enter_vm_ip() {
    log "INFO" "Please provide an IP address that falls within the range of $subnet."
    while true; do
        read -p "Enter the IP that juju-vm will used : " juju_vm_ip
        if validate_ip "$juju_vm_ip"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

check_vm_ip() {
    local subnet_mode=$(maas admin interfaces read $machineID | jq -r '.[].links' | jq -r '.[] | select(.subnet.name=="'"$subnet"'") | .mode')
    log "INFO" "Machine $machineID interfaces mode is $subnet_mode."

    if [[ $subnet_mode != "static" ]]; then
        enter_vm_ip
        update_vm_ip
        return 0
    fi

    local current_vm_ip=$(maas admin interfaces read $machineID | jq -r '.[].links' | jq '.[] | select(.subnet.name=="'"$subnet"'") | .ip_address')
    while true; do
        read -p "The juju-vm is already configured with a static IP $current_vm_ip, do you want to continue using this IP? [y/n]: " confirm
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

update_vm_ip() {
    log "INFO" "Update $juju_vm_ip to juju-vm."
    local interface_name=$(maas admin interfaces read $machineID | jq -r '.[].name')
    local subnet_cidr=$(maas admin subnet read $subnet | jq -r '.cidr')

    # unlink_subnet
    for id in $(maas admin interfaces read $machineID | jq -r '.[].links | .[].id'); do
        maas admin interface unlink-subnet $machineID $interface_name id=$id >>"$TEMP_LOG" 2>&1
    done

    # link_subnet and give static ip
    if ! maas admin interface link-subnet $machineID $interface_name mode=static subnet=$subnet_cidr ip_address=$juju_vm_ip >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update ip $interface_name to machine $machineID."
    fi
}

set_vm_static_ip() {
    if maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q .; then
        machineID=$(maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm") | .system_id')
	subnet=$(maas admin subnet read $(ip -o -4 addr show dev $bridge | awk '{print $4}') | jq -r '.name')
        check_vm_ip
    else
        error_exit "Machine juju-vm not found."
    fi
}
