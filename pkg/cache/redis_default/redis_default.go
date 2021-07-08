package redis_default

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

// 新的代码级灰度
var redisDefaultConn *redis.Pool

const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 24
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 360
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 180 * time.Second
)

// Setup Gary Initialize the Redis instance
func Init() error {
	// 不可重复生成
	if redisDefaultConn == nil {
		redisDefaultConn = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_LIVECLUSTER_RW_HOST"), os.Getenv("REDIS_LIVECLUSTER_RW_PORT")))
				if err != nil {
					return nil, err
				}
				if os.Getenv("REDIS_LIVECLUSTER_RW_PASSWORD") != "" {
					if _, err := c.Do("AUTH", os.Getenv("REDIS_LIVECLUSTER_RW_PASSWORD")); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}
	return nil
}

// 获取直播IM消息【database = 4】
func GetLiveInfoConn() (redis.Conn, error) {
	conn := redisDefaultConn.Get()
	_, err := conn.Do("SELECT", 7)
	if err != nil {
		return conn, err
	}
	return conn, nil
}
