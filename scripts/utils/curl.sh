#!/bin/bash

send_request() {
    local url_path=$1
    local data=$2

    response=$(curl -s --header "Content-Type: application/json" --data "$data" "$otterscale_url$url_path")
    if [[ "$?" != 0 ]]; then
        trap cleanup EXIT
        exit 1
    fi
}

send_statue_data() {
    local phase=$1
    local message=$2

    local data=$(cat <<EOF
{
"phase": "$phase",
"message": "$message"
}
EOF
)

    send_request "/otterscale.environment.v1.EnvironmentService/UpdateStatus" "$data"
}

send_config_data() {
    ## MAAS
    local maas_url="http://$current_ip:5240/MAAS"
    local maas_key=$(su "$username" -c "juju show-credentials --show-secrets --client | grep maas-oauth | awk '{print \$2}'")

    ## Get juju config
    local current_controller=$(su "$username" -c "juju controllers --format json | jq -r '.\"current-controller\"'")
    local show_controller=$(su "$username" -c "current_controller=\$(juju controllers --format json | jq -r '.\"current-controller\"'); juju show-controller \$current_controller --show-password --format=json")

    ## Juju
    local juju_controller_addresses=$(echo $show_controller | jq -r '."'"$current_controller"'"."details"."api-endpoints"' | tr '\n' ' ' | sed 's/ \+/ /g' | grep -v '^ *\[[0-9a-fA-F:]\+.*')
    local juju_username=$(echo $show_controller | jq -r '."'"$current_controller"'"."account"."user"')
    local juju_password=$(echo $show_controller | jq -r '."'"$current_controller"'"."account"."password"')
    local juju_ca_cert=$(echo $show_controller | jq -r '."'"$current_controller"'"."details"."ca-cert"')
    local juju_cloud_name="maas-cloud"
    local juju_cloud_region="default"

    local data=$(cat <<EOF
{"maas_url": "$maas_url",
"maas_key": "$maas_key",
"maas_version": "$MAAS_API_VERSION",
"juju_controller_addresses": $juju_controller_addresses,
"juju_username": "$juju_username",
"juju_password": "$juju_password",
"juju_ca_cert": $(echo "$juju_ca_cert" | jq -sRr '@json'),
"juju_cloud_name": "$juju_cloud_name",
"juju_cloud_region": "$juju_cloud_region",
"juju_charmhub_api_url": "$JUJU_CHARMHUB_API_URL"}
EOF
)

    send_request "/otterscale.environment.v1.EnvironmentService/UpdateConfig" "$data"
}
