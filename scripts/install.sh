#!/bin/bash

##
# Host requirment
MIN_MEMORY_GB=8
MIN_DISK_GB=100
OTTERSCALE_OS="24.04"
OTTERSCALE_MAAS_VERSION="2.0"
OTTERSCALE_BASE_IMAGE="ubuntu@24.04"

##
# Host LXD config
LXD_STORAGE_SIZE_GB="60GB"
LXD_CORES=2
LXD_MEMORY_MB=4096
LXD_DISK_GB=50G

##
# Install packages
APT_PACKAGES="jq openssh-server bridge-utils openvswitch-switch"
SNAP_PACKAGES="core24 maas maas-test-db juju lxd microk8s"

##
#Snap version
CORE24_CHANNEL="latest/stable"
MAAS_CHANNEL="3.6/stable"
MAAS_DB_CHANNEL="3.6/stable"
JUJU_CHANNEL="3.6/stable"
LXD_CHANNEL="5.0/stable"
MICROK8S_CHANNEL="1.32/stable"
CONTROLLER_CHARM_CHANNEL="3.5/stable"

##
# Otterscale
OTTERSCALE_MAX_RETRIES=5
OTTERSCALE_CHARMHUB_URL="https://api.charmhub.io"
OTTERSCALE_MAAS_ADMIN_USER="admin"
OTTERSCALE_MAAS_ADMIN_PASS="admin"
OTTERSCALE_MAAS_ADMIN_EMAIL="admin@example.com"
export OTTERSCALE_INSTALL_DIR=$(dirname "$(readlink -f $0)")

##
# LOG
export TEMP_LOG=$(mktemp)
export LOG=$OTTERSCALE_INSTALL_DIR/setup.log
touch "$LOG"
chmod 666 "$LOG"


# ------------------------------------------------------------
#  Functions –  English documentation and log‑message polish
# ------------------------------------------------------------

# Update apt package lists
apt_update() {
    log "INFO" "Executing command apt update..." "APT update"
    if ! apt-get update --fix-missing >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update apt package lists – check network connectivity"
    fi
}

# Install a list of apt packages
apt_install() {
    local PKG_LIST=$1
    log "INFO" "Installing apt packages: $PKG_LIST" "APT Install"
    if ! DEBIAN_FRONTEND=noninteractive apt-get install -y $PKG_LIST >>"$TEMP_LOG" 2>&1; then
        error_exit "APT package installation failed"
    fi
    log "INFO" "Apt packages installed successfully" "APT Install"
}

# Send a JSON payload to the Otterscale endpoint
send_request() {
    local URL_PATH=$1
    local DATA=$2
    if ! curl -s --header "Content-Type: application/json" --data "$DATA" "$OTTERSCALE_ENDPOINT$URL_PATH" > /dev/null 2>&1 ; then
        echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] Failed to execute curl request"
        trap cleanup EXIT
        exit 1
    fi
}

# Report the current phase/status to Otterscale
send_status_data() {
    local PHASE=$1
    local MESSAGE=$2
    local DATA=$(cat <<EOF
{
"phase": "$PHASE",
"message": "$MESSAGE"
}
EOF
)

    send_request "/otterscale.environment.v1.EnvironmentService/UpdateStatus" "$DATA"
}

# ----------------------------------------------------------------
#  Function: send_otterscale_config_data
#  Description: Build a JSON payload with MAAS/Juju configuration
#               and POST it to the Otterscale service.
# ---------------------------------------------------------------
send_otterscale_config_data() {
    local OTTERSCALE_MAAS_ENDPOINT="http://$OTTERSCALE_INTERFACE_IP:5240/MAAS"
    local OTTERSCALE_MAAS_KEY=$(su "$NON_ROOT_USER" -c "juju show-credentials maas-cloud maas-cloud-credential --show-secrets --client | grep maas-oauth | awk '{print \$2}'")
    local OTTERSCALE_CONTROLLER=$(su "$NON_ROOT_USER" -c "juju controllers --format json | jq -r '.\"current-controller\"'")
    local OTTERSCALE_CONTROLLER_DETIAL=$(su "$NON_ROOT_USER" -c "OTTERSCALE_CONTROLLER=\$(juju controllers --format json | jq -r '.\"current-controller\"'); juju show-controller \$OTTERSCALE_CONTROLLER --show-password --format=json")
    local OTTERSCALE_JUJU_ENDPOINTS=$(echo "$OTTERSCALE_CONTROLLER_DETIAL" | jq -r '."'"$OTTERSCALE_CONTROLLER"'"."details"."api-endpoints"' | tr '\n' ' ' | sed 's/ \+/ /g' | grep -v '^ *\[[0-9a-fA-F:]\+.*')
    local OTTERSCALE_JUJU_USERNAME=$(echo "$OTTERSCALE_CONTROLLER_DETIAL" | jq -r '."'"$OTTERSCALE_CONTROLLER"'"."account"."user"')
    local OTTERSCAKE_JUJU_PASSWORD=$(echo "$OTTERSCALE_CONTROLLER_DETIAL" | jq -r '."'"$OTTERSCALE_CONTROLLER"'"."account"."password"')
    local OTTERSCALE_JUJU_CACERT=$(echo "$OTTERSCALE_CONTROLLER_DETIAL" | jq -r '."'"$OTTERSCALE_CONTROLLER"'"."details"."ca-cert"')
    local OTTERSCALE_JUJU_CLOUD_NAME="maas-cloud"
    local OTTERSCALE_JUJU_REGION="default"
    local OTTERSCALE_K8S_ENDPOINT_JSON=$(microk8s kubectl get endpoints -o json | jq '.items[].subsets[]')
    local OTTERSCALE_K8S_ENDPOINT=$(echo "$OTTERSCALE_K8S_ENDPOINT_JSON" | jq -r '.ports[].name')"://"$(echo $OTTERSCALE_K8S_ENDPOINT_JSON | jq -r '.addresses[].ip')":"$(echo $OTTERSCALE_K8S_ENDPOINT_JSON | jq '.ports[].port')
    local OTTERSCALE_MICROK8S_ENCODE_TOKEN=$(base64 -i --wrap=0 "$KUBE_FOLDER/config")
    local DATA=$(cat <<EOF
{
"maas_url": "$OTTERSCALE_MAAS_ENDPOINT",
"maas_key": "$OTTERSCALE_MAAS_KEY",
"maas_version": "$OTTERSCALE_MAAS_VERSION",
"juju_controller": "$OTTERSCALE_CONTROLLER",
"juju_controller_addresses": $OTTERSCALE_JUJU_ENDPOINTS,
"juju_username": "$OTTERSCALE_JUJU_USERNAME",
"juju_password": "$OTTERSCAKE_JUJU_PASSWORD",
"juju_ca_cert": $(echo "$OTTERSCALE_JUJU_CACERT" | jq -sRr '@json'),
"juju_cloud_name": "$OTTERSCALE_JUJU_CLOUD_NAME",
"juju_cloud_region": "$OTTERSCALE_JUJU_REGION",
"juju_charmhub_api_url": "$OTTERSCALE_CHARMHUB_URL",
"micro_k8s_token": "$OTTERSCALE_MICROK8S_ENCODE_TOKEN",
"micro_k8s_host": "$OTTERSCALE_K8S_ENDPOINT"
}
EOF
)

    send_request "/otterscale.environment.v1.EnvironmentService/UpdateConfig" "$DATA"
}

