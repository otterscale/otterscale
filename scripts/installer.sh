#!/bin/bash

## Source env
for envfile in ./env/*.env; do
    source "$envfile"
done

## Import script
for file in ./scripts/*.sh; do
    source "$file"
done

## LOG
export INSTALLER_DIR=$(dirname "$(readlink -f $0)")
export LOG=$INSTALLER_DIR/setup.log
touch $LOG
chmod 666 $LOG

main() {
    ## Validate
    validate_system

    ## Proxy setting
    ask_proxy

    ## Package install
    install_packages
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

    ## OtterScale
    #ottersacle_start

    ## cleanup
    trap cleanup EXIT
}

main "$@"
