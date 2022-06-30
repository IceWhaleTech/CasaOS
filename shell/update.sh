#!/bin/bash
###
 # @Author:  LinkLeong link@icewhale.com
 # @Date: 2022-06-30 10:08:33
 # @LastEditors: LinkLeong
 # @LastEditTime: 2022-06-30 19:46:12
 # @FilePath: /CasaOS/shell/update.sh
 # @Description:
### 

((EUID)) && sudo_cmd="sudo"

# SYSTEM INFO
readonly UNAME_M="$(uname -m)"

# CasaOS PATHS
readonly CASA_REPO=LinkLeong/casaos-alpha
readonly CASA_UNZIP_TEMP_FOLDER=/tmp/casaos
readonly CASA_BIN=casaos
readonly CASA_BIN_PATH=/usr/bin/casaos
readonly CASA_CONF_PATH=/etc/casaos.conf
readonly CASA_SERVICE_PATH=/etc/systemd/system/casaos.service
readonly CASA_HELPER_PATH=/usr/share/casaos/shell/
readonly CASA_USER_CONF_PATH=/var/lib/casaos/conf/
readonly CASA_DB_PATH=/var/lib/casaos/db/
readonly CASA_TEMP_PATH=/var/lib/casaos/temp/
readonly CASA_LOGS_PATH=/var/log/casaos/
readonly CASA_PACKAGE_EXT=".tar.gz"
readonly CASA_RELEASE_API="https://api.github.com/repos/${CASA_REPO}/releases"
readonly CASA_OPENWRT_DOCS="https://github.com/IceWhaleTech/CasaOS-OpenWrt"

readonly COLOUR_RESET='\e[0m'
readonly aCOLOUR=(
    '\e[38;5;154m' # green  	| Lines, bullets and separators
    '\e[1m'        # Bold white	| Main descriptions
    '\e[90m'       # Grey		| Credits
    '\e[91m'       # Red		| Update notifications Alert
    '\e[33m'       # Yellow		| Emphasis
)

Target_Arch=""
Target_Distro="debian"
Target_OS="linux"
Casa_Tag=""


#######################################
# Custom printing function
# Globals:
#   None
# Arguments:
#   $1 0:OK   1:FAILED  2:INFO  3:NOTICE
#   message
# Returns:
#   None
#######################################

Show() {
    # OK
    if (($1 == 0)); then
        echo -e "${aCOLOUR[2]}[$COLOUR_RESET${aCOLOUR[0]}  OK  $COLOUR_RESET${aCOLOUR[2]}]$COLOUR_RESET $2"
    # FAILED
    elif (($1 == 1)); then
        echo -e "${aCOLOUR[2]}[$COLOUR_RESET${aCOLOUR[3]}FAILED$COLOUR_RESET${aCOLOUR[2]}]$COLOUR_RESET $2"
    # INFO
    elif (($1 == 2)); then
        echo -e "${aCOLOUR[2]}[$COLOUR_RESET${aCOLOUR[0]} INFO $COLOUR_RESET${aCOLOUR[2]}]$COLOUR_RESET $2"
    # NOTICE
    elif (($1 == 3)); then
        echo -e "${aCOLOUR[2]}[$COLOUR_RESET${aCOLOUR[4]}NOTICE$COLOUR_RESET${aCOLOUR[2]}]$COLOUR_RESET $2"
    fi
}

Warn() {
    echo -e "${aCOLOUR[3]}$1$COLOUR_RESET"
}

