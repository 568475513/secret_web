package cache

import (
	"abs/pkg/cache/redis_im"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"

	"abs/pkg/cache/alive_static"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/cache/redis_gray"
)

var RedisConn *redis.Pool

/* 区分两种使用场景：
- 高频调用的场景，需要尽量压榨redis的性能：
调高MaxIdle的大小，该数目小于maxActive，由于作为一个缓冲区一样的存在，扩大缓冲区自然没有问题
调高MaxActive，考虑到服务端的支持上限，尽量调高
IdleTimeout由于是高频使用场景，设置短一点也无所谓，需要注意的一点是MaxIdle设置的长了，队列中的过期连接可能会增多，这个时候IdleTimeout也要相应变化
- 低频调用的场景，调用量远未达到redis的负载，稳定性为重：
MaxIdle可以设置的小一些
IdleTimeout相应地设置小一些
MaxActive随意，够用就好，容易检测到异常 */
const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 30
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 1000
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 200 * time.Millisecond
)

// Setup Initialize the Redis instance
func Init() {
	fmt.Println(">开始初始化缓存连接池...")
	// 默认集群[暂时没用就别开！！！]
	// if err := defaultInit(); err != nil {
	// 	log.Fatal(err)
	// }

	// 直播redis
	if err := redis_alive.Init(); err != nil {
		log.Fatal(err)
	}
	// 直播静态redis
	if err := alive_static.Init(); err != nil {
		log.Fatal(err)
	}
	// 不知道新旧的请问abner!!!
	// 灰度控制【新】
	if err := redis_gray.Init(); err != nil {
		log.Fatal(err)
	}
	// 灰度控制【旧】
	if err := redis_gray.InitOldGary(); err != nil {
		log.Fatal(err)
	}
	// 灰度控制【直播专用】
	if err := redis_gray.InitSpecialGary(); err != nil {
		log.Fatal(err)
	}
	// IM【直播专用】
	if err := redis_im.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(">>>初始化缓存连接池完成")
}

// Job Cmd Setup Initialize the Redis instance
func InitJob() {
	fmt.Println(">开始初始化Job缓存连接池...")
	// 直播redis
	if err := redis_alive.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(">>>初始化Job缓存连接池完成")
}

// 集群redis
func defaultInit() error {
	RedisConn = &redis.Pool{
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
	return nil
}

// 以下是示例使用 =========================================
// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
