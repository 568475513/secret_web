package redis_elive

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"strconv"
	"time"
)

var EliveRedisConnPool *redis.Pool

type EliveRedisConn struct {
	conn                redis.Conn
	accessTimeListLimit int // 限制队列长度，0表示不限制
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
	limit, _ := strconv.Atoi(os.Getenv("ACCESS_TIME_LIST_LIMIT"))
	eliveRedisConn := &EliveRedisConn{
		conn:                conn,
		accessTimeListLimit: limit,
	}
	return eliveRedisConn, nil
}

// 丢进处理最近查看时间的队列
func (c *EliveRedisConn) PushToUpdateAccessTimeQueue(key string, data []byte) error {
	if c.accessTimeListLimit == 0 {
		return nil
	}
	// 达到数量限制就停止丢队列
	listLen, err := redis.Int(c.conn.Do("LLEN", key))
	if listLen >= c.accessTimeListLimit {
		errStr := fmt.Sprintf("curent access time list length is out of limit, curent length value is : %d", listLen)
		return fmt.Errorf(errStr)
	}
	_, err = c.conn.Do("LPUSH", key, data)
	return err
}

func (c *EliveRedisConn) Close() {
	c.conn.Close()
}
