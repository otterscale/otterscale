#!/bin/bash

ask_proxy() {
    while true; do
        read -p "Do you need to set HTTP/HTTPS proxy ? (Y/n): " confirm
        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            read -p "Enter HTTP proxy: " http_proxy
            read -p "Enter HTTPS proxy: " https_proxy
            add_proxy
            log "INFO" "Set http&https proxy finish."
            break
        elif [[ "$confirm" =~ ^[Nn]$ ]]; then
            break
        else
            echo "Invalid input. Please enter y or n."
        fi
    done
}

set_proxy() {
    local proxy_type="$1"
    local proxy_value="$2"
    local existing_proxy=""

    if [[ -f "$PROXY_FILE" ]]; then
        existing_proxy=$(grep -P 'Acquire::$proxy_type::Proxy\s+"[^"]*";' $PROXY_FILE | sed -E 's/Acquire::$proxy_type::Proxy\s+"([^"]*)";/\1/')
    fi

    if [[ "$proxy_value" != "$existing_proxy" ]]; then
        log "INFO" "Set $proxy_type proxy to $PROXY_FILE..."
        echo "Acquire::$proxy_type::Proxy \"$proxy_value\";" | tee -a "$PROXY_FILE" > /dev/null
    fi

    log "INFO" "Set snap system $proxy_type proxy..."
    snap set system proxy.$proxy_type="$proxy_value"
}

add_proxy() {
    if [ -n "$http_proxy" ]; then
        set_proxy "http" "$http_proxy"
    fi

    if [ -n "$https_proxy" ]; then
        set_proxy "https" "$https_proxy"
    fi
}

del_proxy() {
    rm -f $PROXY_FILE
    snap unset system proxy.http
    snap unset system proxy.https
}
