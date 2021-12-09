package redis_elive

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

var EliveRedisConnPool *redis.Pool

type EliveRedisConn struct {
	conn redis.Conn
}

const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 36
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 600
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 180 * time.Second
)

// Setup Alive Initialize the Redis instance
func Init() error {
	if EliveRedisConnPool == nil {
		EliveRedisConnPool = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_EALIVE_RW_HOST"), os.Getenv("REDIS_EALIVE_RW_PORT")))
				if err != nil {
					return nil, err
				}
				if os.Getenv("REDIS_EALIVE_RW_PASSWORD") != "" {
					if _, err := c.Do("AUTH", os.Getenv("REDIS_EALIVE_RW_PASSWORD")); err != nil {
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

func GetEliveRedisConn() (*EliveRedisConn, error) {
	conn := EliveRedisConnPool.Get()

	eliveRedisConn := &EliveRedisConn{
		conn: conn,
	}
	return eliveRedisConn, nil
}

// 丢进处理最近查看时间的队列
func (conn *EliveRedisConn) PushToUpdateAccessTimeQueue(key string, data []byte) error {
	_, err := conn.conn.Do("LPUSH", key, data)
	return err
}

func (c *EliveRedisConn) Close() {
	c.conn.Close()
}
