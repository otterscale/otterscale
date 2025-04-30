#!/bin/bash

get_juju_account() {
    local maas_controller=$(juju show-controller --show-password maas-controller --format=json)
    local controller_user=$(echo $maas_controller | jq -r '."maas-controller"."account"."user"')
    local controller_password=$(echo $maas_controller | jq -r '."maas-controller"."account"."password"')
    local controller_ca_cert=$(echo $maas_controller | jq -r '."maas-controller"."details"."ca-cert"')
}

get_maas_oath() {
    local maas_oath=$(juju show-credentials maas-one anyuser --show-secrets --client | grep maas-oauth | awk '{print $2}')
}

ottersacle_start() {
    get_juju_account 
    get_maas_oath

    ## Start service
}
