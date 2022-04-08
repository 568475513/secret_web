package secret

import (
	"abs/internal/server/repository/prevent"
	"abs/internal/server/rules/validator"
	"abs/pkg/app"
	"abs/pkg/enums"
	"github.com/gin-gonic/gin"
)

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
	err = u.InsertUserPreventInfo()
	if err != nil {
		app.FailWithMessage("录入用户数据异常", enums.ERROR, c)
	}
	app.OkWithMessage("ok", c)
}
