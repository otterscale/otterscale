#!/bin/bash

update_microk8s_config() {
    kubefolder="/home/$username/.kube"

    ## Add user group
    usermod -aG microk8s "$username"

    ## Create folder
    if [ ! -d "$kubefolder" ]; then
        mkdir -p "$kubefolder"
    fi
    chown "$username":"$username" "$kubefolder"
}

enable_microk8s_option() {
    local IPADDR=$(ip -4 -j route get 2.2.2.2 | jq -r '.[] | .prefsrc')
    if microk8s status --wait-ready >/dev/null 2>&1; then
        log "INFO" "microk8s is ready."
        microk8s config > "$kubefolder/config"
        chown "$username":"$username" "$kubefolder/config"

        log "INFO" "Enable microk8s dns."
        microk8s enable dns >>"$TEMP_LOG" 2>&1
        log "INFO" "Enable microk8s hostpath-storage."
        microk8s enable hostpath-storage >>"$TEMP_LOG" 2>&1
        log "INFO" "Enable microk8s metallb."
        microk8s enable metallb:$IPADDR-$IPADDR >>"$TEMP_LOG" 2>&1
    fi
}

extend_microk8s_cert() {
    log "INFO" "Refresh microk8s certificate."
    local SNAP="/snap/microk8s/current"
    local SNAP_DATA="/var/snap/microk8s/current"
    local OPENSSL_CONF="/snap/microk8s/current/etc/ssl/openssl.cnf"

    if ! ${SNAP}/usr/bin/openssl req -new -sha256 -key ${SNAP_DATA}/certs/server.key -out ${SNAP_DATA}/certs/server.csr -config ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out server.csr)."
    fi

    if ! ${SNAP}/usr/bin/openssl x509 -req -sha256 -in ${SNAP_DATA}/certs/server.csr -CA ${SNAP_DATA}/certs/ca.crt -CAkey ${SNAP_DATA}/certs/ca.key -CAcreateserial -out ${SNAP_DATA}/certs/server.crt -days 3650 -extensions v3_ext -extfile ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out server.crt)."
    fi

    if ! ${SNAP}/usr/bin/openssl req -new -sha256 -key ${SNAP_DATA}/certs/front-proxy-client.key -out ${SNAP_DATA}/certs/front-proxy-client.csr -config <(sed '/^prompt = no/d' ${SNAP_DATA}/certs/csr.conf) -subj "/CN=front-proxy-client" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out front-proxy-client.csr)."
    fi

    if ! ${SNAP}/usr/bin/openssl x509 -req -sha256 -in ${SNAP_DATA}/certs/front-proxy-client.csr -CA ${SNAP_DATA}/certs/front-proxy-ca.crt -CAkey ${SNAP_DATA}/certs/front-proxy-ca.key -CAcreateserial -out ${SNAP_DATA}/certs/front-proxy-client.crt -days 3650 -extensions v3_ext -extfile ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate out front-proxy-client.crt)."
    fi
}
