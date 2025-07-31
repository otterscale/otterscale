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
        echo "No non-root user found"
	exit 1
    fi

    NON_ROOT_USER=$(basename "$USER_HOME")
}


confirm_remove() {
    read -p "Confirm remove otters (y/n): " confirm
    if [[ "$confirm" =~ ^[Yy]$ ]]; then
        echo "Start uninstall otterscale..."
    elif [[ "$confirm" =~ ^[Nn]$ ]]; then
	break
    else
        echo "Invalid input. Please enter y or n."
    fi
}

remove_juju_machine() {
    
    NON_ROOT_USER
}

main() {
    confirm_remove

    #dd_ceph_disk
    remove_juju_machine
    remove_snap
    remove_apt
    remove_juju_file
}

main


