package user

import (
	secret "abs/models/secret"
	"abs/pkg/logging"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	UserId               string        `json:"user_id"`
	UserIp               string        `json:"user_ip"`
	UserName             string        `json:"user_name"`
	UserDnsPreventDomain string        `json:"user_dns_prevent_domain"`
	UserPrice            float64       `json:"user_price"`
	PreventSwitch        int           `json:"prevent_switch"`
	PreventInfo          []UserPrevent `json:"prevent_info"`
}

type UserPrevent struct {
	PreventName string `json:"prevent_name"`
	PreventType int    `json:"prevent_type"`
	PreventNum  int    `json:"prevent_num"`
}

//获取用户id并注册
func (u *User) GetUserOnlyId() *User {

	user_id := uuid.NewV4()
	u.UserId = user_id.String()
	u.UserPrice = 80
	u.UserDnsPreventDomain = "https://" + u.UserId + ".privacy.prisecurity.com/dns-query"
	secret.RegisterUser(u.UserId, u.UserDnsPreventDomain)
	return u
}

//获取用户信息
func (u *User) GetUserInfo() (*User, error) {

	ui, err := secret.GetUserInfo(u.UserId, u.UserIp)
	if err != nil || ui == nil {
		logging.Error(err)
		return nil, err
	}
	//用户基本信息赋值
	u.UserId = ui.UserId
	u.UserName = ui.UserName
	u.UserIp = ui.UserIp
	u.UserPrice = ui.UserPrice
	u.UserDnsPreventDomain = ui.UserDnsPreventDomain
	//获取拦截类型
	d, err := secret.GetAllDomainType()
	//获取用户拦截信息
	pi, err := secret.GetPreventCountByUserId(u.UserId, u.UserIp)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	for _, v := range pi {
		u.PreventInfo = append(u.PreventInfo, UserPrevent{PreventName: d[v.PreventName].DomainName, PreventNum: v.PreventNum, PreventType: v.PreventName})
	}
	return u, nil
}