# 0 Check_exist
Check_Exist() {
    #Create Dir
    Show 2 "Create Folders."
    ${sudo_cmd} mkdir -p ${CASA_HELPER_PATH}
    ${sudo_cmd} mkdir -p ${CASA_LOGS_PATH}
    ${sudo_cmd} mkdir -p ${CASA_USER_CONF_PATH}
    ${sudo_cmd} mkdir -p ${CASA_DB_PATH}
    ${sudo_cmd} mkdir -p ${CASA_TEMP_PATH}

   
    Show 2 "Start cleaning up the old version."

    ${sudo_cmd} rm -rf /usr/lib/systemd/system/casaos.service
    
    ${sudo_cmd} rm -rf /lib/systemd/system/casaos.service

    if [[ -f "/casaOS/server/conf/conf.ini" ]]; then
        ${sudo_cmd} cp -rf /casaOS/server/conf/conf.ini ${CASA_CONF_PATH}
        ${sudo_cmd} cp -rf /casaOS/server/conf/*.json ${CASA_USER_CONF_PATH}
    fi

    if [[ -d "/casaOS/server/db" ]]; then
        ${sudo_cmd} cp -rf /casaOS/server/db/* ${CASA_DB_PATH}
    fi

    Show 0 "Clearance completed."

}

# 1 Check Arch
Check_Arch() {
    case $UNAME_M in
    *aarch64*)
        Target_Arch="arm64"
        ;;
    *64*)
        Target_Arch="amd64"
        ;;
    *armv7*)
        Target_Arch="arm-7"
        ;;
    *)
        Show 1 "Aborted, unsupported or unknown architecture: $UNAME_M"
        exit 1
        ;;
    esac
    Show 0 "Your hardware architecture is : $UNAME_M"
}



#Download CasaOS Package
Download_CasaOS() {
    Show 2 "Downloading CasaOS for ${Target_OS}/${Target_Arch}..."
    Net_Getter="curl -fsSLk"
    Casa_Package="${Target_OS}-${Target_Arch}-casaos${CASA_PACKAGE_EXT}"
    if [[ ! -n "$version" ]]; then
        Casa_Tag="$(${Net_Getter} ${CASA_RELEASE_API}/latest | grep -o '"tag_name": ".*"' | sed 's/"//g' | sed 's/tag_name: //g')"
    elif [[ $version == "pre" ]]; then
        Casa_Tag="$(${net_getter} ${CASA_RELEASE_API} | grep -o '"tag_name": ".*"' | sed 's/"//g' | sed 's/tag_name: //g' | sed -n '1p')"
    else
        Casa_Tag="$version"
    fi
    Casa_Package_URL="https://github.com/${CASA_REPO}/releases/download/${Casa_Tag}/${Casa_Package}"
    echo
    # Remove Temp File
    ${sudo_cmd} rm -rf "$PREFIX/tmp/${Casa_Package}"
    # Download Package
    ${Net_Getter} "${Casa_Package_URL}" >"$PREFIX/tmp/${Casa_Package}"
    if [[ $? -ne 0 ]]; then
        Show 1 "Download failed, Please check if your internet connection is working and retry."
        exit 1
    else
        Show 0 "Download successful!"
    fi
    #Extract CasaOS Package
    Show 2 "Extracting..."
    case "${Casa_Package}" in
    *.zip) ${sudo_cmd} unzip -o "$PREFIX/tmp/${Casa_Package}" -d "$PREFIX/tmp/" ;;
    *.tar.gz) ${sudo_cmd} tar -xzf "$PREFIX/tmp/${Casa_Package}" -C "$PREFIX/tmp/" ;;
    esac
    #Setting Executable Permissions
    ${sudo_cmd} chmod +x "$PREFIX${CASA_UNZIP_TEMP_FOLDER}/${CASA_BIN}"

}

#Install Addons
Install_Addons() {
    Show 2 "Installing CasaOS Addons"
    ${sudo_cmd} cp -rf "$PREFIX${CASA_UNZIP_TEMP_FOLDER}/shell/11-usb-mount.rules" "/etc/udev/rules.d/"
    ${sudo_cmd} cp -rf "$PREFIX${CASA_UNZIP_TEMP_FOLDER}/shell/usb-mount@.service" "/etc/systemd/system/"
    sync
}

#Clean Temp Files
Clean_Temp_Files() {
    Show 0 "Clean..."
    ${sudo_cmd} rm -rf "$PREFIX${CASA_UNZIP_TEMP_FOLDER}"
    sync
}

#Install CasaOS
Install_CasaOS() {
    Show 2 "Installing..."

    # Install Bin
    ${sudo_cmd} mv -f $PREFIX${CASA_UNZIP_TEMP_FOLDER}/${CASA_BIN} ${CASA_BIN_PATH}

    # Install Helper
    if [[ -d ${CASA_HELPER_PATH} ]]; then
        ${sudo_cmd} rm -rf ${CASA_HELPER_PATH}*
    fi
    ${sudo_cmd} cp -rf $PREFIX${CASA_UNZIP_TEMP_FOLDER}/shell/* ${CASA_HELPER_PATH}
    #Setting Executable Permissions
    ${sudo_cmd} chmod +x $PREFIX${CASA_HELPER_PATH}*

    # Install Conf
    if [[ ! -f ${CASA_CONF_PATH} ]]; then
        if [[ -f $PREFIX${CASA_UNZIP_TEMP_FOLDER}/conf/conf.ini.sample ]]; then
            ${sudo_cmd} mv -f $PREFIX${CASA_UNZIP_TEMP_FOLDER}/conf/conf.ini.sample ${CASA_CONF_PATH}
        else
            ${sudo_cmd} mv -f $PREFIX${CASA_UNZIP_TEMP_FOLDER}/conf/conf.conf.sample ${CASA_CONF_PATH}
        fi

    fi
    sync

    if [[ ! -x "$(command -v ${CASA_BIN})" ]]; then
        Show 1 "Installation failed, please try again."
        exit 1
    else
        Show 0 "CasaOS Successfully installed."
    fi
}

#Generate Service File
Generate_Service() {
    if [ -f ${CASA_SERVICE_PATH} ]; then
        Show 2 "Try stop CasaOS system service."
        # Stop before generation
        if [[ $(systemctl is-active ${CASA_BIN} &>/dev/null) ]]; then
            ${sudo_cmd} systemctl stop ${CASA_BIN}
        fi
    fi
    Show 2 "Create system service for CasaOS."

    ${sudo_cmd} tee ${CASA_SERVICE_PATH} >/dev/null <<EOF
				[Unit]
				Description=CasaOS Service
				StartLimitIntervalSec=0

				[Service]
				Type=simple
				LimitNOFILE=15210
				Restart=always
				RestartSec=1
				User=root
				ExecStart=${CASA_BIN_PATH} -c ${CASA_CONF_PATH}

				[Install]
				WantedBy=multi-user.target
EOF
    Show 0 "CasaOS service Successfully created."
}

# Start CasaOS
Start_CasaOS() {
    Show 2 "Create a system startup service for CasaOS."
    $sudo_cmd systemctl daemon-reload
    $sudo_cmd systemctl enable ${CASA_BIN}
}

Check_Arch

# Step 7: Download CasaOS
Check_Exist
Download_CasaOS

# Step 8: Install Addon
Install_Addons

# Step 9: Install CasaOS
Install_CasaOS

# Step 10: Generate_Service
Generate_Service

# Step 11: Start CasaOS
Start_CasaOS
Clean_Temp_Files