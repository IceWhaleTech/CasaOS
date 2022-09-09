#!/bin/bash

# copy to /casaOS/util/shell path
# chmod 755

log="logger -t usb-mount.sh -s "

ACTION=$1

DEVBASE=$2

DEVICE="/dev/${DEVBASE}"

# See if this drive is already mounted, and if so where
MOUNT_POINT=$(lsblk -l -p -o name,mountpoint | grep ${DEVICE} | awk '{print $2}')

do_mount() {

  if [ -n "${MOUNT_POINT}" ]; then
    ${log} "Warning: ${DEVICE} is already mounted at ${MOUNT_POINT}"
    exit 1
  fi

  # Get info for this drive: $ID_FS_LABEL and $ID_FS_TYPE
  eval $(blkid -o udev ${DEVICE} | grep -i -e "ID_FS_LABEL" -e "ID_FS_TYPE")

  #ID_FS_LABEL=新加卷
  #ID_FS_LABEL_ENC=新加卷
  #ID_FS_TYPE=ntfs

  # Figure out a mount point to use
  # LABEL=${ID_FS_LABEL}
  LABEL=${DEVBASE}
  if grep -q " /DATA/USB_Storage_${LABEL} " /etc/mtab; then
    # Already in use, make a unique one
    LABEL+="_${DEVBASE}"
  fi
  DEV_LABEL="${LABEL}"

  # Use the device name in case the drive doesn't have label
  if [ -z ${DEV_LABEL} ]; then
    DEV_LABEL="${DEVBASE}"
  fi


 MOUNT_POINT="/DATA/USB_Storage_${DEV_LABEL}"

  ${log} "Mount point: ${MOUNT_POINT}"

  mkdir -p ${MOUNT_POINT}


  # MOUNT_POINT="/DATA/USB_Storage1"
  # arr=("/DATA/USB_Storage1" "/DATA/USB_Storage2" "/DATA/USB_Storage3" "/DATA/USB_Storage4" "/DATA/USB_Storage5" "/DATA/USB_Storage6" "/DATA/USB_Storage7" "/DATA/USB_Storage8" "/DATA/USB_Storage9" "/DATA/USB_Storage10" "/DATA/USB_Storage11" "/DATA/USB_Storage12")
  # for folder in ${arr[@]}; do
  #   #如果文件夹不存在，创建文件夹
  #   if [ ! -d "$folder" ]; then
  #     mkdir -p ${folder}
  #     MOUNT_POINT=$folder
  #     break
  #   fi
  # done

  # ${log} "Mount point: ${MOUNT_POINT}"

  

  #  # Global mount options
  #  OPTS="rw,relatime"
  #
  #  # File system type specific mount options
  #  if [[ ${ID_FS_TYPE} == "vfat" ]]; then
  #    OPTS+=",users,gid=100,umask=000,shortname=mixed,utf8=1,flush"
  #  fi

  #  if ! mount -o ${OPTS} ${DEVICE} ${MOUNT_POINT}; then
  #    ${log} "Error mounting ${DEVICE} (status = $?)"
  #    rmdir "${MOUNT_POINT}"
  #    exit 1
  #  else
  #    # Track the mounted drives
  #    echo "${MOUNT_POINT}:${DEVBASE}" | cat >>"/var/log/usb-mount.track"
  #  fi
  #
  #  ${log} "Mounted ${DEVICE} at ${MOUNT_POINT}"

  case ${ID_FS_TYPE} in
  vfat)
    mount -t vfat -o rw,relatime,users,gid=100,umask=000,shortname=mixed,utf8=1,flush ${DEVICE} ${MOUNT_POINT}
    ;;
  ext[2-4])
    mount -o noatime ${DEVICE} ${MOUNT_POINT} >/dev/null 2>&1
    ;;
  exfat)
    mount -t exfat ${DEVICE} ${MOUNT_POINT} >/dev/null 2>&1
    ;;
  ntfs)
    ntfs-3g ${DEVICE} ${MOUNT_POINT}
    ;;
  iso9660)
    mount -t iso9660 ${DEVICE} ${MOUNT_POINT}
    ;;
  *)
    /bin/rmdir "${MOUNT_POINT}"
    exit 0
    ;;
  esac
}

do_umount() {

  if [[ -z ${MOUNT_POINT} ]]; then
    ${log} "Warning: ${DEVICE} is not mounted"
  else
    #/bin/kill -9 $(lsof ${MOUNT_POINT})
    umount -l ${DEVICE}
    ${log} "Unmounted ${DEVICE} from ${MOUNT_POINT}"
    if [ "`ls -A ${MOUNT_POINT}`" = "" ]; then
      /bin/rm -fr "${MOUNT_POINT}"
    fi
    sed -i.bak "\@${MOUNT_POINT}@d" /var/log/usb-mount.track
  fi

}

case "${ACTION}" in
add)
  do_mount
  ;;
remove)
  do_umount
  ;;
*)
  exit 1
  ;;
esac
