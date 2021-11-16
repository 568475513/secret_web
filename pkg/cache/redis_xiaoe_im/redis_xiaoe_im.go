package redis_xiaoe_im

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var XiaoEImRedisConn *redis.Pool

const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 36
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 600
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 180 * time.Second
)

// Setup XiaoEIm Initialize the Redis instance
func Init() error {
	if XiaoEImRedisConn == nil {
		XiaoEImRedisConn = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_XIAOEIM2_RW_HOST"), os.Getenv("REDIS_XIAOEIM2_RW_PORT")))
				if err != nil {
					return nil, err
				}
				if os.Getenv("REDIS_XIAOEIM2_RW_PASSWORD") != "" {
					if _, err := c.Do("AUTH", os.Getenv("REDIS_XIAOEIM2_RW_PASSWORD")); err != nil {
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

// 快捷查询方法
// 带key和database
func Get(key string, database string) ([]byte, error) {
	conn := XiaoEImRedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("SELECT", database)
	if err != nil {
		return nil, err
	}
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// 请记得释放回连接池
// 获取小鹅im redis主连接
func GetConn() (redis.Conn, error) {
	conn := XiaoEImRedisConn.Get()

	database := os.Getenv("REDIS_XIAOEIM2_DATABASE")
	if database != "" {
		_, err := conn.Do("SELECT", database)
		if err != nil {
			return conn, err
		}
	}

	return conn, nil
}
