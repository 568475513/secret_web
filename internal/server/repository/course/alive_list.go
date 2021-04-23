//获取直播课程列表逻辑
package course

import (
	//内部包
	"abs/models/alive"
	"abs/pkg/logging"
	"abs/pkg/util"
	"abs/service"

	//系统标准包
	"errors"
	"time"
)

type ListInfo struct {
	AppId            string
	UserId           string
	UniversalUnionId string
}

//根据直播开时间获取直播列表
func (l *ListInfo) GetALiveListByTime(startTime time.Time, endTime time.Time) ([]*alive.Alive, error) {
	var (
		err       error
		aliveList []*alive.Alive
	)
	//时间范围限定为3天以内，防止查询范围太大导致慢查询
	timeRange := endTime.Unix() - startTime.Unix()
	if timeRange <= 0 || timeRange > 3600*24*3 {
		return nil, errors.New("startTime or endTime error")
	}
	aliveList, err = alive.GetAliveListByZbStartTime(l.AppId,
		startTime.Format("2006-01-02 15:04:05"),
		endTime.Format("2006-01-02 15:04:05"),
		[]string{"*"})
	if err != nil {
		return nil, err
	}
	return aliveList, nil
}

//从给定的直播列表筛选出当前用户已订阅的直播
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
