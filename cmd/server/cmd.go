package server

import (
	"github.com/spf13/cobra"

	"abs/models"
	"abs/pkg/cache"
	"abs/pkg/conf"
	"abs/pkg/logging"
	"abs/internal/server/rules"
)

var port int
var isInternal bool
var isManage bool

// Cmd run http server
var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Run absgo server",
	Long:  `Run absgo server`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

// init server cmd
func init() {
	Cmd.Flags().IntVar(&port, "port", 8080, "listen port")
	Cmd.Flags().BoolVar(&isInternal, "internal", false, "internal service")
	Cmd.Flags().BoolVar(&isManage, "manage", false, "manage service")

	// 初始化各项服务
	// 配置加载
	conf.Init()
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
}