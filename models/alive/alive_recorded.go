package alive

import "github.com/jinzhu/gorm"

type AliveVodRetweet struct {
	AppId     string `json:"app_id"`
	AliveId   string `json:"alive_id"`
	TaskId    string `json:"task_id"`
	TaskState int    `json:"task_state"`
}

func GetRecordedRetweetTaskInfo(appId, aliveId, fields string) (*AliveVodRetweet, error) {
	var a AliveVodRetweet
	err := db.Table("t_alive_vod_retweet").Select(fields).Where("app_id=? and alive_id=? ", appId, aliveId).First(&a).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &a, nil
}
