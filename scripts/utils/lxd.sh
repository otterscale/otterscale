#!/bin/bash

generate_lxd_config() {
    cat > $lxd_file <<EOF
config:
  core.https_address: '[::]:8443'
  core.trust_password: password
storage_pools:
- config:
    size: $LXD_STORAGE_SIZE
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
      parent: $bridge
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
    lxd_file=$INSTALLER_DIR/lxd-config.yaml
    generate_lxd_config

    log "INFO" "Initializing LXD with bridge $bridge..."
    if ! cat $lxd_file | lxd init --preseed >>$TEMP_LOG 2>&1; then
        error_exit "LXD initialization failed."
    fi
    log "INFO" "LXD initialized successfully"

    rm -f "$lxd_file"
}
