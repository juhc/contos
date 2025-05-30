#!/bin/bash
set -euo pipefail

# Function for displaying an error and exiting
function error() {
    echo "ERROR: $1"
    exit 1
}

# Check arguments count
if [[ $# -ne 4 ]]; then
    echo "Usage: $0 <rootfs-update.tar.gz> <path_to_system_a> <path_to_system_b> <update_conf>"
    exit 1
fi

# Input parameters
UPDATE_ARCHIVE="$1"
SYSTEM_A="$2"
SYSTEM_B="$3"
UPDATE_CONF="$4"

# Check if update archive exists and is readable
[[ ! -f "$UPDATE_ARCHIVE" ]] && error "Update archive '$UPDATE_ARCHIVE' not found"
[[ ! -r "$UPDATE_ARCHIVE" ]] && error "Update archive '$UPDATE_ARCHIVE' not readable"

# Check if update config exists and is readable/writable
[[ ! -f "$UPDATE_CONF" ]] && error "Update config '$UPDATE_CONF' not found"
[[ ! -r "$UPDATE_CONF" || ! -w "$UPDATE_CONF" ]] && error "Update config '$UPDATE_CONF' is not readable or writable"

# Determine which system is currently active (A or B)
ACTIVE=$(grep '^ACTIVE=' "$UPDATE_CONF" | cut -d= -f2)

if [[ "$ACTIVE" == "A" ]]; then
    CURRENT="$SYSTEM_A"
    TARGET="$SYSTEM_B"
    NEXT="B"
elif [[ "$ACTIVE" == "B" ]]; then
    CURRENT="$SYSTEM_B"
    TARGET="$SYSTEM_A"
    NEXT="A"
else
    error "Invalid ACTIVE value in update.conf (must be A or B)"
fi

# Check that system directories exist and are writable
[[ ! -d "$CURRENT" ]] && error "Current rootfs path '$CURRENT' does not exist"
[[ ! -d "$TARGET" ]] && error "Target rootfs path '$TARGET' does not exist"
[[ ! -w "$TARGET" ]] && error "No write access to target rootfs '$TARGET'"

# Clean the target rootfs directory before extracting the new version
echo "[*] Cleaning target rootfs: $TARGET"
rm -rf "${TARGET:?}/"*

# Extract the new rootfs archive into the target directory
echo "[*] Extracting new rootfs to $TARGET"
tar -xzf "$UPDATE_ARCHIVE" -C "$TARGET"

# Basic sanity check: ensure the extracted rootfs has an init binary
if [[ ! -f "$TARGET/sbin/init" && ! -f "$TARGET/init" ]]; then
    error "Extracted rootfs missing /sbin/init or /init — invalid image"
fi

# Update the active partition flag in the config
echo "[*] Updating active rootfs to $NEXT"
echo "ACTIVE=$NEXT" > "$UPDATE_CONF"

# Done!
echo "Update complete. Reboot to apply changes."