# Execute a Juju command as the non‑root user and abort on failure
juju_cmd() {
    local CMD=$1
    local MSG=$2
    log "INFO" "Execute Juju command: $CMD" "$MSG"
    if ! execute_as_user "$NON_ROOT_USER" "$CMD"; then
        error_exit "Failed $MSG"
    fi
}

# Write the Juju cloud definition file (executed as the non‑root user)
generate_clouds_yaml() {
    log "INFO" "Generating Juju cloud definition file ($JUJU_CLOUD)" "Juju clouds"
    su "$NON_ROOT_USER" -c 'cat > $JUJU_CLOUD <<EOF
clouds:
  maas-cloud:
    type: maas
    description: Metal As A Service
    auth-types: [oauth1]
    endpoint: http://$OTTERSCALE_INTERFACE_IP:5240/MAAS/api/2.0/
    regions:
      default: {}
EOF'
}

# Write the Juju credential definition file (executed as the non‑root user)
generate_credentials_yaml() {
     log "INFO" "Generating Juju credential file ($JUJU_CREDENTIAL)" "Juju credentials"
    su "$NON_ROOT_USER" -c 'cat > $JUJU_CREDENTIAL <<EOF
credentials:
  maas-cloud:
    maas-cloud-credential:
      auth-type: oauth1
      maas-oauth: $APIKEY
EOF'
}

# Register the Maas cloud with Juju (if it does not already exist)
juju_clouds() {
    log "INFO" "Configuring Juju clouds..." "JuJu clouds"
    generate_clouds_yaml

    if su "$NON_ROOT_USER" -c 'juju clouds 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Juju cloud maas-cloud already exists – skipping creation" "JuJu clouds"
    else
        juju_cmd "juju add-cloud maas-cloud $JUJU_CLOUD --client --debug" "add juju cloud"
    fi
}

# Register the Maas credential with Juju (if it does not already exist)
juju_credentials() {
    log "INFO" "Configuring Juju credentials..." "JuJu credentials"
    generate_credentials_yaml

    if su "$NON_ROOT_USER" -c 'juju credentials 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Juju credential for maas-cloud already exists – skipping creation" "JuJu credentials"
    else
        juju_cmd "juju add-credential maas-cloud -f $JUJU_CREDENTIAL --client --debug" "add juju credential"
    fi
}

# Check whether a machine named 'juju‑vm' exists in MAAS
is_machine_exist() {
    if maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q . > /dev/null 2>&1; then
        return 0
    fi
    return 1
}

# Check whether the 'juju‑vm' machine is already deployed
is_machine_deployed() {
    if [ $(maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | jq -r '.status_name') == Deployed ]; then
        return 0
    fi
    return 1
}

# Bootstrap Juju on the previously provisioned MAAS machine
bootstrap_juju() {
    su "$NON_ROOT_USER" -c 'mkdir -p ~/.local/share/juju'
    su "$NON_ROOT_USER" -c 'mkdir -p ~/otterscale-tmp'

    export JUJU_CLOUD=/home/$NON_ROOT_USER/otterscale-tmp/cloud.yaml
    export JUJU_CREDENTIAL=/home/$NON_ROOT_USER/otterscale-tmp/credential.yaml
    export OTTERSCALE_INTERFACE_IP=$OTTERSCALE_INTERFACE_IP
    export APIKEY=$APIKEY

    juju_clouds
    juju_credentials

    rm -rf /home/"$NON_ROOT_USER"/otterscale-tmp
    unset JUJU_CLOUD
    unset JUJU_CREDENTIAL
    unset APIKEY

    bootstrap_cmd="juju bootstrap maas-cloud maas-cloud-controller --bootstrap-base=$OTTERSCALE_BASE_IMAGE"
    bootstrap_config="--config default-base=$OTTERSCALE_BASE_IMAGE --controller-charm-channel=$CONTROLLER_CHARM_CHANNEL"
    bootstrap_machine="--to juju-vm"

    if ! is_machine_exist; then
        error_exit "Juju bootstrap failed – juju‑vm not found in MAAS"
    fi

    if is_machine_deployed; then
        log "INFO" "Juju is already been bootstrapped – skipping" "JuJu bootstrap"
    else
        log "INFO" "Bootstrapping Juju (this may take several minutes)..." "JuJu bootstrap"
        juju_cmd "$bootstrap_cmd $bootstrap_config $bootstrap_machine --debug" "juju bootstrap"
        log "INFO" "MAAS and Juju setup completed successfully!" "Finished bootstrap"
    fi
}

# Deploy a Kubernetes cluster via Juju
juju_add_k8s() {
    if execute_as_user "$NON_ROOT_USER" "juju show-cloud cos-k8s --debug"; then
        log "INFO" "K8s cloud cos‑k8s already exists – skipping" "JuJu cloud"
    else
        juju_cmd "juju add-k8s cos-k8s --controller maas-cloud-controller --client --debug"
    fi

    if execute_as_user "$NON_ROOT_USER" "juju show-model cos"; then
        log "INFO" "Model cos already exists – skipping" "JuJu model"
    else
        juju_cmd "juju add-model cos cos-k8s --debug" "execute juju add-model"
        juju_cmd "juju deploy cos-lite --trust --debug" "juju deploy cos-lite"
    fi

    juju_cmd "juju config -m cos prometheus metrics_retention_time=180d --debug" "update metric retention time to 180 days"
    juju_cmd "juju config -m cos prometheus maximum_retention_size=70% --debug" "update max retention size to 70%"
    juju_cmd "juju offer cos.grafana:grafana-dashboard global-grafana --debug" "offer grafana-dashboard"
    juju_cmd "juju offer cos.prometheus:receive-remote-write global-prometheus --debug" "offer prometheus-receive-remote-write"
}

