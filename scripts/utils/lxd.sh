#!/bin/bash

generate_lxd_config() {
    cat > $lxd_file <<EOF
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
}

# Enhanced LXD initialization
init_lxd() {
    lxd_file=$OTTERSCALE_INSTALL_DIR/lxd-config.yaml
    generate_lxd_config

    log "INFO" "Initializing LXD with bridge $OTTERSCALE_BRIDGE_NAME..."
    if ! cat $lxd_file | lxd init --preseed >>$TEMP_LOG 2>&1; then
        error_exit "LXD initialization failed."
    else
        log "INFO" "LXD initialized successfully"
        rm -f "$lxd_file"
    fi
}