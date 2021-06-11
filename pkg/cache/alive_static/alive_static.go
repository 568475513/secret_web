package alive_static

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var AliveStaticRedisConn *redis.Pool

const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 16
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 400
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 180 * time.Second
)

// Setup Static Initialize the Redis instance
func Init() error {
	AliveStaticRedisConn = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_ALIVESTATIC_RW_HOST"), os.Getenv("REDIS_ALIVESTATIC_RW_PORT")))
			if err != nil {
				return nil, err
			}
			if os.Getenv("REDIS_ALIVESTATIC_RW_PASSWORD") != "" {
				if _, err := c.Do("AUTH", os.Getenv("REDIS_ALIVESTATIC_RW_PASSWORD")); err != nil {
					c.Close()
					return nil, err
				}
			}
			// 设定默认数据库[但是放入连接池里面不能重置DataBase]
			if os.Getenv("ALIVE_STATIC_REDIS_DATABASE") != "" {
				if _, err := c.Do("SELECT", os.Getenv("ALIVE_STATIC_REDIS_DATABASE")); err != nil {
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
	return nil
}

func GetStaticRedisCon() (redis.Conn, error) {

	conn := AliveStaticRedisConn.Get()
	_, err := conn.Do("SELECT", 6)
	if err != nil {
		return conn, err
	}
	return conn, err
}

// HSetNx a key/value
func HsetNxString(key, hashKey string, data interface{}, time int) error {
	conn := AliveStaticRedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("HSETNX", key, hashKey, data)
	if err != nil {
		return err
	}

	if time != 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}
