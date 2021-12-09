#!/bin/bash

#######################################
# Custom printing function
# Globals:
#   None
# Arguments:
#   $1 0:OK   1:FAILED
#   message
# Returns:
#   None
#######################################

readonly CASA_PATH=/casaOS/server
readonly casa_bin="casaos"

version=""

usage() {
  cat <<-EOF
		Usage: tool.sh [options]
		Valid options are:
		    -r <version>            verison of casaos
		    -h                      show this help message and exit
	EOF
  exit $1
}

show() {
  local color=("$@") output grey green red reset
  if [[ -t 0 || -t 1 ]]; then
    output='\e[0m\r\e[J' grey='\e[90m' green='\e[32m' red='\e[31m' reset='\e[0m'
  fi
  local left="${grey}[$reset" right="$grey]$reset"
  local ok="$left$green  OK  $right " failed="$left${red}FAILED$right " info="$left$green INFO $right "
  # Print color array from index $1
  Print() {
    [[ $1 == 1 ]]
    for ((i = $1; i < ${#color[@]}; i++)); do
      output+=${color[$i]}
    done
    echo -ne "$output$reset"
  }

  if (($1 == 0)); then
    output+=$ok
    color+=('\n')
    Print 1

  elif (($1 == 1)); then
    output+=$failed
    color+=('\n')
    Print 1

  elif (($1 == 2)); then
    output+=$info
    color+=('\n')
    Print 1
  fi
}

run_external_script() {
  assist.sh
}

update() {
  trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; return 1' ERR

  # Not every platform has or needs sudo (https://termux.com/linux.html)
  ((EUID)) && sudo_cmd="sudo"

  target_os="unsupported"
  target_arch="unknown"
  install_path="/usr/local/bin"

  # Fall back to /usr/bin if necessary
  if [[ ! -d $install_path ]]; then
    install_path="/usr/bin"
  fi

  #########################
  # Which OS and version? #
  #########################
  casa_tmp_folder="casaos"

  casa_dl_ext=".tar.gz"

  # NOTE: `uname -m` is more accurate and universal than `arch`
  # See https://en.wikipedia.org/wiki/Uname
  unamem="$(uname -m)"
  case $unamem in
  *aarch64*)
    target_arch="arm64"
    ;;
  *64*)
    target_arch="amd64"
    ;;
  *86*)
    target_arch="386"
    ;;
  *armv5*)
    target_arch="arm-5"
    ;;
  *armv6*)
    target_arch="arm-6"
    ;;
  *armv7*)
    target_arch="arm-7"
    ;;
  *)
    show 1 "Aborted, unsupported or unknown architecture: $unamem"
    return 2
    ;;
  esac

  unameu="$(tr '[:lower:]' '[:upper:]' <<<$(uname))"
  if [[ $unameu == *DARWIN* ]]; then
    target_os="darwin"
  elif [[ $unameu == *LINUX* ]]; then
    target_os="linux"
  elif [[ $unameu == *FREEBSD* ]]; then
    target_os="freebsd"
  elif [[ $unameu == *NETBSD* ]]; then
    target_os="netbsd"
  elif [[ $unameu == *OPENBSD* ]]; then
    target_os="openbsd"
  else
    show 1 "Aborted, unsupported or unknown OS: $uname"
    return 6
  fi

  ########################
  # Download and extract #
  ########################
  show 2 "Downloading CasaOS for $target_os/$target_arch..."
  if type -p curl >/dev/null 2>&1; then
    net_getter="curl -fsSL"
  elif type -p wget >/dev/null 2>&1; then
    net_getter="wget -qO-"
  else
    show 1 "Aborted, could not find curl or wget"
    return 7
  fi

  casa_file="${target_os}-$target_arch-casaos$casa_dl_ext"
  casa_tag="$(${net_getter} https://api.github.com/repos/IceWhaleTech/CasaOS/releases/latest | grep -o '"tag_name": ".*"' | sed 's/"//g' | sed 's/tag_name: //g')"
  casa_url="https://github.com/IceWhaleTech/CasaOS/releases/download/$casa_tag/$casa_file"
  show 2 "$casa_url"
  # Use $PREFIX for compatibility with Termux on Android
  rm -rf "$PREFIX/tmp/$casa_file"

  ${net_getter} "$casa_url" >"$PREFIX/tmp/$casa_file"

  show 2 "Extracting..."
  case "$casa_file" in
  *.zip) unzip -o "$PREFIX/tmp/$casa_file" -d "$PREFIX/tmp/" ;;
  *.tar.gz) tar -xzf "$PREFIX/tmp/$casa_file" -C "$PREFIX/tmp/" ;;
  esac

  chmod +x "$PREFIX/tmp/$casa_tmp_folder/$casa_bin"

  #stop service
  show 2 "Putting CasaOS in $install_path (may require password)"
  $sudo_cmd mv -f "$PREFIX/tmp/$casa_tmp_folder/$casa_bin" "$install_path/"
  show 2 "Putting CasaOS Shell file in $CASA_PATH (may require password)"
  #check shell folder
  local casa_shell_path=$CASA_PATH/shell

  if [[ -d $casa_shell_path ]]; then
    rm -rf $casa_shell_path
  fi

  $sudo_cmd mv -f $PREFIX/tmp/$casa_tmp_folder/shell "$CASA_PATH/shell"

  # remove tmp files
  $sudo_cmd rm -rf $PREFIX/tmp/$casa_tmp_folder

  if type -p $casa_bin >/dev/null 2>&1; then
    trap ERR
    run_external_script
    #   $sudo_cmd systemctl start casaos
     $sudo_cmd systemctl restart casaos
    show 0 "CasaOS Successfully updated."
    return 0
  else
    show 1 "Something went wrong, CasaOS is not in your path"
    trap ERR
    return 1
  fi
}

while getopts ":s:l:S:L:i:e:a:b:w:p:G:D:oOuUfgrczh" arg; do
  case "$arg" in
  r)
    version=$OPTARG
    update
    ;;
  h)
    usage 0
    ;;
  esac
done
