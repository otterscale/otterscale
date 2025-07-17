#!/bin/bash

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
    ## If cloud exist, try use update
    if su "$NON_ROOT_USER" -c 'juju clouds 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "JuJu cloud maas-cloud already exists, skipping created..." "JuJu clouds"
    else
        if ! execute_non_user_cmd "$NON_ROOT_USER" "juju add-cloud maas-cloud $JUJU_CLOUD" "add cloud"; then
            error_exit "Failed to add Juju cloud"
        fi
    fi
}

juju_credentials() {
    if su "$NON_ROOT_USER" -c 'juju credentials 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "JuJu Credential maas-cloud already exists, skipping created..." "JuJu credentials"
    else
        if ! execute_non_user_cmd "$NON_ROOT_USER" "juju add-credential maas-cloud -f $JUJU_CREDENTIAL" "add credential"; then
            error_exit "Failed to add Juju credentials"
        fi
    fi
}

set_juju_config() {
    su "$NON_ROOT_USER" -c 'mkdir -p ~/.local/share/juju'
    su "$NON_ROOT_USER" -c 'mkdir -p ~/ottersacle'

    export JUJU_CLOUD=/home/$NON_ROOT_USER/ottersacle/cloud.yaml
    export JUJU_CREDENTIAL=/home/$NON_ROOT_USER/ottersacle/credential.yaml
    export OTTERSCALE_INTERFACE_IP=$OTTERSCALE_INTERFACE_IP
    export APIKEY=$APIKEY

    log "INFO" "Configuring Juju clouds..." "JuJu clouds"
    generate_clouds_yaml
    juju_clouds

    log "INFO" "Configuring Juju credentials..." "JuJu credentials"
    generate_credentials_yaml
    juju_credentials

    log "INFO" "Juju configuration completed" "JuJu prepare"
    rm -rf /home/$NON_ROOT_USER/ottersacle
    unset JUJU_CLOUD
    unset JUJU_CREDENTIAL
    unset APIKEY
}

# Juju bootstrap with validation
bootstrap_juju() {
    log "INFO" "Juju bootstrap, it will take a few minutes..." "JuJu bootstrap"
    local bootstrap_cmd="juju bootstrap maas-cloud maas-cloud-controller --bootstrap-base=$OTTERSCALE_BASE_IMAGE"
    local bootstrap_config="--config default-base=$BASE_IMAGE --controller-charm-channel=$CONTROLLER_CHARM_CHANNEL"
    local bootstrap_machine="--to juju-vm"

    if ! maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q . >/dev/null 2>&1; then
        error_exit "Juju bootstrap failed, do not found juju-vm in MAAS"
    fi

    if [ $(maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | jq -r '.status_name') == Deployed ]; then
        log "INFO" "Already juju bootstarp, skipping..."
    else
        if ! execute_non_user_cmd "$NON_ROOT_USER" "$bootstrap_cmd $bootstrap_config $bootstrap_machine --debug" "juju bootstrap"; then
            error_exit "Juju bootstrap failed"
        else
            log "INFO" "MAAS and Juju setup completed successfully!" "Finished bootstrap"
        fi
    fi
}

create_scope() {
    log "INFO" "Create juju default scope" "JuJu scope"
    execute_non_user_cmd "$NON_ROOT_USER" "juju add-model default" "create default model"
}

cross_scope() {
    if ! execute_non_user_cmd "$NON_ROOT_USER" "juju offer grafana:grafana-dashboard grafana-dashboard" "JuJu offer grafana-dashboard"; then
        error_exit "Failed execute juju offer grafana-dashboard"
    fi

    if ! execute_non_user_cmd "$NON_ROOT_USER" "juju offer prometheus:receive-remote-write prometheus-receive-remote-write" "JuJu offer prometheus-receive-remote-write"; then
        error_exit "Failed execute juju offer prometheus-receive-remote-write"
    fi
}

# JuJu K8S
juju_add_k8s() {
    if execute_non_user_cmd "$NON_ROOT_USER" "juju show-cloud cos-k8s | grep -q ." "check juju cloud"; then
        log "INFO" "cos-k8s already exist, skipping..."
    else
        log "INFO" "Juju add-k8s to maas-cloud-controller" "JuJu K8S"
        if ! execute_non_user_cmd "$NON_ROOT_USER" "juju add-k8s cos-k8s --controller maas-cloud-controller --client --debug" "execute juju add-k8s"; then
            error_exit "Failed execute juju add-k8s"
        fi
    fi

    if execute_non_user_cmd "$NON_ROOT_USER" "juju show-model cos | grep -q ." "check juju model"; then
        log "INFO" "cos model already exist, skipping..."
    else
        log "INFO" "Juju add-model cos" "JuJu K8S"
        if ! execute_non_user_cmd "$NON_ROOT_USER" "juju add-model cos cos-k8s --debug" "execute juju add-model"; then
            error_exit "Failed execute juju add-model"
        fi

        log "INFO" "Juju deploy cos-lite" "JuJu K8S"
        if ! execute_non_user_cmd "$NON_ROOT_USER" "juju deploy cos-lite --trust --debug" "juju deploy cos-lite"; then
            error_exit "Failed execute juju deploy cos-lite"
        fi
        juju_config_k8s	
    fi
}

juju_config_k8s() {
    log "INFO" "Juju config prometheus" "JuJu config"
    execute_non_user_cmd "$NON_ROOT_USER" "juju config prometheus metrics_retention_time=180d --debug" "update metric retention time to 180 days"
    execute_non_user_cmd "$NON_ROOT_USER" "juju config prometheus maximum_retention_size=60% --debug" "update max retention size to 60%"
}
