#!/bin/bash

set -euo pipefail

ROOTFS_IMAGE="contos-rootfs.tar.gz"
MOUNT_DIR="/mnt/contos"
BOOT_PARTITION="/dev/sda1"
ROOT_PARTITION="/dev/sda2"
DISK="/dev/sda"

print_info() {
  echo -e "\e[1;34m[INFO]\e[0m $1"
}

print_error() {
  echo -e "\e[1;31m[ERROR]\e[0m $1"
}

check_root() {
  if [[ $EUID -ne 0 ]]; then
    print_error "Запустите скрипт от root"
    exit 1
  fi
}

confirm_disk() {
  lsblk
  echo
  read -rp "Введите диск для установки (например, /dev/sda): " DISK
  echo "ВНИМАНИЕ: Все данные на $DISK будут удалены!"
  read -rp "Подтвердите удаление (yes/no): " CONFIRM
  [[ "$CONFIRM" != "yes" ]] && exit 1
}

partition_disk() {
  print_info "Разметка диска $DISK"
  wipefs -a "$DISK"
  parted -s "$DISK" mklabel gpt
  parted -s "$DISK" mkpart ESP fat32 1MiB 513MiB
  parted -s "$DISK" set 1 esp on
  parted -s "$DISK" mkpart primary ext4 513MiB 100%

  BOOT_PARTITION="${DISK}1"
  ROOT_PARTITION="${DISK}2"

  mkfs.vfat -F32 "$BOOT_PARTITION"
  mkfs.ext4 "$ROOT_PARTITION"
}

mount_partitions() {
  mkdir -p "$MOUNT_DIR"
  mount "$ROOT_PARTITION" "$MOUNT_DIR"
  mkdir -p "$MOUNT_DIR/boot/efi"
  mount "$BOOT_PARTITION" "$MOUNT_DIR/boot/efi"
}

install_rootfs() {
  print_info "Распаковка rootfs"
  tar -xzf "$ROOTFS_IMAGE" -C "$MOUNT_DIR"
}

generate_fstab() {
  print_info "Генерация fstab"
  UUID_ROOT=$(blkid -s UUID -o value "$ROOT_PARTITION")
  UUID_BOOT=$(blkid -s UUID -o value "$BOOT_PARTITION")

  cat > "$MOUNT_DIR/etc/fstab" <<EOF
UUID=$UUID_ROOT / ext4 defaults 0 1
UUID=$UUID_BOOT /boot/efi vfat defaults 0 2
EOF
}

set_hostname() {
  echo "contos" > "$MOUNT_DIR/etc/hostname"
}

install_bootloader() {
  print_info "Установка загрузчика"
  mount --bind /dev "$MOUNT_DIR/dev"
  mount --bind /proc "$MOUNT_DIR/proc"
  mount --bind /sys "$MOUNT_DIR/sys"

  chroot "$MOUNT_DIR" grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=contos
  chroot "$MOUNT_DIR" grub-mkconfig -o /boot/grub/grub.cfg
}

cleanup() {
  umount -R "$MOUNT_DIR"
  print_info "Установка завершена. Вы можете перезагрузить систему."
}

main() {
  check_root
  confirm_disk
  partition_disk
  mount_partitions
  install_rootfs
  generate_fstab
  set_hostname
  install_bootloader
  cleanup
}

main "$@"
