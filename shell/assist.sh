#!/bin/bash



# add in v0.2.5

readonly CASA_DEPANDS="curl smartmontools parted fdisk ntfs-3g"

version_0_2_5() {
  install_depends "$CASA_DEPANDS"
}
version_0_2_11() {
  sysctl -w net.core.rmem_max=2500000
}

#Install Depends
install_depends() {
    ((EUID)) && sudo_cmd="sudo"
    if [[ ! -x "$(command -v '$1')" ]]; then
        packagesNeeded=$1
        if [ -x "$(command -v apk)" ]; then
            $sudo_cmd apk add --no-cache $packagesNeeded
        elif [ -x "$(command -v apt-get)" ]; then
            $sudo_cmd apt-get -y -q install $packagesNeeded
        elif [ -x "$(command -v dnf)" ]; then
            $sudo_cmd dnf install $packagesNeeded
        elif [ -x "$(command -v zypper)" ]; then
            $sudo_cmd zypper install $packagesNeeded
        fi
    fi
}

version_0_2_5

version_0_2_11
