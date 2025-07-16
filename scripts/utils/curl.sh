#!/bin/bash

send_request() {
    local URL_PATH=$1
    local DATA=$2

    local RESPONSE=$(curl -s --header "Content-Type: application/json" --data "$DATA" "$otterscale_url$URL_PATH")
    if [[ "$?" != 0 ]]; then
        echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] Failed execute curl request"
        trap cleanup EXIT
        exit 1
    fi
}

send_status_data() {
    local PHASE=$1
    local MESSAGE=$2
    local DATA=$(cat <<EOF
{
"phase": "$PHASE",
"message": "$MESSAGE"
}
EOF
)

    send_request "/otterscale.environment.v1.EnvironmentService/UpdateStatus" "$DATA"
}

send_config_data() {
    local OTTERSCALE_MAAS_ENDPOINT="http://$OTTERSCALE_INTERFACE_IP:5240/MAAS"
    local OTTERSCALE_MAAS_KEY=$(su "NON_ROOT_USER" -c "juju show-credentials maas-cloud maas-cloud-credential --show-secrets --client | grep maas-oauth | awk '{print \$2}'")
    local OTTERSCALE_CONTROLLER=$(su "NON_ROOT_USER" -c "juju controllers --format json | jq -r '.\"current-controller\"'")
    local OTTERSCALE_CONTROLLER_DETIAL=$(su "NON_ROOT_USER" -c "OTTERSCALE_CONTROLLER=\$(juju controllers --format json | jq -r '.\"current-controller\"'); juju show-controller \$OTTERSCALE_CONTROLLER --show-password --format=json")
    local OTTERSCALE_JUJU_ENDPOINTS=$(echo $OTTERSCALE_CONTROLLER_DETIAL | jq -r '."'"$OTTERSCALE_CONTROLLER"'"."details"."api-endpoints"' | tr '\n' ' ' | sed 's/ \+/ /g' | grep -v '^ *\[[0-9a-fA-F:]\+.*')
    local OTTERSCALE_JUJU_USERNAME=$(echo $OTTERSCALE_CONTROLLER_DETIAL | jq -r '."'"$OTTERSCALE_CONTROLLER"'"."account"."user"')
    local OTTERSCAKE_JUJU_PASSWORD=$(echo $OTTERSCALE_CONTROLLER_DETIAL | jq -r '."'"$OTTERSCALE_CONTROLLER"'"."account"."password"')
    local OTTERSCALE_JUJU_CACERT=$(echo $OTTERSCALE_CONTROLLER_DETIAL | jq -r '."'"$OTTERSCALE_CONTROLLER"'"."details"."ca-cert"')
    local OTTERSCALE_JUJU_CLOUD_NAME="maas-cloud"
    local OTTERSCALE_JUJU_REGION="default"
    local OTTERSCALE_K8S_TOKEN=$(microk8s kubectl describe secret -n kube-system otters-secret | grep -E '^token' | cut -f2 -d ':' |  tr -d ' ')
    local OTTERSCALE_K8S_ENDPOINT_JSON=$(microk8s kubectl get endpoints -o json | jq '.items[].subsets[]')
    local OTTERSCALE_K8S_ENDPOINT=$(echo $OTTERSCALE_K8S_ENDPOINT_JSON | jq -r '.ports[].name')"://"$(echo $OTTERSCALE_K8S_ENDPOINT_JSON | jq -r '.addresses[].ip')":"$(echo $OTTERSCALE_K8S_ENDPOINT_JSON | jq '.ports[].port')
    local DATA=$(cat <<EOF
{"maas_url": "$OTTERSCALE_MAAS_ENDPOINT",
"maas_key": "$OTTERSCALE_MAAS_KEY",
"maas_version": "$OTTERSCALE_MAAS_VERSION",
"juju_controller": "$OTTERSCALE_CONTROLLER",
"juju_controller_addresses": $OTTERSCALE_JUJU_ENDPOINTS,
"juju_username": "$OTTERSCALE_JUJU_USERNAME",
"juju_password": "$OTTERSCAKE_JUJU_PASSWORD",
"juju_ca_cert": $(echo "$OTTERSCALE_JUJU_CACERT" | jq -sRr '@json'),
"juju_cloud_name": "$OTTERSCALE_JUJU_CLOUD_NAME",
"juju_cloud_region": "$OTTERSCALE_JUJU_REGION",
"juju_charmhub_api_url": "$OTTERSCALE_CHARMHUB_URL",
"microk8s_token": "$OTTERSCALE_K8S_TOKEN",
"microk8s_host": "$OTTERSCALE_K8S_ENDPOINT"}
EOF
)

    send_request "/otterscale.environment.v1.EnvironmentService/UpdateConfig" "$DATA"
}
