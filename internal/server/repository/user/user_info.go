package user

import (
	secret "abs/models/secret"
	"abs/pkg/logging"
	cache "github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

type User struct {
	UserId               string        `json:"user_id"`
	UserIp               string        `json:"user_ip"`
	UserName             string        `json:"user_name"`
	UserDnsPreventDomain string        `json:"user_dns_prevent_domain"`
	UserPrice            float64       `json:"user_price"`
	PreventSwitch        int           `json:"prevent_switch"`
	PreventInfo          []UserPrevent `json:"prevent_info"`
	RegisterId           string        `json:"register_id"`
	PreventWeekData      interface{}   `json:"prevent_week_data"`
}

type UserPrevent struct {
	PreventName string `json:"prevent_name"`
	PreventType int    `json:"prevent_type"`
	PreventNum  int    `json:"prevent_num"`
}

var Cache *cache.Cache

//获取用户id并注册
func (u *User) GetUserOnlyId() *User {

	user_id := uuid.NewV4()
	u.UserId = user_id.String()
	u.UserPrice = 80
	u.UserDnsPreventDomain = "https://" + u.UserId + ".privacy.prisecurity.com/dns-query"
	secret.RegisterUser(u.UserId, u.UserDnsPreventDomain, u.RegisterId)
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
	u.RegisterId = ui.RegisterId
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
	//t := time.Now().Weekday().String()
	//如果当天是周天则返回周报信息
	//if t == "Sunday" {
	var s []interface{}
	for _, v := range d {
		r, t := Cache.Get(u.UserId + "_" + strconv.Itoa(v.DomainType))
		if t == true && r != nil {
			s = append(s, r)
			Cache.Delete(u.UserId + "_" + strconv.Itoa(v.DomainType))
		}
	}
	u.PreventWeekData = s
	//}
	return u, nil
}

//获取用户周报数据
func (u *User) WeekGetUserData() (err error) {

	Cache = cache.New(5*time.Minute, 60*time.Second)
	// 获取用户id列表
	err, ids := secret.GetUserId()
	if err != nil {
		logging.Error(err)
	}
	//获取拦截类型
	d, err := secret.GetAllDomainType()
	for _, v := range ids {
		re, err := secret.SelectUserDataTime(v)
		if err != nil {
			logging.Error(err)
		}
		for _, v2 := range re[v] {
			dn := d[v2.DomainType].DomainName
			v2.DomainName = dn
			Cache.Set(v+"_"+strconv.Itoa(v2.DomainType), map[string]interface{}{"count": v2.Count, "name": v2.DomainName}, time.Hour*24)
		}
	}
	return
}