# Centralised logger – also forwards status to Otterscale
log() {
    local LOG_LEVEL=$1
    local MESSAGE=$2
    local PHASE=$3
    local TIMESTAMP="$(date '+%Y-%m-%d %H:%M:%S')"
    echo "$TIMESTAMP [$LOG_LEVEL] $MESSAGE" | tee -a "$OTTERSCALE_INSTALL_DIR/setup.log"
    send_status_data "$PHASE" "$MESSAGE"
}

# Abort execution, dump temporary log and clean up
error_exit() {
    local MESSAGE=$1
    log "ERROR" "$MESSAGE" "ERROR"
    if [ -s "$TEMP_LOG" ]; then
        log "DEBUG" "Full error output:" "ERROR"
        cat "$TEMP_LOG" | while read line; do log "DEBUG" "$line" "ERROR"; done
    fi
    trap cleanup EXIT
    exit 1
}

# Remove temporary files created during the run
cleanup() {
    echo "Cleaning up temporary files..."
    rm -f "$TEMP_LOG"
    log "INFO" "Temporary log file removed, cleanup finished" "Cleanup"
}

# Generate LXD pre‑seed configuration file
generate_lxd_config() {
    log "INFO" "Generating LXD pre‑seed configuration ($lxd_file)" "LXD init"
cat > "$lxd_file" <<EOF
config:
  core.https_address: '[::]:8443'
  core.trust_password: password
storage_pools:
- config:
    size: $LXD_STORAGE_SIZE_GB
  description: ""
  name: default
  driver: zfs
profiles:
- name: default
  config:
    boot.autostart: "true"
  description: ""
  devices:
    eth0:
      name: eth0
      nictype: bridged
      parent: $OTTERSCALE_BRIDGE_NAME
      type: nic
    root:
      path: /
      pool: default
      type: disk
projects: []
cluster: null
EOF

    log "INFO" "LXD pre‑seed file generated" "LXD init"
}

# Initialise LXD using the pre‑seed file
init_lxd() {
    lxd_file=$OTTERSCALE_INSTALL_DIR/lxd-config.yaml
    generate_lxd_config

    log "INFO" "Initialising LXD with bridge $OTTERSCALE_BRIDGE_NAME..." "LXD init"
    if ! cat "$lxd_file" | lxd init --preseed >>"$TEMP_LOG" 2>&1; then
        error_exit "LXD initialisation failed"
    else
        log "INFO" "LXD initialised successfully" "LXD init"
        rm -f "$lxd_file"
    fi
}

# Update an existing MAAS boot source (or create a new one if none exists)
update_boot_source(){
    MAAS_BOOT_SOURCE_ID=$(echo "$MAAS_BOOT_SOURCES" | jq -r '.[0].id')
    MAAS_BOOT_SELECTION_ID=$(maas admin boot-source-selections read "$MAAS_BOOT_SOURCE_ID" | jq -r '.[0].id')
    log "INFO" "Modifying existing boot source (ID: $MAAS_BOOT_SOURCE_ID)" "MAAS boot image"

    if [ "$MAAS_BOOT_SOURCE_COUNT" -gt 1 ]; then
        log "INFO" "Removing duplicate boot sources..." "MAAS boot image"
        for id in $(echo "$sources" | jq -r '.[].id' | tail -n +2); do
            execute_cmd "maas admin boot-source delete $id" "remove maas boot source $id"
        done
    fi

    log "INFO" "Updating boot source to Ubuntu Noble (24.04) amd64..." "MAAS boot image"
    if ! maas admin boot-source-selection update "$MAAS_BOOT_SOURCE_ID" "$MAAS_BOOT_SELECTION_ID" \
        release=noble \
        arches=amd64 \
        subarches="*" \
        labels="*" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to update maas boot source"
    fi
}

create_boot_source() {
    # No existing source - create new one
    log "INFO" "Creating new boot source for Ubuntu Noble (24.04) amd64..." "MAAS boot image"
    if ! maas admin boot-sources create \
        url="http://images.maas.io/ephemeral-v3/stable/" \
        keyring_filename="/usr/share/keyrings/ubuntu-cloudimage-keyring.gpg" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to create maas boot source"
    fi

    # Get the new source ID
    MAAS_BOOT_SOURCE_ID=$(maas admin boot-sources read | jq -r '.[0].id')

    # Create selection for Jammy amd64
    if ! maas admin boot-source-selections create "$MAAS_BOOT_SOURCE_ID" \
        release=noble \
        arches=amd64 \
        subarches="*" \
        labels="*" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed create maas boot source selection"
    fi
}

# Launch MAAS image import and wait for it to finish
start_import() {
    log "INFO" "Starting download MAAS boot image..." "MAAS boot image"
    maas admin boot-resources stop-import >>"$TEMP_LOG" 2>&1
    sleep 10

    execute_cmd "maas admin boot-resources import" "start maas image download"
    sleep 10

    log "INFO" "Waiting for image download to complete..." "MAAS boot image"
    while true; do
        if [ $(maas admin boot-resources is-importing | jq -r) != "true" ]; then
            break
        fi
        sleep 10
    done
}

# Choose between creating a new source or updating an existing one, then import images
download_maas_img() {
    log "INFO" "Configuring MAAS boot sources..." "MAAS boot image"

    MAAS_BOOT_SOURCES=$(maas admin boot-sources read)
    MAAS_BOOT_SOURCE_COUNT=$(echo "$MAAS_BOOT_SOURCES" | jq '. | length')
    if [ "$MAAS_BOOT_SOURCE_COUNT" -gt 0 ]; then
        update_boot_source
    else
        create_boot_source
    fi

    start_import
    set_config "commissioning_distro_series" "noble"
    set_config "default_distro_series" "noble"
    set_config "default_osystem" "ubuntu"
    log "INFO" "MAAS images downloaded successfully" "MAAS boot image"
}

