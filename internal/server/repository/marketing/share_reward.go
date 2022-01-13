package marketing

import (
	"fmt"

	"abs/pkg/enums"
	"abs/pkg/logging"
	"abs/service"
)

type ShareReward struct {
	AppId       string
	ResourceId  string
	ShareUserId string
	UserId      string
}

// AddToShareRewardList 助力数据入队列
func (sr *ShareReward) AddToShareRewardList() {
	var err error

	s := service.ShareRewardService{
		AppId:        sr.AppId,
		ResourceId:   sr.ResourceId,
		ResourceType: enums.ResourceTypeLive,
		ShareUserId:  sr.ShareUserId,
		UserId:       sr.UserId,
		ShareType:    5,
	}
	_, err = s.RequestShareRewardToList()
	if err != nil {
		logging.Error(fmt.Sprintf("request xe.sharereward.assists.add fails: %s, app_id: %s, alive_id: %s, share_user_id: %s, user_id: %s",
			err.Error(), s.AppId, s.ResourceId, s.ShareUserId, s.UserId))
	}
}
