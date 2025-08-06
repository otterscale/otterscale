#!/bin/bash

juju_cmd() {
    local CMD=$1
    local MSG=$2
    log "INFO" "Execute command: $CMD" "$MSG"
    if ! execute_non_user_cmd "$NON_ROOT_USER" "$CMD" "$MSG"; then
        error_exit "Failed $MSG"
    fi
}

generate_clouds_yaml() {
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

generate_credentials_yaml() {
    su "$NON_ROOT_USER" -c 'cat > $JUJU_CREDENTIAL <<EOF
credentials:
  maas-cloud:
    maas-cloud-credential:
      auth-type: oauth1
      maas-oauth: $APIKEY
EOF'
}

juju_clouds() {
    log "INFO" "Configuring Juju clouds..." "JuJu clouds"
    generate_clouds_yaml

    if su "$NON_ROOT_USER" -c 'juju clouds 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "JuJu cloud maas-cloud already exists, skipping created..." "JuJu clouds"
    else
        juju_cmd "juju add-cloud maas-cloud $JUJU_CLOUD" "add juju cloud"
    fi
}

juju_credentials() {
    log "INFO" "Configuring Juju credentials..." "JuJu credentials"
    generate_credentials_yaml

    if su "$NON_ROOT_USER" -c 'juju credentials 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "JuJu Credential maas-cloud already exists, skipping created..." "JuJu credentials"
    else
        juju_cmd "juju add-credential maas-cloud -f $JUJU_CREDENTIAL" "add credential"
    fi
}

is_machine_exist() {
    if maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q . >/dev/null 2>&1; then
        return 0
    fi
    return 1
}

is_machine_deployed() {
    if [ $(maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | jq -r '.status_name') == Deployed ]; then
	return 0
    fi
    return 1
}

# Juju bootstrap with validation
bootstrap_juju() {
    su "$NON_ROOT_USER" -c 'mkdir -p ~/.local/share/juju'
    su "$NON_ROOT_USER" -c 'mkdir -p ~/otterscale'

    export JUJU_CLOUD=/home/$NON_ROOT_USER/otterscale/cloud.yaml
    export JUJU_CREDENTIAL=/home/$NON_ROOT_USER/otterscale/credential.yaml
    export OTTERSCALE_INTERFACE_IP=$OTTERSCALE_INTERFACE_IP
    export APIKEY=$APIKEY

    juju_clouds
    juju_credentials

    rm -rf /home/$NON_ROOT_USER/otterscale
    unset JUJU_CLOUD
    unset JUJU_CREDENTIAL
    unset APIKEY
	
    bootstrap_cmd="juju bootstrap maas-cloud maas-cloud-controller --bootstrap-base=$OTTERSCALE_BASE_IMAGE"
    bootstrap_config="--config default-base=$BASE_IMAGE --controller-charm-channel=$CONTROLLER_CHARM_CHANNEL"
    bootstrap_machine="--to juju-vm"

    if ! is_machine_exist; then
        error_exit "Juju bootstrap failed, do not found juju-vm in MAAS"
    fi

    if is_machine_deployed; then
        log "INFO" "JuJu had already bootstrap, skipping..."
    else
        log "INFO" "Juju bootstrap, it will take a few minutes..." "JuJu bootstrap"
        juju_cmd "$bootstrap_cmd $bootstrap_config $bootstrap_machine --debug" "juju bootstrap"
        log "INFO" "MAAS and Juju setup completed successfully!" "Finished bootstrap"
    fi
}

juju_add_k8s() {
    if execute_non_user_cmd "$NON_ROOT_USER" "juju show-cloud cos-k8s > /dev/null 2>&1" "check juju cloud if cos-k8s exist"; then
        log "INFO" "cos-k8s already exist, skipping..." "JuJu cloud"
    else
        juju_cmd "juju add-k8s cos-k8s --controller maas-cloud-controller --client --debug" "execute juju add-k8s"
    fi

    if execute_non_user_cmd "$NON_ROOT_USER" "juju show-model cos > /dev/null 2>&1" "check juju model if cos exist"; then
        log "INFO" "cos model already exist, skipping..." "JuJu model"
    else
        juju_cmd "juju add-model cos cos-k8s --debug" "execute juju add-model"
        juju_cmd "juju deploy cos-lite --trust --debug" "juju deploy cos-lite"
    fi

    juju_cmd "juju config prometheus metrics_retention_time=180d --debug" "update metric retention time to 180 days"
    juju_cmd "juju config prometheus maximum_retention_size=60% --debug" "update max retention size to 60%"
    juju_cmd "juju offer grafana:grafana-dashboard global-grafana" "offer grafana-dashboard"
    juju_cmd "juju offer prometheus:receive-remote-write global-prometheus" "offer prometheus-receive-remote-write"
}
