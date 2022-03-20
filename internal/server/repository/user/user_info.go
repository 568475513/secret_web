package user

import (
	secret "abs/models/secret"
	"abs/pkg/logging"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	UserId        string           `json:"user_id"`
	UserName      string           `json:"user_name"`
	PreventSwitch int              `json:"prevent_switch"`
	PreventInfo   []secret.Prevent `json:"prevent_info"`
}

//获取用户id并注册
func (u *User) GetUserOnlyId() *User {

	user_id := uuid.NewV4()
	u.UserId = user_id.String()
	secret.RegisterUser(u.UserId)
	return u
}

//获取用户信息
func (u *User) GetUserInfo() (*User, error) {

	ui, err := secret.GetUserInfo(u.UserId)
	if err != nil || ui == nil {
		logging.Error(err)
		return nil, err
	}
	//用户基本信息赋值
	u.UserId = ui.UserId
	u.UserName = ui.UserName

	//获取用户拦截信息
	pi, err := secret.GetPreventCountByUserId(u.UserId)
	if err != nil || pi == nil {
		logging.Error(err)
		return nil, err
	}
	for _, v := range pi {
		u.PreventInfo = append(u.PreventInfo, v)
	}
	return u, nil
}
