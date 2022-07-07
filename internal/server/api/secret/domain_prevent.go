package secret

import (
	"abs/internal/server/repository/prevent"
	"abs/internal/server/rules/validator"
	secret "abs/models/secret"
	"abs/pkg/app"
	"abs/pkg/enums"
	"abs/pkg/logging"
	"abs/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"time"
)

var GoCache *cache.Cache

func init() {
	GoCache = cache.New(cache.NoExpiration, 60*time.Second)
}

//用户拦截信息接口
func DomainPrevent(c *gin.Context) {
	var (
		err error
		req validator.DomainPreventRule
		u   prevent.U
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	u.UserId = req.UserId
	u.UserIp = req.UserIp
	u.DomainType = req.DomainType
	u.Domain = req.Domain
	u.DomainTag = req.DomainTag
	u.DomainSource = req.DomainSource
	u.DomainSourceInfo = req.DomainSourceInfo
	u.RiskLevel = req.RiskLevel
	u.IsPrevent = req.IsPrevent
	err = u.InsertUserPreventInfo()
	if err != nil {
		app.FailWithMessage("录入用户数据异常", enums.ERROR, c)
	}
	app.OkWithMessage("ok", c)
}

//违规收集类APP推送接口
func CollectPush(c *gin.Context) {
	var (
		err error
		req validator.CollectPushRule
	)
	if err = app.ParseRequest(c, &req); err != nil {
		return
	}
	//如果用户拦截到违规收集app数据则发送推送
	if !GetUserFirstPushV(req.UserId, req.DomainTag) {
		err = GoCache.Add(req.UserId+":"+req.DomainTag, 1, cache.NoExpiration)
		if err != nil {
			logging.Error(err)
		}
		//获取用户配置信息
		c, err := secret.GetUserConfig(req.UserId)
		if err != nil {
			logging.Error(err)
		}
		Isbuy := false
		if c.IsBuy == 1 && c.ExpiredAt.Unix() > time.Now().Unix() {
			Isbuy = true
		}
		//获取用户注册id
		ui, err := secret.GetUserInfo(req.UserId, "")
		if err != nil {
			logging.Error(err)
		}
		if Isbuy {
			msg := "系统检测到" + req.DomainTag + "正在运行，该应用\n" +
				"曾被国家通报存在违规收集个人信息行为，\n" +
				"您已开启隐私安全模式，可以安全使用该应用。"
			err := util.SendPushMsg(ui.RegisterId, req.Url, msg)
			if err != nil {
				logging.Error(err)
			}
		} else {
			msg := "系统检测到" + req.DomainTag + "正在运行，该应用\n" +
				"曾被国家通报存在违规收集个人信息行为，\n" +
				"您可以开启隐私安全模式进行拦截，点击查看详情。"
			err := util.SendPushMsg(ui.RegisterId, req.Url, msg)
			if err != nil {
				logging.Error(err)
			}
		}
	}
	if err != nil {
		app.FailWithMessage("提送失败", enums.ERROR, c)
	}
	app.OkWithMessage("ok", c)
}

func GetUserFirstPushV(userId, domainTag string) (in bool) {

	r, t := GoCache.Get(userId + ":" + domainTag)
	if t == true && r != nil {
		if r == 1 {
			in = true
		}
	}
	return
}
