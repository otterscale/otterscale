#!/bin/bash

update_microk8s_config() {
    kubefolder="/home/$username/.kube"

    ## Add user group
    log "INFO" "Add $username to group microk8s"
    usermod -aG microk8s "$username"

    ## Create folder
    if [ ! -d "$kubefolder" ]; then
        mkdir -p "$kubefolder"
    fi
    chown "$username":"$username" "$kubefolder"

    ## Update calico-node env
    log "INFO" "Update microk8s calico daemonset environment IP_AUTODETECTION_METHOD to $bridge_name"
    if ! microk8s kubectl set env -n kube-system daemonset.apps/calico-node -c calico-node IP_AUTODETECTION_METHOD="interface=$bridge_name"; then
        error_exit "Failed update microk8s calico env IP_AUTODETECTION_METHOD."
    fi
}

enable_microk8s_option() {
    local IPADDR=$(ip -4 -j route get 2.2.2.2 | jq -r '.[] | .prefsrc')
    if microk8s status --wait-ready >/dev/null 2>&1; then
        log "INFO" "microk8s is ready."
        microk8s config > "$kubefolder/config"
        chown "$username":"$username" "$kubefolder/config"

        log "INFO" "Enable microk8s dns"
        microk8s enable dns >>"$TEMP_LOG" 2>&1;
        log "INFO" "Enable microk8s hostpath-storage"
        microk8s enable hostpath-storage >>"$TEMP_LOG" 2>&1;
        log "INFO" "Enable microk8s metallb"
	microk8s enable metallb:$IPADDR-$IPADDR >>"$TEMP_LOG" 2>&1;
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

generate_sa_yaml() {
    cat > $SA_PATH <<EOF
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: otters-sa
  namespace: kube-system
EOF
}

generate_rbac_yaml() {
    cat > $RBAC_PATH <<EOF
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: otters
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: otters-sa
  namespace: kube-system
EOF
}

generate_secret_yaml() {
    cat > $SECRET_PATH <<EOF
---
apiVersion: v1
kind: Secret
metadata:
  name: otters-secret
  namespace: kube-system
  annotations:
    kubernetes.io/service-account.name: otters-sa
type: kubernetes.io/service-account-token
EOF
}

apply_yaml() {
    local YAML_FILE=$1
    if microk8s kubectl apply -f $YAML_FILE >/dev/null 2>&1; then
        log "INFO" "Success apply $YAML_FILE"
        rm $YAML_FILE
    else
        error_exit "Failed microk8s kubectl apply $YAML_FILE"
    fi
}

create_k8s_token() {
    export SA_PATH=$INSTALLER_DIR/otters_sa.yaml
    export RBAC_PATH=$INSTALLER_DIR/otters_rbac.yaml
    export SECRET_PATH=$INSTALLER_DIR/otter_secret.yaml

    log "INFO" "Gererate service account"
    generate_sa_yaml
    apply_yaml $SA_PATH

    log "INFO" "Gererate cluster role binding"
    generate_rbac_yaml
    apply_yaml $RBAC_PATH

    log "INFO" "Gererate secret"
    generate_secret_yaml
    apply_yaml $SECRET_PATH

    unset SA_PATH
    unset RBAC_PATH
    unset SECRET_PATH
}
