//获取直播课程列表逻辑
package course

import (
	//内部包
	"abs/models/alive"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"

	//第三方包
	"github.com/gomodule/redigo/redis"

	//系统标准包
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ListInfo struct {
	AppId            string
	UserId           string
	UniversalUnionId string
}

const (
	aliveListByTimeCacheKey  = "alive_list_by_time:%s:%s" //根据直播开时间直播列表缓存key
	aliveListByTimeCacheTime = "30"                       //根据直播开时间直播列表缓存时间，单位s
)

//根据直播开时间获取直播列表
func (l *ListInfo) GetALiveListByTime(startTime time.Time, endTime time.Time) ([]*alive.Alive, error) {
	var (
		err          error
		aliveList    []*alive.Alive
		startTimeStr = startTime.Format(util.TIME_LAYOUT)
		endTimeStr   = endTime.Format(util.TIME_LAYOUT)
	)
	conn, _ := redis_alive.GetSubBusinessConn()
	defer conn.Close()

	//时间范围限定为3天以内，防止查询范围太大导致慢查询
	timeRange := endTime.Unix() - startTime.Unix()
	if timeRange <= 0 || timeRange > 3600*24*3 {
		return nil, errors.New("startTime or endTime error")
	}

	//去缓存读取数据
	tempStr := strings.Join([]string{startTimeStr, endTimeStr}, "")
	md5Str := fmt.Sprintf("%x", md5.Sum([]byte(tempStr)))
	cacheKey := fmt.Sprintf(aliveListByTimeCacheKey, l.AppId, md5Str)
	cacheData, err := redis.Bytes(conn.Do("get", cacheKey))
	if err == nil {
		if err = util.JsonDecode(cacheData, &aliveList); err != nil {
			logging.Error(err)
			logging.LogToEs("GetALiveListByTime", aliveList)
		}
		return aliveList, nil
	}

	//无缓存则读数据库
	aliveList, err = alive.GetAliveListByZbStartTime(l.AppId, startTimeStr, endTimeStr, []string{"*"})
	if err != nil {
		return nil, err
	}

	//写入缓存
	if value, err := util.JsonEncode(aliveList); err == nil {
		if _, err = conn.Do("SET", cacheKey, value, "EX", aliveListByTimeCacheTime); err != nil {
			logging.Error(err)
		}
	}
	return aliveList, nil
}

//从给定的直播列表筛选出当前用户已订阅的直播，并且按直播开始日期分组
func (l *ListInfo) GetSubscribedALiveList(aliveList []*alive.Alive) map[string][]*alive.Alive {
	var (
		result     = make(map[string][]*alive.Alive)
		aliveIds   []string
		filterList = make(map[string]*alive.Alive)
	)
	for _, aliveInfo := range aliveList {
		aliveIds = append(aliveIds, aliveInfo.Id)
		filterList[aliveInfo.Id] = aliveInfo
	}
	subscribedAliveIds, err := service.GetMultipleSubscribe(l.AppId, l.UniversalUnionId, aliveIds)
	if err == nil && len(subscribedAliveIds) > 0 {
		for _, aliveId := range subscribedAliveIds {
			aliveInfo, ok := filterList[aliveId]
			if ok {
				zbStartDate := aliveInfo.ZbStartAt.Time.Format(util.DATE_LAYOUT)
				result[zbStartDate] = append(result[zbStartDate], aliveInfo)
			}
		}
	} else if err != nil {
		logging.Error(err)
	}
	return result
}
