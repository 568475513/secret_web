package secret

import (
	"abs/internal/server/repository/prevent"
	"abs/internal/server/repository/user"
	"abs/internal/server/rules/validator"
	"abs/pkg/app"
	"abs/pkg/enums"
	"github.com/gin-gonic/gin"
)

//用户拦截信息接口
func UserPreventInfo(c *gin.Context) {
	var (
		err error
		req validator.SecretUserInfoRule
		u   prevent.U
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	u.UserId = req.UserId
	u.UserIp = req.UserIp
	u.DomainType = req.DomainType
	u.Page = req.Page
	u.PageSize = req.PageSize
	ps, err := u.GetPreventById()
	if err != nil {
		app.FailWithMessage("获取用户数据异常", enums.ERROR, c)
	}
	app.OkWithData(ps, c)
}

//用户拦截信息列表接口
func UserPreventInfoList(c *gin.Context) {
	var (
		err error
		req validator.SecretUserInfoListRule
		u   prevent.U
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	u.UserId = req.UserId
	u.HighRisk = req.HighRisk
	u.Page = req.Page
	u.PageSize = req.PageSize
	ps, err := u.GetPreventListById()
	if err != nil {
		app.FailWithMessage("获取用户数据异常", enums.ERROR, c)
	}
	app.OkWithData(ps, c)
}

//用户拦截信息分类接口
func UserPreventInfoClassify(c *gin.Context) {
	var (
		err error
		req validator.SecretUserInfoClassifyRule
		u   prevent.U
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	u.UserId = req.UserId
	ps, err := u.GetPreventClassifyById()
	if err != nil {
		app.FailWithMessage("获取用户数据异常", enums.ERROR, c)
	}
	app.OkWithData(ps, c)
}

//用户拦截信息分类详情接口
func UserPreventInfoClassifyDetail(c *gin.Context) {
	var (
		err error
		req validator.SecretUserInfoClassifyDetailRule
		u   prevent.U
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	u.UserId = req.UserId
	u.DomainTag = req.DomainTag
	ps, err := u.GetPreventClassifyDetailById()
	if err != nil {
		app.FailWithMessage("获取用户数据异常", enums.ERROR, c)
	}
	app.OkWithData(ps, c)
}

//用户拦截分类开关接口
func UserPreventClassifySwitch(c *gin.Context) {
	var (
		err error
		req validator.SecretUserClassifySwitchRule
		u   prevent.PreventSwitch
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	u.UserId = req.UserId
	u.IsBusMonitor = req.IsBusMonitor
	u.IsCollectInfo = req.IsCollectInfo
	u.IsSpy = req.IsSpy
	u.IsLargeData = req.IsLargeData
	err = u.UpdateUserPreventSwitch()
	if err != nil {
		app.FailWithMessage("更新用户开关数据异常", enums.ERROR, c)
	}
	app.OkWithData(u, c)
}

//获取用户配置信息列表接口
func GetUserConfigList(c *gin.Context) {
	var (
		err error
		u   user.UC
	)
	uc, err := u.GetUserConfList()
	if err != nil {
		app.FailWithMessage("获取用户配置", enums.ERROR, c)
	}
	app.OkWithData(uc, c)
}

//获取购买会员回调接口
func UserBuy(c *gin.Context) {
	var (
		err error
		req validator.UserBuyRule
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	t, err := user.UserBuyVip(req.UserId, req.ValidTime)
	if err != nil {
		app.FailWithMessage("用户购买会员失败", enums.ERROR, c)
	}
	app.OkWithData(t, c)
}

//用户拦截分类开关接口
func GetPreventList(c *gin.Context) {
	var (
		err error
		req validator.UserPreventListRule
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	s := map[int][]map[string]string{
		1: {
			//////////////////////////////////////////间谍开始////////////////////////////////////////
			0: {
				"AppName": "mSpy",
				"AppType": "间谍软件",
				"AppRank": "1",
				"AppLink": "",
			},
			1: {
				"AppName": "eyeZy",
				"AppType": "间谍软件",
				"AppRank": "2",
				"AppLink": "",
			},
			2: {
				"AppName": "Spyera",
				"AppType": "间谍软件",
				"AppRank": "3",
				"AppLink": "",
			},
			3: {
				"AppName": "FlexiSPY",
				"AppType": "间谍软件",
				"AppRank": "4",
				"AppLink": "",
			},
			4: {
				"AppName": "Qustodio",
				"AppType": "间谍软件",
				"AppRank": "5",
				"AppLink": "",
			},
			5: {
				"AppName": "SpyBubble",
				"AppType": "间谍软件",
				"AppRank": "6",
				"AppLink": "",
			},
			6: {
				"AppName": "TheWiSpy",
				"AppType": "间谍软件",
				"AppRank": "7",
				"AppLink": "",
			},
			7: {
				"AppName": "Spyic",
				"AppType": "间谍软件",
				"AppRank": "8",
				"AppLink": "",
			},
			8: {
				"AppName": "wondershare",
				"AppType": "间谍软件",
				"AppRank": "9",
				"AppLink": "",
			},
			9: {
				"AppName": "MobiStealth",
				"AppType": "间谍软件",
				"AppRank": "10",
				"AppLink": "",
			},
			10: {
				"AppName": "AppMia",
				"AppType": "间谍软件",
				"AppRank": "11",
				"AppLink": "",
			},
			11: {
				"AppName": "ESET Parental Control",
				"AppType": "间谍软件",
				"AppRank": "12",
				"AppLink": "",
			},
			12: {
				"AppName": "Highster Mobile",
				"AppType": "间谍软件",
				"AppRank": "13",
				"AppLink": "",
			},
			13: {
				"AppName": "Spyzie",
				"AppType": "间谍软件",
				"AppRank": "14",
				"AppLink": "",
			},
			14: {
				"AppName": "Notron Family",
				"AppType": "间谍软件",
				"AppRank": "15",
				"AppLink": "",
			},
			15: {
				"AppName": "EasySpyApp",
				"AppType": "间谍软件",
				"AppRank": "16",
				"AppLink": "",
			},
			////////////////////////////////////////间谍结束////////////////////////////////////////
		},
		2: {
			////////////////////////////////////////监控员工软件开始////////////////////////////////////////
			0: {
				"AppName": "Scalefusion",
				"AppType": "企业监控软件",
				"AppRank": "1",
				"AppLink": "",
			},
			1: {
				"AppName": "HubStaff",
				"AppType": "企业监控软件",
				"AppRank": "2",
				"AppLink": "",
			},
			2: {
				"AppName": "Workfolio",
				"AppType": "企业监控软件",
				"AppRank": "3",
				"AppLink": "",
			},
			3: {
				"AppName": "Atto-Work",
				"AppType": "企业监控软件",
				"AppRank": "4",
				"AppLink": "",
			},
			4: {
				"AppName": "OnTheClock",
				"AppType": "企业监控软件",
				"AppRank": "5",
				"AppLink": "",
			},
			5: {
				"AppName": "The Team Tracker",
				"AppType": "企业监控软件",
				"AppRank": "6",
				"AppLink": "",
			},
			6: {
				"AppName": "QuickBooks",
				"AppType": "企业监控软件",
				"AppRank": "7",
				"AppLink": "",
			},
			7: {
				"AppName": "Homebase",
				"AppType": "企业监控软件",
				"AppRank": "8",
				"AppLink": "",
			},
			8: {
				"AppName": "Employee Tracker Lite",
				"AppType": "企业监控软件",
				"AppRank": "9",
				"AppLink": "",
			},
			9: {
				"AppName": "ezClocker-Employee",
				"AppType": "企业监控软件",
				"AppRank": "10",
				"AppLink": "",
			},
			10: {
				"AppName": "Justworks Hours",
				"AppType": "企业监控软件",
				"AppRank": "11",
				"AppLink": "",
			},
			11: {
				"AppName": "Sling-Employee",
				"AppType": "企业监控软件",
				"AppRank": "12",
				"AppLink": "",
			},
			12: {
				"AppName": "Humanity-Employee",
				"AppType": "企业监控软件",
				"AppRank": "13",
				"AppLink": "",
			},
			13: {
				"AppName": "Deputy-Shift",
				"AppType": "企业监控软件",
				"AppRank": "14",
				"AppLink": "",
			},
			14: {
				"AppName": "Employee Link-Timesheet",
				"AppType": "企业监控软件",
				"AppRank": "15",
				"AppLink": "",
			},
			15: {
				"AppName": "busybusy",
				"AppType": "企业监控软件",
				"AppRank": "16",
				"AppLink": "",
			},
			////////////////////////////////////////监控员工软件结束////////////////////////////////////////
		},
		3: {
			////////////////////////////////////////违规收集软件开始////////////////////////////////////////
			0: {
				"AppName": "爱卡汽车",
				"AppType": "违规收集软件",
				"AppRank": "1",
				"AppLink": "",
			},
			1: {
				"AppName": "营创书院",
				"AppType": "违规收集软件",
				"AppRank": "1",
				"AppLink": "",
			},
			////////////////////////////////////////违规收集软件结束////////////////////////////////////////
		},
		4: {
			////////////////////////////////////////大数据滥收集软件开始////////////////////////////////////////
			0: {
				"AppName": "Google统计",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			1: {
				"AppName": "友盟统计",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			2: {
				"AppName": "极光统计",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			3: {
				"AppName": "MobTech",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			4: {
				"AppName": "GrowingIO",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			5: {
				"AppName": "TalkingData",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			6: {
				"AppName": "阿里妈妈",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			7: {
				"AppName": "字节跳动",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			8: {
				"AppName": "快手",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			9: {
				"AppName": "百度",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			10: {
				"AppName": "腾讯",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			11: {
				"AppName": "美团",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			12: {
				"AppName": "淘宝",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			13: {
				"AppName": "今日头条",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			14: {
				"AppName": "苏宁易购",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			15: {
				"AppName": "真快乐",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			16: {
				"AppName": "京东",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			17: {
				"AppName": "拼多多",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			18: {
				"AppName": "58同城",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			19: {
				"AppName": "携程旅行",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			20: {
				"AppName": "同程旅行",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			21: {
				"AppName": "知乎",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			22: {
				"AppName": "新浪微博",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			23: {
				"AppName": "唯品会",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			24: {
				"AppName": "得物",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			25: {
				"AppName": "Bilibili",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			26: {
				"AppName": "西瓜视频",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			27: {
				"AppName": "小红书",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			28: {
				"AppName": "探探",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			29: {
				"AppName": "陌陌",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			30: {
				"AppName": "Soul",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			31: {
				"AppName": "链家",
				"AppType": "大数据滥收集软件",
				"AppRank": "0",
				"AppLink": "",
			},
			////////////////////////////////////////大数据滥收集软件结束////////////////////////////////////////
		},
	}
	app.OkWithData(s[req.PrevemtType], c)
}

//用户投诉接口
func UserComplain(c *gin.Context) {
	var (
		err error
		req validator.UserComplainRule
		u   user.UserComplain
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	u.UserId = req.UserId
	u.ComplainType = req.ComplainType
	u.ComplainMsg = req.ComplainMsg
	u.ComplainContact = req.ComplainContact
	err = u.InsertUserComplainData()
	if err != nil {
		app.FailWithMessage("获取用户配置", enums.ERROR, c)
	}
	app.OK(c)
}
