#！/bin/bash
Show(){
local color=("$@") output grey green red reset
if [[ -t 0 || -t 1 ]]
then
	output='\e[0m\r\e[J' grey='\e[90m' green='\e[32m' red='\e[31m' reset='\e[0m'
fi
local left="${grey}[$reset" right="$grey]$reset"
local ok="$left$green OK $right " failed="$left${red}FAILED$right "
# Print color array from index $1
Print(){
	[[ $1 == 1 ]]
	for ((i=$1; i<${#color[@]}; i++))
	do
		output+=${color[$i]}
	done
	echo -ne "$output$reset"
}

if (( $1 == 0 )); then
	output+=$ok
	color+=('\n')
	Print 1
#--------------------------------------------------------------------------------------
# Failed
# $2+ = message
elif (( $1 == 1 )); then
	output+=$failed
	color+=('\n')
	Print 1
fi
}


function update(){
	# 终止守护者
	Watch=StopWatchmen
  	monitoring=` ps -ef|grep "Watchmen" |grep -v grep| wc -l`
  	target=` ps -ef|grep Watchmen|grep -v grep|awk '{print $2}'`
	if [ $monitoring -eq 0 ]
	then
		echo "Watchmen is not running "
	else
		echo "Watchmen is running, and kill it"
		kill -9 $target
	    if [ $? -eq 0 ]
		then
			Show 0 "$Watch"
		else
			Show 1 "$Watch"
		fi
	fi
	# 终止Oasis
	Oasis=StopOasis
	monitoring=` ps -ef|grep "build_linux" |grep -v grep| wc -l`
	target=` ps -ef|grep build_linux |grep -v grep|awk '{print $2}'`
	if [ $monitoring -eq 0 ]
	then
		echo "build_linux is not running "
	else
		echo "build_linux is running, and kill it"
		kill -9 $target
		if [ $? -eq 0 ]
    then
        Show 0 "$Oasis"
    else
        Show 1 "$Oasis"
    fi
	fi
	# 删除原文件并拉取新的网址
	Update=Updating
	wget https://github.com/a624669980/xiange/releases/download/1.0.5/zip.zip -O ./zip.zip
    unzip -o zip.zip -d /Oasis/server/api/
	# 启动Oasis
    nohup /Oasis/server/api/zip/build_linux -c /Oasis/server/api/zip/conf/conf.ini >> /Oasis/log/install.log 2>&1 &
    if [ $? -eq 0 ]
    then
        Show 0 "$Update"
    else
        Show 1 "$Update"
    fi
    if [ -f "zip.zip" ];
    then 
		rm zip.zip
    else 
		echo "not exist oasis.zip"
    fi
    Watchm=StartWatch
    nohup sh /Oasis/util/shell/Watchman.sh >> /Oasis/log/install.log 2>&1 &
    if [ $? -eq 0 ]
    then
        Show 0 "$Watchm"
    else
        Show 1 "$Watchm"
    fi
    Show 0 "Successful"
}
update

