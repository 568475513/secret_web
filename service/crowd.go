package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"abs/pkg/app"
	"abs/pkg/enums"
)

const (
	// 查看封禁状态接口
	GetUserBlackStatusUrl = "crowd_v2/get_black_status"
	// 查看封禁BuzServer版本
	GetUserBlackStatusUrlByBuz = "crowd/get_user_resource_auth"
	// 超时设置ms
	timeoutCrowd = 700
)

// 用户黑名单服务
type CrowdService struct {
	// br BaseRequest
	AppId  string
	UserId string
}

// 用户黑名单信息
type UserBlackInfo struct {
	PermissionVisit   int `json:"permission_visit"`   // 禁止访问
	PermissionComment int `json:"permission_comment"` // 禁止评论
	PermissionBuy     int `json:"permission_buy"`     // 禁止购买
}

type CrowdUserResponse struct {
	app.Response
	Data UserBlackInfo `json:"data"`
}

// 请求获取黑名单
// func (c *CrowdService) RequestCrowdServiceOld(methods string, route string, data map[string]interface{}) (result map[string]interface{}, err error) {
// 	bT := time.Now()
// 	done := make(chan bool)
// 	go func(data map[string]interface{}) {
// 		defer close(done)
// 		result, err = c.br.Request(methods, fmt.Sprintf("%s%s", os.Getenv("LB_CT_BUZSERVER_IN"), route), data)
// 	}(data)

// 	select {
// 	// 监听发送请求是否超时, 如果超时，则记录数据.
// 	case <-time.After(timeoutCrowd * time.Millisecond):
// 		err = errors.New(fmt.Sprintf("-- RequestCrowdService Timeout Url：%s%s", os.Getenv("LB_CT_BUZSERVER_IN"), route))
// 		fmt.Printf("黑名单请求url: %s%s[RequestCrowdService Timeout] - %d\n", os.Getenv("LB_CT_BUZSERVER_IN"), route, timeoutCrowd)
// 		c.br.RecordFailData(data)
// 		return
// 	// 数据在规定时间内已经请求业务侧.
// 	case <-done:
// 		eT := time.Since(bT)
// 		fmt.Printf("黑名单请求url: %s%s[RequestCrowdService cos time] - %s\n", os.Getenv("LB_CT_BUZSERVER_IN"), route, eT)
// 		return
// 	}
// }

// 发送请求用户黑名单
func (c *CrowdService) GetCrowdUserInfo() (UserBlackInfo, error) {
	var result CrowdUserResponse
	request := Post(fmt.Sprintf("%s%s", os.Getenv("LB_CT_BUZSERVER_IN"), GetUserBlackStatusUrlByBuz))
	request.SetParams(map[string]interface{}{
		"user_id": c.UserId,
		"app_id":  c.AppId,
	})
	request.SetHeader("Content-Type", "application/json")
	request.SetTimeout(timeoutCrowd * time.Millisecond)
	err := request.ToJSON(&result)
	if err != nil {
		return UserBlackInfo{}, err
	}
	if result.Code != enums.SUCCESS {
		return UserBlackInfo{}, errors.New(fmt.Sprintf("请求用户黑名单错误：%s", result.Msg))
	}
	return result.Data, nil
}
