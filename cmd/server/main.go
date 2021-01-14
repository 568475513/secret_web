package server

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fvbock/endless"

	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/RichardKnop/machinery/v2"

	"abs/pkg/conf"
	"abs/cmd/server/routers"
)

// 启动服务
func main() {
	gin.SetMode(conf.ServerConf.RunMode)
	// 路由初始化
	routersInit := routers.InitRouter()
	// 允许读取的最大时间
	readTimeout := conf.ServerConf.ReadTimeout * time.Second
	// 允许写入的最大时间
	writeTimeout := conf.ServerConf.WriteTimeout * time.Second
	// 端口
	endPoint := fmt.Sprintf(":%d", conf.ServerConf.HttpPort)
	// 请求头的最大字节数
	maxHeaderBytes := 1 << 20

	// 本地调试直接这么用
	//server := &http.Server{
	//	Addr:           endPoint,
	//	Handler:        routersInit,
	//	ReadTimeout:    readTimeout,
	//	WriteTimeout:   writeTimeout,
	//	MaxHeaderBytes: maxHeaderBytes,
	//}
	//err := server.ListenAndServe()

	// 平滑启动，须编译环境【生产环境用下面的代码】
	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	endless.DefaultReadTimeOut = readTimeout
	endless.DefaultWriteTimeOut = writeTimeout
	endless.DefaultMaxHeaderBytes = maxHeaderBytes
	server := endless.NewServer(endPoint, routersInit)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d, Server Listening %s", syscall.Getpid(), endPoint)
	}
	err := server.ListenAndServe()

	// Server Start Error!!!
	if err != nil {
		log.Printf("Server err: %v", err)
	}
	// 关闭连接，这一步后期再考虑
	// models.CloseDB()
	log.Printf("[info] Start Http Server Listening %s", endPoint)
}

func startServer() (*machinery.Server, error) {
	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	// Create server instance
	broker := redisbroker.New(cnf, "localhost:6379", "", "", 0)
	backend := redisbackend.New(cnf, "localhost:6379", "", "", 0)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)

	// Register tasks
	tasks := map[string]interface{}{
		"add":               exampletasks.Add,
		"multiply":          exampletasks.Multiply,
		"sum_ints":          exampletasks.SumInts,
		"sum_floats":        exampletasks.SumFloats,
		"concat":            exampletasks.Concat,
		"split":             exampletasks.Split,
		"panic_task":        exampletasks.PanicTask,
		"long_running_task": exampletasks.LongRunningTask,
	}

	return server, server.RegisterTasks(tasks)
}