# Initialise MAAS unless it has already been set up
init_maas() {
    # Check if MAAS admin user already exists
    if maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "INFO" "MAAS is already initialized (user $OTTERSCALE_MAAS_ADMIN_USER exist). Skipping initialization..."
        return 0
    fi

    log "INFO" "Initializing MAAS..." "MAAS init"
    execute_cmd "maas init region+rack --database-uri maas-test-db:/// --maas-url http://$OTTERSCALE_INTERFACE_IP:5240/MAAS" "maas initialization"
}

# Create the MAAS admin user (if not already present)
create_maas_admin() {
    log "INFO" "Creating MAAS admin user..." "MAAS init"
    if maas apikey --username "$OTTERSCALE_MAAS_ADMIN_USER" >/dev/null 2>&1; then
        log "INFO" "Admin user '$OTTERSCALE_MAAS_ADMIN_USER' already exists. Using existing credentials" "MAAS init"
    else
        execute_cmd "maas createadmin --username $OTTERSCALE_MAAS_ADMIN_USER --password $OTTERSCALE_MAAS_ADMIN_PASS --email $OTTERSCALE_MAAS_ADMIN_EMAIL" "create MAAS admin user"
    fi
    log "INFO" "MAAS web url: http://$OTTERSCALE_INTERFACE_IP:5240/MAAS" "MAAS init"
    log "INFO" "MAAS Username: $OTTERSCALE_MAAS_ADMIN_USER" "MAAS init"
    log "INFO" "MAAS Password: $OTTERSCALE_MAAS_ADMIN_PASS" "MAAS init"
}

# Login to MAAS, retrying on failure
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
    if [ -z "$maas_current_dns" ]; then
        dns_value="$CURRENT_DNS"
        log "INFO" "MAAS upstream DNS not set, will use system DNS: $dns_value" "MAAS config"
    fi
}

set_config() {
    local NAME=$1
    local VALUE=$2
    log "INFO" "Update Config: $NAME, Value: $VALUE" "MAAS config update"
    execute_cmd "maas admin maas set-config name=$NAME value=$VALUE" "set maas $NAME config"
}

# Update MAAS configuration values
update_maas_config() {
    get_maas_dns
    set_config "upstream_dns" "$dns_value"
    set_config "boot_images_auto_import" "false"
    set_config "enable_http_proxy" "false"
    set_config "enable_analytics" "false"
    set_config "network_discovery" "disabled"
    set_config "release_notifications" "false"
}

