#!/bin/bash
# 用法:
# 1. 设置当前脚本可执行权限 chmod +x
# 2. 执行 ./alived.sh start(启动服务)
if [ $# -lt 1 ]; then
  echo "command:
  start       Start server.
  stop        Stop server.
  restart     Restart server.
  startjob    Start job.
  stopjob     Stop job.
  restartjob  Restart job."
  exit 1
fi
dir=`exec pwd`
# 启动服务
if [ $1 == "start" ]; then
    # source ~/.profile
    echo "start..."
    cd ${dir}
    #   nohup go run main.go > ./runtime/run.log 2>&1 &
    # rm -f ./absGo
    go build -tags=jsoniter -o absGo -ldflags "-w -s"
    if [ $# -eq 2 ] && [ $2 == "-d" ]; then
        nohup ./absGo server > ./runtime/run.log 2>&1 &
        # tail -n 20 ./runtime/run.log
        # tail -f ./runtime/run.log
    else
        ./absGo server
    fi
    echo "start success!"
  # 平滑重启服务
  elif [ $1 == "restart" ]; then
    echo "restart..."
    cd ${dir}
    rm -f ./absGo
    go build -tags=jsoniter -o absGo -ldflags "-w -s"
    ps aux | grep "absGo server" | grep -v grep | awk '{print $2}' | xargs kill -9
    if [ $# -eq 2 ] && [ $2 == "-d" ]; then
        nohup ./absGo server  > ./runtime/run.log 2>&1 &
        # tail -n 20 ./runtime/run.log
    else
        ./absGo server
    fi
    echo "restart success!"
  # 停止服务
  elif [ $1 == "stop" ]; then
    echo "stop..."
    ps aux | grep "absGo server" | grep -v grep | awk '{print $2}' | xargs kill
    echo "stop success!"
  # job服务系列
  elif [ $1 == "startjob" ]; then
    echo "start job..."
    cd ${dir}
    # rm ./absGo
    # go build -tags=jsoniter -o absGo -ldflags "-w -s"
    if [ $# -eq 2 ] && [ $2 == "-d" ]; then
        nohup ./absGo job > ./runtime/runjob.log 2>&1 &
    else
        ./absGo job
    fi
    echo "start job success!"
  elif [ $1 == "restartjob" ]; then
    echo "restart job..."
    cd ${dir}
    rm -f ./absGo
    go build -tags=jsoniter -o absGo -ldflags "-w -s"
    ps aux | grep "absGo job" | grep -v grep | awk '{print $2}' | xargs kill -9
    if [ $# -eq 2 ] && [ $2 == "-d" ]; then
        nohup ./absGo job > ./runtime/runjob.log 2>&1 &
    else
        ./absGo job
    fi
    echo "restart job success!"
  elif [ $1 == "stopjob" ]; then
    echo "stop job..."
    ps aux | grep "absGo job" | grep -v grep | awk '{print $2}' | xargs kill
    echo "stop job success!"
fi
