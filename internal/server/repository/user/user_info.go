package user

import (
	secret "abs/models/secret"
	"abs/pkg/logging"
	"fmt"
	cache "github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	jpushclient "github.com/ylywyn/jpush-api-go-client"
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

type UserV2 struct {
	UserId               string `json:"user_id"`
	UserIp               string `json:"user_ip"`
	UserName             string `json:"user_name"`
	UserDnsPreventDomain string `json:"user_dns_prevent_domain"`
	RegisterId           string `json:"register_id"`
	UserConfig           UC     `json:"user_config"`
}

type UC struct {
	IsBuy         int `json:"is_buy"`
	IsBusMonitor  int `json:"is_bus_monitor"`
	IsLargeData   int `json:"is_large_data"`
	IsSpy         int `json:"is_spy"`
	IsCollectInfo int `json:"is_collect_info"`
}

type UserPrevent struct {
	PreventName string `json:"prevent_name"`
	PreventType int    `json:"prevent_type"`
	PreventNum  int    `json:"prevent_num"`
}

type UserPrice struct {
	Price int
	Count int
}

const appKey = "11066b2bfdf825c774968dce"
const secretKey = "d245c0ece98da21888765fa6"

var Cache *cache.Cache

func init() {
	Cache = cache.New(5*time.Minute, 60*time.Second)
}

//获取用户id并注册
func (u *User) GetUserOnlyId() *User {

	user_id := uuid.NewV4()
	u.UserId = user_id.String()
	u.UserPrice = 80.00000
	u.UserDnsPreventDomain = "https://privacy.prisecurity.com/dns-query/" + u.UserId
	secret.RegisterUser(u.UserId, u.UserDnsPreventDomain, u.RegisterId, u.UserPrice)
	return u
}

//获取用户id并注册
func (u *UserV2) GetUserOnlyId() *UserV2 {

	user_id := uuid.NewV4()
	u.UserId = user_id.String()
	u.UserDnsPreventDomain = "https://privacy.prisecurity.com/dns-query/" + u.UserId
	secret.RegisterUserV2(u.UserId, u.UserDnsPreventDomain, u.RegisterId)
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

	//判断用户是否存在registerId
	if u.RegisterId != "" && ui.RegisterId == "" {
		err := secret.UpdateUserRegisterId(u.UserId, u.RegisterId)
		if err != nil {
			logging.Error(err)
		}
	} else {
		u.RegisterId = ui.RegisterId
	}

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
	t := time.Now().Weekday().String()
	//如果当天是周天则返回周报信息
	if t == "Sunday" {
		var s []interface{}
		for _, v := range d {
			r, t := Cache.Get(u.UserId + "_" + strconv.Itoa(v.DomainType))
			if t == true && r != nil {
				s = append(s, r)
				Cache.Delete(u.UserId + "_" + strconv.Itoa(v.DomainType))
			}
		}
		u.PreventWeekData = s
	}
	return u, nil
}

//获取用户信息
func (u *UserV2) GetUserInfo() (*UserV2, error) {

	ui, err := secret.GetUserInfo(u.UserId, u.UserIp)
	if err != nil || ui == nil {
		logging.Error(err)
		return nil, err
	}
	//用户基本信息赋值
	u.UserId = ui.UserId
	u.UserName = ui.UserName
	u.UserIp = ui.UserIp

	//判断用户是否存在registerId
	if u.RegisterId != "" && ui.RegisterId == "" {
		err := secret.UpdateUserRegisterId(u.UserId, u.RegisterId)
		if err != nil {
			logging.Error(err)
		}
	} else {
		u.RegisterId = ui.RegisterId
	}

	uc, err := secret.GetUserConfig(u.UserId)
	if err != nil || ui == nil {
		logging.Error(err)
		return nil, err
	}
	u.UserConfig.IsBuy = uc.IsBuy
	u.UserConfig.IsBusMonitor = uc.IsBusMonitor
	u.UserConfig.IsLargeData = uc.IsLargeData
	u.UserConfig.IsSpy = uc.IsSpy
	u.UserConfig.IsCollectInfo = uc.IsCollectInfo
	return u, nil
}

//获取用户周报数据
func (u *User) WeekGetUserData() (err error) {

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
			dn := d[v2.DomainType].DomainName + "风险"
			v2.DomainName = dn
			Cache.Set(v+"_"+strconv.Itoa(v2.DomainType), map[string]interface{}{"count": v2.Count, "name": v2.DomainName}, time.Hour*24)
		}
	}
	return
}

//获取每日用户数据
func (u *User) GetUserPriceDay() (err error) {
	var (
		pf     jpushclient.Platform
		ad     jpushclient.Audience
		op     jpushclient.Option
		notice jpushclient.Notice
	)
	pf.Add(jpushclient.IOS)
	op.ApnsProduction = true
	// 获取用户id以及registerId列表
	err, ids := secret.GetUserIdAndRegisterID()
	if err != nil {
		logging.Error(err)
	}
	for k, v := range ids {
		count, err := secret.GetCountByUserId(k)
		if err != nil {
			logging.Error(err)
		}
		price := ((v.UserPrice - 50) / 50) * 100
		ad.SetID([]string{v.RegisterId})
		notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: fmt.Sprintf("您的个人信息反追踪能力已提升 %d %s，过去24小时内已为您拦截%d条，点击查看详情", int(price), "%", count.PreventNum)})
		payload := jpushclient.NewPushPayLoad()
		payload.SetPlatform(&pf)
		payload.SetAudience(&ad)
		payload.SetNotice(&notice)
		payload.SetOptions(&op)
		bytes, _ := payload.ToBytes()
		fmt.Printf("%s\r\n", string(bytes))
		c := jpushclient.NewPushClient(secretKey, appKey)
		str, err := c.Send(bytes)
		if err != nil {
			fmt.Printf("err:%s", err.Error())
		} else {
			fmt.Printf("ok:%s", str)
		}
	}
	return
}

//获取每日用户数据
func (u *User) GetUserDataWeekPush() (err error) {
	var (
		pf     jpushclient.Platform
		ad     jpushclient.Audience
		op     jpushclient.Option
		notice jpushclient.Notice
	)
	pf.Add(jpushclient.IOS)
	op.ApnsProduction = true
	// 获取用户id以及registerId列表
	err, ids := secret.GetUserIdAndRegisterID()
	if err != nil {
		logging.Error(err)
	}
	for _, v := range ids {
		ad.SetID([]string{v.RegisterId})
		notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: "您本周的安全报告已生成，请点击查看详情"})
		payload := jpushclient.NewPushPayLoad()
		payload.SetPlatform(&pf)
		payload.SetAudience(&ad)
		payload.SetNotice(&notice)
		payload.SetOptions(&op)
		bytes, _ := payload.ToBytes()
		fmt.Printf("%s\r\n", string(bytes))
		c := jpushclient.NewPushClient(secretKey, appKey)
		str, err := c.Send(bytes)
		if err != nil {
			fmt.Printf("err:%s", err.Error())
		} else {
			fmt.Printf("ok:%s", str)
		}
	}
	return
}
