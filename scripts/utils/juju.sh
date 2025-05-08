#!/bin/bash

execute_juju_command() {
    local username="$1"
    local command="$2"
    local description="$3"
    if ! su "$username" -c "${command} >$LOG 2>&1"; then
        log "WARN" "Failed to $description, check $LOG for details."
	return 1
    fi
    return 0
}

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
        if ! execute_juju_command "$username" "juju update-cloud maas-cloud -f $JUJU_CLOUD" "update cloud"; then
            log "WARN" "Failed to update cloud, removing and recreating..."
            execute_juju_command "$username" "juju remove-cloud maas-cloud || true" "remove cloud"
            if ! execute_juju_command "$username" "juju add-cloud maas-cloud $JUJU_CLOUD" "add cloud"; then
                error_exit "Failed to add Juju cloud."
            fi
        fi
    else
        if ! execute_juju_command "$username" "juju add-cloud maas-cloud $JUJU_CLOUD" "add cloud"; then
            error_exit "Failed to add Juju cloud."
        fi
    fi
}

juju_credentials() {
    ## If credential exist, try use update
    if su "$username" -c 'juju credentials 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Credential already exists, updating..."
        if ! execute_juju_command "$username" "juju update-credential maas-cloud maas-cloud-credential -f $JUJU_CREDENTIAL" "update credential"; then
            log "WARN" "Failed to update credential, removing and recreating..."
            execute_juju_command "$username" "juju remove-credential maas-cloud maas-cloud-credential || true" "remove credential"
            if ! execute_juju_command "$username" "juju add-credential maas-cloud -f $JUJU_CREDENTIAL" "add credential"; then
                error_exit "Failed to add Juju credentials."
            fi
        fi
    else
        if ! execute_juju_command "$username" "juju add-credential maas-cloud -f $JUJU_CREDENTIAL" "add credential"; then
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
    local bootstrap_cmd="juju bootstrap maas-cloud maas-cloud-controller --bootstrap-base=ubuntu@22.04"
    local bootstrap_config="--config default-base=ubuntu@22.04"
    local bootstrap_machine="--to juju-vm"

    if ! maas admin machines read | jq -r '.[] | select(.hostname=="juju-vm")' | grep -q . >/dev/null 2>&1; then
        error_exit "Juju bootstrap failed, do not found juju-vm in MAAS."
    fi

    if [ -n "$http_proxy" ] && [ -n "$https_proxy" ]; then
	export http_proxy=$http_proxy
	export https_proxy=$https_proxy
	bootstrap_config="$bootstrap_config --config http-proxy=$http_proxy --config snap-http-proxy=$http_proxy"
	bootstrap_config="$bootstrap_config --config https-proxy=$https_proxy --config snap-https-proxy=$https_proxy"
    fi

    if ! execute_juju_command "$username" "$bootstrap_cmd $bootstrap_config $bootstrap_machine --debug" "juju bootstrap"; then
        error_exit "Juju bootstrap with proxy failed."
    else
        unset http_proxy
        unset https_proxy
        log "INFO" "MAAS and Juju setup completed successfully!"
    fi
}

create_scope() {
    if ! execute_juju_command "$username" "juju models --format=json | jq '.\"models\" | select(.[].\"short-name\"==\"default\")'" "check default model"; then
        execute_juju_command "$username" "juju create-model default" "create default model"
    fi
}
