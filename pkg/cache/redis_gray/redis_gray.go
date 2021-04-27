package redis_gray

import (
	"abs/pkg/logging"
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

// 新的代码级灰度
var redisGrayConn *redis.Pool

// 老的灰度redis
var redisGrayOldConn *redis.Pool

// 直播专用灰度redis
var redisGraySpecialConn *redis.Pool

const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 24
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 250
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 180 * time.Second
)

// Setup Gary Initialize the Redis instance
func Init() error {
	// 不可重复生成
	if redisGrayConn == nil {
		redisGrayConn = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_CODEGRAY_RW_HOST"), os.Getenv("REDIS_CODEGRAY_RW_PORT")))
				if err != nil {
					return nil, err
				}
				if os.Getenv("REDIS_CODEGRAY_RW_PASSWORD") != "" {
					if _, err := c.Do("AUTH", os.Getenv("REDIS_CODEGRAY_RW_PASSWORD")); err != nil {
						c.Close()
						return nil, err
					}
				}
				// 设定默认数据库，默认为0不用设置
				// if os.Getenv("xxx") != "" {
				// 	if _, err := c.Do("SELECT", os.Getenv("xxx")); err != nil {
				// 		c.Close()
				// 		return nil, err
				// 	}
				// }
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

// Setup Gary Old Initialize the Redis instance
func InitOldGary() error {
	// 不可重复生成
	if redisGrayOldConn == nil {
		redisGrayOldConn = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_DATA_RW_HOST"), os.Getenv("REDIS_DATA_RW_PORT")))
				if err != nil {
					return nil, err
				}
				if os.Getenv("REDIS_DATA_RW_PASSWORD") != "" {
					if _, err := c.Do("AUTH", os.Getenv("REDIS_DATA_RW_PASSWORD")); err != nil {
						c.Close()
						return nil, err
					}
				}
				//设定默认数据库，默认为0不用设置
				if os.Getenv("REDIS_DATA_RW_DATABASE") != "" {
					if _, err := c.Do("SELECT", os.Getenv("REDIS_DATA_RW_DATABASE")); err != nil {
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

// 搞个直播专用灰度redis实例
func InitSpecialGary() error {
	// 不可重复生成
	if redisGraySpecialConn == nil {
		redisGraySpecialConn = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_ALIVECODEGRAY_RW_HOST"), os.Getenv("REDIS_ALIVECODEGRAY_RW_PORT")))
				if err != nil {
					return nil, err
				}
				if os.Getenv("REDIS_ALIVECODEGRAY_RW_PASSWORD") != "" {
					if _, err := c.Do("AUTH", os.Getenv("REDIS_ALIVECODEGRAY_RW_PASSWORD")); err != nil {
						c.Close()
						return nil, err
					}
				}
				//设定默认数据库，默认为0不用设置
				if os.Getenv("REDIS_ALIVECODEGRAY_RW_DATABASE") != "" {
					if _, err := c.Do("SELECT", os.Getenv("REDIS_ALIVECODEGRAY_RW_DATABASE")); err != nil {
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

// 判断店铺是否在灰度名单
// redis 实例需要根据O端配置进行设置 $gray_id 灰度名单项目id请查看O端设置
func InGrayShopNew(garyKey, appId string) bool {
	if garyKey == "" || appId == "" {
		return false
	}
	conn := redisGrayConn.Get()
	defer conn.Close()

	// 全网打开
	if replyAll, _ := redis.Bool(conn.Do("SISMEMBER", garyKey, "*")); replyAll {
		return replyAll
	}

	// 指定查询
	reply, err := redis.Bool(conn.Do("SISMEMBER", garyKey, appId))
	if err != nil {
		logging.Error(fmt.Sprintf("注意！！！InGrayShopNew有错误：%s", err.Error()))
		return false
	}

	return reply
}

// 判断店铺是否在灰度名单【旧】
func InGrayShop(garyKey, appId string) bool {
	if garyKey == "" || appId == "" {
		return false
	}
	conn := redisGrayOldConn.Get()
	defer conn.Close()

	// 全网打开
	if replyAll, _ := redis.Bool(conn.Do("SISMEMBER", garyKey, "*")); replyAll {
		return replyAll
	}

	// 指定查询
	reply, err := redis.Bool(conn.Do("SISMEMBER", garyKey, appId))
	if err != nil {
		logging.Error(fmt.Sprintf("注意！！！InGrayShop有错误：%s", err.Error()))
		return false
	}

	return reply
}

// 判断店铺是否在灰度名单【直播专用】
func InGrayShopSpecial(garyKey, appId string) bool {
	if garyKey == "" || appId == "" {
		return false
	}
	conn := redisGraySpecialConn.Get()
	defer conn.Close()

	// 全网打开
	if replyAll, _ := redis.Bool(conn.Do("SISMEMBER", garyKey, "*")); replyAll {
		return replyAll
	}

	// 指定查询
	reply, err := redis.Bool(conn.Do("SISMEMBER", garyKey, appId))
	if err != nil {
		logging.Error(fmt.Sprintf("注意！！！InGrayShopSpecial有错误：%s", err.Error()))
		return false
	}

	return reply
}
