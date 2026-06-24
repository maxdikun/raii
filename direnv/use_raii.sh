#!/usr/bin/env bash
# Add this file to your direnv stdlib, e.g.:
#   mkdir -p ~/.config/direnv/lib
#   cp use_raii.sh ~/.config/direnv/lib/use_raii.sh
#
# Then in your .envrc:
#   use raii
#   # or with a custom config:
#   use raii ./my-config.toml

use_raii() {
    local config="${1:-raii.toml}"
    local owner=$$

    # Start resources and register this shell as an owner.
    # A background watchdog will automatically call stop when this shell exits.
    raii start --config "$config" --owner "$owner"
}
