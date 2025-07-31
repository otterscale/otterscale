#!/bin/bash

## Current directory
export OTTERSCALE_INSTALL_DIR=$(dirname "$(readlink -f $0)")

## Source env
source $OTTERSCALE_INSTALL_DIR/static.env

## Import script
for file in $OTTERSCALE_INSTALL_DIR/utils/*.sh; do
    source "$file"
done

## LOG
export TEMP_LOG=$(mktemp)
export LOG=$OTTERSCALE_INSTALL_DIR/setup.log
touch $LOG
chmod 666 $LOG

## without parameter
if [[ $# -eq 0 ]]; then
    while true; do
        read -p "Please enter otterscale endpoint (default is http://127.0.0.1:8299): " USER_INPUT_ENDPOINT
        OTTERSCALE_ENDPOINT=${USER_INPUT_ENDPOINT:-http://127.0.0.1:8299}
        if validate_url "$OTTERSCALE_ENDPOINT"; then
            break
        else
            echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] URL $OTTERSCALE_ENDPOINT is invalid, please try again" | tee -a $OTTERSCALE_INSTALL_DIR/setup.log
        fi
    done
fi

## with parameter
while [ $# -gt 0 ]; do
    case $1 in
        --config=* | config=*)
            OTTERSCALE_CONFIG_PATH="${1#*=}"
            if [ ! -f $OTTERSCALE_CONFIG_PATH ]; then
                echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] Config file $OTTERSCALE_CONFIG_PATH not found, please try again" | tee -a $OTTERSCALE_INSTALL_DIR/setup.log
                exit 1
            fi

            source $OTTERSCALE_CONFIG_PATH
            if ! validate_url "$OTTERSCALE_ENDPOINT"; then
                exit 1
            fi
            ;;
        -h | --help | help)
            echo "Usage: sudo bash install.sh [options]"
            echo ""
            echo "Options:"
            echo "  -h | --help | help     Show this help message"
            echo "  --config= | config=    Specific the configuration file to use"
            echo ""
            echo "Example"
            echo "  sudo bash install.sh"
            echo "  sudo bash install.sh config=FILEPATH"
            exit 0
            ;;
        *)
            echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] Invalid option: $1, please try again" | tee -a $OTTERSCALE_INSTALL_DIR/setup.log
            exit 1
            ;;
    esac
    shift
done

main() {
    validate_system
    apt_update
    apt_install "$APT_PACKAGES"
    snap_install
    config_modules

    check_bridge
    check_maas
    bootstrap_juju
    check_microk8s
    juju_add_k8s
    send_otterscale_config_data

    trap cleanup EXIT
    log "INFO" "Otterscale install finished" "Otterscale"
}

main "$@"