enter_dhcp_start_ip() {
    while true; do
        read -p "Enter DHCP start IP: " MAAS_DHCP_START_IP
        if validate_ip "$MAAS_DHCP_START_IP"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

enter_dhcp_end_ip() {
    while true; do
        read -p "Enter DHCP end IP: " MAAS_DHCP_END_IP
        if validate_ip "$MAAS_DHCP_END_IP"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

update_fabric_dns() {
    local FABRIC_DNS=$(maas admin subnet read "$MAAS_NETWORK_SUBNET" | jq -r '.dns_servers[]')
    log "INFO" "Update dns $CURRENT_DNS to fabric $MAAS_NETWORK_SUBNET" "MAAS config update"

    if [[ "$FABRIC_DNS" =~ "$CURRENT_DNS" ]]; then
        log "INFO" "Current dns already existed, skipping..." "MAAS config update"
    elif [ ! -n "$maas_current_dns" ]; then
        if [ -z "$FABRIC_DNS" ]; then
            dns_value="$CURRENT_DNS"
        else
            dns_value="$FABRIC_DNS $CURRENT_DNS"
        fi
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
    if ! maas admin ipranges create type=dynamic start_ip="$MAAS_DHCP_START_IP" end_ip="$MAAS_DHCP_END_IP" >>"$TEMP_LOG" 2>&1; then
        log "WARN" "Please confirm if address is within subnet $MAAS_NETWORK_SUBNET, or maybe it conflicts with an existing IP address or range" "MAAS config update"
        error_exit "Failed to create DHCP range"
    fi
}

update_dhcp_config() {
    local ENABLED=$1
    log "INFO" "Set MAAS VLAN DHCP to $ENABLED" "MAAS config update"
    if ! maas admin vlan update $FABRIC_ID $VLAN_TAG dhcp_on=$ENABLED primary_rack=$PRIMARY_RACK >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed to set MAAS DHCP to $ENABLED"
    fi
}

get_dhcp_subnet_and_ip() {
    if [ -z "$MAAS_DHCP_START_IP" ]; then
        enter_dhcp_start_ip
    fi

    if [ -z "$MAAS_DHCP_END_IP" ]; then
        enter_dhcp_end_ip
    fi
}

# ----------------------------------------------------------------
#  Function: enable_maas_dhcp
#  Description: Configure DHCP settings in MAAS if no dynamic ranges exist.
# ----------------------------------------------------------------
enable_maas_dhcp() {
    if [ $(maas admin ipranges read | jq '. | length') -ne 0 ]; then
        log "INFO" "MAAS already has dynamic IP ranges – skipping DHCP configuration" "MAAS config update"
        return 0
    fi

    log "INFO" "Configuring MAAS DHCP..." "MAAS config update"
    MAAS_NETWORK_SUBNET=$OTTERSCALE_CIDR
    while true; do
        get_dhcp_subnet_and_ip
        if check_ip_range ; then
            update_fabric_dns
            get_fabric
            create_dhcp_iprange
            update_dhcp_config "True"
            log "INFO" "DHCP configuration completed" "MAAS config update"
            break
        else
            if [ -n "$MAAS_NETWORK_SUBNET" ] && [ -n "$MAAS_DHCP_START_IP" ] && [ -n "$MAAS_DHCP_END_IP" ]; then
                break
            fi
        fi
    done
}

# Create an LXD project for MAAS and adjust the default profile
create_maas_lxd_project() {
    if ! lxc project list --format json | jq --exit-status '.[] | select(.name == "maas")' >>"$TEMP_LOG" 2>&1; then
        lxc project create maas >>"$TEMP_LOG" 2>&1
        log "INFO" "Create LXD project maas" "MAAS lxd create"
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
            log "INFO" "Available resources - Cores: $AVAILABLE_CORES, Memory: ${AVAILABLE_MEMORY_GB}GB, Disk: ${AVAILABLE_DISK_GB}GB" "MAAS vmhost create"
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
    log "INFO" "Machine $MACHINE_ID renamed to $MACHINE_NAME" "Machine rename"
}

wait_commissioning() {
    log "INFO" "Waiting for the machine to transition from commissioning to ready state" "MAAS prepare machine"
    while true; do
        local MACHINE_STATUS=$(maas admin machine read $JUJU_MACHINE_ID | jq -r '.status_name')
        case $MACHINE_STATUS in
            "Ready")
                log "INFO" "Machine $JUJU_MACHINE_ID created successfully" "MAAS prepare machine"
                log "INFO" "Machine juju-vm id is $JUJU_MACHINE_ID" "MAAS prepare machine"
                rename_machine "$JUJU_MACHINE_ID" "juju-vm"
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
        log "INFO" "Creating VM from kvm lxd id $VM_HOST_ID..." "MAAS prepare machine"
        JUJU_MACHINE_ID=$(maas admin vm-host compose "$VM_HOST_ID" cores="$LXD_CORES" memory="$LXD_MEMORY_MB" disk=1:size="$LXD_DISK_GB" | jq -r '.system_id')
        if [[ -z $JUJU_MACHINE_ID ]]; then
            error_exit "Failed create vm host from kvm lxd id $VM_HOST_ID"
        else
            wait_commissioning
        fi
    fi
}

enter_vm_ip() {
    log "INFO" "Please provide an IP address that falls within the range of $MAAS_NETWORK_SUBNET" "Enter VM IP"
    while true; do
        read -p "Enter the IP that juju-vm will used: " JUJU_VM_IP
        if validate_ip "$JUJU_VM_IP"; then
            break
        fi
        echo "Invalid IP format. Please try again."
    done
}

set_static_vm_ip() {
    if [ ! -z "$OTTERSCALE_CONFIG_JUJU_IP" ]; then
        JUJU_VM_IP="$OTTERSCALE_CONFIG_JUJU_IP"
        log "INFO" "Using static IP from config: $JUJU_VM_IP" "VM IP"
    else
        log "INFO" "Prompting user for static VM IP" "VM IP"
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
    if [ ! -z "$OTTERSCALE_CONFIG_JUJU_IP" ]; then
        if [ "$OTTERSCALE_CONFIG_JUJU_IP" != "$CURRENT_JUJU_IP" ]; then
            JUJU_VM_IP="$OTTERSCALE_CONFIG_JUJU_IP"
            update_vm_ip
        fi
    else
        ask_user_type_vm_ip
    fi
}

check_vm_ip() {
    local JUJU_MACHINE_SUBNET_MODE=$(maas admin interfaces read $JUJU_MACHINE_ID | jq -r '.[].links' | jq -r '.[] | select(.subnet.name=="'"$MAAS_NETWORK_SUBNET"'") | .mode')
    log "INFO" "Machine $JUJU_MACHINE_ID interfaces mode is $JUJU_MACHINE_SUBNET_MODE" "Check VM IP"
    if [[ $JUJU_MACHINE_SUBNET_MODE != "static" ]]; then
        set_static_vm_ip
    else
        check_modify_vm_ip
    fi
}

update_vm_ip() {
    log "INFO" "Update $JUJU_VM_IP to juju-vm." "Update Juju-vm IP"
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

prepare_microk8s_config() {
    usermod -aG microk8s "$NON_ROOT_USER"
    log "INFO" "Added $NON_ROOT_USER to microk8s group" "MicroK8s config"

    KUBE_FOLDER="/home/$NON_ROOT_USER/.kube"
    if [ ! -d "$KUBE_FOLDER" ]; then
        mkdir -p "$KUBE_FOLDER"
    fi
    chown "$NON_ROOT_USER":"$NON_ROOT_USER" "$KUBE_FOLDER"

    log "INFO" "Update microk8s calico daemonset environment IP_AUTODETECTION_METHOD to $OTTERSCALE_BRIDGE_NAME" "Update microk8s calico"
    if [ -z $OTTERSCALE_BRIDGE_NAME ]; then
        error_exit "Bridge name is empty."
    fi
    if ! microk8s kubectl set env -n kube-system daemonset.apps/calico-node -c calico-node IP_AUTODETECTION_METHOD="interface=$OTTERSCALE_BRIDGE_NAME" >/dev/null 2>&1; then
        error_exit "Failed update microk8s calico env IP_AUTODETECTION_METHOD."
    fi
}

enable_microk8s_option() {
    if microk8s status --wait-ready >/dev/null 2>&1; then
        log "INFO" "microk8s is ready." "MicroK8S config"
        microk8s config > "$KUBE_FOLDER/config"
        chown "$NON_ROOT_USER":"$NON_ROOT_USER" "$KUBE_FOLDER/config"

        execute_cmd "microk8s enable dns" "enable microk8s dns"
        execute_cmd "microk8s enable hostpath-storage" "enable microk8s hostpath-storage"
        execute_cmd "microk8s enable metallb:$OTTERSCALE_INTERFACE_IP-$OTTERSCALE_INTERFACE_IP" "enable microk8s metallb"
    fi
}

extend_microk8s_cert() {
    log "INFO" "Refresh microk8s certificate to 3650 days" "MicroK8S certificate update"
    SNAP_SSL="/snap/microk8s/current/usr/bin/openssl"
    SNAP_DATA="/var/snap/microk8s/current"
    OPENSSL_CONF="/snap/microk8s/current/etc/ssl/openssl.cnf"

    log "INFO" "Generating server CSR for microk8s certificate" "MicroK8s cert"
    if ! ${SNAP}/usr/bin/openssl req -new -sha256 -key ${SNAP_DATA}/certs/server.key -out ${SNAP_DATA}/certs/server.csr -config ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out server.csr)."
    fi

    log "INFO" "Signing server certificate (3650 days)" "MicroK8s cert"
    if ! ${SNAP}/usr/bin/openssl x509 -req -sha256 -in ${SNAP_DATA}/certs/server.csr -CA ${SNAP_DATA}/certs/ca.crt -CAkey ${SNAP_DATA}/certs/ca.key -CAcreateserial -out ${SNAP_DATA}/certs/server.crt -days 3650 -extensions v3_ext -extfile ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out server.crt)."
    fi

    log "INFO" "Generating front‑proxy CSR" "MicroK8s cert"
    if ! ${SNAP}/usr/bin/openssl req -new -sha256 -key ${SNAP_DATA}/certs/front-proxy-client.key -out ${SNAP_DATA}/certs/front-proxy-client.csr -config <(sed '/^prompt = no/d' ${SNAP_DATA}/certs/csr.conf) -subj "/CN=front-proxy-client" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out front-proxy-client.csr)."
    fi

    log "INFO" "Signing front‑proxy certificate (3650 days)" "MicroK8s cert"
    if ! ${SNAP}/usr/bin/openssl x509 -req -sha256 -in ${SNAP_DATA}/certs/front-proxy-client.csr -CA ${SNAP_DATA}/certs/front-proxy-ca.crt -CAkey ${SNAP_DATA}/certs/front-proxy-ca.key -CAcreateserial -out ${SNAP_DATA}/certs/front-proxy-client.crt -days 3650 -extensions v3_ext -extfile ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate out front-proxy-client.crt)."
    fi

    log "INFO" "MicroK8s certificates refreshed successfully" "MicroK8s cert"
}

