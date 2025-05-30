#!/bin/bash
set -e

IMAGES="$1"
ISO="${IMAGES}/iso"
ROOTFS="${IMAGES}/rootfs.tar"
DISK_IMG="${IMAGES}/contos.img"
ISO_IMG="${IMAGES}/contos.iso"

echo ">>> [INFO] Creating output directories..."
mkdir -p "${ISO}"
mkdir -p "${IMAGES}/mnt"

echo ">>> [INFO] Extracting rootfs..."
tar -xf "${ROOTFS}" -C "${ISO}"

echo ">>> [INFO] Creating ISO image..."
xorriso -as mkisofs \
  -o "${ISO_IMG}" \
  -isohybrid-mbr /usr/lib/syslinux/isohdpfx.bin \
  -c boot/boot.cat \
  -b boot/syslinux/isolinux.bin \
  -no-emul-boot -boot-load-size 4 -boot-info-table \
  "${ISO}"

echo ">>> [INFO] Setting up loop device..."
LOOP_ISO=$(losetup --find --show "${ISO_IMG}")
if [ -z "$LOOP_ISO" ]; then
    echo "ERROR: Failed to setup loop device for ISO."
    exit 1
fi

echo ">>> [INFO] Mounting ISO..."
mount "${LOOP_ISO}" "${IMAGES}/mnt"

# Пример дальнейших действий:
# cp -r "${IMAGES}/mnt/somefile" /some/target/location

echo ">>> [INFO] Cleaning up..."
umount "${IMAGES}/mnt"
losetup -d "${LOOP_ISO}"
rm -rf "${IMAGES}/mnt"
rm -rf "${ISO}"

echo ">>> [SUCCESS] post_image.sh completed successfully."
