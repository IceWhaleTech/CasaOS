#!/bin/bash
###
 # @Author: LinkLeong link@icewhale.com
 # @Date: 2021-12-06 17:12:32
 # @LastEditors: LinkLeong
 # @LastEditTime: 2022-06-27 14:23:15
 # @FilePath: /CasaOS/shell/tools.sh
 # @Description: 
 # @Website: https://www.casaos.io
 # Copyright (c) 2022 by icewhale, All Rights Reserved. 
### 

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


run_external_script() {
  assist.sh
}

update() {
  curl -fsSL https://get.icewhale.io/casaos.sh | bash
  run_external_script
}