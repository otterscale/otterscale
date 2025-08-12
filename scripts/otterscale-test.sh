#!/usr/bin/env bash
# otterscale_test.sh
#
# Purpose:  Repeatedly install and uninstall the OtterScale package.
#           Intended for production‑level testing.
# Usage:    curl -L <repo‑url>/otterscale_test.sh | sudo bash -s -- <iterations>
#           where <iterations> is a positive integer.

set -euo pipefail

# --------------------------------------------------------------------
# Logging helpers
# --------------------------------------------------------------------
log() {
    local level="$1"; shift
    printf '[%s] %s: %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$level" "$*"
}

log_info()    { log "INFO"    "$@"; }
log_error()   { log "ERROR"   "$@"; }
log_success() { log "SUCCESS" "$@"; }

# --------------------------------------------------------------------
# Argument validation
# --------------------------------------------------------------------
if [[ $# -ne 1 ]]; then
    log_error "Usage: sudo bash $0 <iteration-count>"
    exit 1
fi

ITERATIONS="$1"

if ! [[ "$ITERATIONS" =~ ^[1-9][0-9]*$ ]]; then
    log_error "Iteration count must be a positive integer"
    exit 1
fi

# --------------------------------------------------------------------
# Core functions
# --------------------------------------------------------------------
install_otterscale() {
    log_info "Starting installation (attempt $1)"
    curl -fsSL \
        https://raw.githubusercontent.com/openhdc/otterscale/refs/heads/main/scripts/install.sh |
        sudo bash -s -- --config=~/install.cfg
    log_success "Installation succeeded (attempt $1)"
}

uninstall_otterscale() {
    log_info "Starting uninstallation (attempt $1)"
    curl -fsSL \
        https://raw.githubusercontent.com/openhdc/otterscale/refs/heads/main/scripts/unstall.sh |
        sudo bash -s --
    log_success "Uninstallation succeeded (attempt $1)"
}

# --------------------------------------------------------------------
# Main loop
# --------------------------------------------------------------------
for ((i=1; i<=ITERATIONS; i++)); do
    log_info "=== Cycle $i of $ITERATIONS ==="
    install_otterscale "$i"
    uninstall_otterscale "$i"
    echo
done

log_success "All $ITERATIONS cycles completed without error"
exit 0