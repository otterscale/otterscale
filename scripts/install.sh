#!/bin/bash

## Current directory
export OTTERSCALE_INSTALL_DIR=$(dirname "$(readlink -f $0)")

## Source env
for envfile in $OTTERSCALE_INSTALL_DIR/env/*.env; do
    source "$envfile"
done

## Import script
for file in $OTTERSCALE_INSTALL_DIR/utils/*.sh; do
    source "$file"
done

## LOG
export TEMP_LOG=$(mktemp)
export LOG=$OTTERSCALE_INSTALL_DIR/setup.log
touch $LOG
chmod 666 $LOG

main() {
    ## Validate
    validate_system

    ## Package install
    apt_update
    apt_install "$APT_PACKAGES"
    snap_install

    ## Host network
    select_bridge
    
    ## MAAS init and login
    init_maas
    create_maas_admin
    login_maas

    ## User ssh-key
    set_sshkey

    ## MAAS configure
    update_maas_dns
    update_maas_config
    download_maas_img
    enable_maas_dhcp

    ## Create LXD
    init_lxd
    create_maas_lxd_project
    create_lxd_vm
    create_vm_from_maas
    set_vm_static_ip

    ## Bootstrap
    set_juju_config
    bootstrap_juju

    ## Create default model
    create_scope

    ## Config microk8s
    update_microk8s_config
    enable_microk8s_option
    extend_microk8s_cert

    ## Add juju-k8s
    juju_add_k8s
    juju_config_k8s

    ## Create cluster token
    create_k8s_token

    ## Send config to otterscale
    send_config_data

    ## cleanup
    trap cleanup EXIT

    log "INFO" "Otterscale install finished" "Otterscale"
}

if [[ $# -eq 0 ]]; then
    error_exit "URL must be provided as a parameter"
fi
while [ $# -gt 0 ]; do
    case $1 in
        url=*)
            otterscale_url="${1#*=}"
            validate_url "$otterscale_url"
            ;;
        *)
            error_exit "Invalid option: $1"
            ;;
    esac
    shift
done

main "$@"
