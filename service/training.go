package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"abs/pkg/app"
	e "abs/pkg/enums"
)

type TrainingReq struct {
	AppId  string
	UserId string
}

const (
	trainingAuthCheck = "/privilege/employee/auth.check"
	// 超时设置ms
	trainingAuthTimeOut = 1000
)

type TrainingResp struct {
	app.Response
	Data TrainingSer `json:"data"`
}

type TrainingSer struct {
	Isauth bool `json:"is_auth"`
}

// 配置中心获取配置，传入fields中需要获取的字段。
func (c *TrainingReq) AuthCheck() (TrainingResp, error) {
	var result TrainingResp
	request := Post(fmt.Sprintf("%s%s?app_id=%s", os.Getenv("LB_PF_TRAININGAPI_IN"), trainingAuthCheck, c.AppId))
	request.SetHeader("Content-Type", "application/json")
	request.SetHeader("App-Id", c.AppId)
	request.SetHeader("User-Id", c.UserId)
	request.SetTimeout(trainingAuthTimeOut * time.Millisecond)
	err := request.ToJSON2(&result)
	if err != nil {
		return TrainingResp{}, err
	}
	if result.Code != e.SUCCESS {
		return TrainingResp{}, errors.New(fmt.Sprintf("请求训练营鉴权错误：%s", result.Msg))
	}
	return result, nil
}
