#!/bin/bash

#add in v0.2.3
version_0_2_3() {
  ((EUID)) && sudo_cmd="sudo"
  $sudo_cmd cp -rf /casaOS/server/shell/11-usb-mount.rules /etc/udev/rules.d/
  $sudo_cmd chmod +x /casaOS/server/shell/usb-mount.sh
  $sudo_cmd cp -rf /casaOS/server/shell/usb-mount@.service /etc/systemd/system/

}

version_0_2_3
