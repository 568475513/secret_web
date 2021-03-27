package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"abs/models/user"
	"abs/pkg/app"
	e "abs/pkg/enums"
)

type UserService struct {
	AppId  string
	UserId string
}

const (
	// 用户基本信息接口
	userInfoUrl = "xe.user.user_info/1.0.0"
	// 超时设置ms
	userServiceTime = 500
)

type UserServiceResponse struct {
	app.Response
	Data user.User `json:"data"`
}

// 获取用户数据
func (userService *UserService) RequestUserInfo() (UserServiceResponse, error) {
	var result UserServiceResponse
	request := Post(fmt.Sprintf("%s%s", os.Getenv("LB_SP_USERSERVICE_IN"), userInfoUrl))
	request.SetParams(map[string]string{
		"gray_type": "0",
		"gray_id":   "1",
		"app_id":    userService.AppId,
		"user_id":   userService.UserId,
	})
	request.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	request.SetTimeout(userServiceTime * time.Millisecond)
	err := request.ToJSON(&result)
	if err != nil {
		return result, err
	}
	if result.Code != e.SUCCESS {
		return result, errors.New(fmt.Sprintf("请求用户数据错误：%s[appId:%s,userId:%s]", result.Msg, userService.AppId, userService.UserId))
	}
	return result, nil
}
