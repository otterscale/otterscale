#!/bin/bash

export INSTALLER_DIR=$(dirname "$(readlink -f $0)")

## Source env
for envfile in $INSTALLER_DIR/env/*.env; do
    source "$envfile"
done

## Import script
for file in $INSTALLER_DIR/utils/*.sh; do
    source "$file"
done

## LOG
export TEMP_LOG=$(mktemp)
export LOG=$INSTALLER_DIR/setup.log
touch $LOG
chmod 666 $LOG

main() {
    ## Validate
    validate_system

    ## Proxy setting
    ask_proxy

    ## Package install
    apt_update
    apt_install $APT_PACKAGES
    install_snaps

    ## Host network
    select_bridge
    
    ## MAAS init and login
    init_maas
    create_maas_admin
    login_maas

    ## User ssh-key
    set_sshkey

    ## MAAS configure
    update_dns
    update_img_autosync
    update_proxy
    download_maas_img
    enable_maas_dhcp

    ## Create JuJu VM
    init_lxd
    create_maas_lxd_project
    create_lxd_vm
    create_vm_from_maas
    set_vm_static_ip

    ## JuJu
    set_juju_config
    bootstrap_juju

    ## Create default model
    create_scope

    ## Install otterscale
    local deb_file=$(ls $INSTALLER_DIR/packages/ | grep deb | head -n 1)
    apt_install $INSTALLER_DIR/packages/$deb_file
    start_service "ottersacle"
    enable_service "ottersacle"
    log "INFO" "OtterScale has been launched, you can access via http://localhost:5090"

    ## cleanup
    trap cleanup EXIT
}

main "$@"
