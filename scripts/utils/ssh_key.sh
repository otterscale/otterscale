#!/bin/bash

find_first_non_user() {
    local USER_HOME=""
    for USER in $(ls /home); do
        if [ -d "/home/$USER" ]; then
            USER_HOME="/home/$USER"
            break
        fi
    done

    if [ -z "$USER_HOME" ]; then
        error_exit "No non-root user found for SSH key setup"
    fi

    username=$(basename "$USER_HOME")
}

generate_ssh_key() {
    if [[ ! -f "/home/NON_ROOT_USER/.ssh/id_rsa" ]]; then
        if ! su "NON_ROOT_USER" -c 'mkdir -p $HOME/.ssh; ssh-keygen -q -t rsa -N "" -f "$HOME/.ssh/id_rsa" >>"$LOG" 2>&1'; then
            error_exit "SSH key generation failed"
        fi
    fi

    chown -R "NON_ROOT_USER:NON_ROOT_USER" "/home/NON_ROOT_USER/.ssh"
    chmod 600 "/home/NON_ROOT_USER/.ssh/id_rsa"
    chmod 644 "/home/NON_ROOT_USER/.ssh/id_rsa.pub"
}

add_key_to_maas() {
    if [[ $(maas admin sshkeys read | jq -r 'length') -eq 0 ]]; then
        if ! maas admin sshkeys create key="$(cat "/home/NON_ROOT_USER/.ssh/id_rsa.pub")" >>"$TEMP_LOG" 2>&1; then
            error_exit "Failed to add SSH key to MAAS"
        fi
    fi
}

set_sshkey() {
    log "INFO" "Setting up SSH keys..." "OS SSH"
    find_first_non_user
    generate_ssh_key
    add_key_to_maas
    log "INFO" "SSH keys setup completed" "OS SSH"
}
