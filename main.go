package main

import (
	"github.com/spf13/cobra"
	// _ "go.uber.org/automaxprocs" // 根据容器配额设置maxprocs【如果是容器启动请打开注释！！！】

	"abs/cmd/server"
)

// 反模式设计启动
func main() {
	root := cobra.Command{Use: "abs-go"}
	root.AddCommand(
		server.Cmd,
	)

	// 执行
	root.Execute()
}
