#!/bin/bash
echo "注意！:【打包需要GO环境版本 > 1.15.*】"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter -o absGo -ldflags "-w -s"
chmod +x ./absGo