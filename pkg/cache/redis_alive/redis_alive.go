package redis_alive

import (
	"fmt"
	"time"
	"os"

	"github.com/gomodule/redigo/redis"
)

var AliveRedisConn *redis.Pool

const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 32
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 5000
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 180 * time.Second
)

// Setup Alive Initialize the Redis instance
func Init() error {
	if AliveRedisConn == nil {
		AliveRedisConn = &redis.Pool{
			MaxIdle: maxIdle,
			MaxActive: maxActive,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_LIVEBUSINESS_RW_HOST"), os.Getenv("REDIS_LIVEBUSINESS_RW_PORT")))
				if err != nil {
					return nil, err
				}
				if os.Getenv("REDIS_LIVEBUSINESS_RW_PASSWORD") != "" {
					if _, err := c.Do("AUTH", os.Getenv("REDIS_LIVEBUSINESS_RW_PASSWORD")); err != nil {
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
	conn := AliveRedisConn.Get()
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
// 获取直播业务主连接【database = 4】
func GetLiveBusinessConn() (redis.Conn, error) {
	conn := AliveRedisConn.Get()

	database := os.Getenv("REDIS_LIVEBUSINESS_ALIVE_DATABASE")
	if database != "" {
		_, err := conn.Do("SELECT", database)
		if err != nil {
			return conn, err
		}
	}

	return conn, nil
}

// 获取直播业务营销活动相关连接【database = 1】
func GetLiveMarketingConn() (redis.Conn, error) {
	conn := AliveRedisConn.Get()

	_, err := conn.Do("SELECT", 1)
	if err != nil {
		return conn, err
	}

	return conn, nil
}

// 获取直播业务评论互动相关连接【database = 2】
func GetLiveInteractConn() (redis.Conn, error) {
	conn := AliveRedisConn.Get()

	_, err := conn.Do("SELECT", 2)
	if err != nil {
		return conn, err
	}

	return conn, nil
}

// 获取直播次级业务连接【database = 5】
func GetSubBusinessConn() (redis.Conn, error) {
	conn := AliveRedisConn.Get()

	_, err := conn.Do("SELECT", 5)
	if err != nil {
		return conn, err
	}

	return conn, nil
}