# 构建Docker镜像使用
# 任何问题，联系基础架构@sparklizhang

FROM golang:1.17 as builder

WORKDIR /app

ENV  GOPROXY=https://goproxy.cn,http://goproxy.xiaoe-tools.com,direct GO111MODULE=on GOOS=linux GOARCH=amd64

# 缓存处理，如gomod gosum未更改则不会重新拉取
COPY go.mod go.mod
COPY go.sum go.sum
RUN  go mod download

# 清理本地二进制包
RUN  go clean

COPY . .

# 构建二进制文件命令,替换为自身程序的构建命令
RUN  go build -tags=jsoniter -o main -ldflags "-w -s"
RUN  go build -tags=jsoniter -o job  -ldflags "-w -s" ./cmd/job/cmd.go
RUN chmod 755 main job
# 为了缩小镜像体积，做分层处理
FROM centos:7

WORKDIR /app

COPY --from=builder /app/main ./main
COPY --from=builder /app/job ./job

# 启动命令，多行参数使用,隔开
# 如原启动命令 ./main run -p 8888，则以下启动命令为 ENTRYPOINT ["./main","run","-p","8888"]