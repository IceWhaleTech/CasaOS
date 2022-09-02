#!/bin/bash
###
# @Author: LinkLeong link@icewhale.org
# @Date: 2022-08-25 11:41:22
 # @LastEditors: LinkLeong
 # @LastEditTime: 2022-08-31 17:54:17
 # @FilePath: /CasaOS/build/scripts/setup/service.d/casaos/debian/setup-casaos.sh
# @Description:

# @Website: https://www.casaos.io
# Copyright (c) 2022 by icewhale, All Rights Reserved.
###

set -e

APP_NAME="casaos"

# copy config files
CONF_PATH=/etc/casaos
OLD_CONF_PATH=/etc/casaos.conf
CONF_FILE=${CONF_PATH}/${APP_NAME}.conf
CONF_FILE_SAMPLE=${CONF_PATH}/${APP_NAME}.conf.sample


if [ -f "${OLD_CONF_PATH}" ]; then
    echo "copy old conf"
    cp "${OLD_CONF_PATH}" "${CONF_FILE}"
fi
if [ ! -f "${CONF_FILE}" ]; then
    echo "Initializing config file..."
    cp -v "${CONF_FILE_SAMPLE}" "${CONF_FILE}"
fi

if systemctl is-active "${APP_NAME}.service" &>/dev/null ;then
    echo "server started"
else
    # enable and start service
    systemctl daemon-reload

    echo "Enabling service..."
    systemctl enable --force --no-ask-password "${APP_NAME}.service"

    echo "Starting service..."
    systemctl start --force --no-ask-password "${APP_NAME}.service"
fi