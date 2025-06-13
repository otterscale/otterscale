#!/bin/bash

generate_clouds_yaml() {
    su "$username" -c 'cat > $JUJU_CLOUD <<EOF
clouds:
  maas-cloud:
    type: maas
    description: Metal As A Service
    auth-types: [oauth1]
    endpoint: http://$BRIDGE_IP:5240/MAAS/api/2.0/
    regions:
      default: {}
EOF'
}

generate_credentials_yaml() {
    su "$username" -c 'cat > $JUJU_CREDENTIAL <<EOF
credentials:
  maas-cloud:
    maas-cloud-credential:
      auth-type: oauth1
      maas-oauth: $APIKEY
EOF'
}

juju_clouds() {
    ## If cloud exist, try use update
    if su "$username" -c 'juju clouds 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Cloud already exists, updating configuration..."
        if ! execute_non_user_cmd "$username" "juju update-cloud maas-cloud -f $JUJU_CLOUD" "update cloud"; then
            log "WARN" "Failed to update cloud, removing and recreating..."
            execute_non_user_cmd "$username" "juju remove-cloud maas-cloud || true" "remove cloud"
            if ! execute_non_user_cmd "$username" "juju add-cloud maas-cloud $JUJU_CLOUD" "add cloud"; then
                error_exit "Failed to add Juju cloud."
            fi
        fi
    else
        if ! execute_non_user_cmd "$username" "juju add-cloud maas-cloud $JUJU_CLOUD" "add cloud"; then
            error_exit "Failed to add Juju cloud."
        fi
    fi
}

juju_credentials() {
    ## If credential exist, try use update
    if su "$username" -c 'juju credentials 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Credential already exists, updating..."
        if ! execute_non_user_cmd "$username" "juju update-credential maas-cloud maas-cloud-credential -f $JUJU_CREDENTIAL" "update credential"; then
            log "WARN" "Failed to update credential, removing and recreating..."
            execute_non_user_cmd "$username" "juju remove-credential maas-cloud maas-cloud-credential || true" "remove credential"
            if ! execute_non_user_cmd "$username" "juju add-credential maas-cloud -f $JUJU_CREDENTIAL" "add credential"; then
                error_exit "Failed to add Juju credentials."
            fi
        fi
    else
        if ! execute_non_user_cmd "$username" "juju add-credential maas-cloud -f $JUJU_CREDENTIAL" "add credential"; then
            error_exit "Failed to add Juju credentials."
        fi
    fi
}

set_juju_config() {
    su "$username" -c 'mkdir -p ~/.local/share/juju'
    su "$username" -c 'mkdir -p ~/ottersacle'

    export JUJU_CLOUD=/home/$username/ottersacle/cloud.yaml
    export JUJU_CREDENTIAL=/home/$username/ottersacle/credential.yaml
    export BRIDGE_IP=$BRIDGE_IP
    export APIKEY=$APIKEY

    log "INFO" "Configuring Juju clouds..."
    generate_clouds_yaml
    juju_clouds

    log "INFO" "Configuring Juju credentials..."
    generate_credentials_yaml
    juju_credentials

    log "INFO" "Juju configuration completed"
    rm -rf /home/$username/ottersacle
    unset JUJU_CLOUD
    unset JUJU_CREDENTIAL
    unset BRIDGE_IP
    unset APIKEY
}

# Juju bootstrap with validation
bootstrap_juju() {
    log "INFO" "Juju bootstrap, it will take few minutes..."
    local bootstrap_cmd="juju bootstrap maas-cloud maas-cloud-controller --bootstrap-base=$BASE_IMAGE"
    local bootstrap_config="--config default-base=$BASE_IMAGE --controller-charm-channel=$CONTROLLER_CHARM_CHANNEL"
    local bootstrap_machine="--to juju-vm"

    if ! maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q . >/dev/null 2>&1; then
        error_exit "Juju bootstrap failed, do not found juju-vm in MAAS."
    fi

    if ! execute_non_user_cmd "$username" "$bootstrap_cmd $bootstrap_config $bootstrap_machine --debug" "juju bootstrap"; then
        error_exit "Juju bootstrap failed."
    else
        log "INFO" "MAAS and Juju setup completed successfully!"
    fi
}

create_scope() {
    log "INFO" "Create juju default scope"
    execute_non_user_cmd "$username" "juju create-model default" "create default model"
}

cross_scope() {
    execute_non_user_cmd "$username" "juju offer grafana:grafana-dashboard grafana-dashboard" "JuJu offer grafana-dashboard"
    execute_non_user_cmd "$username" "juju offer prometheus:receive-remote-write prometheus-receive-remote-write" "JuJu offer prometheus-receive-remote-write"
}

# JuJu K8S
juju_add_k8s() {
    log "INFO" "Juju add-k8s to maas-cloud-controller"
    if ! execute_non_user_cmd "$username" "juju add-k8s cos-k8s --controller maas-cloud-controller --client --debug" "execute juju add-k8s"; then
        error_exit "Failed execute juju add-k8s"
    fi

    log "INFO" "Juju add-model cos"
    if ! execute_non_user_cmd "$username" "juju add-model cos cos-k8s --debug" "execute juju add-model"; then
        error_exit "Failed execute juju add-model"
    fi

    log "INFO" "Juju deploy cos-lite"
    if ! execute_non_user_cmd "$username" "juju deploy cos-lite --trust --debug" "juju deploy cos-lite"; then
        error_exit "Failed execute juju deploy cos-lite"
    fi
}

juju_config_k8s() {
    log "INFO" "Juju config prometheus"
    execute_non_user_cmd "$username" "juju config prometheus metrics_retention_time=180d --debug" "update metric retention time to 180 days"
    execute_non_user_cmd "$username" "juju config prometheus maximum_retention_size=60% --debug" "update max retention size to 60%"
}
