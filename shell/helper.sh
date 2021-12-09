#!/bin/bash

# 获取系统信息
GetSysInfo() {
  if [ -s "/etc/redhat-release" ]; then
    SYS_VERSION=$(cat /etc/redhat-release)
  elif [ -s "/etc/issue" ]; then
    SYS_VERSION=$(cat /etc/issue)
  fi
  SYS_INFO=$(uname -a)
  SYS_BIT=$(getconf LONG_BIT)
  MEM_TOTAL=$(free -m | grep Mem | awk '{print $2}')
  CPU_INFO=$(getconf _NPROCESSORS_ONLN)

  echo -e ${SYS_VERSION}
  echo -e Bit:${SYS_BIT} Mem:${MEM_TOTAL}M Core:${CPU_INFO}
  echo -e ${SYS_INFO}
}

#获取网卡信息
GetNetCard() {
  if [ "$1" == "1" ]; then
    if [ -d "/sys/devices/virtual/net" ]; then
      ls /sys/devices/virtual/net
    fi
  else
    if [ -d "/sys/devices/virtual/net" ] && [ -d "/sys/class/net" ]; then
      ls /sys/class/net/ | grep -v "$(ls /sys/devices/virtual/net/)"
    fi
  fi
}

#查看网卡状态
#param 网卡名称
CatNetCardState() {
  if [ -e "/sys/class/net/$1/operstate" ]; then
    cat /sys/class/net/$1/operstate
  fi
}

#获取docker根目录
GetDockerRootDir() {
  if hash docker 2>/dev/null; then
    docker info | grep 'Docker Root Dir' | awk -F ':' '{print $2}'
  else
    echo ""
  fi
}

#删除安装应用文件夹
#param 需要删除的文件夹路径
DelAppConfigDir() {
  if [ -d $1 ]; then
    rm -fr $1
  fi
}

#zerotier本机已加入的网络
#result start,end,sectors
GetLocalJoinNetworks() {
  zerotier-cli listnetworks -j
}

#移除挂载点,删除已挂在的文件夹
UMountPorintAndRemoveDir() {
  DEVICE=$1
  MOUNT_POINT=$(mount | grep ${DEVICE} | awk '{ print $3 }')
  if [[ -z ${MOUNT_POINT} ]]; then
    ${log} "Warning: ${DEVICE} is not mounted"
  else
    umount -l ${DEVICE}
    ${log} "Unmounted ${DEVICE} from ${MOUNT_POINT}"
    /bin/rmdir "${MOUNT_POINT}"
    sed -i.bak "\@${MOUNT_POINT}@d" /var/log/usb-mount.track
  fi
}

#格式化fat32磁盘
#param 需要格式化的目录 /dev/sda1
#param 格式
FormatDisk() {
  if [ "$2" == "fat32" ]; then
    mkfs.vfat -F 32 $1
  elif [ "$2" == "ntfs" ]; then
    mkfs.ntfs $1
  elif [ "$2" == "ext4" ]; then
    mkfs.ext4 -F $1
  elif [ "$2" == "exfat" ]; then
    mkfs.exfat $1
  else
    mkfs.ext4 -F $1
  fi
}

#删除分区
#param 路径   /dev/sdb
#param 删除分区的区号
DelPartition() {
  fdisk $1 <<EOF
  d
  $2
  wq
EOF
}

#添加分区只有一个分区
#param 路径   /dev/sdb
#param 要挂载的目录
AddPartition() {

  DelPartition $1
  parted -s $1 mklabel gpt

  parted -s $1 mkpart primary ext4 0 100%

  mkfs.ext4 $11

  partprobe $1

  #  mount $11 $2

}

#磁盘类型
GetDiskType() {
  fdisk $1 -l | grep Disklabel | awk -F: '{print $2}'
}

#获磁盘的插入路径
#param 路径 /dev/sda
GetPlugInDisk() {
  fdisk -l | grep 'Disk' | grep 'sd' | awk -F , '{print substr($1,11,3)}'
}

#获取磁盘状态
#param 磁盘路径
GetDiskHealthState() {
  smartctl -H $1 | grep "SMART Health Status" | awk -F ":" '{print$2}'
}

#获取磁盘字节数量和扇区数量
#param 磁盘路径  /dev/sda
#result bytes
#result sectors
GetDiskSizeAndSectors() {
  fdisk $1 -l | grep "/dev/sda:" | awk -F, 'BEGIN {OFS="\n"}{print $2,$3}' | awk '{print $1}'
}

#获取磁盘分区数据扇区
#param 磁盘路径  /dev/sda
#result start,end,sectors
GetPartitionSectors() {
  fdisk $1 -l | grep "/dev/sda[1-9]" | awk 'BEGIN{OFS=","}{print $1,$2,$3,$4}'
}

#检查没有使用的挂载点删除文件夹
AutoRemoveUnuseDir() {
  DIRECTORY="/mnt/"
  dir=$(ls -l $DIRECTORY | awk '/^d/ {print $NF}')
  for i in $dir; do

    path="$DIRECTORY$i"
    mountStr=$(mountpoint $path)
    notMountpoint="is not a mountpoint"
    if [[ $mountStr =~ $notMountpoint ]]; then
      if [ "$(ls -A $path)" = "" ]; then
        rm -fr $path
      else
        echo "$path is not empty"
      fi
    fi
  done
}

#重载samba服务
ReloadSamba() {
  /etc/init.d/smbd reload
}

# $1=sda1
# $2=volume{1}
do_mount() {
  DEVBASE=$1
  DEVICE="${DEVBASE}"
  # See if this drive is already mounted, and if so where
  MOUNT_POINT=$(mount | grep ${DEVICE} | awk '{ print $3 }')

  if [ -n "${MOUNT_POINT}" ]; then
    ${log} "Warning: ${DEVICE} is already mounted at ${MOUNT_POINT}"
    exit 1
  fi

  # Get info for this drive: $ID_FS_LABEL and $ID_FS_TYPE
  eval $(blkid -o udev ${DEVICE} | grep -i -e "ID_FS_LABEL" -e "ID_FS_TYPE")

  LABEL=$2
  if grep -q " /media/${LABEL} " /etc/mtab; then
    # Already in use, make a unique one
    LABEL+="-${DEVBASE}"
  fi
  DEV_LABEL="${LABEL}"

  # Use the device name in case the drive doesn't have label
  if [ -z ${DEV_LABEL} ]; then
    DEV_LABEL="${DEVBASE}"
  fi

  MOUNT_POINT="/media/${DEV_LABEL}"

  ${log} "Mount point: ${MOUNT_POINT}"

  mkdir -p ${MOUNT_POINT}

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
