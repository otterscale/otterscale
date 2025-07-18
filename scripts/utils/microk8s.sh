#!/bin/bash

prepare_microk8s_config() {
    usermod -aG microk8s "$NON_ROOT_USER"

    KUBE_FOLDER="/home/$NON_ROOT_USER/.kube"
    if [ ! -d "$KUBE_FOLDER" ]; then
        mkdir -p "$KUBE_FOLDER"
    fi
    chown "$NON_ROOT_USER":"$NON_ROOT_USER" "$KUBE_FOLDER"
}

update_microk8s_cni_interface() {
    log "INFO" "Update microk8s calico daemonset environment IP_AUTODETECTION_METHOD to $OTTERSCALE_BRIDGE_NAME" "MicroK8S config"
    if ! microk8s kubectl set env -n kube-system daemonset.apps/calico-node -c calico-node IP_AUTODETECTION_METHOD="interface=$OTTERSCALE_BRIDGE_NAME" >> "$TEMP_LOG" 2>&1; then
        error_exit "Failed update microk8s calico env IP_AUTODETECTION_METHOD"
    fi
}

enable_microk8s_option() {
    if microk8s status --wait-ready >/dev/null 2>&1; then
        log "INFO" "microk8s is ready." "MicroK8S config"
        microk8s config > "$KUBE_FOLDER/config"
        chown "$NON_ROOT_USER":"$NON_ROOT_USER" "$KUBE_FOLDER/config"

        log "INFO" "Enable microk8s dns" "MicroK8S config"
        if ! microk8s enable dns >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed enable microk8s dns"
        fi
        log "INFO" "Enable microk8s hostpath-storage" "MicroK8S config"
	if ! microk8s enable hostpath-storage >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed enable microk8s hostpath-storage"
	fi

        log "INFO" "Enable microk8s metallb" "MicroK8S config"
        if ! microk8s enable metallb:$OTTERSCALE_INTERFACE_IP-$OTTERSCALE_INTERFACE_IP >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed enable microk8s metallb"
        fi
    fi
}

extend_microk8s_cert() {
    log "INFO" "Refresh microk8s certificate to 3650 days" "MicroK8S certificate update"
    local SNAP="/snap/microk8s/current"
    local SNAP_DATA="/var/snap/microk8s/current"
    local OPENSSL_CONF="/snap/microk8s/current/etc/ssl/openssl.cnf"

    if ! ${SNAP}/usr/bin/openssl req -new -sha256 -key ${SNAP_DATA}/certs/server.key -out ${SNAP_DATA}/certs/server.csr -config ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out server.csr)"
    fi

    if ! ${SNAP}/usr/bin/openssl x509 -req -sha256 -in ${SNAP_DATA}/certs/server.csr -CA ${SNAP_DATA}/certs/ca.crt -CAkey ${SNAP_DATA}/certs/ca.key -CAcreateserial -out ${SNAP_DATA}/certs/server.crt -days 3650 -extensions v3_ext -extfile ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out server.crt)"
    fi

    if ! ${SNAP}/usr/bin/openssl req -new -sha256 -key ${SNAP_DATA}/certs/front-proxy-client.key -out ${SNAP_DATA}/certs/front-proxy-client.csr -config <(sed '/^prompt = no/d' ${SNAP_DATA}/certs/csr.conf) -subj "/CN=front-proxy-client" >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate (out front-proxy-client.csr)"
    fi

    if ! ${SNAP}/usr/bin/openssl x509 -req -sha256 -in ${SNAP_DATA}/certs/front-proxy-client.csr -CA ${SNAP_DATA}/certs/front-proxy-ca.crt -CAkey ${SNAP_DATA}/certs/front-proxy-ca.key -CAcreateserial -out ${SNAP_DATA}/certs/front-proxy-client.crt -days 3650 -extensions v3_ext -extfile ${SNAP_DATA}/certs/csr.conf >>"$TEMP_LOG" 2>&1; then
        error_exit "Failed extend microk8s certificate out front-proxy-client.crt)"
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
        log "INFO" "Success apply $YAML_FILE" "MicroK8S config"
        rm $YAML_FILE
    else
        error_exit "Failed microk8s kubectl apply $YAML_FILE"
    fi
}

create_k8s_token() {
    if ! microk8s kubectl get secret otters-secret -n kube-system >/dev/null 2>&1; then
        export SA_PATH=$OTTERSCALE_INSTALL_DIR/otters_sa.yaml
        export RBAC_PATH=$OTTERSCALE_INSTALL_DIR/otters_rbac.yaml
        export SECRET_PATH=$OTTERSCALE_INSTALL_DIR/otter_secret.yaml

        log "INFO" "Gererate service account" "MicroK8S create token"
        generate_sa_yaml
        apply_yaml $SA_PATH

        log "INFO" "Gererate cluster role binding" "MicroK8S create token"
        generate_rbac_yaml
        apply_yaml $RBAC_PATH

        log "INFO" "Gererate secret" "MicroK8S create token"
        generate_secret_yaml
        apply_yaml $SECRET_PATH

        unset SA_PATH
        unset RBAC_PATH
        unset SECRET_PATH
    fi
}

check_microk8s() {
    prepare_microk8s_config
    update_microk8s_cni_interface
    enable_microk8s_option
    extend_microk8s_cert
}
