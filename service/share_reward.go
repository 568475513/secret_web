package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"abs/pkg/app"
	e "abs/pkg/enums"
)

type ShareRewardService struct {
	AppId        string `json:"app_id"`
	ResourceId   string `json:"resource_id"`
	ResourceType int    `json:"resource_type"`
	ShareUserId  string `json:"share_user_id"`
	UserId       string `json:"user_id"`
	ShareType    int    `json:"share_type"`
}

const (
	//助力数据入队列接口
	shareRewardInsertListUrl = "_alive/share_reward/xe.sharereward.assists.add/1.0.0"
	//超时设置ms
	shareRewardInsertListTimeOut = 1000
)

// RequestShareRewardToList 助力数据入队列
func (srs *ShareRewardService) RequestShareRewardToList() (app.Response, error) {
	var result app.Response
	request := Post(fmt.Sprintf("%s/%s", os.Getenv("LB_PF_SHAREREWARD_IN"), shareRewardInsertListUrl))
	request.SetParams(srs)
	request.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	request.SetTimeout(shareRewardInsertListTimeOut * time.Millisecond)
	err := request.ToJSON(&result)
	if err != nil {
		return result, err
	}
	if result.Code != e.SUCCESS {
		return result, errors.New(
			fmt.Sprintf("请求助力数据入队列错误：%s[appId:%s,aliveId:%s,shareUserId:%s,userId:%s]",
				result.Msg,
				srs.AppId,
				srs.ResourceId,
				srs.ShareUserId,
				srs.UserId))
	}
	return result, nil
}
