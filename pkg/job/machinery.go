package job

import (
	"fmt"
	"os"
	"strconv"

	backendsiface "github.com/RichardKnop/machinery/v1/backends/iface"
	redisbackend "github.com/RichardKnop/machinery/v1/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v1/brokers/redis"
	"github.com/RichardKnop/machinery/v1/config"
	lockiface "github.com/RichardKnop/machinery/v1/locks/iface"
	"github.com/RichardKnop/machinery/v2"
)

var Machinery *machinery.Server

// Machinery队列服务
func MachineryStartServer(queue string) {
	fmt.Println(">开始启动Machinery异步队列服务...")
	// 队列消息过期时间
	rExpirein, _ := strconv.Atoi(os.Getenv("JOB_RESULTSEXPIREIN"))
	cnf := &config.Config{
		DefaultQueue: queue,
		ResultsExpireIn: rExpirein,
		Redis: &config.RedisConfig{
			MaxIdle:                16,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	// Create server instance
	var lock lockiface.Lock
	var backend backendsiface.Backend
	// 入列队列
	broker := redisbroker.New(
		cnf, 
		fmt.Sprintf("%s:%s", os.Getenv("REDIS_LIVEBUSINESS_RW_HOST"), 
		os.Getenv("REDIS_LIVEBUSINESS_RW_PORT")), 
		os.Getenv("REDIS_LIVEBUSINESS_RW_PASSWORD"), 
		"", 
		14,
	)
	// 结果队列
	backend = redisbackend.New(
		cnf, 
		fmt.Sprintf("%s:%s", os.Getenv("REDIS_LIVEBUSINESS_RW_HOST"), 
		os.Getenv("REDIS_LIVEBUSINESS_RW_PORT")), 
		os.Getenv("REDIS_LIVEBUSINESS_RW_PASSWORD"), 
		"", 
		15,
	)
	Machinery = machinery.NewServer(cnf, broker, backend, lock)
	fmt.Println(">>>Machinery异步队列服务启动完成")
}