select_bridge() {
    while true; do
        echo "Available network bridges:"
        for i in "${!bridges[@]}"; do
            echo "$i) ${bridges[$i]}"
        done

        read -p "Select bridge (0-${#bridges[@]}): " choice
        case $choice in
            [0-9]*)
                if [ $choice -le ${#bridges[@]} ]; then
                    OTTERSCALE_BRIDGE_NAME=${bridges[$((choice-1))]}
                    log "INFO" "User selected bridge: $OTTERSCALE_BRIDGE_NAME" "Network"
                    return
                fi
                ;;
        esac
    done
}

prompt_bridge_creation() {
    read -p "No network bridge found, should be provisioned programmatically (Y/N)? " CONFIRM_CREATE
    if [[ "$CONFIRM_CREATE" == "Y" || "$CONFIRM_CREATE" == "y" ]]; then
        log "INFO" "User chose to auto‑create bridge" "Network"
        get_default_interface
        get_default_gateway
        get_default_cidr
        get_default_dns $CURRENT_INTERFACE
        backup_netplan
        create_netplan
        apply_netplan
    else
        log "INFO" "User asked to create bridge manually; exiting" "Network"
        exit 0
    fi
}

# ----------------------------------------------------------------
#  Function: check_bridge
#  Description: Ensure the network bridge exists, otherwise prompt/create it.
# ----------------------------------------------------------------
check_bridge() {
    OTTERSCALE_BRIDGE_NAME="br-otters"

    if ip link show "$OTTERSCALE_BRIDGE_NAME" &>/dev/null; then
        log "INFO" "Default bridge $OTTERSCALE_BRIDGE_NAME exist" "Check network"
    else
        log "INFO" "Bridge $OTTERSCALE_BRIDGE_NAME not found – will attempt to create it" "Network"

        bridges=($(brctl show 2>/dev/null | awk 'NR>1 {print $1}' | grep -v '^$'))
        if [ ${#bridges[@]} -eq 0 ]; then
            prompt_bridge_creation
        else
            select_bridge
        fi
    fi

    OTTERSCALE_CIDR=$(ip -o -4 addr show dev "$OTTERSCALE_BRIDGE_NAME" | awk '{print $4}')    
    OTTERSCALE_INTERFACE_IP=$(echo $OTTERSCALE_CIDR | cut -d'/' -f1)
    OTTERSCALE_INTERFACE_IP_MASK=$(echo $OTTERSCALE_CIDR | cut -d'/' -f2)
    get_default_dns $OTTERSCALE_BRIDGE_NAME
}

backup_netplan() {
    NETPLAN_FILE=$(find /etc/netplan/ -name *yaml | head -n 1)
    log "INFO" "Backing up current netplan file ($NETPLAN_FILE)" "OS network"
    cp $NETPLAN_FILE $NETPLAN_FILE-$(date '+%Y%m%d%H%M%S')
}

get_default_interface() {
    CURRENT_INTERFACE=$(ip route show default | awk '{print $5}')
    if [ -z $CURRENT_INTERFACE ]; then
        error_exit "Default route network interface is empty"
    fi
    log "INFO" "Detected default interface: $CURRENT_INTERFACE" "OS network"
}

get_default_gateway() {
    CURRENT_GATEWAY=$(ip route show default | awk '{print $3}')
    if [ -z $CURRENT_GATEWAY ]; then
        error_exit "Default gateway is empty"
    fi
    log "INFO" "Detected default gateway: $CURRENT_GATEWAY" "OS network"
}

get_default_cidr() {
    CURRENT_CIDR=$(ip -o -4 addr show dev "$CURRENT_INTERFACE" | awk '{print $4}')
    if [ -z $CURRENT_CIDR ]; then
        error_exit "Interface $CURRENT_INTERFACE CIDR is empty"
    fi
    log "INFO" "Detected CIDR for $CURRENT_INTERFACE: $CURRENT_CIDR" "OS network"
}

get_default_dns() {
    local INTERFACE=$1
    CURRENT_DNS=$(resolvectl -i $INTERFACE | grep "Current DNS Server" | awk '{print $4}' | paste -sd, -)
    if [ -z "$CURRENT_DNS" ]; then
        log "WARN" "No dns found for $INTERFACE, falling back to 8.8.8.8" "OS network"
        CURRENT_DNS="8.8.8.8"
    else
        log "INFO" "Detected DNS server(s) for $INTERFACE: $CURRENT_DNS" "OS network"
    fi
}

create_netplan() {
    cat > "$NETPLAN_FILE" <<EOF
network:
  version: 2
  renderer: networkd
  ethernets:
    $CURRENT_INTERFACE:
      dhcp4: no
      dhcp6: no
  bridges:
    br-otters:
      interfaces: [$CURRENT_INTERFACE]
      addresses: [$CURRENT_CIDR]
      routes:
      - to: default
        via: $CURRENT_GATEWAY
      nameservers:
        addresses: [$CURRENT_DNS]
EOF
    chmod 600 /etc/netplan/*.yaml
}

apply_netplan() {
    systemctl stop NetworkManager >/dev/null 2>&1
    systemctl disable NetworkManager >/dev/null 2>&1

    systemctl restart systemd-networkd >/dev/null 2>&1
    systemctl enable systemd-networkd >/dev/null 2>&1

    if netplan apply; then
        log "INFO" "Netplan applied successfully" "OS network"
	sleep 10
    else
        error_exit "Failed to apply netplan configuration"
    fi
}

get_snap_channel() {
    local SNAP_NAME=$1
    snap list | grep "^${SNAP_NAME}[[:space:]]" | awk '{print $4}'
}

retry_snap_install() {
    local SNAP_NAME="$1"
    local MAX_RETRIES="$2"
    local SNAP_OPTION="$3"
    local RETRIES=0

    while [ $RETRIES -lt $MAX_RETRIES ]; do
        log "INFO" "Installing snap $SNAP_NAME... (Attempt $((RETRIES)))" "Snap install"
        if snap install $SNAP_NAME $SNAP_OPTION >>"$TEMP_LOG" 2>&1; then
            break
        fi
        log "WARN" "Failed to install snap $SNAP_NAME. Retrying... (Attempt $((RETRIES)))" "Snap install"
        RETRIES=$((RETRIES+1))
    done

    if [ $RETRIES -eq $MAX_RETRIES ]; then
        error_exit "Failed to install snap $SNAP_NAME after $MAX_RETRIES attempts"
    fi
}

retry_snap_refresh() {
    local SNAP_NAME="$1"
    local SNAP_CHANNEL="$2"
    local MAX_RETRIES="$3"
    local RETRIES=0

    while [ $RETRIES -lt $MAX_RETRIES ]; do
        log "INFO" "Refreshing snap $SNAP_NAME to $MAX_RETRIES... (Attempt $((RETRIES)))" "Snap refresh"
        if snap refresh $SNAP_NAME --channel=$MAX_RETRIES >>"$TEMP_LOG" 2>&1; then
            break
        fi
        log "WARN" "Failed to refresh snap $snSNAP_NAMEap to $MAX_RETRIES. Retrying... (Attempt $((RETRIES)))" "Snap refresh"
        RETRIES=$((RETRIES+1))
    done

    if [ $RETRIES -eq $MAX_RETRIES ]; then
        error_exit "Failed to refresh snap $SNAP_NAME to $MAX_RETRIES after $MAX_RETRIES attempts"
    fi
}

install_or_update_snap() {
    local SNAP_NAME=$1
    local SNAP_CHANNEL=$2
    if snap list | grep -q "^${SNAP_NAME}[[:space:]]"; then
        if [[ $(get_snap_channel "$SNAP_NAME") != "$SNAP_CHANNEL" ]]; then
            retry_snap_refresh "$SNAP_NAME" "$SNAP_CHANNEL" "$OTTERSCALE_MAX_RETRIES"
        fi
    else
        if [[ $SNAP_NAME == "microk8s" ]]; then
            retry_snap_install "$SNAP_NAME" "$OTTERSCALE_MAX_RETRIES" "--classic --channel=$SNAP_CHANNEL"
        else
            retry_snap_install "$SNAP_NAME" "$OTTERSCALE_MAX_RETRIES" "--channel=$SNAP_CHANNEL"
        fi
    fi
}

snap_install() {
    declare -A SNAP_CHANNELS
    SNAP_CHANNELS[core24]=$CORE24_CHANNEL
    SNAP_CHANNELS[maas]=$MAAS_CHANNEL
    SNAP_CHANNELS[maas-test-db]=$MAAS_DB_CHANNEL
    SNAP_CHANNELS[juju]=$JUJU_CHANNEL
    SNAP_CHANNELS[lxd]=$LXD_CHANNEL
    SNAP_CHANNELS[microk8s]=$MICROK8S_CHANNEL

    for SNAP_NAME in $SNAP_PACKAGES; do
        install_or_update_snap "$SNAP_NAME" "${SNAP_CHANNELS[$SNAP_NAME]}"
    done

    log "INFO" "Holding all snaps..." "Snap hold"
    snap refresh --hold >>"$TEMP_LOG" 2>&1
}

disable_ipv6() {
    log "INFO" "Sysctl disable ipv6, it will resume after reboot" "OS config"
    sysctl -w net.ipv6.conf.all.disable_ipv6=1 >/dev/null 2>&1
    sysctl -w net.ipv6.conf.default.disable_ipv6=1 >/dev/null 2>&1
}

find_first_non_user() {
    local USER_HOME=""
    for USER in $(ls /home); do
        if [ -d "/home/$USER" ]; then
            USER_HOME="/home/$USER"
            break
        fi
    done

    if [ -z "$USER_HOME" ]; then
        error_exit "No non-root user found for SSH key setup"
    fi

    NON_ROOT_USER=$(basename "$USER_HOME")
    log "INFO" "Non-root user is $NON_ROOT_USER"
}

generate_ssh_key() {
    if [[ ! -f "/home/$NON_ROOT_USER/.ssh/id_rsa" ]]; then
        if ! su "$NON_ROOT_USER" -c 'mkdir -p $HOME/.ssh; ssh-keygen -q -t rsa -N "" -f "$HOME/.ssh/id_rsa" >>"$LOG" 2>&1'; then
            error_exit "SSH key generation failed"
        fi
    fi

    chown -R "$NON_ROOT_USER:$NON_ROOT_USER" "/home/$NON_ROOT_USER/.ssh"
    chmod 600 "/home/$NON_ROOT_USER/.ssh/id_rsa"
    chmod 644 "/home/$NON_ROOT_USER/.ssh/id_rsa.pub"
}

add_key_to_maas() {
    if [[ $(maas admin sshkeys read | jq -r 'length') -eq 0 ]]; then
        if ! maas admin sshkeys create key="$(cat "/home/$NON_ROOT_USER/.ssh/id_rsa.pub")" >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed to add SSH key to MAAS"
        fi
    fi
}

execute_as_user() {
    local USER=$1
    shift
    local CMD="$*"
    if ! su "$USER" -c "$CMD" >>"$TEMP_LOG" 2>&1; then
        log "WARN" "Failed to $2, check $LOG for details" "Non-root cmd"
        return 1
    fi
    return 0
}

execute_cmd() {
    local CMD=$1
    local MSG=$2
    log "INFO" "Execute command: $CMD" "$MSG"
    if ! $CMD >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed $MSG"
    else
        log "DEBUG" "Command succeeded: $CMD" "$MSG"
    fi
    return 0
}

check_root() {
    [ "$(id -u)" -eq 0 ] || error_exit "This script must be run as root"
}

check_os() {
    local OS_ID=$(lsb_release -si)
    if [ "$OS_ID" != "Ubuntu" ]; then
        error_exit "This script requires Ubuntu $OTTERSCALE_OS."
    fi
}

check_memory() {
    local OS_MEMORY_GB=$(free -g | awk '/Mem:/ {print $2}')
    if [ "$OS_MEMORY_GB" -lt "$MIN_MEMORY_GB" ]; then
        error_exit "Insufficient memory. Minimum required: ${MIN_MEMORY_GB}GB, Available: ${OS_MEMORY_GB}GB"
    fi
}

check_disk() {
    local OS_DISK_AVAILABLE_GB=$(df -BG / | awk 'NR==2 {print $4}' | tr -d 'G')
    if [ "$OS_DISK_AVAILABLE_GB" -lt "$MIN_DISK_GB" ]; then
        error_exit "Insufficient disk space. Minimum required: ${MIN_DISK_GB}GB, Available: ${OS_DISK_AVAILABLE_GB}GB"
    fi
}

config_modules() {
    local MODULE=rbd
    local MODULES_FILE="/etc/modules"
    if ! grep -q "^$MODULE$" "$MODULES_FILE"; then
        echo "$MODULE" >> "$MODULES_FILE"
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

validate_url() {
    local URL=$1
    local IP=$(echo "$URL" | awk -F '[/:]' '{print $4}')
    local PORT=$(echo "$URL" | awk -F '[/:]' '{print $5}')

    if ! validate_ip $IP; then
        log "WARN" "Invalid URL – IP part $IP is malformed" "Input validation"
        error_exit "Invalid IP format: $IP"
    fi

    if ! validate_port $PORT; then
    log "WARN" "Invalid URL – port part $PORT is out of range" "Input validation"
        error_exit "Invalid Port format: $PORT"
    fi
}

validate_ip() {
    local IP=$1
    if [[ ! $IP =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        return 1
    else
        return 0
    fi
}

validate_port() {
    local PORT=$1
    if [[ ! $PORT =~ ^[0-9]+$ ]]; then
        return 1
    fi

    if [[ "$PORT" -lt 1 || "$PORT" -gt 65535 ]]; then
        return 1
    fi
    return 0
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
    if is_ip_in_network $MAAS_DHCP_START_IP $network $mask_dotted; then
        if is_ip_in_network $MAAS_DHCP_END_IP $network $mask_dotted; then
            log "INFO" "IP range $MAAS_DHCP_START_IP to $MAAS_DHCP_END_IP is within the network $MAAS_NETWORK_SUBNET"
            return 0
        else
            log "WARN" "End IP $MAAS_DHCP_END_IP is not in the network $MAAS_NETWORK_SUBNET"
            return 1
        fi
    else
        log "WARN" "Start IP $MAAS_DHCP_START_IP is not in the network $MAAS_NETWORK_SUBNET"
        return 1
    fi
}

## without parameter
if [[ $# -eq 0 ]]; then
    while true; do
        read -p "Please enter otterscale endpoint (default is http://127.0.0.1:8299): " USER_INPUT_ENDPOINT
        OTTERSCALE_ENDPOINT=${USER_INPUT_ENDPOINT:-http://127.0.0.1:8299}
        if validate_url "$OTTERSCALE_ENDPOINT"; then
            break
        else
            echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] URL $OTTERSCALE_ENDPOINT is invalid, please try again" | tee -a $OTTERSCALE_INSTALL_DIR/setup.log
        fi
    done
fi

## with parameter
while [ $# -gt 0 ]; do
    case $1 in
        --url=* | url=*)
            OTTERSCALE_ENDPOINT="${1#*=}"
            if ! validate_url "$OTTERSCALE_ENDPOINT"; then
                exit 1
            fi
	    ;;
        --config=* | config=*)
            OTTERSCALE_CONFIG_PATH="${1#*=}"
            if [ ! -f $OTTERSCALE_CONFIG_PATH ]; then
                echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] Config file $OTTERSCALE_CONFIG_PATH not found, please try again" | tee -a $OTTERSCALE_INSTALL_DIR/setup.log
                exit 1
            fi
            source $OTTERSCALE_CONFIG_PATH
            ;;
        -h | --help | help)
            echo "Usage: sudo bash install.sh [options]"
            echo ""
            echo "Options:"
            echo "  -h | --help | help     Show this help message"
            echo "  --url=    | url=       Specific Otters endpoint"
            echo "  --config= | config=    Specific the configuration file to use"
            echo ""
            echo "Example"
            echo "  sudo bash install.sh"
            echo "  sudo bash install.sh url=http://IP:PORT"
            exit 0
            ;;
        *)
            echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] Invalid option: $1, please try again" | tee -a $OTTERSCALE_INSTALL_DIR/setup.log
            exit 1
            ;;
    esac
    shift
done

main() {
    log "INFO" "Starting Otterscale installation..." "Otterscale"
    log "INFO" "OTTERSCALE_ENDPOINT=$OTTERSCALE_ENDPOINT" "Otterscale"

    ##
    # Validate environment
    check_root
    check_os
    check_memory
    check_disk
    disable_ipv6
    log "INFO" "All pre‑checks passed, proceeding with installation." "Otterscale"

    ##
    # Package install
    apt_update
    apt_install "$APT_PACKAGES"
    snap_install

    ##
    # Network bridge check
    check_bridge

    ##
    # MAAS
    # 1. Init, create, and login
    # 2. Generate ssh
    # 3. Configure maas
    # 4. Lxd init
    init_maas
    create_maas_admin
    login_maas
    find_first_non_user
    generate_ssh_key
    add_key_to_maas
    update_maas_config
    download_maas_img
    enable_maas_dhcp
    init_lxd
    create_maas_lxd_project
    create_lxd_vm
    create_vm_from_maas
    #set_vm_static_ip

    ##
    # Juju
    # 1. Bootstrap
    # 2. Deploy microk8s
    # 3. Juju config
    bootstrap_juju
    prepare_microk8s_config
    enable_microk8s_option
    extend_microk8s_cert
    juju_add_k8s
    config_modules

    ##
    # Finished
    send_otterscale_config_data

    trap cleanup EXIT
    log "INFO" "Otterscale install finished" "Otterscale"
}

main "$@"
