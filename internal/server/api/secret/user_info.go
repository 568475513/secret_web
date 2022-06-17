package secret

import (
	"abs/internal/server/repository/prevent"
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
