package user

import (
	secret "abs/models/secret"
	"abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
	"fmt"
	cache "github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	jpushclient "github.com/ylywyn/jpush-api-go-client"
	"math"
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
	Count                int    `json:"count"`
	CreatedAt            int    `json:"created_at"`
	UserConfig           UC     `json:"user_config"`
}

type UC struct {
	IsBuy         int    `json:"is_buy"`
	IsBusMonitor  int    `json:"is_bus_monitor"`
	IsLargeData   int    `json:"is_large_data"`
	IsSpy         int    `json:"is_spy"`
	IsCollectInfo int    `json:"is_collect_info"`
	ExpiredAt     string `json:"expired_at"`
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

type UserComplain struct {
	UserId          string `json:"user_id"`
	ComplainType    int    `json:"complain_type"`
	ComplainMsg     string `json:"complain_msg"`
	ComplainContact string `json:"complain_contact"`
}

type Expired struct {
	ExpiredAt string `json:"expired_at"`
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
	u.UserDnsPreventDomain = "https://test.privacy.prisecurity.com/dns-query/" + u.UserId
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
	u.CreatedAt = int(math.Ceil(float64(time.Now().Unix()-ui.CreatedAt.Unix()) / float64(24*60*60)))
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

	//获取用户配置
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
	u.UserConfig.ExpiredAt = uc.ExpiredAt.Format("2006-01-02 15:04:05")

	//获取用户拦截天数以及拦截次数
	c, err := secret.GetCountByUserId(u.UserId)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	u.Count = c.PreventNum

	return u, nil
}

//获取所有用户配置信息
func (u *UC) GetUserConfList() (map[string][]string, error) {

	ui, err := secret.GetAllUserConfigList()
	if err != nil || ui == nil {
		logging.Error(err)
		return nil, err
	}
	return ui, nil
}

//用户购买会员回调
func UserBuyVip(userId string, year int) (ti Expired, err error) {

	e := Expired{}
	t := time.Time{}
	s, err := secret.GetUserConfig(userId)
	if err != nil {
		logging.Error(err)
		return e, err
	}
	if time.Now().Unix() > s.ExpiredAt.Unix() {
		t = time.Now().AddDate(year, 0, 0)
	} else {
		t = s.ExpiredAt.AddDate(year, 0, 0)
	}
	err = secret.UpdateUserVipExpired(userId, t)
	if err != nil {
		logging.Error(err)
		return e, err
	}
	e.ExpiredAt = t.Format("2006-01-02 15:04:05")
	return e, nil
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
func (u *User) GetUserPriceDay2() (err error) {

	// 获取用户id以及registerId列表
	err, ids := secret.GetUidRid()
	if err != nil {
		logging.Error(err)
	}
	for k, v := range ids {
		c, err := secret.GetUserConfig(k)
		if err != nil {
			logging.Error(err)
		}
		Isbuy := false
		if c.IsBuy == 1 && c.ExpiredAt.Unix() > time.Now().Unix() {
			Isbuy = true
		}
		d, err := secret.GetPreventCountByUserIdLastDay(k)
		if err != nil {
			logging.Error(err)
		}
		count := 0
		spy := 0
		buz := 0
		col := 0
		lar := 0
		if Isbuy {
			for k, v := range d {
				count = count + v
				if k == enums.IsSpyCode {
					spy = v
				}
				if k == enums.IsBusMonitorCode {
					buz = v
				}
				if k == enums.IsCollectInfoCode {
					col = v
				}
				if k == enums.IsLargeDataCode {
					lar = v
				}
			}
			msg := fmt.Sprintf("过去24小时内已为您拦截%d条隐私监控行为，\n"+
				"其中间谍软件拦截%d条，企业监控拦截%d条，\n"+
				"违规收集信息拦截%d条，大数据滥收集拦截%d条\n"+
				"其中拦截最多的应用是mSpy，点击查看详情", count, spy, buz, col, lar)
			err := util.SendPushMsg(v.RegisterId, "", msg)
			if err != nil {
				logging.Error(err)
			}
		} else {
			for k, v := range d {
				count = count + v
				if k == enums.IsSpyCode {
					spy = v
				}
				if k == enums.IsBusMonitorCode {
					buz = v
				}
				if k == enums.IsCollectInfoCode {
					col = v
				}
				if k == enums.IsLargeDataCode {
					lar = v
				}
			}
			msg := fmt.Sprintf("过去24小时内已发现%d条隐私监控行为，\n"+
				"其中间谍软件%d条，企业监控%d条，\n"+
				"违规收集信息%d条，大数据滥收集%d条，\n"+
				"您可以开启隐私安全模式进行拦截，点击查看详情", count, spy, buz, col, lar)
			err := util.SendPushMsg(v.RegisterId, "", msg)
			if err != nil {
				logging.Error(err)
			}
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

//获取所有用户配置信息
func (u *UserComplain) InsertUserComplainData() error {

	err := secret.InsertUserComplain(u.UserId, u.ComplainMsg, u.ComplainContact, u.ComplainType)
	if err != nil {
		logging.Error(err)
		return err
	}
	return nil
}
