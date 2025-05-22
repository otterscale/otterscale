#!/bin/bash

find_first_non_user() {
    local user_home=""
    for user in $(ls /home); do
        if [ -d "/home/$user" ]; then
            user_home="/home/$user"
            break
        fi
    done

    if [ -z "$user_home" ]; then
        error_exit "No non-root user found for SSH key setup"
    fi

    username=$(basename "$user_home")
}

generate_ssh_key() {
    if [[ ! -f "/home/$username/.ssh/id_rsa" ]]; then
        if ! su "$username" -c 'mkdir -p $HOME/.ssh; ssh-keygen -q -t rsa -N "" -f "$HOME/.ssh/id_rsa" >"$LOG" 2>&1'; then
            error_exit "SSH key generation failed."
        fi
    fi

    chown -R "$username:$username" "/home/$username/.ssh"
    chmod 600 "/home/$username/.ssh/id_rsa"
    chmod 644 "/home/$username/.ssh/id_rsa.pub"
}

add_key_to_maas() {
    if [[ $(maas admin sshkeys read | jq -r 'length') -eq 0 ]]; then
        if ! maas admin sshkeys create key="$(cat "/home/$username/.ssh/id_rsa.pub")" >"$TEMP_LOG" 2>&1; then
            error_exit "Failed to add SSH key to MAAS."
        fi
    fi
}

set_sshkey() {
    log "INFO" "Setting up SSH keys..."
    find_first_non_user
    generate_ssh_key
    add_key_to_maas
    log "INFO" "SSH keys setup completed"
}
