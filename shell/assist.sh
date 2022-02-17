#!/bin/bash

#add in v0.2.3
version_0_2_3() {
  ((EUID)) && sudo_cmd="sudo"
  $sudo_cmd cp -rf /casaOS/server/shell/11-usb-mount.rules /etc/udev/rules.d/
  $sudo_cmd chmod +x /casaOS/server/shell/usb-mount.sh
  $sudo_cmd cp -rf /casaOS/server/shell/usb-mount@.service /etc/systemd/system/

}

# add in v0.2.5

readonly CASA_DEPANDS="curl smartmontools parted fdisk partprobe ntfs-3g"

version_0_2_5() {
  install_depends "$CASA_DEPANDS"
}


#Install Depends
install_depends() {
    ((EUID)) && sudo_cmd="sudo"
    if [[ ! -x "$(command -v '$1')" ]]; then
        show 2 "Install the necessary dependencies: $1"
        packagesNeeded=$1
        if [ -x "$(command -v apk)" ]; then
            $sudo_cmd apk add --no-cache $packagesNeeded
        elif [ -x "$(command -v apt-get)" ]; then
            $sudo_cmd apt-get -y -q install $packagesNeeded
        elif [ -x "$(command -v dnf)" ]; then
            $sudo_cmd dnf install $packagesNeeded
        elif [ -x "$(command -v zypper)" ]; then
            $sudo_cmd zypper install $packagesNeeded
        else
            show 1 "Package manager not found. You must manually install: $packagesNeeded"
        fi
    fi
}

version_0_2_3

version_0_2_5
