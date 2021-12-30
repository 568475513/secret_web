package course

import (
	"abs/pkg/cache/redis_elive"
	"abs/pkg/logging"
	"encoding/json"
	"time"
)

// 空结构体
type EliveInfo struct{}

type AccessTimeTaskBean struct {
	appId           string
	aliveId         string
	accessTimestamp int64
}

const (
	// 更新最近查看时间的list key名
	accessTimeListKey = "elive_access_time_update_list"
)

// 将最近查看时间的更新任务丢进Redis队列异步处理
func (c *EliveInfo) UpdateAccessTimeToQueue(appId, aliveId, userId string, userType uint) {
	cur := time.Now().Format("2006-01-02 15:04:05")
	var role int
	if userType == 1 {
		// 讲师身份
		role = 2
	} else {
		// 学员身份
		role = 1
	}
	jsonData := map[string]interface{}{
		"app_id":      appId,
		"alive_id":    aliveId,
		"user_id":     userId,
		"access_time": cur,
		"role":        role,
	}
	jsonStrData, _ := json.Marshal(jsonData)
	conn, _ := redis_elive.GetEliveRedisConn()
	defer conn.Close()
	err := conn.PushToUpdateAccessTimeQueue(accessTimeListKey, jsonStrData)
	if err != nil {
		logging.Error(err)
	}
}