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
	err = user.UserBuyVip(req.UserId, req.ValidTime)
	if err != nil {
		app.FailWithMessage("用户购买会员失败", enums.ERROR, c)
	}
	app.OK(c)
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
				"AppName": "EasySpyApp",
				"AppType": "间谍软件",
				"AppRank": "3",
				"AppLink": "",
			},
			3: {
				"AppName": "麦苗守护",
				"AppType": "间谍软件",
				"AppRank": "4",
				"AppLink": "",
			},
		},
		2: {
			0: {
				"AppName": "HubStaff",
				"AppType": "监控员工软件",
				"AppRank": "1",
				"AppLink": "",
			},
			1: {
				"AppName": "Atto-Work",
				"AppType": "监控员工软件",
				"AppRank": "2",
				"AppLink": "",
			},
			3: {
				"AppName": "OnTheClock",
				"AppType": "监控员工软件",
				"AppRank": "3",
				"AppLink": "",
			},
		},
		3: {
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
		},
		4: {
			0: {
				"AppName": "美团",
				"AppType": "大数据滥收集软件",
				"AppRank": "1",
				"AppLink": "",
			},
			1: {
				"AppName": "淘宝",
				"AppType": "大数据滥收集软件",
				"AppRank": "1",
				"AppLink": "",
			},
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
