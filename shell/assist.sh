#!/bin/bash

#update to v0.2.3
version_0_2_3(){
  ((EUID)) && sudo_cmd="sudo"

#copy file to path
  if [ ! -s "/etc/udev/rules.d/11-usb-mount.rules" ]; then
    $sudo_cmd cp /casaOS/server/shell/11-usb-mount.rules /etc/udev/rules.d/
  fi

  if [ ! -s "/casaOS/util/shell/usb-mount.sh" ]; then
    $sudo_cmd cp /casaOS/server/shell/usb-mount.sh /casaOS/util/shell/
    $sudo_cmd chmod +x /casaOS/util/shell/usb-mount.sh
  fi
  if [ ! -s "/etc/systemd/system/cp /casaOS/server/shell/usb-mount@.service" ]; then
     $sudo_cmd cp /casaOS/server/shell/usb-mount@.service /etc/systemd/system/
  fi




}

version_0_2_3