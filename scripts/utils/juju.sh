#!/bin/bash

execute_juju_command() {
    local username="$1"
    local command="$2"
    local description="$3"
    if ! su "$username" -c "$command >\"\$LOG\" 2>&1"; then
        log "WARN" "Failed to $description, check \$LOG for details."
	return 1
    fi
    return 0
}

create_juju_folder() {
    su "$username" -c 'mkdir -p ~/.local/share/juju'
}

generate_clouds_yaml() {
    su "$username" -c 'cat > $INSTALLER_DIR/cloud.yaml <<EOF
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

juju_clouds() {
    export BRIDGE_IP=$BRIDGE_IP
    generate_clouds_yaml
    unset BRIDGE_IP

    ## If cloud exist, try use update
    if su "$username" -c 'juju clouds 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Cloud already exists, updating configuration..."
        if ! execute_juju_command "$username" "juju update-cloud maas-cloud -f cloud.yaml" "update cloud"; then
            log "WARN" "Failed to update cloud, removing and recreating..."
            execute_juju_command "$username" "juju remove-cloud maas-cloud || true" "remove cloud"
            if ! execute_juju_command "$username" "juju add-cloud maas-cloud cloud.yaml" "add cloud"; then
                error_exit "Failed to add Juju cloud."
            fi
        fi
    else
        if ! execute_juju_command "$username" "juju add-cloud maas-cloud cloud.yaml" "add cloud"; then
            error_exit "Failed to add Juju cloud."
        fi
    fi
}

generate_credentials_yaml() {
    su "$username" -c 'cat > $INSTALLER_DIR/credential.yaml <<EOF
credentials:
  maas-cloud:
    maas-cloud-credential:
      auth-type: oauth1
      maas-oauth: $APIKEY
EOF'
}

juju_credentials() {
    export APIKEY=$APIKEY
    generate_credentials_yaml
    unset APIKEY

    ## If credential exist, try use update
    if su "$username" -c 'juju credentials 2>/dev/null | grep -q "^maas-cloud[[:space:]]"'; then
        log "WARN" "Credential already exists, updating..."
        if ! execute_juju_command "$username" "juju update-credential maas-cloud maas-cloud-credential -f credential.yaml" "update credential"; then
            log "WARN" "Failed to update credential, removing and recreating..."
            execute_juju_command "$username" "juju remove-credential maas-cloud maas-cloud-credential || true" "remove credential"
            if ! execute_juju_command "$username" "juju add-credential maas-cloud -f credential.yaml" "add credential"; then
                error_exit "Failed to add Juju credentials."
            fi
        fi
    else
        if ! execute_juju_command "$username" "juju add-credential maas-cloud -f credential.yaml" "add credential"; then
            error_exit "Failed to add Juju credentials."
        fi
    fi
}

set_juju_config() {
    create_juju_folder

    log "INFO" "Configuring Juju clouds..."
    juju_clouds

    log "INFO" "Configuring Juju credentialss..."
    juju_credentials

    log "INFO" "Juju configuration completed"
}

bootstrap_finished() {
    unset http_proxy
    unset https_proxy

    log "INFO" "MAAS and Juju setup completed successfully!"
}

# Juju bootstrap with validation
bootstrap_juju() {
    log "INFO" "Bootstrap Juju..."
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
        bootstrap_finished
    fi
}
