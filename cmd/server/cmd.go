package server

import (
	"github.com/spf13/cobra"

	"abs/internal/server/rules"
	"abs/models"
	"abs/pkg/cache"
	"abs/pkg/conf"
	"abs/pkg/job"
	"abs/pkg/logging"
)

var port int
var env string
var queue string

// Cmd run http server
var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Run absgo server",
	Long:  `Run absgo server`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化各项服务
		initStep()
		// 启动服务
		main()
	},
}

// init server cmd
func init() {
	Cmd.Flags().IntVar(&port, "port", 9090, "listen port")
	Cmd.Flags().StringVar(&env, "env", "local", "conf environmental science")
	Cmd.Flags().StringVar(&queue, "queue", "abs_machinery_tasks", "job default queue")
}

// initStep server
func initStep() {
	// 初始化各项服务
	// 配置加载
	conf.Init(env)
	// 设定真实端口
	conf.ServerConf.HttpPort = port
	// 自定义日志
	logging.InitLog()
	// 调用链路日志
	logging.InitZipkin()
	// 请求以及错误日志
	logging.InitEs()
	// 启动各数据库连接
	models.Init()
	// 启动相关redis
	cache.Init()
	// 初始化验证器
	rules.InitVali()
	// 初始化队列服务
	job.MachineryStartServer(queue)
}