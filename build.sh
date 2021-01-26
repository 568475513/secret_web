#!/bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter -o absGo -ldflags "-w -s"
chmod +x ./absGo